// go test -v ./internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler/
package warcraftlogsBuildsTemporalScheduler

import (
	"fmt"
	"testing"
	"time"

	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

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

// TestScheduleIDValidation tests the schedule and workflow ID formats
func TestScheduleIDValidation(t *testing.T) {
	testCases := []struct {
		name        string
		className   string
		validClass  bool
		validTiming bool
		description string
	}{
		{
			name:        "Valid Priest Schedule",
			className:   "Priest",
			validClass:  true,
			validTiming: true,
			description: "Known class with valid schedule",
		},
		{
			name:        "Valid DeathKnight Schedule",
			className:   "DeathKnight",
			validClass:  true,
			validTiming: true,
			description: "Known class with valid schedule",
		},
		{
			name:        "Invalid Class Schedule",
			className:   "InvalidClass",
			validClass:  false,
			validTiming: false,
			description: "Unknown class should be invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing case: %s - %s", tc.name, tc.description)

			// Test if class is in valid time slot
			isValid := isClassInValidTimeSlot(tc.className)
			assert.Equal(t, tc.validClass, isValid,
				"Class %s validity should be %v", tc.className, tc.validClass)

			if tc.validClass {
				// Get scheduled time
				hour := getScheduledHour(tc.className)
				day := getScheduledDay(tc.className)

				assert.NotZero(t, hour, "Hour should not be zero for valid class")
				assert.NotZero(t, day, "Day should not be zero for valid class")
				t.Logf("Class %s scheduled for day %d at %d:00",
					tc.className, day, hour)
			}
		})
	}
}

// TestFilterConfigForClass tests the configuration filtering for specific classes
func TestFilterConfigForClass(t *testing.T) {
	baseConfig := &workflows.Config{
		Specs: []workflows.ClassSpec{
			{ClassName: "Priest", SpecName: "Shadow"},
			{ClassName: "Priest", SpecName: "Holy"},
			{ClassName: "Druid", SpecName: "Balance"},
			{ClassName: "Mage", SpecName: "Frost"},
		},
		Rankings: workflows.RankingsConfig{
			MaxRankingsPerSpec: 150,
		},
	}

	testCases := []struct {
		name             string
		config           *workflows.Config
		className        string
		expectedSpecsLen int
		description      string
	}{
		{
			name:             "Filter Priest Specs",
			config:           baseConfig,
			className:        "Priest",
			expectedSpecsLen: 2,
			description:      "Should return only Priest specs",
		},
		{
			name:             "Filter Druid Specs",
			config:           baseConfig,
			className:        "Druid",
			expectedSpecsLen: 1,
			description:      "Should return only Druid spec",
		},
		{
			name:             "Non-existent Class",
			config:           baseConfig,
			className:        "Warlock",
			expectedSpecsLen: 0,
			description:      "Should return empty spec list",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing case: %s - %s", tc.name, tc.description)

			filtered := filterConfigForClass(tc.config, tc.className)
			assert.Equal(t, tc.expectedSpecsLen, len(filtered.Specs),
				"Expected %d specs for class %s, got %d",
				tc.expectedSpecsLen, tc.className, len(filtered.Specs))

			// Verify all specs belong to the correct class
			for _, spec := range filtered.Specs {
				assert.Equal(t, tc.className, spec.ClassName,
					"Spec %s should belong to class %s",
					spec.SpecName, tc.className)
			}
		})
	}
}

// TestScheduleValidation tests the schedule validation logic
func TestScheduleValidation(t *testing.T) {
	testCases := []struct {
		name          string
		className     string
		expectedHour  int
		expectedDay   int
		shouldBeValid bool
		description   string
	}{
		{
			name:          "Valid Tuesday Early Schedule",
			className:     "DeathKnight",
			expectedHour:  2,
			expectedDay:   2,
			shouldBeValid: true,
			description:   "DeathKnight should be scheduled for Tuesday 2 AM",
		},
		{
			name:          "Valid Tuesday Late Schedule",
			className:     "Hunter",
			expectedHour:  7,
			expectedDay:   2,
			shouldBeValid: true,
			description:   "Hunter should be scheduled for Tuesday 7 AM",
		},
		{
			name:          "Valid Wednesday Schedule",
			className:     "Priest",
			expectedHour:  2,
			expectedDay:   3,
			shouldBeValid: true,
			description:   "Priest should be scheduled for Wednesday 2 AM",
		},
		{
			name:          "Invalid Class",
			className:     "InvalidClass",
			expectedHour:  0,
			expectedDay:   0,
			shouldBeValid: false,
			description:   "Invalid class should not have a valid schedule",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing case: %s - %s", tc.name, tc.description)

			// Verify class validity
			isValid := isClassInValidTimeSlot(tc.className)
			assert.Equal(t, tc.shouldBeValid, isValid,
				"Class %s validity should be %v", tc.className, tc.shouldBeValid)

			if tc.shouldBeValid {
				// Verify scheduled time
				hour := getScheduledHour(tc.className)
				day := getScheduledDay(tc.className)

				assert.Equal(t, tc.expectedHour, hour,
					"Expected hour %d for class %s", tc.expectedHour, tc.className)
				assert.Equal(t, tc.expectedDay, day,
					"Expected day %d for class %s", tc.expectedDay, tc.className)
			}
		})
	}
}

// TestCronExpressionGeneration tests the generation of cron expressions
func TestCronExpressionGeneration(t *testing.T) {
	testCases := []struct {
		name             string
		className        string
		expectedCronExpr string
		description      string
	}{
		{
			name:             "DeathKnight Cron",
			className:        "DeathKnight",
			expectedCronExpr: "0 2 * * 2",
			description:      "Tuesday 2 AM schedule",
		},
		{
			name:             "Hunter Cron",
			className:        "Hunter",
			expectedCronExpr: "0 7 * * 2",
			description:      "Tuesday 7 AM schedule",
		},
		{
			name:             "Priest Cron",
			className:        "Priest",
			expectedCronExpr: "0 2 * * 3",
			description:      "Wednesday 2 AM schedule",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing case: %s - %s", tc.name, tc.description)

			hour := getScheduledHour(tc.className)
			cronExpr := fmt.Sprintf("0 %d * * %d", hour, getScheduledDay(tc.className))

			assert.Equal(t, tc.expectedCronExpr, cronExpr,
				"Expected cron expression %s for class %s",
				tc.expectedCronExpr, tc.className)
		})
	}
}
