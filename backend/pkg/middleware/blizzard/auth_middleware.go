package blizzard

import (
	"log"
	"net/http"
	"wowperf/internal/services/blizzard/auth"

	"github.com/gin-gonic/gin"
)

// BattleNetMiddleware handles the verification of Battle.net tokens
type BattleNetMiddleware struct {
	authService *auth.BattleNetAuthService
}

func NewBattleNetMiddleware(authService *auth.BattleNetAuthService) *BattleNetMiddleware {
	return &BattleNetMiddleware{
		authService: authService,
	}
}

// RequireValidToken checks that the user has a valid Battle.net token
func (m *BattleNetMiddleware) RequireValidToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		if userID == 0 {
			log.Printf("No user ID found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
				"code":  "user_not_authenticated",
			})
			c.Abort()
			return
		}

		// Get the token
		token, err := m.authService.GetUserToken(c.Request.Context(), userID)
		if err != nil {
			log.Printf("Error getting token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Battle.net token not found or invalid",
				"code":  "token_invalid",
			})
			c.Abort()
			return
		}

		// Set it in context
		c.Set("blizzard_token", token)

		// Debug log
		log.Printf("Token set in context for user %d", userID)

		c.Next()
	}
}

// RequireBattleNetAccount checks that the user has a linked Battle.net account
func (m *BattleNetMiddleware) RequireBattleNetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
				"code":  "user_not_found",
			})
			c.Abort()
			return
		}

		// Check if the account is linked
		status, err := m.authService.GetUserBattleNetStatus(c.Request.Context(), userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check Battle.net status",
				"code":  "status_check_failed",
			})
			c.Abort()
			return
		}

		if !status.Linked {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Battle.net account not linked",
				"code":  "account_not_linked",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
