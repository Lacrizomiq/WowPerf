package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// CSRF is a middleware that protects against CSRF attacks
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		csrfMiddleware := csrf.Protect(
			[]byte(os.Getenv("CSRF_SECRET")),
			csrf.Secure(true),
			csrf.HttpOnly(true),
		)

		csrfHandler := csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		}))

		csrfHandler.ServeHTTP(c.Writer, c.Request)
	}
}

// CSRFToken returns the CSRF token
func CSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("csrf_token", csrf.Token(c.Request))
		c.Next()
	}
}
