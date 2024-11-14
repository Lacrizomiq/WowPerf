package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// CSRFConfig struct for CSRF middleware configuration
type CSRFConfig struct {
	isDevelopment bool
	Domain        string
	AllowedHosts  []string
}

// CSRFMiddleware is the CSRF middleware available for the app
var CSRFMiddleware func(http.Handler) http.Handler

// InitCSRFMiddleware initializes the CSRF middleware
func InitCSRFMiddleware() {
	secret := os.Getenv("CSRF_SECRET")
	if secret == "" {
		log.Fatal("CSRF_SECRET is not set")
	}

	CSRFMiddleware = csrf.Protect(
		[]byte(secret),
		csrf.Secure(true), // Always true because we use HTTPS
		csrf.Path("/"),
		csrf.Domain(".localhost"),           // Add . to the domain to allow subdomains
		csrf.SameSite(csrf.SameSiteLaxMode), // Lax mode to allow cross-domain requests
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("_csrf"),
		csrf.TrustedOrigins([]string{
			"https://localhost",
			"https://api.localhost",
		}),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CSRF Error Details - Headers: %v", r.Header)
			log.Printf("CSRF Error Details - Method: %s", r.Method)
			log.Printf("CSRF Error Details - URL: %s", r.URL.String())

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

func NewCSRFHandler() gin.HandlerFunc {
	// Initialize the CSRF middleware if it's not already done
	if CSRFMiddleware == nil {
		InitCSRFMiddleware()
	}

	return func(c *gin.Context) {
		log.Printf("🔒 CSRF Middleware - Method: %s, Path: %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("🔑 Headers: %v", c.Request.Header)

		// For OPTIONS requests, return immediately with a 200 status
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}

		// Skip for other safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			c.Next()
			return
		}

		// Verify the CSRF token
		token := c.GetHeader("X-CSRF-Token")
		log.Printf("CSRF Token: %s", token)

		csrfHandler := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("✅ CSRF validation passed")
			c.Next()
		}))
		csrfHandler.ServeHTTP(c.Writer, c.Request)
	}
}

// GetCSRFToken generates and returns a new CSRF token
func GetCSRFToken() gin.HandlerFunc {
	if CSRFMiddleware == nil {
		InitCSRFMiddleware()
	}

	return func(c *gin.Context) {
		log.Printf("📝 Generating CSRF token for %s", c.Request.Host)
		handler := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := csrf.Token(r)

			log.Printf("📤 Token CSRF generated: %s", token)

			// Adding explicit CORS headers for the token response
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Expose-Headers", "X-CSRF-Token")
			c.Header("X-CSRF-Token", token)
			c.JSON(http.StatusOK, gin.H{
				"token":  token,
				"header": "X-CSRF-Token",
			})
		}))
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// GetTokenFromRequest gets the CSRF token from the request
func GetTokenFromRequest(c *gin.Context, handler func(token string)) {
	if CSRFMiddleware == nil {
		InitCSRFMiddleware()
	}

	wrapped := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := csrf.Token(r)
		handler(token)
	}))
	wrapped.ServeHTTP(c.Writer, c.Request)
}
