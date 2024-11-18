// csrf.go
package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// CSRFConfig struct for CSRF middleware configuration
type Config struct {
	Domain         string
	AllowedOrigins []string
	Environment    string
}

// CSRFMiddleware is the CSRF middleware available for the app
var CSRFMiddleware func(http.Handler) http.Handler

// InitCSRFMiddleware initializes the CSRF middleware with the provided configuration
func InitCSRFMiddleware(config Config) {
	secret := os.Getenv("CSRF_SECRET")
	if secret == "" {
		log.Fatal("CSRF_SECRET is not set")
	}

	// Log CSRF configuration in local environment
	if config.Environment == "local" {
		log.Printf("🔒 Initializing CSRF middleware with config:")
		log.Printf("   Domain: %s", config.Domain)
		log.Printf("   Allowed Origins: %v", config.AllowedOrigins)
	}

	CSRFMiddleware = csrf.Protect(
		[]byte(secret),
		csrf.Secure(true),
		csrf.Path("/"),
		csrf.Domain(config.Domain),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("_csrf"),
		csrf.TrustedOrigins([]string{"localhost"}),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CSRF Error Details - Headers: %v", r.Header)
			log.Printf("CSRF Error Details - Method: %s", r.Method)
			log.Printf("CSRF Error Details - URL: %s", r.URL.String())
			log.Printf("CSRF Error Details - Domain: %s", config.Domain)
			log.Printf("CSRF Error Details - Allowed Origins: %v", config.AllowedOrigins)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)

			errorReason := csrf.FailureReason(r)
			log.Printf("CSRF Error: %v", errorReason)

			json.NewEncoder(w).Encode(gin.H{
				"error":   "CSRF validation failed",
				"code":    "INVALID_CSRF_TOKEN",
				"details": errorReason.Error(),
			})
		})),
	)
}

// NewCSRFHandler creates a new CSRF middleware handler for Gin
func NewCSRFHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if CSRFMiddleware == nil {
			log.Fatal("CSRF middleware not initialized")
		}

		// Skip CSRF check for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" ||
			c.Request.Method == "OPTIONS" || strings.HasSuffix(c.Request.URL.Path, "/csrf-token") {
			c.Next()
			return
		}

		// Ajout des headers CORS pour les requêtes CSRF
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
			c.Header("Access-Control-Expose-Headers", "X-CSRF-Token")
		}

		handler := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		}))
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// GetCSRFToken generates and returns a new CSRF token
func GetCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if CSRFMiddleware == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "CSRF middleware not initialized"})
			return
		}

		handler := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := csrf.Token(r)
			if token == "" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
				return
			}

			// Configuration CORS spécifique pour le token
			origin := c.GetHeader("Origin")
			if origin != "" && isAllowedOrigin(origin, c.GetString("allowed_origins")) {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Access-Control-Expose-Headers", "X-CSRF-Token, Set-Cookie")
			}

			// Définir le cookie avec SameSite=None et Secure=true
			c.SetSameSite(http.SameSiteNoneMode)
			c.SetCookie("_csrf", token, 3600, "/", "localhost", true, true)

			c.Header("X-CSRF-Token", token)
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		}))
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// Helper function to check if origin is allowed
func isAllowedOrigin(origin string, allowedOrigins string) bool {
	origins := strings.Split(allowedOrigins, ",")
	for _, allowed := range origins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}

// GetTokenFromRequest gets the CSRF token from the request
func GetTokenFromRequest(c *gin.Context, handler func(token string)) {
	if CSRFMiddleware == nil {
		log.Fatal("CSRF middleware not initialized")
	}

	wrapped := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := csrf.Token(r)
		handler(token)
	}))
	wrapped.ServeHTTP(c.Writer, c.Request)
}
