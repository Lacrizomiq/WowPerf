package googleauth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"wowperf/internal/models"
	"wowperf/internal/services/auth"
	googleauthService "wowperf/internal/services/auth/google"

	"github.com/gin-gonic/gin"
)

// GoogleAuthHandler gère les endpoints Google OAuth
type GoogleAuthHandler struct {
	service     *googleauthService.GoogleAuthService
	authService *auth.AuthService
}

// NewGoogleAuthHandler initialise un nouveau handler GoogleAuthHandler
func NewGoogleAuthHandler(
	service *googleauthService.GoogleAuthService,
	authService *auth.AuthService,
) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		service:     service,
		authService: authService,
	}
}

// InitiateGoogleAuth initie le flux OAuth Google
// GET /auth/google/login
func (h *GoogleAuthHandler) InitiateGoogleAuth(c *gin.Context) {
	log.Printf("Initiating Google OAuth flow")

	// Genérer l'URL d'autorisation et l'état CSRF
	authURL, state, err := h.service.GetAuthURL()
	if err != nil {
		log.Printf("Failed to generate Google auth URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initiate Google Authentication",
			"code":  "google_auth_init_failed",
		})
		return
	}

	// Stocker l'état dans un cookie temporaire sécurisé (10 minutes)
	c.SetCookie(
		"google_oauth_state",              // name
		state,                             // value
		int((10 * time.Minute).Seconds()), // maxAge (10 minutes)
		"/",                               // Path
		"",                                // domain (sera défini automatiquement)
		true,                              // secure (HTTPS uniquement)
		true,                              // httpOnly (pas accessible en JS)
	)

	log.Printf("Google Oauth state cookie set, redirecting to Google")

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// HandleGoogleCallback traite le callback OAuth Google
// GET /auth/google/callback
func (h *GoogleAuthHandler) HandleGoogleCallback(c *gin.Context) {
	log.Printf("Starting Google OAuth callback processing")

	// 1. Récupérer et valider les paramètres
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// Gérer le cas où l'utilisateur refuse l'autorisation
	if errorParam != "" {
		h.handleOAuthError(c, errorParam, c.Query("error_description"))
		return
	}

	if code == "" || state == "" {
		log.Printf("Missing code or state in Google callback")
		h.redirectToFrontendWithError(c, "invalid_callback", "Invalid Google callback parameters")
		return
	}

	// 2. Vérifier l'état CSRF
	if err := h.validateState(c, state); err != nil {
		h.handleGenericError(c, err, "CSRF state validation")
		return
	}

	// 3. Échange du code contre un token
	token, err := h.service.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		h.redirectToFrontendWithError(c, "token_exchange_failed", "Failed to exchange authorization code")
		return
	}

	// 4. Récupérer les informations utilisateur
	userInfo, err := h.service.GetUserInfoWithRetry(c.Request.Context(), token)
	if err != nil {
		h.redirectToFrontendWithError(c, "user_info_failed", "Failed to get user information from Google")
		return
	}

	// 5. Logique de loging/signup/linking
	authResult, err := h.service.ProcessUserAuthentication(userInfo)
	if err != nil {
		log.Printf("Failed to process user authentication: %v", err)
		h.redirectToFrontendWithError(c, "auth_processing_failed", "Failed to process authentication")
		return
	}

	log.Printf("Authentication successful: method=%s, new_user=%t, user_id=%d",
		authResult.Method, authResult.IsNewUser, authResult.User.ID)

	// 6. Générer JWT et définir cookies
	if err := h.setAuthenticationCookies(c, authResult.User); err != nil {
		log.Printf("Failed to set authentication cookies: %v", err)
		h.redirectToFrontendWithError(c, "cookie_setting_failed", "Failed to set authentication")
		return
	}

	// 7. Redirection finale vers le frontend
	h.redirectToFrontendWithSuccess(c, authResult)
}

// validateState valide l'état CSRF
func (h *GoogleAuthHandler) validateState(c *gin.Context, receivedState string) error {
	// Récupérer l'état stocké dans le cookie
	storedState, err := c.Cookie("google_oauth_state")
	if err != nil {
		return fmt.Errorf("state cookie not found: %w", err)
	}

	// Supprimer le cookie d'état après utilisation
	c.SetCookie("google_oauth_state", "", -1, "/", "", true, true)

	// Comparer les états
	if receivedState != storedState {
		return fmt.Errorf("state mismatch")
	}

	return nil
}

// setAuthenticationCookies définit les cookies JWT
func (h *GoogleAuthHandler) setAuthenticationCookies(c *gin.Context, user *models.User) error {
	log.Printf("Setting authentication cookies for user ID: %d", user.ID)

	// Générer access token via AuthService existant
	accessToken, err := h.authService.GenerateToken(user.ID, h.authService.TokenExpiration)
	if err != nil {
		return fmt.Errorf("failed to generate access token: %w", err)
	}

	// Générer refresh token via AuthService existant
	refreshToken, err := h.authService.GenerateRefreshToken(user.ID)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Définir les cookies via AuthService existant
	h.authService.SetAuthCookies(c, accessToken, refreshToken)

	log.Printf("Authentication cookies set successfully for user: %s", user.Email)
	return nil
}

// redirectToFrontendWithSuccess redirige vers le frontend avec succès
func (h *GoogleAuthHandler) redirectToFrontendWithSuccess(c *gin.Context, result *googleauthService.AuthResult) {
	frontendURL := h.service.GetFrontendURL()

	// Toujours rediriger vers /auth/callback pour une meilleur UX
	redirectPath := "/auth/callback"

	// Ajouter query param pour indiquer si nouvel utilisateur
	var queryParams string
	if result.IsNewUser {
		queryParams = "?new_user=true"
	}

	finalURL := frontendURL + redirectPath + queryParams
	log.Printf("Redirecting to frontend callback: %s (method: %s)", finalURL, result.Method)

	c.Redirect(http.StatusSeeOther, finalURL)
}

// redirectToFrontendWithError redirige vers le frontend avec erreur
func (h *GoogleAuthHandler) redirectToFrontendWithError(c *gin.Context, errorCode, errorMessage string) {
	frontendURL := h.service.GetFrontendURL()
	errorPath := h.service.GetErrorPath()

	// Encoder les paramètres pour URL
	encodedCode := strings.ReplaceAll(errorCode, " ", "%20")
	encodedMessage := strings.ReplaceAll(errorMessage, " ", "%20")

	finalURL := fmt.Sprintf("%s%s?error=%s&message=%s",
		frontendURL, errorPath, encodedCode, encodedMessage) // ✅ CORRECT

	log.Printf("Redirecting to frontend with error: %s", finalURL)
	c.Redirect(http.StatusSeeOther, finalURL)
}

// RegisterRoutes enregistre les routes Google OAuth
func (h *GoogleAuthHandler) RegisterRoutes(router *gin.Engine) {
	googleAuth := router.Group("/auth/google")
	{
		// Route publique pour initier l'authentification
		googleAuth.GET("/login", h.InitiateGoogleAuth)

		// Route de callback publique
		googleAuth.GET("/callback", h.HandleGoogleCallback)
	}
}
