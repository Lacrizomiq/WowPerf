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
		// Skip for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" ||
			c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// CORS management for CSRF
		origin := c.GetHeader("Origin")
		if origin != "" {
			if isAllowedOrigin(origin, config.AllowedOrigins) {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
				c.Header("Access-Control-Expose-Headers", "X-CSRF-Token")
			}
		}

		handler := csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		}))
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// GetCSRFToken generate a new CSRF token
func GetCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		if token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
			return
		}

		// Consistent cookie configuration
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(
			"csrf",              // name
			token,               // value
			3600,                // maxAge (1 heure)
			"/",                 // path
			os.Getenv("DOMAIN"), // domain
			true,                // secure (HTTPS only)
			false,               // httpOnly (false car besoin d'accès JS)
		)
		c.Header("X-CSRF-Token", token)

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func handleCSRFError(w http.ResponseWriter, r *http.Request) {
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
