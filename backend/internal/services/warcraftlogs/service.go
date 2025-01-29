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

	// Rate limiting constants
	maxPointsPerHour    = 18000
	workflowQuotaPoints = 14000
	bufferPoints        = 2000
	safetyPoints        = 2000
)

// WarcraftLogsClientService is a struct to manage the Warcraft Logs client
type WarcraftLogsClientService struct {
	Client          *Client
	rateLimiter     *RateLimiter
	refreshInterval time.Duration
}

// RateLimiter handles API rate limiting with optimized point tracking
type RateLimiter struct {
	mu            sync.RWMutex
	maxPoints     float64
	currentPoints float64
	resetTime     time.Time
	lastUpdate    time.Time
	pointsPerHour float64

	// Workflow specific tracking
	workflowQuota   float64  // 14000 points
	workflowPoints  sync.Map // map[string]float64
	activeWorkflows int32

	// API client reference
	client *Client
}

// RateLimitInfo provides current rate limit status
type RateLimitInfo struct {
	RemainingPoints float64
	PointsPerHour   int
	ResetIn         time.Duration
	NextRefresh     time.Time
}

func NewWarcraftLogsClientService() (*WarcraftLogsClientService, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	rateLimiter := &RateLimiter{
		workflowQuota: workflowQuotaPoints,
		maxPoints:     maxPointsPerHour,
		client:        client,
	}

	service := &WarcraftLogsClientService{
		Client:          client,
		refreshInterval: defaultRefreshInterval,
		rateLimiter:     rateLimiter,
	}

	if err := service.updateRateLimit(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize rate limit: %w", err)
	}

	go service.startPeriodicRateLimitCheck()

	return service, nil
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

// ConsumePoint consumes a point from the rate limiter
func (r *RateLimiter) ConsumePoint() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentPoints++
}

// ReserveWorkflowPoints attempts to reserve points for a workflow
func (r *RateLimiter) ReserveWorkflowPoints(workflowID string, points float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentPoints+points > r.workflowQuota {
		return fmt.Errorf("workflow quota exceeded")
	}

	if existingPoints, ok := r.workflowPoints.Load(workflowID); ok {
		points += existingPoints.(float64)
	}
	r.workflowPoints.Store(workflowID, points)
	return nil
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

// waitForAvailablePoints waits for points to be available with exponential backoff
func (s *WarcraftLogsClientService) waitForAvailablePoints(ctx context.Context) error {
	info := s.rateLimiter.GetRateLimitInfo()
	remaining := info.RemainingPoints

	log.Printf("[DEBUG] Points status - Remaining: %.2f, Required: 1.00, Reset in: %v",
		remaining,
		info.ResetIn.Round(time.Second))

	if remaining >= 1.0 {
		return nil
	}

	// Implement exponential backoff for retries
	backoff := minWaitTime
	for retries := 0; retries < maxRetries; retries++ {
		// Wait using backoff
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}

		// Update and check again
		if err := s.updateRateLimit(ctx); err != nil {
			return fmt.Errorf("failed to update rate limit: %w", err)
		}

		info = s.rateLimiter.GetRateLimitInfo()
		if info.RemainingPoints >= 1.0 {
			return nil
		}

		// Double the backoff for next iteration
		backoff *= 2
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

	// Only trigger async update if needed
	if time.Since(s.rateLimiter.lastUpdate) > s.refreshInterval {
		go func() {
			if err := s.updateRateLimit(context.Background()); err != nil {
				log.Printf("[WARN] Async rate limit update failed: %v", err)
			}
		}()
	}

	return response, nil
}

// GetWorkflowQuotaRemaining returns remaining points available for workflows
func (r *RateLimiter) GetWorkflowQuotaRemaining() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.workflowQuota - r.currentPoints
}

// ReleaseWorkflowPoints releases points reserved for a workflow
func (r *RateLimiter) ReleaseWorkflowPoints(workflowID string) {
	if points, ok := r.workflowPoints.Load(workflowID); ok {
		r.mu.Lock()
		r.currentPoints -= points.(float64)
		r.mu.Unlock()
		r.workflowPoints.Delete(workflowID)
	}
}

func (s *WarcraftLogsClientService) GetRateLimiter() *RateLimiter {
	return s.rateLimiter
}
