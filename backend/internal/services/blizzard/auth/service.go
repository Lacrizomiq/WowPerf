package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"wowperf/internal/models"

	"encoding/json"

	"github.com/go-redis/redis/v8"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

const (
	// API endpoints
	authEndpoint      = "https://oauth.battle.net"
	tokenEndpoint     = authEndpoint + "/token"
	authorizeEndpoint = authEndpoint + "/authorize"
	userInfoEndpoint  = authEndpoint + "/userinfo"

	// Token settings
	TokenRefreshThreshold = 5 * time.Minute // Refresh token if expiring within 5 minutes
)

// BattleNetUserInfo represents the user info returned from the Battle.net API
type BattleNetUserInfo struct {
	Sub       string `json:"sub"`
	ID        int    `json:"id"`
	BattleTag string `json:"battletag"`
}

// BattleNetAuthService handles all Battle.net OAuth2 authentication operations
type BattleNetAuthService struct {
	db          *gorm.DB
	oauthConfig *oauth2.Config
	region      string
	redisClient *redis.Client
}

// NewBattleNetAuthService creates a new Battle.net authentication service
func NewBattleNetAuthService(db *gorm.DB, redisClient *redis.Client) (*BattleNetAuthService, error) {
	// Create OAuth2 config with Battle.net specific settings
	config := &oauth2.Config{
		ClientID:     os.Getenv("BLIZZARD_CLIENT_ID"),
		ClientSecret: os.Getenv("BLIZZARD_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("BLIZZARD_REDIRECT_URL"),
		Scopes: []string{
			"openid",      // Required for user info
			"wow.profile", // Required for WoW profile access
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizeEndpoint,
			TokenURL: tokenEndpoint,
		},
	}

	return &BattleNetAuthService{
		db:          db,
		oauthConfig: config,
		region:      os.Getenv("BLIZZARD_REGION"),
		redisClient: redisClient,
	}, nil
}

// storeOAuthState stores the OAuth state in Redis
func (s *BattleNetAuthService) storeOAuthState(ctx context.Context, state *OAuthState) error {
	// Generate Redis key
	redisKey := fmt.Sprintf("oauth_state:%s", state.State)

	// Marshal state to JSON
	stateJSON, err := state.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Store in Redis with expiration
	duration := time.Until(state.ExpiresAt)
	if err := s.redisClient.Set(ctx, redisKey, stateJSON, duration).Err(); err != nil {
		return fmt.Errorf("failed to store OAuth state: %w", err)
	}

	return nil
}

// getAndValidateOAuthState retrieves and validates the OAuth state from Redis
func (s *BattleNetAuthService) getAndValidateOAuthState(ctx context.Context, stateParam string) (*OAuthState, error) {
	// Get state from Redis
	redisKey := fmt.Sprintf("oauth_state:%s", stateParam)
	stateJSON, err := s.redisClient.Get(ctx, redisKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve OAuth state: %w", err)
	}

	// Delete state immediately to prevent replay attacks
	s.redisClient.Del(ctx, redisKey)

	// Unmarshal and validate state
	state, err := UnmarshalOAuthState(stateJSON)
	if err != nil {
		return nil, fmt.Errorf("invalid OAuth state: %w", err)
	}

	return state, nil
}

// InitiateAuth starts the OAuth2 flow by generating the authorization URL
func (s *BattleNetAuthService) InitiateAuth(ctx context.Context, userID uint) (string, error) {
	// Generate a new OAuth state
	state, err := NewOAuthState(userID)
	if err != nil {
		return "", fmt.Errorf("failed to create OAuth state: %w", err)
	}

	// Store the state in Redis
	if err := s.storeOAuthState(ctx, state); err != nil {
		return "", fmt.Errorf("failed to store OAuth state: %w", err)
	}

	// Generate the authorization URL
	return s.oauthConfig.AuthCodeURL(state.State), nil
}

// ExchangeCodeForToken exchanges an authorization code for access and refresh tokens
func (s *BattleNetAuthService) ExchangeCodeForToken(ctx context.Context, code, stateParam string) (*oauth2.Token, uint, error) {
	// Verify and get the state
	state, err := s.getAndValidateOAuthState(ctx, stateParam)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid OAuth state: %w", err)
	}

	// Exchange code for token
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, state.UserID, nil
}

