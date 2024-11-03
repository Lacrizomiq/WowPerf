package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"wowperf/internal/services/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuth(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from the cookie
		tokenCookie, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
			// Verify the signature method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return authService.JWTSecret, nil
		})

		if err != nil {
			// if token is expired, try to refresh
			var validationError *jwt.ValidationError
			if errors.As(err, &validationError) && validationError.Errors&jwt.ValidationErrorExpired != 0 {
				// Try to refresh
				if err := authService.RefreshToken(c); err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired, please login again"})
					c.Abort()
					return
				}
				// Continue with new token
				c.Next()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Verify if token is blacklist
		blacklisted, err := authService.IsTokenBlacklisted(tokenCookie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token status"})
			c.Abort()
			return
		}
		if blacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		// Extract and validate claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Verify the expirations
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
					c.Abort()
					return
				}
			}

			// Extract userID
			userID, ok := claims["user_id"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}

			// Store in the context
			c.Set("user_id", uint(userID))
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
