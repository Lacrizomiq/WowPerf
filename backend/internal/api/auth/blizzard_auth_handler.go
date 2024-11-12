package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"wowperf/internal/models"
	"wowperf/internal/services/auth"
	"wowperf/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type BlizzardAuthHandler struct {
	BlizzardAuthService *auth.BlizzardAuthService
	AuthService         *auth.AuthService
}

func NewBlizzardAuthHandler(blizzardAuthService *auth.BlizzardAuthService, authService *auth.AuthService) *BlizzardAuthHandler {
	return &BlizzardAuthHandler{
		BlizzardAuthService: blizzardAuthService,
		AuthService:         authService,
	}
}

// HandleBattleNetLogin initiates the OAuth2 flow for Battle.net login
func (h *BlizzardAuthHandler) HandleBattleNetLogin(c *gin.Context) {
	log.Printf("Starting Battle.net login process...")

	state := generateRandomState()
	log.Printf("Generated state: %s", state)

	// Set cookie with all attributes in a single SetCookie call
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"oauth_state", // name
		state,         // value
		3600,          // maxAge
		"/",           // path
		"",            // domain
		false,         // secure (false for development)
		true,          // httpOnly
	)

	log.Printf("Cookie settings applied: Path=/, SameSite=Lax, Secure=false, HttpOnly=true")

	url := h.BlizzardAuthService.GetAuthorizationURL(state)
	log.Printf("Generated authorization URL with state parameter: %s", url)

	c.JSON(http.StatusOK, gin.H{"url": url})
	log.Printf("Battle.net login process initiated successfully")
}

// HandleBattleNetCallback handles the Battle.net OAuth callback
// HandleBattleNetCallback handles the Battle.net OAuth callback
func (h *BlizzardAuthHandler) HandleBattleNetCallback(c *gin.Context) {
	log.Printf("Starting Battle.net callback handling with URL: %s", c.Request.URL.String())
	log.Printf("Request Method: %s", c.Request.Method)
	log.Printf("Request Headers: %+v", c.Request.Header)

	// Get code and state from query parameters
	code := c.Query("code")
	state := c.Query("state")

	log.Printf("Received OAuth parameters - code: [REDACTED], state: %s", state)

	if code == "" || state == "" {
		log.Printf("Missing required parameters - code present: %v, state present: %v",
			code != "", state != "")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	// Get stored state from cookie
	storedState, err := c.Cookie("oauth_state")
	if err != nil {
		log.Printf("Error retrieving oauth_state cookie: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing state cookie"})
		return
	}
	log.Printf("Retrieved stored state from cookie: %s", storedState)

	if state != storedState {
		log.Printf("State mismatch - Received: %s, Stored: %s", state, storedState)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}
	log.Printf("State validation successful")

	// Échanger le code contre un token
	log.Printf("Exchanging code for token...")
	token, err := h.BlizzardAuthService.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}
	log.Printf("Token exchange successful")

	// Get user info from Battle.net
	log.Printf("Getting user info from Battle.net...")
	userInfo, err := h.BlizzardAuthService.GetUserInfo(token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	log.Printf("Successfully got user info for BattleTag: %s", userInfo.BattleTag)

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	log.Printf("Found user ID in context: %v", userID)

	// Link account
	log.Printf("Linking Battle.net account...")
	if err := h.BlizzardAuthService.LinkBattleNetAccount(userID.(uint), userInfo, token); err != nil {
		log.Printf("Failed to link account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Account successfully linked")

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	log.Printf("Cleared oauth_state cookie")

	response := gin.H{
		"message":   "Battle.net account linked successfully",
		"linked":    true,
		"battleTag": userInfo.BattleTag,
	}

	log.Printf("Sending success response: %+v", response)
	c.JSON(http.StatusOK, response)
	log.Printf("OAuth flow completed successfully")
}

// GetLinkStatus returns the status of the Battle.net link
func (h *BlizzardAuthHandler) GetLinkStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := h.AuthService.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"linked":    user.BattleNetID != nil,
		"battleTag": user.BattleTag,
	})
}

// UnlinkBattleNetAccount unlinks the Battle.net account
func (h *BlizzardAuthHandler) UnlinkBattleNetAccount(c *gin.Context) {
	userID := c.GetUint("user_id")

	updates := map[string]interface{}{
		"battle_net_id":            nil,
		"battle_tag":               nil,
		"encrypted_token":          nil,
		"battle_net_refresh_token": nil,
		"battle_net_token_type":    nil,
		"battle_net_expires_at":    nil,
		"battle_net_scopes":        nil,
	}

	if err := h.AuthService.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink Battle.net account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Battle.net account unlinked successfully"})
}

// RegisterRoutes registers Battle.net OAuth routes
func (h *BlizzardAuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth/battle-net")
	{
		protected := auth.Group("")
		protected.Use(middleware.JWTAuth(h.AuthService))
		{
			protected.GET("/link", h.HandleBattleNetLogin)
			protected.GET("/callback", h.HandleBattleNetCallback)
			protected.GET("/status", h.GetLinkStatus)
			protected.DELETE("/unlink", h.UnlinkBattleNetAccount)
		}
	}
}

// generateRandomState generates a random state string
func generateRandomState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
