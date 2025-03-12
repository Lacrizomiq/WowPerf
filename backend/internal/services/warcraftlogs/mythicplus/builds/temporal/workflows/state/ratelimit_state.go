// ratelimit_state.go
package warcraftlogsBuildsTemporalWorkflowsState

import "time"

// RateLimitState tracks rate limiting state
type RateLimitState struct {
	RemainingPoints float64
	RetryCount      int32
	LastCheckTime   time.Time
	ResetTime       time.Time
}

// NewRateLimitState creates a new rate limit state
func NewRateLimitState() *RateLimitState {
	return &RateLimitState{
		RemainingPoints: 18000,
		LastCheckTime:   time.Now(),
		ResetTime:       time.Now().Add(time.Hour),
	}
}
