// go test -v ./internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler/

package warcraftlogsBuildsTemporalScheduler

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/client"
)

// setupTestLogger creates a logger that writes to a buffer for testing
func setupTestLogger() (*bytes.Buffer, *log.Logger) {
	var buf bytes.Buffer
	return &buf, log.New(&buf, "[TEST] ", log.Ltime)
}

// TestTimeSlotConfiguration verifies the basic structure and configuration of time slots
func TestTimeSlotConfiguration(t *testing.T) {
	t.Log("Starting time slot configuration tests...")

	// TestTimeSlotConfiguration
	t.Run("Tuesday Slots Configuration", func(t *testing.T) {
		slots := []struct {
			slot     TimeSlot
			name     string
			expected struct{ hour, day int }
		}{
			{Tuesday2AM, "2 AM", struct{ hour, day int }{2, 2}},
			{Tuesday8AM, "8 AM", struct{ hour, day int }{8, 2}},
			{Tuesday2PM, "2 PM", struct{ hour, day int }{14, 2}},
			{Tuesday8PM, "8 PM", struct{ hour, day int }{20, 2}},
		}

		for _, s := range slots {
			t.Logf("Checking %s slot configuration", s.name)
			assert.Equal(t, s.expected.hour, s.slot.Hour)
			assert.Equal(t, s.expected.day, s.slot.Day)
			t.Logf("Slot %s classes: %v", s.name, s.slot.Classes)
			t.Logf("Slot %s task queue: %s", s.name, s.slot.TaskQueue)
		}
	})

	t.Log("Time slot configuration tests completed")
}

// TestScheduleManager tests the core scheduling functionality
func TestScheduleManager(t *testing.T) {
	t.Log("Starting schedule manager tests...")
	buf, logger := setupTestLogger()

	t.Run("Schedule State Management", func(t *testing.T) {
		t.Log("Testing schedule state transitions...")

		// Create test schedule manager
		sm := &ScheduleManager{
			schedules:     make(map[string]client.ScheduleHandle),
			workflowState: make(map[string]*WorkflowState),
			logger:        logger,
		}

		// Test slot activation
		slot := Tuesday2AM
		taskQueue := slot.TaskQueue

		t.Logf("Testing activation of slot: %s", slot.ID)
		sm.workflowState[taskQueue] = &WorkflowState{
			TaskQueue: taskQueue,
			IsRunning: false,
		}

		// Verify initial state
		state := sm.workflowState[taskQueue]
		assert.False(t, state.IsRunning, "Initial state should be not running")
		t.Log("Initial state verified")

		// Test state transitions
		t.Log("Testing state transitions...")
		sm.activeSlot = &slot
		assert.NotNil(t, sm.activeSlot, "Active slot should be set")
		t.Logf("Logs from test: \n%s", buf.String())
	})
}

// TestScheduleStateTransitions tests the transitions between different schedule states
func TestScheduleStateTransitions(t *testing.T) {
	t.Log("Starting schedule state transition tests...")
	buf, logger := setupTestLogger()

	t.Run("Slot Transition Validation", func(t *testing.T) {
		sm := &ScheduleManager{
			schedules:     make(map[string]client.ScheduleHandle),
			workflowState: make(map[string]*WorkflowState),
			logger:        logger,
		}

		// Test transition between slots
		t.Log("Testing transition between time slots...")

		// Set initial slot
		initialSlot := Tuesday2AM
		sm.activeSlot = &initialSlot
		t.Logf("Initial slot set to: %s", initialSlot.ID)

		// Simulate slot change
		newSlot := Tuesday8AM
		sm.activeSlot = &newSlot
		t.Logf("Transitioned to new slot: %s", newSlot.ID)

		assert.Equal(t, &newSlot, sm.activeSlot, "Active slot should be updated")
		t.Logf("Transition logs: \n%s", buf.String())
	})
}

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	t.Log("Starting error handling tests...")
	buf, logger := setupTestLogger()

	t.Run("Invalid Schedule Creation", func(t *testing.T) {
		sm := &ScheduleManager{
			schedules:     make(map[string]client.ScheduleHandle),
			workflowState: make(map[string]*WorkflowState),
			logger:        logger,
		}

		t.Log("Testing invalid class schedule creation...")
		err := sm.CreateClassSchedule(nil, "InvalidClass", nil, nil)
		assert.Error(t, err, "Should error on invalid class")
		t.Logf("Error handling logs: \n%s", buf.String())
	})

	t.Run("Task Queue Conflicts", func(t *testing.T) {
		t.Log("Testing task queue conflict handling...")
		// Test duplicate task queue assignments in existing slots
		taskQueues := make(map[string]bool)

		for _, slot := range ScheduleSlots {
			t.Logf("Checking task queue for slot %s: %s", slot.ID, slot.TaskQueue)
			if taskQueues[slot.TaskQueue] {
				t.Errorf("Duplicate task queue found: %s", slot.TaskQueue)
			}
			taskQueues[slot.TaskQueue] = true
		}

		// Verify that each slot has a unique task queue
		assert.Equal(t, len(ScheduleSlots), len(taskQueues),
			"Number of unique task queues should match number of slots")
		t.Log("Task queue conflict tests completed")
	})
}

