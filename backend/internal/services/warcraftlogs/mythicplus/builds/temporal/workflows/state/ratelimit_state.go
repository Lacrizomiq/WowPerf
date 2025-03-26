// ratelimit_state.go
package warcraftlogsBuildsTemporalWorkflowsState

import (
	"time"

	warcraftlogsTypes "wowperf/internal/services/warcraftlogs/types"
)

// RateLimitState stores rate limit state for Temporal
type RateLimitState struct {
	RemainingPoints float64   // Points remaining from last check
	ResetTime       time.Time // Predicted reset time
	LastCheckTime   time.Time // Last check time
}

// NewRateLimitState creates a new state with default values
func NewRateLimitState() *RateLimitState {
	return &RateLimitState{
		RemainingPoints: 18000,
		ResetTime:       time.Now().Add(time.Hour),
		LastCheckTime:   time.Now(),
	}
}

// UpdateFrom updates the state from the RateLimitInfo
func (s *RateLimitState) UpdateFrom(info *warcraftlogsTypes.RateLimitInfo) {
	s.RemainingPoints = info.RemainingPoints
	s.ResetTime = time.Now().Add(info.ResetIn)
	s.LastCheckTime = time.Now()
}
