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
		log.Printf("🔒 CSRF Config:")
		log.Printf("   Domain: %s", config.Domain)
		log.Printf("   Secure: %v", secure)
		log.Printf("   Origins: %v", config.AllowedOrigins)
	}

	CSRFMiddleware = csrf.Protect(
		[]byte(config.Secret),
		csrf.Secure(secure),
		csrf.Path("/"),
		csrf.Domain(config.Domain),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("_csrf"),
		csrf.TrustedOrigins(config.AllowedOrigins),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("🚫 CSRF Validation Failed:")
			log.Printf("   Origin: %s", r.Header.Get("Origin"))
			log.Printf("   Referer: %s", r.Header.Get("Referer"))
			log.Printf("   Token: %s", r.Header.Get("X-CSRF-Token"))
			// Log the cookie header for debugging
			log.Printf("🍪 Cookies: %s", r.Header.Get("Cookie"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(gin.H{
				"error": "CSRF validation failed",
				"code":  "INVALID_CSRF_TOKEN",
				"debug": gin.H{
					"origin":  r.Header.Get("Origin"),
					"referer": r.Header.Get("Referer"),
					"cookie":  r.Header.Get("Cookie"),
				},
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

		// Create a custom ResponseWriter to prevent multiple writes
		crw := &customResponseWriter{ResponseWriter: c.Writer}
		c.Writer = crw

		handler := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		}))

		handler.ServeHTTP(crw, c.Request)
	}
}

// customResponseWriter wraps gin.ResponseWriter to prevent multiple writes
type customResponseWriter struct {
	gin.ResponseWriter
	written bool
}

func (w *customResponseWriter) WriteHeader(code int) {
	if !w.written {
		w.written = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.written = true
		return w.ResponseWriter.Write(b)
	}
	return len(b), nil
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
