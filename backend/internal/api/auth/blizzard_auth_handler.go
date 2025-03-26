// Package auth handles authentication-related functionality

package auth

/*

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"wowperf/internal/models"
	auth "wowperf/internal/services/auth"
	"wowperf/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// BlizzardAuthHandler handles Battle.net OAuth2 authentication endpoints
type BlizzardAuthHandler struct {
	blizzardAuthService *auth.BlizzardAuthService
	authService         *auth.AuthService
}

// NewBlizzardAuthHandler creates a new Battle.net authentication handler
func NewBlizzardAuthHandler(blizzardAuthService *auth.BlizzardAuthService, authService *auth.AuthService) *BlizzardAuthHandler {
	return &BlizzardAuthHandler{
		blizzardAuthService: blizzardAuthService,
		authService:         authService,
	}
}

// RegisterRoutes registers the Battle.net authentication routes
func (h *BlizzardAuthHandler) RegisterRoutes(router *gin.Engine) {
	battleNet := router.Group("/auth/battle-net")
	{
		// Public routes
		battleNet.GET("/link", h.InitiateAuth)

		// JWT protected routes
		protected := battleNet.Group("")
		protected.Use(middleware.JWTAuth(h.authService))
		{
			protected.GET("/callback", h.HandleCallback)
			protected.GET("/status", h.GetLinkStatus)
			protected.POST("/unlink", h.UnlinkAccount)
		}
	}
}

// InitiateAuth initiates the Battle.net OAuth flow
func (h *BlizzardAuthHandler) InitiateAuth(c *gin.Context) {
	// Generate random state using standard library
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate state",
			"code":  "state_generation_failed",
		})
		return
	}
	state := base64.URLEncoding.EncodeToString(b)[:32]

	// Store state in a secure cookie
	c.SetCookie(
		"oauth_state",
		state,
		3600, // 1 hour
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	authURL := h.blizzardAuthService.GetAuthorizationURL(state)
	c.JSON(http.StatusOK, gin.H{
		"url": authURL,
	})
}

// HandleCallback processes the OAuth2 callback from Battle.net
func (h *BlizzardAuthHandler) HandleCallback(c *gin.Context) {
	log.Printf("Starting Battle.net callback handling with URL: %s", c.Request.URL.String())
	log.Printf("Request Method: %s", c.Request.Method)
	log.Printf("Request Headers: %v", c.Request.Header)

	// 1. Get and validate state parameter to prevent CSRF
	code := c.Query("code")
	state := c.Query("state")
	storedState, err := c.Cookie("oauth_state")

	log.Printf("Received OAuth parameters - code: [REDACTED], state: %s", state)
	log.Printf("Retrieved stored state from cookie: %s", storedState)

	if err != nil || state != storedState {
		log.Printf("State validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid OAuth state",
			"code":  "invalid_state",
		})
		return
	}

	log.Printf("State validation successful")

	// 2. Exchange authorization code for tokens
	log.Printf("Exchanging code for token...")
	token, err := h.blizzardAuthService.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to exchange authorization code",
			"code":    "token_exchange_failed",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Token exchange successful")

	// 3. Use the access token to get user info from Battle.net
	log.Printf("Getting user info from Battle.net...")
	userInfo, err := h.blizzardAuthService.GetUserInfo(token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get Battle.net user info",
			"code":    "user_info_failed",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Successfully got user info for BattleTag: %s", userInfo.BattleTag)

	// 4. Get the authenticated user's ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
			"code":  "user_not_found",
		})
		return
	}

	log.Printf("Found user ID in context: %v", userID)

	// 5. Link the Battle.net account to the user's account
	log.Printf("Linking Battle.net account...")
	err = h.blizzardAuthService.LinkBattleNetAccount(userID.(uint), userInfo, token)
	if err != nil {
		log.Printf("Failed to link Battle.net account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"linked":  false,
			"error":   "Failed to link Battle.net account",
			"code":    "link_failed",
			"details": err.Error(),
		})
		return
	}

	// 6. Send successful response with token information
	log.Printf("Battle.net account successfully linked")
	c.JSON(http.StatusOK, gin.H{
		"linked":    true,
		"battleTag": userInfo.BattleTag,
		"message":   "Battle.net account successfully linked",
		"code":      "link_successful",
		"expiresIn": token.Expiry.Sub(time.Now()).Seconds(),
		"scope":     token.TokenType,
	})

	// Clear the oauth_state cookie as it's no longer needed
	c.SetCookie("oauth_state", "", -1, "/", "", true, true)
}

// GetLinkStatus checks if the user has a linked Battle.net account
func (h *BlizzardAuthHandler) GetLinkStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.authService.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"linked": false,
			"error":  "User not found",
			"code":   "user_not_found",
		})
		return
	}

	// Check if the token is expired and needs refresh
	if user.BattleTag != nil && user.BattleNetExpiresAt.Before(time.Now()) {
		// Token is expired, attempt to refresh
		if user.BattleNetRefreshToken != "" {
			newToken, err := h.blizzardAuthService.RefreshToken(c.Request.Context(), user.BattleNetRefreshToken)
			if err != nil {
				log.Printf("Failed to refresh token: %v", err)
				// Continue with the response, but log the error
			} else {
				// Update the user's tokens
				err = h.blizzardAuthService.UpdateUserBattleNetTokens(&user, newToken)
				if err != nil {
					log.Printf("Failed to update tokens: %v", err)
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"linked":    user.BattleTag != nil,
		"battleTag": user.BattleTag,
	})
}

// UnlinkAccount removes the Battle.net link from the user's account
func (h *BlizzardAuthHandler) UnlinkAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Clear all Battle.net related fields
	updates := map[string]interface{}{
		"battle_net_id":            nil,
		"battle_tag":               nil,
		"battle_net_refresh_token": nil,
		"battle_net_expires_at":    nil,
		"battle_net_token_type":    nil,
		"battle_net_scopes":        nil,
		"encrypted_token":          nil,
	}

	if err := h.authService.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to unlink Battle.net account",
			"code":    "unlink_failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Battle.net account successfully unlinked",
		"code":    "unlink_successful",
	})
}

*/
