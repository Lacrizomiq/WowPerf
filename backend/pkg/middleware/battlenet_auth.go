package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"wowperf/internal/models"
	"wowperf/internal/services/auth"
)

// BlizzardAuthMiddleware handles Battle.net token validation and refresh
type BlizzardAuthMiddleware struct {
	db                  *gorm.DB
	blizzardAuthService *auth.BlizzardAuthService
}

// NewBlizzardAuthMiddleware creates a new Battle.net authentication middleware
func NewBlizzardAuthMiddleware(db *gorm.DB, blizzardAuthService *auth.BlizzardAuthService) *BlizzardAuthMiddleware {
	return &BlizzardAuthMiddleware{
		db:                  db,
		blizzardAuthService: blizzardAuthService,
	}
}

// RequireValidToken ensures a valid Battle.net token exists before proceeding
func (m *BlizzardAuthMiddleware) RequireValidToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		var user models.User
		if err := m.db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if Battle.net account is linked
		if user.BattleNetID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Battle.net account not linked"})
			c.Abort()
			return
		}

		// Get the current token
		accessToken, err := user.GetBattleNetToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt token"})
			c.Abort()
			return
		}

		token := &oauth2.Token{
			AccessToken:  accessToken,
			TokenType:    user.BattleNetTokenType,
			RefreshToken: user.BattleNetRefreshToken,
			Expiry:       user.BattleNetExpiresAt,
		}

		// Check if token needs refresh (expires within 5 minutes)
		if time.Until(token.Expiry) < 5*time.Minute {
			newToken, err := m.blizzardAuthService.RefreshToken(context.Background(), token.RefreshToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token"})
				c.Abort()
				return
			}

			// Update user's token
			if err := m.blizzardAuthService.UpdateUserBattleNetTokens(&user, newToken); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
				c.Abort()
				return
			}

			token = newToken
		}

		// Validate scopes
		if !m.blizzardAuthService.ValidateScopes(&user) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient Battle.net permissions"})
			c.Abort()
			return
		}

		// Store the validated token in context for use in handlers
		c.Set("blizzard_token", token)
		c.Set("battle_net_user", &user)
		c.Next()
	}
}
