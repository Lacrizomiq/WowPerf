// pkg/middleware/cache_middleware.go
package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter is the response writer for the cache middleware
type responseWriter struct {
	body *bytes.Buffer
	gin.ResponseWriter
}

// Write capture the response body
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// generateDefaultCacheKey generates the default cache key
func generateDefaultCacheKey(prefix string, c *gin.Context) string {
	// Get all the query parameters
	params := c.Request.URL.Query()

	// Sort parameters for consistent key generation
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build cache key with path and sorted parameters
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s:route:%s", prefix, c.Request.URL.Path))

	// Add each parameter to the cache key
	for _, k := range keys {
		v := params.Get(k)
		sb.WriteString(fmt.Sprintf(":%s:%s", k, v))
	}

	return sb.String()
}

// CacheMiddleware is the middleware for the cache with advanced options
func (cm *CacheManager) CacheMiddleware(routeConfig RouteConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !routeConfig.Enabled {
			c.Next()
			return
		}

		start := time.Now()
		ctx := c.Request.Context()

		// Generate the cache key
		cacheKey := generateDefaultCacheKey(cm.config.KeyPrefix, c)

		// Try to get the value from the cache
		var cachedResponse interface{}
		err := cm.config.Cache.Get(ctx, cacheKey, &cachedResponse)

		if err == nil {
			// Cache hit
			if cm.config.Metrics {
				duration := time.Since(start)
				cm.UpdateMetrics(true, duration, nil)
			}
			c.JSON(http.StatusOK, cachedResponse)
			c.Abort()
			return
		}

		// Cache miss, capture the response
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = writer

		c.Next()

		// Cache the new response
		if c.Writer.Status() == http.StatusOK {
			var response interface{}
			if err := json.Unmarshal(writer.body.Bytes(), &response); err == nil {
				expiration := routeConfig.Expiration
				if expiration == 0 {
					expiration = cm.config.Expiration
				}

				err = cm.config.Cache.Set(ctx, cacheKey, response, expiration)
				if cm.config.Metrics {
					cm.UpdateMetrics(false, time.Since(start), err)
				}

				// Add tags if necessary
				if len(routeConfig.Tags) > 0 {
					for _, tag := range routeConfig.Tags {
						tagKey := fmt.Sprintf("%s:tag:%s:%s", cm.config.KeyPrefix, tag, cacheKey)
						cm.config.Cache.Set(ctx, tagKey, "", expiration)
					}
				}
			}
		}
	}
}
