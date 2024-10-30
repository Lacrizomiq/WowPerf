package warcraftlogs

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

const (
	requestsPerSecond = 2
	burstLimit        = 5
	requestTimeout    = 10 * time.Second
)

type WarcraftLogsClientService struct {
	Client      *Client
	RateLimiter *rate.Limiter
}

func NewWarcraftLogsClientService() (*WarcraftLogsClientService, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	return &WarcraftLogsClientService{
		Client:      client,
		RateLimiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burstLimit),
	}, nil
}

// MakeRequest fait un appel à l'API avec rate limiting
func (s *WarcraftLogsClientService) MakeRequest(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	if err := s.RateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	return s.Client.MakeGraphQLRequest(query, variables)
}
