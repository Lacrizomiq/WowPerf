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
	LimitPerHour        float64 `json:"limitPerHour"`
	PointsSpentThisHour float64 `json:"pointsSpentThisHour"`
	PointsResetIn       int     `json:"pointsResetIn"`
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

// RateLimiter handles API rate limiting with real-time point tracking
type RateLimiter struct {
	mu          sync.RWMutex
	maxPoints   float64   // Maximum points allowed per hour
	usedPoints  float64   // Currently used points based on last API check
	resetTime   time.Time // Time of next reset based on API
	lastCheck   time.Time // Last time we checked the API
	initialized bool
}

// NewWarcraftLogsClientService creates a new service instance
func NewWarcraftLogsClientService() (*WarcraftLogsClientService, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	rateLimiter := &RateLimiter{
		maxPoints: 18000, // Valeur conservative
		lastCheck: time.Now(),
		resetTime: time.Now().Add(time.Hour),
	}

	service := &WarcraftLogsClientService{
		Client:      client,
		rateLimiter: rateLimiter,
	}

	// Make up to 3 rate limit init
	var initError error
	for i := 0; i < 3; i++ {
		if err := service.initializeRateLimiter(); err != nil {
			initError = err
			log.Printf("[WARN] Rate limit initialization attempt %d failed: %v", i+1, err)
			time.Sleep(time.Second * 2)
			continue
		}
		initError = nil
		break
	}

	// even if it fails, return the service, it will use default value
	if initError != nil {
		log.Printf("[WARN] Using conservative default values due to initialization failure: %v", initError)
	}

	return service, nil
}

// Initialize rate limiter with API data
func (s *WarcraftLogsClientService) initializeRateLimiter() error {
	// Initial API check for rate limit data
	response, err := s.Client.MakeGraphQLRequest(RateLimitQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to get initial rate limit data: %w", err)
	}

	var result struct {
		Data struct {
			RateLimitData RateLimitData `json:"rateLimitData"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("failed to parse initial rate limit data: %w", err)
	}

	data := result.Data.RateLimitData
	s.rateLimiter.initialize(
		data.LimitPerHour,
		data.PointsSpentThisHour,
		data.PointsResetIn,
	)

	log.Printf("[INFO] Rate limiter initialized - Max: %.2f, Used: %.2f, Reset in: %ds",
		data.LimitPerHour,
		data.PointsSpentThisHour,
		data.PointsResetIn)

	return nil
}

// MakeRequest performs a rate-limited API request with real-time point checking
func (s *WarcraftLogsClientService) MakeRequest(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	if query == "" {
		return nil, fmt.Errorf("empty query not allowed")
	}

	// Skip rate checking for rate limit queries to avoid recursion
	if query == RateLimitQuery {
		return s.Client.MakeGraphQLRequest(query, variables)
	}

	// Check current rate limit status
	rateLimitResponse, err := s.Client.MakeGraphQLRequest(RateLimitQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}

	var result struct {
		RateLimitData RateLimitData `json:"rateLimitData"`
	}

	if err := json.Unmarshal(rateLimitResponse, &result); err != nil {
		return nil, fmt.Errorf("failed to parse rate limit data: %w", err)
	}

	// Update rate limiter with latest data
	data := result.RateLimitData
	s.rateLimiter.initialize(
		data.LimitPerHour,
		data.PointsSpentThisHour,
		data.PointsResetIn,
	)

	// Check if we have enough points (considering the 2 points just used for checking)
	if s.rateLimiter.maxPoints-s.rateLimiter.usedPoints < 3.0 { // 2 points for check + 1 minimum for request
		info := s.rateLimiter.GetRateLimitInfo()
		return nil, warcraftlogsTypes.NewQuotaExceededError(info)
	}

	// Make the actual request
	response, err := s.Client.MakeGraphQLRequest(query, variables)
	if err != nil {
		if warcraftlogsTypes.IsRateLimit(err) {
			return nil, warcraftlogsTypes.NewQuotaExceededError(s.rateLimiter.GetRateLimitInfo())
		}
		return nil, err
	}

	return response, nil
}

// initialize sets up the rate limiter with latest API data
func (r *RateLimiter) initialize(limitPerHour, pointsSpent float64, resetIn int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	r.maxPoints = limitPerHour // 18000
	r.usedPoints = pointsSpent // 326
	r.resetTime = now.Add(time.Duration(resetIn) * time.Second)
	r.lastCheck = now
	r.initialized = true

	// Log après mise à jour pour vérification
	log.Printf("[DEBUG] Rate limiter - Max: %.2f, Used: %.2f, Remaining: %.2f, Reset in: %v",
		r.maxPoints,
		r.usedPoints,
		r.maxPoints-r.usedPoints,
		time.Until(r.resetTime))
}

// GetRateLimitInfo returns current rate limit information for monitoring
func (r *RateLimiter) GetRateLimitInfo() *warcraftlogsTypes.RateLimitInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	remaining := r.maxPoints - r.usedPoints
	if remaining < 0 {
		remaining = 0
	}

	if time.Now().After(r.resetTime) {
		remaining = r.maxPoints
	}

	return &warcraftlogsTypes.RateLimitInfo{
		RemainingPoints: remaining,
		ResetIn:         time.Until(r.resetTime),
	}
}

// GetLastCheck returns the time of the last rate limit check (for monitoring)
func (r *RateLimiter) GetLastCheck() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastCheck
}

// GetMaxPoints returns the maximum points per hour (for monitoring)
func (r *RateLimiter) GetMaxPoints() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.maxPoints
}

// GetUsedPoints returns the number of points used (for monitoring)
func (r *RateLimiter) GetUsedPoints() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.usedPoints
}

// GetResetTime returns the time of the next reset (for monitoring)
func (r *RateLimiter) GetResetTime() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.resetTime
}

// GetRateLimiter returns the rate limiter instance (for monitoring)
func (s *WarcraftLogsClientService) GetRateLimiter() *RateLimiter {
	return s.rateLimiter
}
