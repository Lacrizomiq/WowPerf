package auth

import (
	"fmt"
	"log"
	"net/http"
	bnAuth "wowperf/internal/services/blizzard/auth"
	"wowperf/pkg/middleware/blizzard"
	"wowperf/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// BattleNetAuthHandler handles Battle.net authentication endpoints
type BattleNetAuthHandler struct {
	BattleNetAuthService *bnAuth.BattleNetAuthService
	middleware           *blizzard.BattleNetMiddleware
}

func NewBattleNetAuthHandler(
	BattleNetAuthService *bnAuth.BattleNetAuthService,
) *BattleNetAuthHandler {
	return &BattleNetAuthHandler{
		BattleNetAuthService: BattleNetAuthService,
		middleware:           blizzard.NewBattleNetMiddleware(BattleNetAuthService),
	}
}

// RegisterRoutes registers all Battle.net authentication routes
func (h *BattleNetAuthHandler) RegisterRoutes(r *gin.Engine, requireAuth gin.HandlerFunc) {
	battleNet := r.Group("/auth/battle-net")
	{
		// Route publique pour le callback OAuth
		battleNet.GET("/callback", h.HandleCallback)

		// Routes protégées par authentification utilisateur
		authed := battleNet.Group("")
		authed.Use(requireAuth)
		{
			authed.GET("/link", h.InitiateAuth)
			authed.GET("/status", h.GetLinkStatus)
			authed.POST("/unlink", h.UnlinkAccount)
		}

		// Routes nécessitant un compte Battle.net lié
		bnetProtected := authed.Group("")
		bnetProtected.Use(h.middleware.RequireBattleNetAccount())
		bnetProtected.Use(h.middleware.RequireValidToken())
		{
			bnetProtected.GET("/profile", h.GetBattleNetProfile)
		}
	}
}

// InitiateAuth starts the Battle.net OAuth flow
func (h *BattleNetAuthHandler) InitiateAuth(c *gin.Context) {
	state := utils.GenerateRandomString(32)

	// Store state in a secure cookie
	c.SetCookie(
		"oauth_state",
		state,
		3600, // 1 hour expiry
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	// Get authorization URL
	authURL := h.BattleNetAuthService.GetAuthorizationURL(state)

	c.JSON(http.StatusOK, gin.H{
		"url":  authURL,
		"code": "auth_url_generated",
	})
}

// HandleCallback processes the Battle.net OAuth callback
// HandleCallback processes the Battle.net OAuth callback
func (h *BattleNetAuthHandler) HandleCallback(c *gin.Context) {
	log.Printf("Processing Battle.net callback: %s", c.Request.URL.String())

	// Vérifier le state
	state := c.Query("state")
	storedState, err := c.Cookie("oauth_state")
	log.Printf("State verification - Received: %s, Stored: %s, Cookie Error: %v", state, storedState, err)

	if err != nil {
		log.Printf("State cookie not found: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "State cookie not found",
			"code":    "invalid_state",
			"details": err.Error(),
		})
		return
	}

	if state == "" {
		log.Printf("State parameter missing")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "State parameter missing",
			"code":  "invalid_state",
		})
		return
	}

	if state != storedState {
		log.Printf("State mismatch - Received: %s, Expected: %s", state, storedState)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "State mismatch",
			"code":  "invalid_state",
		})
		return
	}

	// Get authorization code
	code := c.Query("code")
	log.Printf("Received authorization code: %s", code)

	// Exchange code for token
	ctx := c.Request.Context()
	token, err := h.BattleNetAuthService.ExchangeCodeForToken(ctx, code)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to exchange code for token",
			"code":    "token_exchange_failed",
			"details": err.Error(),
		})
		return
	}
	log.Printf("Successfully exchanged code for token")

	// Get user info
	userInfo, err := h.BattleNetAuthService.GetUserInfo(ctx, token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get user info",
			"code":    "user_info_failed",
			"details": err.Error(),
		})
		return
	}
	log.Printf("Retrieved user info for BattleTag: %s", userInfo.BattleTag)

	// Clear state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", true, true)

	// Si on a un user ID dans le contexte, on lie le compte
	if userID, exists := c.Get("user_id"); exists {
		log.Printf("Linking account for user ID: %v", userID)
		if err := h.BattleNetAuthService.LinkUserAccount(ctx, token, fmt.Sprint(userID.(uint))); err != nil {
			log.Printf("Account linking failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to link account",
				"code":    "link_failed",
				"details": err.Error(),
			})
			return
		}
		log.Printf("Successfully linked Battle.net account")
	} else {
		log.Printf("No user ID found in context, skipping account linking")
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Battle.net authentication successful",
		"code":      "auth_successful",
		"battleTag": userInfo.BattleTag,
		"linked":    true,
	})
}

// GetLinkStatus returns the current Battle.net link status
func (h *BattleNetAuthHandler) GetLinkStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	status, err := h.BattleNetAuthService.GetUserBattleNetStatus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get Battle.net status",
			"code":    "status_check_failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// UnlinkAccount removes the Battle.net account link
func (h *BattleNetAuthHandler) UnlinkAccount(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.BattleNetAuthService.UnlinkUserAccount(c.Request.Context(), userID); err != nil {
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

// GetBattleNetProfile returns the user's Battle.net profile
func (h *BattleNetAuthHandler) GetBattleNetProfile(c *gin.Context) {
	tokenInterface, exists := c.Get("blizzard_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Battle.net token not found",
			"code":  "token_not_found",
		})
		return
	}

	token, ok := tokenInterface.(*oauth2.Token)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid token format",
			"code":  "invalid_token_format",
		})
		return
	}

	profile, err := h.BattleNetAuthService.GetUserInfo(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get Battle.net profile",
			"code":    "profile_fetch_failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// Private helper methods

func (h *BattleNetAuthHandler) verifyState(c *gin.Context) error {
	state := c.Query("state")
	storedState, err := c.Cookie("oauth_state")

	if err != nil {
		return fmt.Errorf("state cookie not found")
	}

	if state == "" {
		return fmt.Errorf("state parameter missing")
	}

	if state != storedState {
		return fmt.Errorf("state mismatch")
	}

	return nil
}
