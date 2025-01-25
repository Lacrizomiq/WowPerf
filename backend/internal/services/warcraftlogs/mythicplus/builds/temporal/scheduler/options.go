package warcraftlogsBuildsTemporalScheduler

import (
	"time"

	"go.temporal.io/api/enums/v1"
)

// SchedulePolicy defines the core scheduling behavior
type SchedulePolicy struct {
	CronExpression string // Cron expression for the schedule
	TimeZone       string // Timezone for the schedule
	OverlapPolicy  enums.ScheduleOverlapPolicy
}

// RetryPolicy defines how failures are handled
type RetryPolicy struct {
	InitialInterval    time.Duration
	BackoffCoefficient float64
	MaximumInterval    time.Duration
	MaximumAttempts    int
}

// BackfillPolicy defines how missed executions are handled
type BackfillPolicy struct {
	Enabled        bool
	BackfillWindow time.Duration
}

// ScheduleOptions combines all configuration options
type ScheduleOptions struct {
	Policy   SchedulePolicy
	Retry    RetryPolicy
	Backfill BackfillPolicy
	Timeout  time.Duration
	Paused   bool
}

// DefaultScheduleOptions returns the default configuration
func DefaultScheduleOptions() *ScheduleOptions {
	return &ScheduleOptions{
		Policy: SchedulePolicy{
			CronExpression: "0 7 * * 2", // Tuesday 7am
			TimeZone:       "Europe/Paris",
			OverlapPolicy:  enums.SCHEDULE_OVERLAP_POLICY_SKIP,
		},
		Retry: RetryPolicy{
			InitialInterval:    30 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Minute,
			MaximumAttempts:    3,
		},
		Backfill: BackfillPolicy{
			Enabled:        true,
			BackfillWindow: 7 * 24 * time.Hour, // 1 week
		},
		Timeout: 24 * time.Hour,
		Paused:  false,
	}
}
