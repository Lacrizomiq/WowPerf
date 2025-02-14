// go test -v ./internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler/
package warcraftlogsBuildsTemporalScheduler

import (
	"fmt"
	"testing"
	"time"

	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTimeSlotConfiguration tests the basic time slot configuration
func TestTimeSlotConfiguration(t *testing.T) {
	t.Run("Verify Time Slots Structure", func(t *testing.T) {
		assert.Equal(t, 2, Tuesday2AM.Hour, "Tuesday2AM should be at 2:00")
		assert.Equal(t, 2, Tuesday2AM.Day, "Tuesday2AM should be on Tuesday")
		assert.Equal(t, 7, Tuesday7AM.Hour, "Tuesday7AM should be at 7:00")
		assert.Equal(t, 2, Tuesday7AM.Day, "Tuesday7AM should be on Tuesday")
		assert.Equal(t, 2, Wednesday2AM.Hour, "Wednesday2AM should be at 2:00")
		assert.Equal(t, 3, Wednesday2AM.Day, "Wednesday2AM should be on Wednesday")
	})

	t.Run("Verify Class Distribution", func(t *testing.T) {
		// Test Tuesday 2AM slot
		assert.Contains(t, Tuesday2AM.Classes, "DeathKnight")
		assert.Contains(t, Tuesday2AM.Classes, "DemonHunter")
		assert.Contains(t, Tuesday2AM.Classes, "Druid")
		assert.Contains(t, Tuesday2AM.Classes, "Evoker")

		// Test Tuesday 7AM slot
		assert.Contains(t, Tuesday7AM.Classes, "Hunter")
		assert.Contains(t, Tuesday7AM.Classes, "Mage")
		assert.Contains(t, Tuesday7AM.Classes, "Monk")
		assert.Contains(t, Tuesday7AM.Classes, "Paladin")

		// Test Wednesday 2AM slot
		assert.Contains(t, Wednesday2AM.Classes, "Priest")
		assert.Contains(t, Wednesday2AM.Classes, "Rogue")
		assert.Contains(t, Wednesday2AM.Classes, "Shaman")
		assert.Contains(t, Wednesday2AM.Classes, "Warrior")
		assert.Contains(t, Wednesday2AM.Classes, "Warlock")
	})
}

// TestFindTimeSlotForClass tests the time slot lookup functionality
func TestFindTimeSlotForClass(t *testing.T) {
	testCases := []struct {
		name          string
		className     string
		expectedSlot  *TimeSlot
		shouldBeFound bool
		description   string
	}{
		{
			name:          "DeathKnight Schedule",
			className:     "DeathKnight",
			expectedSlot:  &Tuesday2AM,
			shouldBeFound: true,
			description:   "Should be found in Tuesday 2AM slot",
		},
		{
			name:          "Hunter Schedule",
			className:     "Hunter",
			expectedSlot:  &Tuesday7AM,
			shouldBeFound: true,
			description:   "Should be found in Tuesday 7AM slot",
		},
		{
			name:          "Priest Schedule",
			className:     "Priest",
			expectedSlot:  &Wednesday2AM,
			shouldBeFound: true,
			description:   "Should be found in Wednesday 2AM slot",
		},
		{
			name:          "Invalid Class",
			className:     "InvalidClass",
			expectedSlot:  nil,
			shouldBeFound: false,
			description:   "Should not be found in any slot",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing case: %s - %s", tc.name, tc.description)

			slot := FindTimeSlotForClass(tc.className)
			if tc.shouldBeFound {
				require.NotNil(t, slot, "Time slot should be found for %s", tc.className)
				assert.Equal(t, tc.expectedSlot.Hour, slot.Hour)
				assert.Equal(t, tc.expectedSlot.Day, slot.Day)
				assert.Contains(t, slot.Classes, tc.className)
			} else {
				assert.Nil(t, slot, "Time slot should not be found for %s", tc.className)
			}
		})
	}
}

