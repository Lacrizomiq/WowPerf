// pkg/middleware/cache_config.go
package middleware

import (
	"time"
	"wowperf/pkg/cache"

	"github.com/gin-gonic/gin"
)

// CacheConfig is the configuration for the cache middleware
type CacheConfig struct {
	Cache      cache.CacheService
	Expiration time.Duration
	KeyPrefix  string
	Tags       []string
	Metrics    bool
}

// RouteConfig is the configuration for a specific route
type RouteConfig struct {
	Enabled      bool
	Expiration   time.Duration
	Tags         []string
	KeyGenerator func(ctx *gin.Context) string
}

// CacheMetrics is the metrics for the cache middleware
type CacheMetrics struct {
	Hits        int64
	Misses      int64
	Errors      int64
	TotalTime   time.Duration
	LastUpdated time.Time
}
