package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// CSRFConfing holds the configuration for the CSRF middleware
type CSRFConfing struct {
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

	config := &CSRFConfing{
		Secret:   []byte(secret),
		Secure:   false,
		Path:     "/",
		MaxAge:   86400,
		SameSite: csrf.SameSiteLaxMode,
	}

	csrfMiddleware := csrf.Protect(
		config.Secret,
		csrf.Secure(config.Secure), // Set to true in production
		csrf.Path(config.Path),
		csrf.MaxAge(config.MaxAge),
		csrf.SameSite(config.SameSite),
		csrf.HttpOnly(true),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.FieldName("_csrf"),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CSRF error: %v", csrf.FailureReason(r))
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "CSRF token invalid"}`))
		})),
	)

	return func(c *gin.Context) {
		// Skip CSRF check for safe methods
		if isSafeMethod(c.Request.Method) {
			c.Next()
			return
		}

		csrfHandler := csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := csrf.Token(r)
			log.Printf("Generated CSRF Token: %s", token)
			c.Set("csrf_token", token)
			c.Next()

		}))

		csrfHandler.ServeHTTP(c.Writer, c.Request)
	}
}

// GetCSRFToken returns the CSRF token from the context
func GetCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		c.JSON(http.StatusOK, gin.H{"csrf_token": token})
	}
}

// isSafeMethod checks if the HTTP method is safe
func isSafeMethod(method string) bool {
	return method == "GET" || method == "HEAD" || method == "OPTIONS"
}
