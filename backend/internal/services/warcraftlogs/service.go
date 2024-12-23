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
	LimitPerHour        int     `json:"limitPerHour"`
	PointsSpentThisHour float64 `json:"pointsSpentThisHour"`
	PointsResetIn       int     `json:"pointsResetIn"`
}

// RateLimitQuery is the GraphQL query to fetch rate limit data
const RateLimitQuery = `
query getRateLimitData {
    rateLimitData {
        limitPerHour
        pointsSpentThisHour
        pointsResetIn
    }
}`

const (
	requestTimeout         = 10 * time.Second
	checkRateLimitInterval = 5 * time.Minute
)

// WarcraftLogsClientService handles API calls with intelligent rate limiting
type WarcraftLogsClientService struct {
	Client             *Client
	currentPoints      float64
	maxPointsHour      int
	resetTime          time.Time
	lastCheck          time.Time
	mu                 sync.RWMutex
	minPointsThreshold float64       // Threshold before slowing down requests
	refreshInterval    time.Duration // Rate limit refresh interval
}

// NewWarcraftLogsClientService creates a new service instance
func NewWarcraftLogsClientService() (*WarcraftLogsClientService, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	service := &WarcraftLogsClientService{
		Client:             client,
		minPointsThreshold: 10.0,
		refreshInterval:    time.Minute * 5,
	}

	// Initial rate limit check
	if err := service.updateRateLimit(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to update rate limit: %w", err)
	}

	// Start periodic rate limit check
	go service.startPeriodicRateLimitCheck()

	return service, nil
}

// startPeriodicRateLimitCheck starts the periodic rate limit check
func (s *WarcraftLogsClientService) startPeriodicRateLimitCheck() {
	ticker := time.NewTicker(checkRateLimitInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.updateRateLimit(context.Background()); err != nil {
			log.Printf("[ERROR] Failed to update rate limit: %v", err)
		}
	}
}

// getRateLimitInfo returns current rate limit information
func (s *WarcraftLogsClientService) getRateLimitInfo() *warcraftlogsTypes.RateLimitInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &warcraftlogsTypes.RateLimitInfo{
		RemainingPoints: float64(s.maxPointsHour) - s.currentPoints,
		PointsPerHour:   s.maxPointsHour,
		ResetIn:         time.Until(s.resetTime),
		NextRefresh:     s.lastCheck.Add(s.refreshInterval),
	}
}

// shouldUpdateRateLimit determines if rate limit info should be updated
func (s *WarcraftLogsClientService) shouldUpdateRateLimit() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	timeSinceLastCheck := time.Since(s.lastCheck)
	remainingPoints := float64(s.maxPointsHour) - s.currentPoints

	return timeSinceLastCheck >= s.refreshInterval || remainingPoints < s.minPointsThreshold
}

// updateRateLimit updates rate limit information from the API
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

	s.mu.Lock()
	s.currentPoints = result.Data.RateLimitData.PointsSpentThisHour
	s.maxPointsHour = result.Data.RateLimitData.LimitPerHour
	s.resetTime = time.Now().Add(time.Duration(result.Data.RateLimitData.PointsResetIn) * time.Second)
	s.lastCheck = time.Now()
	s.mu.Unlock()

	log.Printf("[DEBUG] Rate limit updated - Max: %d, Current: %.2f, Reset in: %s",
		s.maxPointsHour,
		s.currentPoints,
		time.Until(s.resetTime))

	return nil
}

// waitForPoints waits until enough points are available
func (s *WarcraftLogsClientService) waitForPoints(ctx context.Context, cost float64) error {
	for {
		if s.shouldUpdateRateLimit() {
			if err := s.updateRateLimit(ctx); err != nil {
				return fmt.Errorf("failed to update rate limit: %w", err)
			}
		}

		info := s.getRateLimitInfo()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if info.RemainingPoints >= cost {
				return nil
			}

			if info.RemainingPoints < s.minPointsThreshold {
				return warcraftlogsTypes.NewQuotaExceededError(info)
			}

			waitTime := calculateWaitTime(info, cost)
			log.Printf("[DEBUG] Rate limit approaching threshold. Waiting %v before next request", waitTime)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitTime):
				continue
			}
		}
	}
}

// calculateWaitTime determines how long to wait before the next request
func calculateWaitTime(info *warcraftlogsTypes.RateLimitInfo, cost float64) time.Duration {
	if info.RemainingPoints < cost {
		return info.ResetIn
	}

	// Calculate a proportional wait based on remaining points
	pointRatio := cost / info.RemainingPoints
	baseWait := time.Second * 5
	return time.Duration(float64(baseWait) * pointRatio)
}

// MakeRequest makes a request to the WarcraftLogs API with rate limiting
func (s *WarcraftLogsClientService) MakeRequest(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	const requestCost = 1.0

	if err := s.waitForPoints(ctx, requestCost); err != nil {
		return nil, err
	}

	response, err := s.Client.MakeGraphQLRequest(query, variables)
	if err != nil {
		if warcraftlogsTypes.IsRateLimit(err) {
			go s.updateRateLimit(context.Background())
		}
		return nil, err
	}

	s.mu.Lock()
	s.currentPoints += requestCost
	s.mu.Unlock()

	// Asynchronous rate limit update if needed
	if s.shouldUpdateRateLimit() {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			if err := s.updateRateLimit(ctx); err != nil {
				log.Printf("[WARN] Failed to update rate limit: %v", err)
			}
		}()
	}

	return response, nil
}