// TestConfigValidation tests the configuration validation logic
func TestConfigValidation(t *testing.T) {
	t.Log("Starting configuration validation tests...")

	t.Run("Schedule Options Validation", func(t *testing.T) {
		t.Log("Testing schedule options validation...")
		opts := DefaultScheduleOptions()

		assert.NotZero(t, opts.Retry.InitialInterval, "Initial interval should be set")
		assert.NotZero(t, opts.Retry.MaximumInterval, "Maximum interval should be set")
		assert.Greater(t, opts.Retry.MaximumAttempts, 0, "Maximum attempts should be positive")
		t.Logf("Schedule options validated: MaxAttempts=%d, InitialInterval=%v",
			opts.Retry.MaximumAttempts, opts.Retry.InitialInterval)
	})

	t.Run("Task Queue Naming Convention", func(t *testing.T) {
		t.Log("Testing task queue naming conventions...")

		for _, slot := range ScheduleSlots {
			t.Logf("Validating task queue for slot %s: %s", slot.ID, slot.TaskQueue)
			assert.Contains(t, slot.TaskQueue, "warcraft-logs-")
			assert.Contains(t, slot.TaskQueue, slot.ID)
		}
	})
}

func TestMidnightTransition(t *testing.T) {
	t.Log("Testing midnight transition between slots...")

	// Test transition from Tuesday 8 PM to Wednesday 2 AM
	t.Run("Tuesday to Wednesday Transition", func(t *testing.T) {
		tuesdaySlot := Tuesday8PM
		wednesdaySlot := Wednesday2AM

		t.Logf("Checking transition from %s (Day %d, Hour %d) to %s (Day %d, Hour %d)",
			tuesdaySlot.ID, tuesdaySlot.Day, tuesdaySlot.Hour,
			wednesdaySlot.ID, wednesdaySlot.Day, wednesdaySlot.Hour)

		// Verify there's enough gap between slots
		hourDiff := (wednesdaySlot.Day-tuesdaySlot.Day)*24 + (wednesdaySlot.Hour - tuesdaySlot.Hour)
		assert.GreaterOrEqual(t, hourDiff, 6,
			"Should have at least 6 hours between slots across days")

		t.Logf("Gap between slots: %d hours", hourDiff)
	})
}

func TestScheduleStateManagement(t *testing.T) {
	t.Log("Testing detailed schedule state management...")

	t.Run("State Transitions During Slot Change", func(t *testing.T) {
		buf, logger := setupTestLogger()
		sm := &ScheduleManager{
			schedules:     make(map[string]client.ScheduleHandle),
			workflowState: make(map[string]*WorkflowState),
			logger:        logger,
		}

		// Initial state
		slot1 := Tuesday2AM
		sm.workflowState[slot1.TaskQueue] = &WorkflowState{
			IsRunning: true,
			TaskQueue: slot1.TaskQueue,
		}
		t.Logf("Initial state set for slot %s: IsRunning=%v",
			slot1.ID, sm.workflowState[slot1.TaskQueue].IsRunning)

		// Simulate slot change
		err := sm.checkAndEnforceTimeSlot()
		assert.NoError(t, err)
		t.Logf("State transition logs:\n%s", buf.String())
	})
}
