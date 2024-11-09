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

	isDev := os.Getenv("ENV") != "production"

	CSRFMiddleware = csrf.Protect(
		[]byte(secret),
		csrf.Secure(!isDev),
		csrf.Path("/"),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("_csrf"),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		// Skip csrf for GET, HEAD, OPTIONS
		if c.Request.Method == "GET" ||
			c.Request.Method == "HEAD" ||
			c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		csrfHandler := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		handler := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := csrf.Token(r)
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