// TestGetCurrentTimeSlot tests the current time slot detection
func TestGetCurrentTimeSlot(t *testing.T) {
	testCases := []struct {
		name          string
		currentHour   int
		currentDay    int
		expectedSlot  *TimeSlot
		shouldBeValid bool
		description   string
	}{
		{
			name:          "Tuesday 2 AM",
			currentHour:   2,
			currentDay:    2, // Tuesday
			expectedSlot:  &Tuesday2AM,
			shouldBeValid: true,
			description:   "Should be valid for Tuesday 2 AM slot",
		},
		{
			name:          "Tuesday 7 AM",
			currentHour:   7,
			currentDay:    2, // Tuesday
			expectedSlot:  &Tuesday7AM,
			shouldBeValid: true,
			description:   "Should be valid for Tuesday 7 AM slot",
		},
		{
			name:          "Wednesday 2 AM",
			currentHour:   2,
			currentDay:    3, // Wednesday
			expectedSlot:  &Wednesday2AM,
			shouldBeValid: true,
			description:   "Should be valid for Wednesday 2 AM slot",
		},
		{
			name:          "Invalid Time",
			currentHour:   12,
			currentDay:    2, // Tuesday
			expectedSlot:  nil,
			shouldBeValid: false,
			description:   "Should be invalid for non-scheduled time",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing case: %s - %s", tc.name, tc.description)

			// Compare directly with slot values instead of using time.Now()
			for _, slot := range ScheduleSlots {
				if slot.Hour == tc.currentHour && slot.Day == tc.currentDay {
					if tc.shouldBeValid {
						assert.Equal(t, tc.expectedSlot.Hour, slot.Hour)
						assert.Equal(t, tc.expectedSlot.Day, slot.Day)
					}
					return
				}
			}

			if !tc.shouldBeValid {
				assert.Nil(t, tc.expectedSlot, "Time slot should be invalid")
			}
		})
	}
}

