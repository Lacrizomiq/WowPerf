package warcraftlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	warcraftlogsTypes "wowperf/internal/services/warcraftlogs/types"
)

// RateLimitData represents rate limit information from the WarcraftLogs API
type RateLimitData struct {
	LimitPerHour        float64 `json:"limitPerHour"`        // Maximum points per hour
	PointsSpentThisHour float64 `json:"pointsSpentThisHour"` // Points already spent
	PointsResetIn       int     `json:"pointsResetIn"`       // Seconds until points reset
}

// GraphQL query to fetch rate limit data
const RateLimitQuery = `query {
    rateLimitData {
        limitPerHour
        pointsSpentThisHour
        pointsResetIn
    }
}`

// WarcraftLogsClientService manages interactions with the WarcraftLogs API
type WarcraftLogsClientService struct {
	Client      *Client
	rateLimiter *RateLimiter
}

// RateLimiter handles API rate limiting with precise point tracking
type RateLimiter struct {
	mu          sync.RWMutex
	maxPoints   float64   // Maximum points allowed per hour
	usedPoints  float64   // Currently used points
	resetTime   time.Time // Exact time of next reset
	initialized bool      // Whether the limiter has been initialized
	lastCheck   time.Time // Last time the rate limit was checked
}

// NewWarcraftLogsClientService creates a new service instance
func NewWarcraftLogsClientService() (*WarcraftLogsClientService, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	service := &WarcraftLogsClientService{
		Client: client,
		rateLimiter: &RateLimiter{
			lastCheck: time.Now(),
		},
	}

	// Initialize rate limiter with API state
	if err := service.initializeRateLimiter(); err != nil {
		return nil, fmt.Errorf("failed to initialize rate limiter: %w", err)
	}

	// Start reset timer
	go service.startResetTimer()

	return service, nil
}

// Initialize rate limiter with API data
func (s *WarcraftLogsClientService) initializeRateLimiter() error {
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		response, err := s.Client.MakeGraphQLRequest(RateLimitQuery, nil)
		if err != nil {
			log.Printf("[WARN] Attempt %d: Failed to get rate limit data: %v", attempt+1, err)
			time.Sleep(time.Second * time.Duration(attempt+1))
			continue
		}

		var result struct {
			RateLimitData RateLimitData `json:"rateLimitData"`
		}

		if err := json.Unmarshal(response, &result); err != nil {
			log.Printf("[WARN] Attempt %d: Failed to parse rate limit data: %v", attempt+1, err)
			continue
		}

		data := result.RateLimitData
		if data.LimitPerHour <= 0 {
			if attempt == maxRetries-1 {
				log.Printf("[WARN] Invalid rate limit received from API: %f. Using default values.", data.LimitPerHour)
				s.rateLimiter.initialize(18000, 0, 3600)
				return nil
			}
			continue
		}

		// Valid data received
		s.rateLimiter.initialize(
			data.LimitPerHour,
			data.PointsSpentThisHour,
			data.PointsResetIn,
		)
		log.Printf("[INFO] Rate limiter initialized with API values: limit=%.2f, used=%.2f, reset=%ds",
			data.LimitPerHour,
			data.PointsSpentThisHour,
			data.PointsResetIn)
		return nil
	}

	log.Printf("[WARN] Failed to get valid rate limit after %d attempts. Using default values.", maxRetries)
	s.rateLimiter.initialize(18000, 0, 3600)
	return nil
}

// backgroundRateLimitSync tries to get real values from API
func (s *WarcraftLogsClientService) backgroundRateLimitSync() {
	maxAttempts := 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		time.Sleep(time.Second * time.Duration(attempt+1)) // Exponential backoff

		response, err := s.Client.MakeGraphQLRequest(RateLimitQuery, nil)
		if err != nil {
			log.Printf("[WARN] Background sync attempt %d failed: %v", attempt+1, err)
			continue
		}

		var result struct {
			Data struct {
				RateLimitData RateLimitData `json:"rateLimitData"`
			} `json:"data"`
		}

		if err := json.Unmarshal(response, &result); err != nil {
			log.Printf("[WARN] Failed to parse background sync data: %v", err)
			continue
		}

		data := result.Data.RateLimitData
		if data.LimitPerHour > 0 {
			s.rateLimiter.initialize(
				data.LimitPerHour,
				data.PointsSpentThisHour,
				data.PointsResetIn,
			)
			log.Printf("[INFO] Background sync successful: limit=%.2f", data.LimitPerHour)
			return
		}
	}
}