// GetUserInfo retrieves the user's Battle.net profile information
func (s *BattleNetAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*BattleNetUserInfo, error) {
	// Create an HTTP client with the OAuth2 token
	client := s.oauthConfig.Client(ctx, token)

	// Make request to userinfo endpoint
	resp, err := client.Get(userInfoEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	// Decode the response
	var userInfo BattleNetUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// LinkUserAccount links a Battle.net account to a user account
func (s *BattleNetAuthService) LinkUserAccount(ctx context.Context, token *oauth2.Token, userID string) error {
	log.Printf("Starting LinkUserAccount for userID=%s", userID)

	tx := s.db.Begin()
	if tx.Error != nil {
		log.Printf("Transaction start failed: %v", tx.Error)
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in LinkUserAccount: %v", r)
			tx.Rollback()
		}
	}()

	// Get Battle.net user info
	userInfo, err := s.GetUserInfo(ctx, token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		tx.Rollback()
		return fmt.Errorf("failed to get user info: %w", err)
	}
	log.Printf("Got user info: battleTag=%s, id=%d", userInfo.BattleTag, userInfo.ID)

	battleNetID := fmt.Sprintf("%d", userInfo.ID)

	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		log.Printf("User not found: %v", err)
		tx.Rollback()
		return fmt.Errorf("user not found: %w", err)
	}

	// Chiffrer et sauvegarder uniquement l'access token
	if err := user.SetBattleNetTokens(token.AccessToken, ""); err != nil {
		log.Printf("Failed to set Battle.net token: %v", err)
		tx.Rollback()
		return fmt.Errorf("failed to encrypt token: %w", err)
	}

	// Mise à jour des champs
	user.BattleNetID = battleNetID
	user.BattleTag = userInfo.BattleTag
	user.BattleNetExpiresAt = token.Expiry
	user.BattleNetTokenType = token.TokenType

	// Sauvegarder l'utilisateur
	if err := tx.Save(&user).Error; err != nil {
		log.Printf("Failed to save user: %v", err)
		tx.Rollback()
		return fmt.Errorf("failed to save user: %w", err)
	}

	log.Printf("Successfully linked account and updated user %s", userID)
	return tx.Commit().Error
}

// UnlinkUserAccount removes the Battle.net account link from a user
func (s *BattleNetAuthService) UnlinkUserAccount(ctx context.Context, userID uint) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := map[string]interface{}{
		"battle_net_id":           nil,
		"battle_tag":              nil,
		"encrypted_access_token":  nil,
		"encrypted_refresh_token": nil,
		"battle_net_token_type":   nil,
		"battle_net_expires_at":   nil,
		"battle_net_scopes":       nil,
	}

	if err := tx.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to unlink account: %w", err)
	}

	return tx.Commit().Error
}

// RefreshToken refreshes an expired access token using the refresh token
func (s *BattleNetAuthService) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	// Use OAuth2 config to get a new token
	source := s.oauthConfig.TokenSource(ctx, token)
	newToken, err := source.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

// ValidateToken checks if a token is valid and not expired
func (s *BattleNetAuthService) ValidateToken(ctx context.Context, token *oauth2.Token) error {
	if token == nil || token.AccessToken == "" {
		return fmt.Errorf("invalid token")
	}

	// Check expiration
	if token.Expiry.Before(time.Now()) {
		return fmt.Errorf("token expired")
	}

	// Verify token by making a request to userinfo
	_, err := s.GetUserInfo(ctx, token)
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	return nil
}

// GetUserToken returns a valid OAuth2 token for a user
func (s *BattleNetAuthService) GetUserToken(ctx context.Context, userID uint) (*oauth2.Token, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user has Battle.net linked
	if !user.IsBattleNetLinked() {
		return nil, fmt.Errorf("battle.net account not linked")
	}

	// Get decrypted access token
	accessToken, err := user.GetBattleNetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	token := &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   user.BattleNetTokenType,
		Expiry:      user.BattleNetExpiresAt,
	}

	// Check if token needs refresh
	if time.Until(token.Expiry) < TokenRefreshThreshold {
		refreshToken, err := user.GetBattleNetRefreshToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get refresh token: %w", err)
		}

		newToken, err := s.RefreshToken(ctx, refreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// Update tokens in database
		if err := s.UpdateUserToken(userID, newToken); err != nil {
			return nil, fmt.Errorf("failed to update token: %w", err)
		}

		return newToken, nil
	}

	return token, nil
}

// UpdateUserToken updates the stored tokens for a user
func (s *BattleNetAuthService) UpdateUserToken(userID uint, token *oauth2.Token) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Encrypt and store new tokens
	if err := user.SetBattleNetTokens(token.AccessToken, token.RefreshToken); err != nil {
		return fmt.Errorf("failed to encrypt tokens: %w", err)
	}

	updates := map[string]interface{}{
		"battle_net_expires_at": token.Expiry,
		"battle_net_token_type": token.TokenType,
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user token: %w", err)
	}

	return nil
}

// GetUserBattleNetStatus returns the current Battle.net link status for a user
func (s *BattleNetAuthService) GetUserBattleNetStatus(ctx context.Context, userID uint) (*BattleNetStatus, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &BattleNetStatus{
		Linked:    user.IsBattleNetLinked(),
		BattleTag: user.BattleTag,
	}, nil
}
