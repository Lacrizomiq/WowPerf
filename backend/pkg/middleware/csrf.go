package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// Config contains the CSRF middleware configuration
type Config struct {
	Secret         string
	Domain         string
	AllowedOrigins []string
	Environment    string
}

// CSRFMiddleware is the CSRF middleware instance
var CSRFMiddleware func(http.Handler) http.Handler

// InitCSRFMiddleware initializes the CSRF middleware
func InitCSRFMiddleware(config Config) {
	if config.Secret == "" {
		log.Fatal("CSRF_SECRET is not set")
	}

	secure := config.Environment != "local"

	if config.Environment == "local" {
		log.Printf("🔒 CSRF Config - Domain: %s, Secure: %v", config.Domain, secure)
	}

	CSRFMiddleware = csrf.Protect(
		[]byte(config.Secret),
		csrf.Secure(secure),
		csrf.Path("/"),
		csrf.Domain(config.Domain),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("_csrf"),
		csrf.TrustedOrigins(config.AllowedOrigins),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(gin.H{
				"error": "CSRF validation failed",
				"code":  "INVALID_CSRF_TOKEN",
			})
		})),
	)
}

// NewCSRFHandler returns a Gin middleware handler for CSRF protection
func NewCSRFHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if CSRFMiddleware == nil {
			log.Fatal("CSRF middleware not initialized")
		}

		// Skip CSRF for safe methods
		if isSafeMethod(c.Request.Method) || c.Request.URL.Path == "/api/csrf-token" {
			c.Next()
			return
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
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		}))

		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// Helper function to check if method is safe
func isSafeMethod(method string) bool {
	return method == "GET" ||
		method == "HEAD" ||
		method == "OPTIONS"
}
