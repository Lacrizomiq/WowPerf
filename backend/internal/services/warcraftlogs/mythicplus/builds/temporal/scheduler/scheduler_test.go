// go test -v ./internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler/
package warcraftlogsBuildsTemporalScheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestScheduleTimeFunctions tests the core scheduling time functions
func TestScheduleTimeFunctions(t *testing.T) {
	testCases := []struct {
		name          string
		className     string
		expectedHour  int
		expectedDay   int
		expectedValid bool
	}{
		{
			name:          "DeathKnight Schedule",
			className:     "DeathKnight",
			expectedHour:  2,
			expectedDay:   2,
			expectedValid: true,
		},
		{
			name:          "Hunter Schedule",
			className:     "Hunter",
			expectedHour:  7,
			expectedDay:   2,
			expectedValid: true,
		},
		{
			name:          "Priest Schedule",
			className:     "Priest",
			expectedHour:  2,
			expectedDay:   3,
			expectedValid: true,
		},
		{
			name:          "Invalid Class",
			className:     "InvalidClass",
			expectedHour:  2,
			expectedDay:   2,
			expectedValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test getScheduledHour
			hour := getScheduledHour(tc.className)
			assert.Equal(t, tc.expectedHour, hour,
				"Expected hour %d for class %s", tc.expectedHour, tc.className)

			// Test getScheduledDay
			day := getScheduledDay(tc.className)
			assert.Equal(t, tc.expectedDay, day,
				"Expected day %d for class %s", tc.expectedDay, tc.className)

			// Test isClassInValidTimeSlot
			isValid := isClassInValidTimeSlot(tc.className)
			assert.Equal(t, tc.expectedValid, isValid,
				"Expected validity %v for class %s", tc.expectedValid, tc.className)
		})
	}
}

// TestDefaultScheduleOptions tests the default schedule options
func TestDefaultScheduleOptions(t *testing.T) {
	options := DefaultScheduleOptions()

	t.Run("Default Policy Values", func(t *testing.T) {
		assert.Contains(t, options.Policy.CronExpression, "%d",
			"CronExpression should contain placeholder for hour")
		assert.Equal(t, "Europe/Paris", options.Policy.TimeZone,
			"Default timezone should be Europe/Paris")
	})

	t.Run("Default Retry Policy", func(t *testing.T) {
		assert.Equal(t, time.Second*15, options.Retry.InitialInterval,
			"Initial retry interval should be 15 seconds")
		assert.Equal(t, 2.0, options.Retry.BackoffCoefficient,
			"Backoff coefficient should be 2.0")
		assert.Equal(t, 3, options.Retry.MaximumAttempts,
			"Maximum attempts should be 3")
	})

	t.Run("Default Backfill Policy", func(t *testing.T) {
		assert.True(t, options.Backfill.Enabled,
			"Backfill should be enabled by default")
		assert.Equal(t, 24*time.Hour, options.Backfill.BackfillWindow,
			"Backfill window should be 24 hours")
	})

	t.Run("Default State", func(t *testing.T) {
		assert.False(t, options.Paused,
			"Schedule should not be paused by default")
		assert.Equal(t, 3*time.Hour, options.Timeout,
			"Default timeout should be 3 hours")
	})
}

// TestClassScheduleTimes tests the time slot distribution
func TestClassScheduleTimes(t *testing.T) {
	// Test that we don't have overlapping classes in same time slot
	timeSlots := make(map[string][]string)

	for hour, classes := range classScheduleTimes {
		day := (hour / 24) + 2 // Convert to day
		hour = hour % 24       // Convert to 24-hour format
		slotKey := fmt.Sprintf("Day%d-Hour%d", day, hour)

		t.Logf("Time slot %s has classes: %v", slotKey, classes)
		assert.LessOrEqual(t, len(classes), 5,
			"Time slot %s should not have more than 5 classes", slotKey)

		timeSlots[slotKey] = classes
	}

	// Verify no class appears more than once
	classCount := make(map[string]int)
	for _, classes := range timeSlots {
		for _, class := range classes {
			classCount[class]++
			assert.Equal(t, 1, classCount[class],
				"Class %s appears in multiple time slots", class)
		}
	}
}
