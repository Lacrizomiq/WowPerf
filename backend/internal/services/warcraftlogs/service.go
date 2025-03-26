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
	mu            sync.RWMutex
	maxPoints     float64       // Maximum points allowed per hour
	usedPoints    float64       // Currently used points based on last API check
	resetTime     time.Time     // Time of next reset based on API
	lastCheckTime time.Time     // Last time we checked the rate limit
	checkInterval time.Duration // Minimum interval between checks
	initialized   bool          // If the rate limiter has been initialized
}

// NewWarcraftLogsClientService creates a new service instance
func NewWarcraftLogsClientService() (*WarcraftLogsClientService, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	rateLimiter := &RateLimiter{
		maxPoints:     18000, // Conservative default value
		usedPoints:    0,
		resetTime:     time.Now().Add(time.Hour),
		lastCheckTime: time.Now().Add(-time.Hour), // Force a check at the first request
		checkInterval: time.Minute * 1,            // Check at most once per minute
		initialized:   false,
	}

	service := &WarcraftLogsClientService{
		Client:      client,
		rateLimiter: rateLimiter,
	}

	return service, nil
}

// fetchRateLimitData fetches the rate limit data from the WarcraftLogs API
func (s *WarcraftLogsClientService) fetchRateLimitData() (*RateLimitData, error) {
	response, err := s.Client.MakeGraphQLRequest(RateLimitQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit data: %w", err)
	}

	var result struct {
		Data struct {
			RateLimitData RateLimitData `json:"rateLimitData"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to parse rate limit data: %w", err)
	}

	// Validate and correct values if necessary
	data := result.Data.RateLimitData
	if data.LimitPerHour <= 0 {
		data.LimitPerHour = 18000
		log.Printf("[WARN] Invalid limitPerHour: %.2f, using default: 18000", data.LimitPerHour)
	}
	if data.PointsResetIn <= 0 {
		data.PointsResetIn = 3600 // 1 hour default
		log.Printf("[WARN] Invalid resetIn: %d, using default: 3600", data.PointsResetIn)
	}

	return &data, nil
}

// MakeRequest performs a rate-limited API request with real-time point checking
func (s *WarcraftLogsClientService) MakeRequest(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	if query == "" {
		return nil, fmt.Errorf("empty query not allowed")
	}

	// Avoid recursion for rate limit queries
	if query == RateLimitQuery {
		return s.Client.MakeGraphQLRequest(query, variables)
	}
	// Check if we need to update rate limit info
	shouldCheck := false
	s.rateLimiter.mu.RLock()
	shouldCheck = !s.rateLimiter.initialized || time.Since(s.rateLimiter.lastCheckTime) > s.rateLimiter.checkInterval
	s.rateLimiter.mu.RUnlock()

	if shouldCheck {
		// Make a request to get the latest info
		rateLimitData, err := s.fetchRateLimitData()
		if err != nil {
			log.Printf("[WARN] Failed to get rate limit data: %v, using existing values", err)
		} else {
			s.rateLimiter.mu.Lock()
			s.rateLimiter.maxPoints = rateLimitData.LimitPerHour
			s.rateLimiter.usedPoints = rateLimitData.PointsSpentThisHour
			s.rateLimiter.resetTime = time.Now().Add(time.Duration(rateLimitData.PointsResetIn) * time.Second)
			s.rateLimiter.lastCheckTime = time.Now()
			s.rateLimiter.initialized = true
			s.rateLimiter.mu.Unlock()

			log.Printf("[DEBUG] Rate limit updated - Max: %.2f, Used: %.2f, Remaining: %.2f, Reset in: %v",
				rateLimitData.LimitPerHour,
				rateLimitData.PointsSpentThisHour,
				rateLimitData.LimitPerHour-rateLimitData.PointsSpentThisHour,
				time.Until(s.rateLimiter.resetTime))
		}
	}

	// Check if we have enough points (with a safety margin)
	s.rateLimiter.mu.RLock()
	remaining := s.rateLimiter.maxPoints - s.rateLimiter.usedPoints
	resetIn := time.Until(s.rateLimiter.resetTime)
	s.rateLimiter.mu.RUnlock()

	if remaining < 500.0 {
		info := &warcraftlogsTypes.RateLimitInfo{
			RemainingPoints: remaining,
			ResetIn:         resetIn,
		}
		return nil, warcraftlogsTypes.NewQuotaExceededError(info)
	}

	// Make the actual request
	response, err := s.Client.MakeGraphQLRequest(query, variables)
	if err != nil {
		if warcraftlogsTypes.IsRateLimit(err) {
			s.rateLimiter.mu.RLock()
			info := &warcraftlogsTypes.RateLimitInfo{
				RemainingPoints: s.rateLimiter.maxPoints - s.rateLimiter.usedPoints,
				ResetIn:         time.Until(s.rateLimiter.resetTime),
			}
			s.rateLimiter.mu.RUnlock()
			return nil, warcraftlogsTypes.NewQuotaExceededError(info)
		}
		return nil, err
	}

	return response, nil
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
	return r.lastCheckTime
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
