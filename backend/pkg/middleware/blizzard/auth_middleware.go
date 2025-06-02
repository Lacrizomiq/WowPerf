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

// RequireBattleNetAuth checks that the user has a linked and valid Battle.net account
// This middleware combines the functionality of both RequireBattleNetAccount and RequireValidToken
func (m *BattleNetMiddleware) RequireBattleNetAuth() gin.HandlerFunc {
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

		// Get and validate the token (this checks linking, expiration, and validity in one call)
		token, err := m.authService.GetUserToken(c.Request.Context(), userID)
		if err != nil {
			log.Printf("Error getting/validating token for user %d: %v", userID, err)

			// Distinguish between different error types for better UX
			switch err.Error() {
			case "battle.net account not linked":
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Battle.net account not linked",
					"code":  "account_not_linked",
				})
			case "battle.net token expired - user must re-authenticate":
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Battle.net token expired, please re-authenticate",
					"code":  "token_expired",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to validate Battle.net authentication",
					"code":  "auth_validation_failed",
				})
			}
			c.Abort()
			return
		}

		// Set the validated token in context for use by handlers
		c.Set("blizzard_token", token)

		log.Printf("Valid Battle.net token set in context for user %d", userID)
		c.Next()
	}
}

// DEPRECATED: Use RequireBattleNetAuth() instead
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

// DEPRECATED: Use RequireBattleNetAuth() instead
// RequireBattleNetAccount checks that the user has a linked Battle.net account
func (m *BattleNetMiddleware) RequireBattleNetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
				"code":  "user_not_found",
			})
			c.Abort()
			return
		}

		// Check if the account is linked
		status, err := m.authService.GetUserBattleNetStatus(c.Request.Context(), userID)
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
