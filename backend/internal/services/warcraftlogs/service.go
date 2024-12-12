package warcraftlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// RateLimitData is the data structure for the rate limit data
type RateLimitData struct {
	LimitPerHour        int     `json:"limitPerHour"`
	PointsSpentThisHour float64 `json:"pointsSpentThisHour"`
	PointsResetIn       int     `json:"pointsResetIn"`
}

// RateLimitQuery is the query to get the rate limit data
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

type WarcraftLogsClientService struct {
	Client        *Client
	currentPoints float64
	maxPointsHour int
	resetTime     time.Time
	mu            sync.RWMutex
}

func NewWarcraftLogsClientService() (*WarcraftLogsClientService, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	service := &WarcraftLogsClientService{
		Client: client,
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
			log.Printf("Failed to update rate limit: %v", err)
		}
	}
}

// updateRateLimit updates the rate limit data
func (s *WarcraftLogsClientService) updateRateLimit(ctx context.Context) error {
	response, err := s.Client.MakeGraphQLRequest(RateLimitQuery, nil)
	if err != nil {
		return fmt.Errorf("Failed to get rate limit data: %w", err)
	}

	var result struct {
		Data struct {
			RateLimitData RateLimitData `json:"rateLimitData"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("Failed to unmarshal rate limit data: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentPoints = result.Data.RateLimitData.PointsSpentThisHour
	s.maxPointsHour = result.Data.RateLimitData.LimitPerHour
	s.resetTime = time.Now().Add(time.Duration(result.Data.RateLimitData.PointsResetIn) * time.Second)

	log.Printf("Rate limit updated - Max: %d, Current: %.2f, Reset in: %s",
		s.maxPointsHour,
		s.currentPoints,
		time.Until(s.resetTime))

	return nil
}

// waitForPoints waits for the required points to be available
func (s *WarcraftLogsClientService) waitForPoints(ctx context.Context, cost float64) error {
	for {
		s.mu.RLock()
		remaining := float64(s.maxPointsHour) - s.currentPoints
		resetTime := s.resetTime
		s.mu.RUnlock()

		if remaining >= cost {
			return nil
		}

		waitTime := time.Until(resetTime)
		if waitTime <= 0 {
			// if the reset time has passed, we need to update the rate limit data
			if err := s.updateRateLimit(ctx); err != nil {
				return fmt.Errorf("Failed to update rate limit: %w", err)
			}
			continue
		}

		log.Printf("Rate limit exceeded. Waiting for %s before retrying...", waitTime)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			if err := s.updateRateLimit(ctx); err != nil {
				return fmt.Errorf("Failed to update rate limit: %w", err)
			}
		}
	}
}

// MakeRequest makes a request to the Warcraft Logs API with rate limiting
func (s *WarcraftLogsClientService) MakeRequest(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	// Assuming the cost of the request is 1 point
	const requestCost = 1.0

	if err := s.waitForPoints(ctx, requestCost); err != nil {
		return nil, fmt.Errorf("rate limit wait error: %w", err)
	}

	response, err := s.Client.MakeGraphQLRequest(query, variables)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.currentPoints += requestCost
	s.mu.Unlock()

	// Update the rate limit data
	go func() {
		if err := s.updateRateLimit(ctx); err != nil {
			log.Printf("Failed to update rate limit: %v", err)
		}
	}()

	return response, nil
}
