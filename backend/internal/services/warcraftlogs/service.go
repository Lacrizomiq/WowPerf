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

// RateLimitData represents the rate limit information from the API
type RateLimitData struct {
	LimitPerHour        float64 `json:"limitPerHour"`
	PointsSpentThisHour float64 `json:"pointsSpentThisHour"`
	PointsResetIn       int     `json:"pointsResetIn"`
}

const RateLimitQuery = `query {
	rateLimitData {
			limitPerHour
			pointsSpentThisHour
			pointsResetIn
	}
}`

const (
	defaultRefreshInterval = 5 * time.Minute
	minWaitTime            = 5 * time.Second
	maxRetries             = 3
)

// WarcraftLogsClientService is a struct to manage the Warcraft Logs client
type WarcraftLogsClientService struct {
	Client          *Client
	rateLimiter     *RateLimiter
	refreshInterval time.Duration
}

// RateLimiter is a struct to manage the rate limit
type RateLimiter struct {
	mu            sync.RWMutex
	maxPoints     float64
	currentPoints float64
	resetTime     time.Time
	lastUpdate    time.Time
	pointsPerHour float64
}

func NewWarcraftLogsClientService() (*WarcraftLogsClientService, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	service := &WarcraftLogsClientService{
		Client:          client,
		refreshInterval: defaultRefreshInterval,
		rateLimiter:     &RateLimiter{},
	}

	if err := service.updateRateLimit(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize rate limit: %w", err)
	}

	go service.startPeriodicRateLimitCheck()

	return service, nil
}

// Available checks if points are available
func (r *RateLimiter) Available() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()
	if now.After(r.resetTime) {
		return true
	}

	pointsAvailable := r.maxPoints - r.currentPoints
	return pointsAvailable >= 1.0
}

// UpdateState updates the rate limiter state
func (r *RateLimiter) UpdateState(limitPerHour float64, pointsSpent float64, resetIn int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.maxPoints = limitPerHour
	r.currentPoints = pointsSpent
	r.resetTime = time.Now().Add(time.Duration(resetIn) * time.Second)
	r.lastUpdate = time.Now()
	r.pointsPerHour = limitPerHour

	log.Printf("[DEBUG] Rate limit updated - Max: %.2f, Current: %.2f, Remaining: %.2f, Reset in: %v",
		r.maxPoints,
		r.currentPoints,
		r.maxPoints-r.currentPoints,
		time.Until(r.resetTime).Round(time.Second))
}

// ConsumePoint consumes a point from the rate limiter
func (r *RateLimiter) ConsumePoint() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentPoints++
}

// GetRateLimitInfo returns the current rate limit information
func (r *RateLimiter) GetRateLimitInfo() *warcraftlogsTypes.RateLimitInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	remaining := r.maxPoints - r.currentPoints
	if remaining < 0 {
		remaining = 0
	}

	return &warcraftlogsTypes.RateLimitInfo{
		RemainingPoints: remaining,
		PointsPerHour:   int(r.pointsPerHour),
		ResetIn:         time.Until(r.resetTime),
		NextRefresh:     r.lastUpdate.Add(defaultRefreshInterval),
	}
}

// startPeriodicRateLimitCheck starts the periodic rate limit check
func (s *WarcraftLogsClientService) startPeriodicRateLimitCheck() {
	ticker := time.NewTicker(s.refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.updateRateLimit(context.Background()); err != nil {
			log.Printf("[WARN] Periodic rate limit update failed: %v", err)
		}
	}
}

// updateRateLimit updates the rate limit information from the API
func (s *WarcraftLogsClientService) updateRateLimit(ctx context.Context) error {
	response, err := s.Client.MakeGraphQLRequest(RateLimitQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to get rate limit data: %w", err)
	}

	var result struct {
		Data struct {
			RateLimitData RateLimitData `json:"rateLimitData"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("failed to unmarshal rate limit data: %w", err)
	}

	s.rateLimiter.UpdateState(
		result.Data.RateLimitData.LimitPerHour,
		result.Data.RateLimitData.PointsSpentThisHour,
		result.Data.RateLimitData.PointsResetIn,
	)

	return nil
}

// waitForAvailablePoints waits for points to be available
func (s *WarcraftLogsClientService) waitForAvailablePoints(ctx context.Context) error {
	info := s.rateLimiter.GetRateLimitInfo()
	log.Printf("[DEBUG] Points status - Remaining: %.2f, Required: 1.00, Reset in: %v",
		info.RemainingPoints,
		info.ResetIn.Round(time.Second))

	// Si on a assez de points, on y va
	if info.RemainingPoints >= 1.0 {
		return nil
	}

	// Sinon, on force une mise à jour et on vérifie à nouveau
	if err := s.updateRateLimit(ctx); err != nil {
		return fmt.Errorf("failed to update rate limit: %w", err)
	}

	info = s.rateLimiter.GetRateLimitInfo()
	if info.RemainingPoints >= 1.0 {
		return nil
	}

	// Si toujours pas de points, on attend
	if info.ResetIn > 0 {
		log.Printf("[DEBUG] Waiting %v for rate limit reset", info.ResetIn.Round(time.Second))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(info.ResetIn):
			return nil
		}
	}

	return warcraftlogsTypes.NewQuotaExceededError(info)
}

// MakeRequest makes a request to the API with rate limiting management
func (s *WarcraftLogsClientService) MakeRequest(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	log.Printf("[DEBUG] Making GraphQL request - Variables: %+v", variables)

	if err := s.waitForAvailablePoints(ctx); err != nil {
		return nil, err
	}

	response, err := s.Client.MakeGraphQLRequest(query, variables)
	if err != nil {
		if warcraftlogsTypes.IsRateLimit(err) {
			log.Printf("[WARN] Rate limit hit, forcing update")
			if updateErr := s.updateRateLimit(ctx); updateErr != nil {
				log.Printf("[ERROR] Failed to update rate limit after hit: %v", updateErr)
			}
		}
		return nil, err
	}

	s.rateLimiter.ConsumePoint()

	// Asynchronous update if necessary
	if time.Since(s.rateLimiter.lastUpdate) > s.refreshInterval {
		go func() {
			if err := s.updateRateLimit(context.Background()); err != nil {
				log.Printf("[WARN] Async rate limit update failed: %v", err)
			}
		}()
	}

	return response, nil
}
