package middleware

import (
	"fmt"
	"net/http"
	"time"
	"wowperf/internal/services/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth is a middleware that verifies the JWT token
func JWTAuth(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token cookie is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return authService.JWTSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check expiration
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
					c.Abort()
					return
				}
			}

			// Check if the token is blacklisted
			blacklisted, err := authService.IsTokenBlacklisted(tokenString)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check token status"})
				c.Abort()
				return
			}
			if blacklisted {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
				c.Abort()
				return
			}

			userID, ok := claims["user_id"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}

			c.Set("user_id", uint(userID))
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
