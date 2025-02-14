package warcraftlogsBuildsTemporalScheduler

import (
	"fmt"
	"time"

	"go.temporal.io/api/enums/v1"
)

// TimeSlot represents a specific scheduling time slot for classes
// It defines when certain classes should have their workflows executed
type TimeSlot struct {
	ID          string   // Unique identifier for the slot (e.g., "tuesday-2am")
	Hour        int      // Hour in 24-hour UTC format
	Day         int      // Day of week (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	Classes     []string // List of WoW classes scheduled for this slot
	Description string   // Human-readable description
}

// Pre-defined time slots for class scheduling
// All times are in UTC
var (
	// Tuesday2AM represents the early morning slot on Tuesday (2 AM UTC)
	Tuesday2AM = TimeSlot{
		ID:          "tuesday-2am",
		Hour:        2,
		Day:         2, // Tuesday
		Classes:     []string{"DeathKnight", "DemonHunter", "Druid", "Evoker"},
		Description: "Tuesday 2 AM UTC",
	}

	// Tuesday7AM represents the morning slot on Tuesday (7 AM UTC)
	Tuesday7AM = TimeSlot{
		ID:          "tuesday-7am",
		Hour:        7,
		Day:         2, // Tuesday
		Classes:     []string{"Hunter", "Mage", "Monk", "Paladin"},
		Description: "Tuesday 7 AM UTC",
	}

	// Wednesday2AM represents the early morning slot on Wednesday (2 AM UTC)
	Wednesday2AM = TimeSlot{
		ID:          "wednesday-2am",
		Hour:        2,
		Day:         3, // Wednesday
		Classes:     []string{"Priest", "Rogue", "Shaman", "Warrior", "Warlock"},
		Description: "Wednesday 2 AM UTC",
	}
)

// ScheduleSlots contains all available time slots
var ScheduleSlots = []TimeSlot{
	Tuesday2AM,
	Tuesday7AM,
	Wednesday2AM,
}

// SchedulePolicy defines the core scheduling behavior
type SchedulePolicy struct {
	CronExpression string                      // Cron expression for the schedule
	TimeZone       string                      // Timezone for the schedule
	OverlapPolicy  enums.ScheduleOverlapPolicy // Policy for handling overlapping schedules
}

// RetryPolicy defines how failures are handled
type RetryPolicy struct {
	InitialInterval    time.Duration // Initial retry interval
	BackoffCoefficient float64       // Multiplier for subsequent retries
	MaximumInterval    time.Duration // Maximum retry interval
	MaximumAttempts    int           // Maximum number of retry attempts
}

// BackfillPolicy defines how missed executions are handled
type BackfillPolicy struct {
	Enabled        bool          // Whether backfill is enabled
	BackfillWindow time.Duration // How far back to look for missed executions
}

// ScheduleOptions combines all configuration options
type ScheduleOptions struct {
	Policy   SchedulePolicy // Core scheduling policy
	Retry    RetryPolicy    // Retry handling policy
	Backfill BackfillPolicy // Backfill policy
	Timeout  time.Duration  // Maximum execution time
	Paused   bool           // Whether the schedule starts paused
}

// DefaultScheduleOptions returns the default configuration
func DefaultScheduleOptions() *ScheduleOptions {
	return &ScheduleOptions{
		Policy: SchedulePolicy{
			// %d will be replaced with specific hour for each class
			CronExpression: "0 %d * * 2,3", // Runs on Tuesday and Wednesday in UTC
			TimeZone:       "UTC",          // Changed from Europe/Paris to UTC
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
		Timeout: 3 * time.Hour, // Increased timeout for larger processing windows
		Paused:  false,
	}
}

// FindTimeSlotForClass returns the time slot for a given class
// Returns nil if the class is not scheduled in any slot
func FindTimeSlotForClass(className string) *TimeSlot {
	for i := range ScheduleSlots {
		slot := &ScheduleSlots[i]
		for _, class := range slot.Classes {
			if class == className {
				return slot
			}
		}
	}
	return nil
}

// GetCurrentTimeSlot returns the appropriate time slot for the current time
// Returns nil if current time is not in any defined slot
// All comparisons are done in UTC
func GetCurrentTimeSlot() *TimeSlot {
	now := time.Now().UTC()
	currentHour := now.Hour()
	currentDay := int(now.Weekday())

	for i := range ScheduleSlots {
		slot := &ScheduleSlots[i]
		if slot.Day == currentDay && slot.Hour == currentHour {
			return slot
		}
	}
	return nil
}

// ValidateTimeSlot checks if a time slot configuration is valid
func ValidateTimeSlot(slot *TimeSlot) error {
	if slot == nil {
		return fmt.Errorf("time slot cannot be nil")
	}
	if slot.Hour < 0 || slot.Hour > 23 {
		return fmt.Errorf("invalid hour: %d", slot.Hour)
	}
	if slot.Day < 0 || slot.Day > 6 {
		return fmt.Errorf("invalid day: %d", slot.Day)
	}
	if len(slot.Classes) == 0 {
		return fmt.Errorf("no classes defined for slot %s", slot.ID)
	}
	return nil
}
