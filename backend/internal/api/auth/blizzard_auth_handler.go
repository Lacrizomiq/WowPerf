package auth

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"wowperf/internal/services/auth"
	"wowperf/pkg/middleware"
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

// RegisterRoutes registers Battle.net OAuth routes
func (h *BlizzardAuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth/battle-net")
	{
		// Public route
		auth.GET("/login", h.HandleBattleNetLogin)

		// Protected route
		protected := auth.Group("")
		protected.Use(middleware.JWTAuth(h.AuthService))
		{
			protected.GET("/callback", h.HandleBattleNetCallback)
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
