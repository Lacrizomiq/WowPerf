package warcraftlogsBuildsTemporalScheduler

import (
	"fmt"
	"time"
)

// MainScheduleConfig defines the configuration for the weekly sync
type MainScheduleConfig struct {
	Hour      int    // Hour in 24-hour UTC format
	Day       int    // Day of week (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	TaskQueue string // Task queue name for the workflow
}

// DefaultScheduleConfig provides the default weekly schedule (Tuesday 2 AM UTC)
var DefaultScheduleConfig = MainScheduleConfig{
	Hour:      2, // 2 AM
	Day:       2, // Tuesday
	TaskQueue: "warcraft-logs-sync",
}

// RetryPolicy defines how failures are handled
type RetryPolicy struct {
	InitialInterval    time.Duration // Initial retry interval
	BackoffCoefficient float64       // Multiplier for subsequent retries
	MaximumInterval    time.Duration // Maximum retry interval
	MaximumAttempts    int           // Maximum number of retry attempts
}

// ScheduleOptions combines all configuration options
type ScheduleOptions struct {
	Retry   RetryPolicy   // Retry handling policy
	Timeout time.Duration // Maximum execution time
	Paused  bool          // Whether the schedule starts paused
}

// DefaultScheduleOptions returns the default configuration
func DefaultScheduleOptions() *ScheduleOptions {
	return &ScheduleOptions{
		Retry: RetryPolicy{
			InitialInterval:    time.Minute,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Hour,
			MaximumAttempts:    5,
		},
		Timeout: 24 * time.Hour, // Allowing enough time for sequential processing
		Paused:  false,
	}
}

// ValidateScheduleConfig validates the schedule configuration
func ValidateScheduleConfig(config *MainScheduleConfig) error {
	if config == nil {
		return fmt.Errorf("schedule config cannot be nil")
	}
	if config.Hour < 0 || config.Hour > 23 {
		return fmt.Errorf("invalid hour: %d", config.Hour)
	}
	if config.Day < 0 || config.Day > 6 {
		return fmt.Errorf("invalid day: %d", config.Day)
	}
	if config.TaskQueue == "" {
		return fmt.Errorf("task queue cannot be empty")
	}
	return nil
}
