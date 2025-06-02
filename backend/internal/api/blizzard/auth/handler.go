package auth

import (
	"log"
	"net/http"
	bnAuth "wowperf/internal/services/blizzard/auth"
	"wowperf/pkg/middleware/blizzard"

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
		// Public route for OAuth callback

		// Protected routes requiring user authentication
		authed := battleNet.Group("")
		authed.Use(requireAuth)
		{
			authed.GET("/link", h.InitiateAuth)
			authed.GET("/callback", h.HandleCallback)
			authed.GET("/status", h.GetLinkStatus)
			authed.POST("/unlink", h.UnlinkAccount)
		}

		// Routes requiring a linked Battle.net account
		bnetProtected := authed.Group("")
		bnetProtected.Use(h.middleware.RequireBattleNetAuth()) // Nouveau middleware simplifié
		{
			bnetProtected.GET("/profile", h.GetBattleNetProfile)
		}
	}
}

// InitiateAuth starts the Battle.net OAuth flow
func (h *BattleNetAuthHandler) InitiateAuth(c *gin.Context) {
	// Get user ID from context
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
			"code":  "user_not_authenticated",
		})
		return
	}

	// Initiate the OAuth flow with the user ID
	authURL, err := h.BattleNetAuthService.InitiateAuth(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initiate Battle.net authentication",
			"code":    "auth_initiation_failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":  authURL,
		"code": "auth_url_generated",
	})
}

// HandleCallback processes the Battle.net OAuth callback
func (h *BattleNetAuthHandler) HandleCallback(c *gin.Context) {
	log.Printf("Starting OAuth callback process")

	// 1. Vérifier la présence des paramètres
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required OAuth parameters",
			"code":  "invalid_oauth_params",
		})
		return
	}

	// 2. Récupérer l'utilisateur authentifié
	authenticatedUserID := c.GetUint("user_id")
	if authenticatedUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
			"code":  "user_not_authenticated",
		})
		return
	}

	// 3. Échanger le code et VALIDER que l'état correspond à l'utilisateur authentifié
	token, stateUserID, err := h.BattleNetAuthService.ExchangeCodeForToken(c.Request.Context(), code, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to exchange code for token",
			"code":    "token_exchange_failed",
			"details": err.Error(),
		})
		return
	}

	// 4. SÉCURITÉ : Vérifier que l'utilisateur du state = utilisateur authentifié
	if stateUserID != authenticatedUserID {
		log.Printf("SECURITY ALERT: State user ID (%d) != authenticated user ID (%d)",
			stateUserID, authenticatedUserID)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "OAuth state mismatch - security violation",
			"code":  "state_user_mismatch",
		})
		return
	}

	log.Printf("Token exchange successful: userID=%d, token_type=%s, expires=%v",
		stateUserID, token.TokenType, token.Expiry)

	// 5. Lier le compte Battle.net à l'utilisateur (userID déjà validé)
	if err := h.BattleNetAuthService.LinkUserAccount(c.Request.Context(), token, stateUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to link Battle.net account",
			"code":    "link_failed",
			"details": err.Error(),
		})
		return
	}
	log.Printf("Battle.net account linked: userID=%d", stateUserID)

	// 6. Récupérer les infos utilisateur pour la réponse
	userInfo, err := h.BattleNetAuthService.GetUserInfo(c.Request.Context(), token)
	if err != nil {
		// Le linking a réussi mais impossible de récupérer les infos
		c.JSON(http.StatusOK, gin.H{
			"message": "Battle.net authentication successful",
			"code":    "auth_successful",
			"linked":  true,
		})
		return
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
