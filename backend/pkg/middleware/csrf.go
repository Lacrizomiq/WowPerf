package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// CSRFConfig holds the configuration for the CSRF middleware
type CSRFConfig struct {
	Secret   []byte
	Secure   bool
	Domain   string
	Path     string
	MaxAge   int
	SameSite csrf.SameSiteMode
}

// NewCSRFMiddleware creates a new CSRF middleware
func NewCSRFMiddleware() gin.HandlerFunc {
	secret := os.Getenv("CSRF_SECRET")
	if secret == "" {
		log.Fatal("CSRF_SECRET is not set")
	}

	config := &CSRFConfig{
		Secret:   []byte(secret),
		Secure:   false, // Set to true in production
		Path:     "/",
		MaxAge:   86400,
		SameSite: csrf.SameSiteLaxMode,
	}

	csrfMiddleware := csrf.Protect(
		config.Secret,
		csrf.Secure(config.Secure),
		csrf.Path(config.Path),
		csrf.MaxAge(config.MaxAge),
		csrf.SameSite(config.SameSite),
		csrf.HttpOnly(true),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CSRF error: %v", csrf.FailureReason(r))
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "CSRF token invalid"}`))
		})),
	)

	return func(c *gin.Context) {
		// Skip CSRF check for the token endpoint and safe methods
		if c.Request.URL.Path == "/api/csrf-token" || isSafeMethod(c.Request.Method) {
			c.Next()
			return
		}

		csrfHandler := csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		}))

		csrfHandler.ServeHTTP(c.Writer, c.Request)
	}
}

// GetCSRFToken returns a new CSRF token
func GetCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("CSRF_SECRET")
		csrfMiddleware := csrf.Protect([]byte(secret), csrf.Secure(false))

		csrfHandler := csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := csrf.Token(r)
			c.JSON(http.StatusOK, gin.H{"csrf_token": token})
		}))

		csrfHandler.ServeHTTP(c.Writer, c.Request)
	}
}

func isSafeMethod(method string) bool {
	return method == "GET" || method == "HEAD" || method == "OPTIONS"
}