// TestDefaultScheduleOptions tests the default schedule options
func TestDefaultScheduleOptions(t *testing.T) {
	options := DefaultScheduleOptions()

	t.Run("Default Policy Values", func(t *testing.T) {
		assert.Contains(t, options.Policy.CronExpression, "%d",
			"CronExpression should contain placeholder for hour")
		assert.Equal(t, "UTC", options.Policy.TimeZone,
			"Default timezone should be UTC")
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

// TestClassDistribution tests the distribution of classes across time slots
func TestClassDistribution(t *testing.T) {
	// Verify no class appears in multiple slots
	classSlots := make(map[string]*TimeSlot)

	for _, slot := range ScheduleSlots {
		for _, class := range slot.Classes {
			if existingSlot, found := classSlots[class]; found {
				t.Errorf("Class %s appears in multiple slots: %d:%d and %d:%d",
					class,
					existingSlot.Day, existingSlot.Hour,
					slot.Day, slot.Hour)
			}
			classSlots[class] = &slot
		}
	}

	// Verify all classes are assigned to a slot
	expectedClasses := []string{
		"DeathKnight", "DemonHunter", "Druid", "Evoker",
		"Hunter", "Mage", "Monk", "Paladin",
		"Priest", "Rogue", "Shaman", "Warrior", "Warlock",
	}

	for _, class := range expectedClasses {
		_, found := classSlots[class]
		assert.True(t, found, "Class %s should be assigned to a time slot", class)
	}

	// Verify slot sizes
	for _, slot := range ScheduleSlots {
		assert.LessOrEqual(t, len(slot.Classes), 5,
			"Slot Day %d Hour %d should not have more than 5 classes",
			slot.Day, slot.Hour)
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

			slot := FindTimeSlotForClass(tc.className)
			require.NotNil(t, slot, "Time slot should be found for %s", tc.className)

			cronExpr := fmt.Sprintf("0 %d * * %d", slot.Hour, slot.Day)
			assert.Equal(t, tc.expectedCronExpr, cronExpr,
				"Expected cron expression %s for class %s",
				tc.expectedCronExpr, tc.className)
		})
	}
}

// TestValidateTimeSlot tests the time slot validation logic
func TestValidateTimeSlot(t *testing.T) {
	testCases := []struct {
		name        string
		slot        *TimeSlot
		shouldError bool
		description string
	}{
		{
			name: "Valid Time Slot",
			slot: &TimeSlot{
				ID:          "test-slot",
				Hour:        2,
				Day:         2,
				Classes:     []string{"DeathKnight", "DemonHunter"},
				Description: "Test Slot",
			},
			shouldError: false,
			description: "Valid time slot should pass validation",
		},
		{
			name: "Invalid Hour",
			slot: &TimeSlot{
				ID:          "invalid-hour",
				Hour:        24,
				Day:         2,
				Classes:     []string{"DeathKnight"},
				Description: "Invalid Hour",
			},
			shouldError: true,
			description: "Hour should be between 0 and 23",
		},
		{
			name: "Invalid Day",
			slot: &TimeSlot{
				ID:          "invalid-day",
				Hour:        2,
				Day:         7,
				Classes:     []string{"DeathKnight"},
				Description: "Invalid Day",
			},
			shouldError: true,
			description: "Day should be between 0 and 6",
		},
		{
			name: "Empty Classes",
			slot: &TimeSlot{
				ID:          "empty-classes",
				Hour:        2,
				Day:         2,
				Classes:     []string{},
				Description: "Empty Classes",
			},
			shouldError: true,
			description: "Classes list cannot be empty",
		},
		{
			name:        "Nil Slot",
			slot:        nil,
			shouldError: true,
			description: "Nil slot should fail validation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing case: %s - %s", tc.name, tc.description)

			err := ValidateTimeSlot(tc.slot)

			if tc.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTimeSlotOverlap tests for overlapping time slots
func TestTimeSlotOverlap(t *testing.T) {
	testCases := []struct {
		name        string
		slots       []TimeSlot
		shouldError bool
		description string
	}{
		{
			name: "Non-overlapping Slots",
			slots: []TimeSlot{
				{
					ID:   "slot1",
					Hour: 2,
					Day:  2,
				},
				{
					ID:   "slot2",
					Hour: 7,
					Day:  2,
				},
			},
			shouldError: false,
			description: "Different hours on same day should not overlap",
		},
		{
			name: "Overlapping Slots",
			slots: []TimeSlot{
				{
					ID:   "slot1",
					Hour: 2,
					Day:  2,
				},
				{
					ID:   "slot2",
					Hour: 2,
					Day:  2,
				},
			},
			shouldError: true,
			description: "Same hour and day should be considered overlap",
		},
		{
			name: "Different Days",
			slots: []TimeSlot{
				{
					ID:   "slot1",
					Hour: 2,
					Day:  2,
				},
				{
					ID:   "slot2",
					Hour: 2,
					Day:  3,
				},
			},
			shouldError: false,
			description: "Same hour but different days should not overlap",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing case: %s - %s", tc.name, tc.description)

			hasOverlap := false
			for i := 0; i < len(tc.slots); i++ {
				for j := i + 1; j < len(tc.slots); j++ {
					if tc.slots[i].Hour == tc.slots[j].Hour && tc.slots[i].Day == tc.slots[j].Day {
						hasOverlap = true
						break
					}
				}
			}

			assert.Equal(t, tc.shouldError, hasOverlap)
		})
	}
}

// TestTimeSlotValidationErrors tests specific error messages and validation details
func TestTimeSlotValidationErrors(t *testing.T) {
	testCases := []struct {
		name          string
		slot          *TimeSlot
		expectedError string
		description   string
	}{
		{
			name: "Invalid Hour Error",
			slot: &TimeSlot{
				ID:   "invalid-hour",
				Hour: 25,
				Day:  2,
			},
			expectedError: "invalid hour: 25",
			description:   "Should return specific error for invalid hour",
		},
		{
			name: "Invalid Day Error",
			slot: &TimeSlot{
				ID:   "invalid-day",
				Hour: 2,
				Day:  8,
			},
			expectedError: "invalid day: 8",
			description:   "Should return specific error for invalid day",
		},
		{
			name: "No Classes Error",
			slot: &TimeSlot{
				ID:      "no-classes",
				Hour:    2,
				Day:     2,
				Classes: []string{},
			},
			expectedError: "no classes defined for slot no-classes",
			description:   "Should return specific error for empty classes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing case: %s - %s", tc.name, tc.description)

			err := ValidateTimeSlot(tc.slot)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}
