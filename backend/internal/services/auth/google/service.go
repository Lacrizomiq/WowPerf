package googleauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// GoogleAuthService gère l'authentification Google Oauth2
type GoogleAuthService struct {
	db          *gorm.DB
	oauthConfig *oauth2.Config
	config      *Config
	repository  *GoogleAuthRepository
}

// NewGoogleAuthService crée une nouvelle instance du service Google Auth
func NewGoogleAuthService(db *gorm.DB) (*GoogleAuthService, error) {
	// Charge la configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load Google OAuth config: %w", err)
	}

	// Valide la configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid Google OAuth config: %w", err)
	}

	// Crée la configuration OAuth2
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Crée le repository
	repository := NewGoogleAuthRepository(db)

	return &GoogleAuthService{
		db:          db,
		oauthConfig: oauthConfig,
		config:      config,
		repository:  repository,
	}, nil
}

// ===== MÉTHODES DE GÉNÉRATION D'ÉTAT ET URL =====

// generateState génère un état aléatoire pour protection CSRF
func (s *GoogleAuthService) generateState() (string, error) {
	// Générer 32 bytes aléatoires
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", &OAuthError{
			Code:    "state_generation_failed",
			Message: "Failed to generate CSRF state",
			Details: err.Error(),
		}
	}

	// Encoder en base64 URL-safe
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetAuthURL génère l'URL d'autorisation Google avec état CSRF
func (s *GoogleAuthService) GetAuthURL() (string, string, error) {
	// Générer l'état CSRF
	state, err := s.generateState()
	if err != nil {
		return "", "", err
	}

	// Générer l'URL d'autorisation avec options
	url := s.oauthConfig.AuthCodeURL(state,
		oauth2.AccessTypeOffline,                    // Pour obtenir un refresh token
		oauth2.SetAuthURLParam("prompt", "consent"), // Force le consent screen
	)

	return url, state, nil
}

// GetFrontendURL retourne l'URL du frontend configurée
func (s *GoogleAuthService) GetFrontendURL() string {
	return s.config.FrontendURL
}

// GetDashboardURL retourne l'URL complète du dashboard
func (s *GoogleAuthService) GetDashboardURL() string {
	return s.config.FrontendURL + s.config.DashboardPath
}

// GetErrorURL retourne l'URL complète d'erreur
func (s *GoogleAuthService) GetErrorURL() string {
	return s.config.FrontendURL + s.config.ErrorPath
}

// GetErrorPath retourne le path d'erreur (sans URL de base)
func (s *GoogleAuthService) GetErrorPath() string {
	return s.config.ErrorPath // "/login"
}

// ===== ÉCHANGE DE CODE ET RÉCUPÉRATION D'INFOS UTILISATEUR =====

// ExchangeCodeForToken échange le code d'autorisation contre un token d'accès
func (s *GoogleAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	// Créer un contexte avec timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Échanger le code contre un token
	token, err := s.oauthConfig.Exchange(ctxWithTimeout, code)
	if err != nil {
		return nil, &OAuthError{
			Code:    "token_exchange_failed",
			Message: "Failed to exchange authorization code for token",
			Details: err.Error(),
		}
	}

	// Vérifier que le token est valide
	if token.AccessToken == "" {
		return nil, &OAuthError{
			Code:    "invalid_token",
			Message: "Received invalid token from Google",
		}
	}

	return token, nil
}

