package warcraftlogsBuildsTemporalScheduler

import (
	"time"

	"go.temporal.io/api/enums/v1"
)

// Class scheduling configuration
// Hours are cumulative:
// - Tuesday: 2 AM (2) and 7 AM (7)
// - Wednesday: 2 AM (26) because it's 24 + 2
var classScheduleTimes = map[int][]string{
	// Tuesday (day 2)
	2: {"DeathKnight", "DemonHunter", "Druid", "Evoker"},
	7: {"Hunter", "Mage", "Monk", "Paladin"},

	// Wednesday (day 3)
	26: {"Priest", "Rogue", "Shaman", "Warrior", "Warlock"},
}

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
			// %d will be replaced with specific hour for each class
			CronExpression: "0 %d * * 2,3", // Runs on Tuesday and Wednesday
			TimeZone:       "Europe/Paris",
			OverlapPolicy:  enums.SCHEDULE_OVERLAP_POLICY_SKIP,
		},
		Retry: RetryPolicy{
			InitialInterval:    15 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    5 * time.Minute,
			MaximumAttempts:    3,
		},
		Backfill: BackfillPolicy{
			Enabled:        true,
			BackfillWindow: 24 * time.Hour, // 1 day backfill window
		},
		// Increased timeout to allow for larger processing windows
		Timeout: 3 * time.Hour,
		Paused:  false,
	}
}
