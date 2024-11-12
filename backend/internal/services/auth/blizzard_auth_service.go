// Package auth provides authentication services for the application
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	models "wowperf/internal/models"
)

// BlizzardScopes are the scopes available for Battle.net OAuth2 authentication.
const (
	ScopeOpenID        = "openid"
	ScopeWowProfile    = "wow.profile"
	ScopeWowCharacters = "wow.profile.characters"
)

// RequiredScopes are the scopes required for Battle.net OAuth2 authentication.
var RequiredScopes = []string{
	ScopeOpenID,
	ScopeWowProfile,
}

// BlizzardAuthConfig holds the configuration for Battle.net authentication.
// All fields are required for the service to function properly.
type BlizzardAuthConfig struct {
	ClientID     string // OAuth client ID from Battle.net developer portal
	ClientSecret string // OAuth client secret from Battle.net developer portal
	RedirectURL  string // OAuth redirect URL registered in Battle.net developer portal
	Region       string // Battle.net region (e.g., "eu", "us", etc.)
}

// BlizzardAuthService handles all Battle.net OAuth2 authentication operations.
// It manages token exchange, refresh, and Battle.net API interactions.
type BlizzardAuthService struct {
	db          *gorm.DB       // Database connection for user data persistence
	oauthConfig *oauth2.Config // OAuth2 configuration for Battle.net
	region      string         // Battle.net region for API calls
}

// BattleNetUserInfo represents the user information returned by Battle.net's
// userinfo endpoint. This struct maps directly to the JSON response.
type BattleNetUserInfo struct {
	Sub       string `json:"sub"`       // Unique identifier for the user
	ID        int    `json:"id"`        // Battle.net account ID
	BattleTag string `json:"battletag"` // User's BattleTag
}

// BattleNetTokenInfo represents the OAuth token information received from Battle.net.
// It includes both the token data and its expiration details.
type BattleNetTokenInfo struct {
	AccessToken  string    `json:"access_token"`  // The OAuth2 access token
	TokenType    string    `json:"token_type"`    // Token type (usually "Bearer")
	ExpiresIn    int       `json:"expires_in"`    // Token lifetime in seconds
	RefreshToken string    `json:"refresh_token"` // Token used to refresh access token
	Scope        string    `json:"scope"`         // Granted OAuth scopes
	ExpiresAt    time.Time `json:"expires_at"`    // Calculated token expiration time
}

// NewBlizzardAuthService creates and initializes a new Battle.net authentication service.
// It sets up the OAuth2 configuration with the provided parameters.
//
// Parameters:
//   - db: A pointer to the GORM database connection
//   - config: BlizzardAuthConfig containing OAuth credentials and settings
//
// Returns:
//   - A pointer to the initialized BlizzardAuthService
func NewBlizzardAuthService(db *gorm.DB, config BlizzardAuthConfig) *BlizzardAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       RequiredScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.battle.net/authorize",
			TokenURL: "https://oauth.battle.net/token",
		},
	}

	return &BlizzardAuthService{
		db:          db,
		oauthConfig: oauthConfig,
		region:      config.Region,
	}
}

// GetAuthorizationURL generates the Battle.net OAuth authorization URL.
// The state parameter is used to prevent CSRF attacks.
//
// Parameters:
//   - state: A randomly generated string to verify the OAuth callback
//
// Returns:
//   - The complete authorization URL to redirect the user to
func (s *BlizzardAuthService) GetAuthorizationURL(state string) string {
	s.oauthConfig.Scopes = RequiredScopes
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCodeForToken exchanges an authorization code for an OAuth token.
// This method should be called after receiving the authorization code from Battle.net.
//
// Parameters:
//   - ctx: Context for the request
//   - code: The authorization code received from Battle.net
//
// Returns:
//   - Token: The OAuth2 token if successful
//   - error: Any error that occurred during the exchange
func (s *BlizzardAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	log.Printf("Exchanging code for token with redirect_uri: %s", s.oauthConfig.RedirectURL)

	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, nil
}

// GetUserInfo retrieves the user's Battle.net account information using their OAuth token.
//
// Parameters:
//   - token: A valid OAuth2 token
//
// Returns:
//   - BattleNetUserInfo: User's Battle.net account information
//   - error: Any error that occurred while fetching the information
func (s *BlizzardAuthService) GetUserInfo(token *oauth2.Token) (*BattleNetUserInfo, error) {
	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://oauth.battle.net/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status code %d", resp.StatusCode)
	}

	var userInfo BattleNetUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// RefreshToken refreshes an expired OAuth token using its refresh token.
//
// Parameters:
//   - ctx: Context for the request
//   - refreshToken: The refresh token to use
//
// Returns:
//   - Token: The new OAuth2 token
//   - error: Any error that occurred during refresh
func (s *BlizzardAuthService) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := s.oauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

// UpdateUserBattleNetTokens updates the Battle.net tokens for a user in the database.
// It encrypts the access token before storage.
//
// Parameters:
//   - user: The user whose tokens need updating
//   - token: The new OAuth2 token to store
//
// Returns:
//   - error: Any error that occurred during the update
func (s *BlizzardAuthService) UpdateUserBattleNetTokens(user *models.User, token *oauth2.Token) error {
	if err := user.SetBattleNetToken(token.AccessToken); err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	// Extract scopes from the token
	scope, ok := token.Extra("scope").(string)
	var scopes pq.StringArray
	if ok {
		scopes = pq.StringArray(strings.Split(scope, " "))
	}

	updates := map[string]interface{}{
		"battle_net_refresh_token": token.RefreshToken,
		"battle_net_expires_at":    token.Expiry,
		"battle_net_token_type":    token.TokenType,
		"battle_net_scopes":        scopes,
	}

	if err := s.db.Model(user).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user battle.net tokens: %w", err)
	}

	return nil
}

