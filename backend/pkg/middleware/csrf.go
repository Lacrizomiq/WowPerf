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

type CSRFConfig struct {
	Secret         string
	Domain         string
	AllowedOrigins []string
	Secure         bool
	Environment    string
}

// NewCSRFConfig create a new CSRF config
func NewCSRFConfig(env string) CSRFConfig {
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	return CSRFConfig{
		Secret:         os.Getenv("CSRF_SECRET"),
		Domain:         os.Getenv("DOMAIN"),
		AllowedOrigins: allowedOrigins,
		Secure:         true,
		Environment:    env,
	}
}

// InitCSRFMiddleware initialize the CSRF middleware
func InitCSRFMiddleware(config CSRFConfig) gin.HandlerFunc {
	if config.Secret == "" {
		log.Fatal("CSRF_SECRET must be set")
	}

	csrfMiddleware := csrf.Protect(
		[]byte(config.Secret),
		csrf.Secure(config.Secure),
		csrf.Path("/"),
		csrf.Domain(config.Domain),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("csrf"),
		csrf.ErrorHandler(http.HandlerFunc(handleCSRFError)),
		csrf.TrustedOrigins(config.AllowedOrigins),
	)

	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" ||
			c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		handler := csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		}))

		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// GetCSRFToken generate a new CSRF token
func GetCSRFToken() gin.HandlerFunc {
	config := NewCSRFConfig(os.Getenv("ENVIRONMENT"))
	csrfMiddleware := csrf.Protect(
		[]byte(config.Secret),
		csrf.Secure(config.Secure),
		csrf.Path("/"),
		csrf.Domain(config.Domain),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("csrf"),
	)

	return func(c *gin.Context) {
		var token string
		handler := csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token = csrf.Token(r)
		}))

		handler.ServeHTTP(c.Writer, c.Request)

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func handleCSRFError(w http.ResponseWriter, r *http.Request) {
	log.Printf("CSRF validation failed for request: %s %s", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	json.NewEncoder(w).Encode(gin.H{
		"error": "CSRF validation failed",
		"code":  "INVALID_CSRF_TOKEN",
	})
}

func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}
