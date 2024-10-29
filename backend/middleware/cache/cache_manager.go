// pkg/middleware/cache_manager.go
package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// CacheManager is the manager for the cache middleware
type CacheManager struct {
	config  CacheConfig
	metrics *CacheMetrics
	routes  map[string]RouteConfig
	mu      sync.RWMutex // Mutex for the sync of the routes and metrics
}

// NewCacheManager creates a new cache manager
func NewCacheManager(config CacheConfig) *CacheManager {
	return &CacheManager{
		config:  config,
		metrics: &CacheMetrics{},
		routes:  make(map[string]RouteConfig),
	}
}

// InvalidateByTags invalidates the cache by tags
func (cm *CacheManager) InvalidateByTags(ctx context.Context, tags []string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, tag := range tags {
		pattern := fmt.Sprintf("%s:tag:%s*", cm.config.KeyPrefix, tag)
		keys, err := cm.config.Cache.Keys(ctx, pattern)
		if err != nil {
			return fmt.Errorf("failed to get keys for tag %s: %w", tag, err)
		}
		if err := cm.config.Cache.DeleteMany(ctx, keys); err != nil {
			return fmt.Errorf("failed to delete keys for tag %s: %w", tag, err)
		}
	}
	return nil
}

// UpdateMetrics updates the metrics for the cache
func (cm *CacheManager) UpdateMetrics(hit bool, duration time.Duration, err error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if err != nil {
		cm.metrics.Errors++
		return
	}

	if hit {
		cm.metrics.Hits++
	} else {
		cm.metrics.Misses++
	}

	cm.metrics.TotalTime += duration
	cm.metrics.LastUpdated = time.Now()
}

// GetMetrics returns the metrics for the cache
func (cm *CacheManager) GetMetrics() CacheMetrics {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return *cm.metrics
}
