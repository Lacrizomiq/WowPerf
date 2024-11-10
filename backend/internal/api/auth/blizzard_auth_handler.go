package auth

import (
	"crypto/rand"
	"encoding/base64"
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
	state := generateRandomState()
	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)
	url := h.BlizzardAuthService.GetAuthorizationURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleBattleNetCallback handles the callback from Battle.net OAuth2 flow
func (h *BlizzardAuthHandler) HandleBattleNetCallback(c *gin.Context) {
	state, _ := c.Cookie("oauth_state")
	if state != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	code := c.Query("code")
	token, err := h.BlizzardAuthService.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	userInfo, err := h.BlizzardAuthService.GetUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.BlizzardAuthService.LinkBattleNetAccount(userID.(uint), userInfo, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link Battle.net account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Battle.net account linked successfully"})
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

		// Protected route for Battle.net OAuth that need to be authenticated into WowPerf
		protected := auth.Group("")
		protected.Use(middleware.JWTAuth(h.AuthService))
		{
			protected.GET("/link", h.HandleBattleNetLogin)        // Initiate Battle.net OAuth flow
			protected.GET("/callback", h.HandleBattleNetCallback) // Handle Battle.net OAuth callback
			protected.GET("/status", h.GetLinkStatus)             // Get Battle.net link status
			protected.DELETE("/unlink", h.UnlinkBattleNetAccount) // Unlink Battle.net account
		}
	}
}

// generateRandomState generates a random state string
func generateRandomState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