// initialize sets up the rate limiter with validated API data
func (r *RateLimiter) initialize(limitPerHour, pointsSpent float64, resetIn int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	r.maxPoints = limitPerHour
	r.usedPoints = pointsSpent
	r.resetTime = now.Add(time.Duration(resetIn) * time.Second)
	r.initialized = true
	r.lastCheck = now

	log.Printf("[INFO] Rate limiter initialized - Max: %.2f, Used: %.2f, Reset at: %s",
		r.maxPoints,
		r.usedPoints,
		r.resetTime.Format(time.RFC3339))
}

// startResetTimer manages the automatic reset of points
func (s *WarcraftLogsClientService) startResetTimer() {
	for {
		r := s.rateLimiter
		r.mu.RLock()
		waitDuration := time.Until(r.resetTime)
		r.mu.RUnlock()

		// Wait until next reset
		time.Sleep(waitDuration)

		// Reset points
		r.mu.Lock()
		r.usedPoints = 0
		r.resetTime = time.Now().Add(time.Hour)
		r.lastCheck = time.Now()
		log.Printf("[INFO] Rate limit reset. Next reset at: %s", r.resetTime.Format(time.RFC3339))
		r.mu.Unlock()
	}
}

// HasAvailablePoints checks if points are available for use
func (r *RateLimiter) HasAvailablePoints() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.maxPoints-r.usedPoints >= 1.0
}

// ConsumePoint marks one point as used
func (r *RateLimiter) ConsumePoint() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.usedPoints++
	log.Printf("[DEBUG] Point consumed. Used: %.2f, Remaining: %.2f",
		r.usedPoints, r.maxPoints-r.usedPoints)
}

// GetRateLimitInfo returns current rate limit information
func (r *RateLimiter) GetRateLimitInfo() *warcraftlogsTypes.RateLimitInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	remaining := r.maxPoints - r.usedPoints
	if remaining < 0 {
		remaining = 0
	}

	return &warcraftlogsTypes.RateLimitInfo{
		RemainingPoints: remaining,
		ResetIn:         time.Until(r.resetTime),
	}
}

// GetLastCheck returns the time of the last rate limit check
func (r *RateLimiter) GetLastCheck() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastCheck
}

// GetMaxPoints returns the maximum points per hour
func (r *RateLimiter) GetMaxPoints() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.maxPoints
}

// GetUsedPoints returns the number of points used in current period
func (r *RateLimiter) GetUsedPoints() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.usedPoints
}

// GetResetTime returns the time of the next reset
func (r *RateLimiter) GetResetTime() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.resetTime
}

// MakeRequest performs a rate-limited API request
func (s *WarcraftLogsClientService) MakeRequest(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	if query == "" {
		return nil, fmt.Errorf("empty query not allowed")
	}

	// Skip rate limiting for rate limit queries
	if query == RateLimitQuery {
		return s.Client.MakeGraphQLRequest(query, variables)
	}

	if !s.rateLimiter.HasAvailablePoints() {
		info := s.rateLimiter.GetRateLimitInfo()
		if info.ResetIn > 0 {
			log.Printf("[INFO] Waiting %v for rate limit reset", info.ResetIn.Round(time.Second))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(info.ResetIn):
				if !s.rateLimiter.HasAvailablePoints() {
					return nil, warcraftlogsTypes.NewQuotaExceededError(info)
				}
			}
		} else {
			return nil, warcraftlogsTypes.NewQuotaExceededError(info)
		}
	}

	response, err := s.Client.MakeGraphQLRequest(query, variables)
	if err != nil {
		if warcraftlogsTypes.IsRateLimit(err) {
			log.Printf("[WARN] Unexpected rate limit hit")
			return nil, warcraftlogsTypes.NewQuotaExceededError(s.rateLimiter.GetRateLimitInfo())
		}
		return nil, err
	}

	s.rateLimiter.ConsumePoint()
	return response, nil
}

// GetRateLimiter returns the rate limiter instance
func (s *WarcraftLogsClientService) GetRateLimiter() *RateLimiter {
	return s.rateLimiter
}