// LinkBattleNetAccount links a Battle.net account to an existing user account.
// This operation is performed within a transaction to ensure data consistency.
//
// Parameters:
//   - userID: The ID of the user to link
//   - battleNetInfo: Battle.net account information
//   - token: The OAuth2 token to store
//
// Returns:
//   - error: Any error that occurred during the linking process
func (s *BlizzardAuthService) LinkBattleNetAccount(userID uint, battleNetInfo *BattleNetUserInfo, token *oauth2.Token) error {
	log.Printf("Starting Battle.net account linking for userID: %d, BattleTag: %s", userID, battleNetInfo.BattleTag)

	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Vérification du BattleTag
	var existingUser models.User
	err := tx.Where("battle_tag = ? AND id != ?", battleNetInfo.BattleTag, userID).First(&existingUser).Error

	if err == nil {
		tx.Rollback()
		return fmt.Errorf("battle tag already linked to another account")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("failed to check battle tag uniqueness: %w", err)
	}

	log.Printf("BattleTag %s is available for linking", battleNetInfo.BattleTag)

	// Update Battle.net information
	updates := map[string]interface{}{
		"battle_net_id": battleNetInfo.ID,
		"battle_tag":    battleNetInfo.BattleTag,
	}

	log.Printf("Updating user information with: %+v", updates)
	if err := tx.Model(&user).Updates(updates).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to update user information: %v", err)
		return fmt.Errorf("failed to update user battle.net info: %w", err)
	}

	// Update tokens
	if err := s.UpdateUserBattleNetTokens(&user, token); err != nil {
		tx.Rollback()
		log.Printf("Failed to update tokens: %v", err)
		return fmt.Errorf("failed to update user battle.net tokens: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully linked Battle.net account for user ID %d", userID)
	return nil
}

// ValidateToken checks if a token is valid and not expired.
// It performs both expiration time checking and a test API call.
//
// Parameters:
//   - token: The OAuth2 token to validate
//
// Returns:
//   - bool: True if the token is valid, false otherwise
func (s *BlizzardAuthService) ValidateToken(token *oauth2.Token) bool {
	if token == nil {
		return false
	}

	// Check if token is expired
	if token.Expiry.Before(time.Now()) {
		return false
	}

	// Verify token with an API call
	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://oauth.battle.net/userinfo")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// ValidateScopes checks if the token has the required scopes
func (s *BlizzardAuthService) ValidateScopes(user *models.User) bool {
	return user.HasRequiredScopes(RequiredScopes)
}

// CompleteOAuthFlow manages the entire OAuth process in a single function
func (s *BlizzardAuthService) CompleteOAuthFlow(ctx context.Context, userID uint, code string) (string, error) {
	// 1. Exchange code for a token
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	// 2. Get user information
	userInfo, err := s.GetUserInfo(token)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	// 3. Begin a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// 4. Verify if the BattleTag is already linked
	var existingUser models.User
	err = tx.Where("battle_tag = ? AND id != ?", userInfo.BattleTag, userID).First(&existingUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return "", fmt.Errorf("failed to check battletag uniqueness: %w", err)
	}
	if err == nil {
		tx.Rollback()
		return "", fmt.Errorf("battletag already linked to another account")
	}

	// 5. Update user
	updates := map[string]interface{}{
		"battle_net_id":            userInfo.ID,
		"battle_tag":               userInfo.BattleTag,
		"battle_net_refresh_token": token.RefreshToken,
		"battle_net_expires_at":    token.Expiry,
		"battle_net_token_type":    token.TokenType,
		"battle_net_scopes":        strings.Split(token.TokenType, " "),
	}

	if err := tx.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	// 6. Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userInfo.BattleTag, nil
}
