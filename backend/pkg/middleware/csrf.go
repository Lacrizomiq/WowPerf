package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// CSRF is a middleware that protects against CSRF attacks
func CSRF() gin.HandlerFunc {
	csrfSecret := os.Getenv("CSRF_SECRET")
	if csrfSecret == "" {
		log.Fatal("CSRF_SECRET is not set")
	}

	csrfMiddleware := csrf.Protect(
		[]byte(csrfSecret),
		csrf.Secure(false), // Set to true in production
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.CookieName("csrf_token"),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CSRF error: %v", csrf.FailureReason(r))
			http.Error(w, "CSRF token invalid", http.StatusForbidden)
		})),
	)

	return func(c *gin.Context) {
		csrfHandler := csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := csrf.Token(r)
			log.Printf("Generated CSRF Token: %s", token)
			c.Set("csrf_token", token)
			c.Next()
		}))

		csrfHandler.ServeHTTP(c.Writer, c.Request)
	}
}

// CSRFToken returns the CSRF token
func CSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		c.Header("X-CSRF-Token", token)
		c.JSON(http.StatusOK, gin.H{"csrf_token": token})
	}
}
