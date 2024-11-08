package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"wowperf/internal/services/auth"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// Custom errors for better error handling
var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
	ErrNoToken      = errors.New("no token provided")
)

// JWTAuth middleware handles the JWT token verification and CSRF protection
func JWTAuth(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from cookie
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			handleAuthError(c, ErrNoToken)
			return
		}

		// Verify token blacklist status first
		blacklisted, err := authService.IsTokenBlacklisted(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token status"})
			c.Abort()
			return
		}
		if blacklisted {
			handleAuthError(c, errors.New("token has been revoked"))
			return
		}

		// Parse and validate token
		token, err := parseAndValidateToken(tokenString, authService.JWTSecret)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				// Attempt to refresh the token
				if err := handleTokenRefresh(c, authService); err != nil {
					handleAuthError(c, err)
					return
				}
				// Continue with the new token
				c.Next()
				return
			}
			handleAuthError(c, err)
			return
		}

		// Extract and set user claims
		if err := setUserContext(c, token); err != nil {
			handleAuthError(c, err)
			return
		}

		c.Next()
	}
}

// parseAndValidateToken parses and validates the JWT token
func parseAndValidateToken(tokenString string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		var validationError *jwt.ValidationError
		if errors.As(err, &validationError) && validationError.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}

// setUserContext sets the user context in the request
func setUserContext(c *gin.Context, token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Verify expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return ErrTokenExpired
		}
	}

	// Extract and validate userID
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return errors.New("invalid user ID in token")
	}

	// Set user ID in the context
	c.Set("user_id", uint(userID))
	return nil
}

// handleTokenRefresh attempts to refresh the access token
func handleTokenRefresh(c *gin.Context, authService *auth.AuthService) error {
	if err := authService.RefreshToken(c); err != nil {
		log.Printf("Token refresh failed: %v", err)
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	return nil
}

// handleAuthError handles authentication errors consistently
func handleAuthError(c *gin.Context, err error) {
	var status int
	var message string

	switch {
	case errors.Is(err, ErrNoToken):
		status = http.StatusUnauthorized
		message = "Authentication required"
	case errors.Is(err, ErrTokenExpired):
		status = http.StatusUnauthorized
		message = "Token expired"
	case errors.Is(err, ErrInvalidToken):
		status = http.StatusUnauthorized
		message = "Invalid token"
	default:
		status = http.StatusUnauthorized
		message = "Authentication failed"
	}

	c.JSON(status, gin.H{"error": message})
	c.Abort()
}