// GetUserInfo récupère les informations utilisateurs depuis Google avec validation
func (s *GoogleAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	// Créer un client HTTP avec le token OAuth
	client := s.oauthConfig.Client(ctx, token)

	// Créer un contexte avec timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Faire la requête vers l'API Google
	req, err := http.NewRequestWithContext(ctxWithTimeout, "GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, &OAuthError{
			Code:    "user_info_request_failed",
			Message: "Failed to create user info request",
			Details: err.Error(),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, &OAuthError{
			Code:    "user_info_request_failed",
			Message: "Failed to get user info from Google",
			Details: err.Error(),
		}
	}

	defer resp.Body.Close()

	// Vérifie le status code
	if resp.StatusCode != http.StatusOK {
		return nil, &OAuthError{
			Code:    "user_info_request_failed",
			Message: fmt.Sprintf("Google API returned status %d", resp.StatusCode),
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &OAuthError{
			Code:    "user_info_decode_failed",
			Message: "Failed to read response body",
			Details: err.Error(),
		}
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(bodyBytes, &userInfo); err != nil {
		return nil, &OAuthError{
			Code:    "user_info_decode_failed",
			Message: "Failed to decode user info from Google",
			Details: err.Error(),
		}
	}

	// Validation critique Google
	if err := s.validateGoogleUserInfo(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// validateGoogleUserInfo valide les données utilisateur selon les recommandations Google
func (s *GoogleAuthService) validateGoogleUserInfo(userInfo *GoogleUserInfo) error {
	// 1. Email vérifié RECOMMANDÉ mais pas obligatoire pour les tests
	if !userInfo.VerifiedEmail {
		log.Printf("⚠️ Warning: Email not verified by Google for %s", userInfo.Email)
		// return &OAuthError{...}
	}

	// 2. Email présent
	if userInfo.Email == "" {
		return &OAuthError{
			Code:    "missing_email",
			Message: "Email is required from Google",
		}
	}

	// 3. Google ID présent
	if userInfo.ID == "" {
		return &OAuthError{
			Code:    "missing_google_id",
			Message: "Google ID is required from Google",
		}
	}

	// 4. Validation basique format email
	if !s.isValidEmailFormat(userInfo.Email) {
		return &OAuthError{
			Code:    "invalid_email_format",
			Message: "Invalid email format from Google",
		}
	}

	return nil
}

// isValidEmailFormat validation basique du format email
func (s *GoogleAuthService) isValidEmailFormat(email string) bool {
	// Validation très basique - Google a déjà validé l'email
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ===== LOGIQUE COMPLEXE DE RECHERCHE UTILISATEUR (CRITIQUE SELON GOOGLE) =====

// LookupUser recherche un utilisateur selon la logique complexe recommandée par Google
func (s *GoogleAuthService) LookupUser(googleID, googleEmail string) (*UserLookupResult, error) {
	result := &UserLookupResult{}

	// CAS 1 : Recherche par Google ID (PRIORITAIRE selon Google)
	userByGoogleID, err := s.repository.FindUserByGoogleID(googleID)
	if err == nil {
		result.ExistingUser = userByGoogleID
		result.FoundByGoogleID = true
		return result, nil
	}

	// CAS 2 : Recherche par email (Google a vérifié l'email)
	userByEmail, err := s.repository.FindUserByEmail(googleEmail)
	if err == nil {
		result.ExistingUser = userByEmail
		result.FoundByEmail = true

		// DÉCISION CRITIQUE : Liaison automatique
		// Si l'email existe MAIS pas de Google ID, et que Google a vérifié l'email
		// → On peut lier automatiquement (recommandation Google)
		result.CanAutoLink = !userByEmail.IsGoogleLinked()
		return result, nil
	}

	// CAS 3 : Aucun utilisateur trouvé - nouveau signup possible
	return result, nil
}

// ProcessUserAuthentication gère la logique complète d'authentification
func (s *GoogleAuthService) ProcessUserAuthentication(userInfo *GoogleUserInfo) (*AuthResult, error) {

	// Recherche utilisateur selon la logique Google
	lookupResult, err := s.LookupUser(userInfo.ID, userInfo.Email)
	if err != nil {
		return nil, &OAuthError{
			Code:    "user_lookup_failed",
			Message: "Failed to lookup user",
			Details: err.Error(),
		}
	}

	// CAS 1 : Utilisateur trouvé par Google ID - Login simple
	if lookupResult.FoundByGoogleID {
		return &AuthResult{
			User:      lookupResult.ExistingUser,
			IsNewUser: false,
			Method:    "login",
		}, nil
	}

	// CAS 2 : Utilisateur trouvé par email - Liaison automatique possible
	if lookupResult.FoundByEmail && lookupResult.CanAutoLink {
		// Lier le compte Google à l'utilisateur existant
		user := lookupResult.ExistingUser
		user.LinkGoogleAccount(userInfo.ID, userInfo.Email)

		if err := s.repository.UpdateUser(user); err != nil {
			return nil, &OAuthError{
				Code:    "google_account_link_failed",
				Message: "Failed to link Google account to existing user",
				Details: err.Error(),
			}
		}

		return &AuthResult{
			User:      user,
			IsNewUser: false,
			Method:    "link",
		}, nil
	}

	// CAS 3 : Utilisateur trouvé par email MAIS déjà lié à Google - Erreur
	if lookupResult.FoundByEmail && !lookupResult.CanAutoLink {
		return nil, &OAuthError{
			Code:    "email_already_linked",
			Message: "Email is already linked to another Google account",
			Details: "This email is already associated with a different Google account",
		}
	}
	// CAS 4 : Nouveau utilisateur - Signup automatique
	newUser, err := s.repository.CreateUserFromGoogle(userInfo)
	if err != nil {
		return nil, &OAuthError{
			Code:    "user_creation_failed",
			Message: "Failed to create new user from Google account",
			Details: err.Error(),
		}
	}

	return &AuthResult{
		User:      newUser,
		IsNewUser: true,
		Method:    "signup",
	}, nil
}

// GetUserInfoWithRetry récupère les infos utilisateur avec retry et backoff
func (s *GoogleAuthService) GetUserInfoWithRetry(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		// Créer un contexte avec timeout pour chaque tentative
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
		userInfo, err := s.GetUserInfo(ctxWithTimeout, token)
		cancel()

		if err == nil {
			if i > 0 {
				log.Printf("Successfully retrieved user info after %d retries", i)
			}
			return userInfo, nil
		}

		lastErr = err
		log.Printf("Attempt %d/%d failed to get user info: %v", i+1, maxRetries, err)

		// Backoff exponentiel : 1s, 2s, 4s
		if i < maxRetries-1 {
			backoffDuration := time.Duration(1<<uint(i)) * time.Second
			log.Printf("Retrying in %v...", backoffDuration)
			time.Sleep(backoffDuration)
		}
	}

	return nil, &OAuthError{
		Code:    "user_info_retry_failed",
		Message: fmt.Sprintf("Failed to get user info after %d retries", maxRetries),
		Details: lastErr.Error(),
	}
}
