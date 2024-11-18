package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func DebugMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("🔍 Incoming request:")
		log.Printf("   Method: %s", c.Request.Method)
		log.Printf("   Path: %s", c.Request.URL.Path)
		log.Printf("   Origin: %s", c.Request.Header.Get("Origin"))
		log.Printf("   CSRF Token: %s", c.Request.Header.Get("X-CSRF-Token"))
		log.Printf("   Headers: %v", c.Request.Header)

		c.Next()

		log.Printf("📤 Response status: %d", c.Writer.Status())
	}
}
