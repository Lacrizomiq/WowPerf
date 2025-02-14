package warcraftlogsBuildsTemporalScheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"

	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
)

// ScheduleManager handles Temporal schedules for different classes with time slot enforcement
type ScheduleManager struct {
	client     client.Client
	schedules  map[string]client.ScheduleHandle
	logger     *log.Logger
	mu         sync.RWMutex
	activeSlot *TimeSlot
	monitoring bool
}

// NewScheduleManager creates a new ScheduleManager instance
func NewScheduleManager(temporalClient client.Client, logger *log.Logger) *ScheduleManager {
	sm := &ScheduleManager{
		client:     temporalClient,
		schedules:  make(map[string]client.ScheduleHandle),
		logger:     logger,
		monitoring: false,
	}

	// Start monitoring in a separate goroutine
	sm.startTimeSlotMonitoring(context.Background())

	return sm
}

// CreateClassSchedule creates a schedule for a specific class
func (sm *ScheduleManager) CreateClassSchedule(ctx context.Context, className string, cfg *workflows.Config, opts *ScheduleOptions) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.logger.Printf("[INFO] Creating schedule for class: %s", className)

	// Verify class is scheduled in a valid time slot
	slot := FindTimeSlotForClass(className)
	if slot == nil {
		return fmt.Errorf("class %s is not configured for any time slot", className)
	}

	scheduleID := fmt.Sprintf("warcraft-logs-%s-schedule", className)
	workflowID := fmt.Sprintf("warcraft-logs-%s-workflow-%s", className, time.Now().UTC().Format("2006-01-02T15:04:05Z"))

	// Filter config for just this class
	classConfig := filterConfigForClass(cfg, className)
	workflowParams := workflows.WorkflowParams{
		Config:     classConfig,
		WorkflowID: workflowID,
	}

	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{fmt.Sprintf("0 %d * * %d", slot.Hour, slot.Day)},
			Jitter:          time.Minute,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:                 workflowID,
			Workflow:           workflows.SyncWorkflow,
			TaskQueue:          "warcraft-logs-sync",
			Args:               []interface{}{workflowParams},
			WorkflowRunTimeout: opts.Timeout,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    opts.Retry.InitialInterval,
				BackoffCoefficient: opts.Retry.BackoffCoefficient,
				MaximumInterval:    opts.Retry.MaximumInterval,
				MaximumAttempts:    int32(opts.Retry.MaximumAttempts),
			},
		},
		Paused: true, // Start paused, let the monitor handle activation
	}

	handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
	if err != nil {
		return fmt.Errorf("failed to create schedule for %s: %w", className, err)
	}

	sm.schedules[className] = handle
	sm.logger.Printf("[INFO] Successfully created schedule for class %s in slot %s", className, slot.Description)

	// Initial state check
	currentSlot := GetCurrentTimeSlot()
	if currentSlot != nil && currentSlot.ID == slot.ID {
		if err := handle.Unpause(ctx, client.ScheduleUnpauseOptions{}); err != nil {
			sm.logger.Printf("[WARN] Failed to initially unpause schedule for %s: %v", className, err)
		}
	}

	return nil
}

// checkAndEnforceTimeSlot ensures only classes in the current time slot are active
func (sm *ScheduleManager) checkAndEnforceTimeSlot() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	currentSlot := GetCurrentTimeSlot()

	// Log the time slot change
	sm.logger.Printf("[DEBUG] Time slot check - Current: %+v, Active: %+v",
		currentSlot, sm.activeSlot)

	// If we're outside any slot, pause everything
	if currentSlot == nil {
		sm.logger.Printf("[INFO] Outside of scheduled time slots, pausing all schedules")
		return sm.pauseAllSchedules(context.Background())
	}

	// If we're in a new slot
	if sm.activeSlot == nil || sm.activeSlot.ID != currentSlot.ID {
		sm.logger.Printf("[INFO] Switching to new time slot: %s", currentSlot.Description)

		// Pause all schedules first
		if err := sm.pauseAllSchedules(context.Background()); err != nil {
			return fmt.Errorf("failed to pause schedules during slot switch: %w", err)
		}

		// Then activate only the schedules for the current slot
		for _, className := range currentSlot.Classes {
			if handle, exists := sm.schedules[className]; exists {
				if err := handle.Unpause(context.Background(), client.ScheduleUnpauseOptions{}); err != nil {
					sm.logger.Printf("[ERROR] Failed to unpause schedule for %s: %v", className, err)
				} else {
					sm.logger.Printf("[INFO] Activated schedule for class %s", className)
				}
			}
		}

		sm.activeSlot = currentSlot
	}

	return nil
}

// pauseAllSchedules pauses all active schedules
func (sm *ScheduleManager) pauseAllSchedules(ctx context.Context) error {
	for className, handle := range sm.schedules {
		if err := handle.Pause(ctx, client.SchedulePauseOptions{}); err != nil {
			return fmt.Errorf("failed to pause schedule for %s: %w", className, err)
		}
		sm.logger.Printf("[INFO] Paused schedule for class %s", className)
	}
	return nil
}

// startTimeSlotMonitoring starts the time slot monitoring goroutine
func (sm *ScheduleManager) startTimeSlotMonitoring(ctx context.Context) {
	if sm.monitoring {
		return
	}

	sm.monitoring = true
	ticker := time.NewTicker(time.Minute)

	go func() {
		sm.logger.Printf("[INFO] Starting time slot monitoring")
		defer func() {
			ticker.Stop()
			sm.monitoring = false
			sm.logger.Printf("[INFO] Time slot monitoring stopped")
		}()

		// Initial check
		if err := sm.checkAndEnforceTimeSlot(); err != nil {
			sm.logger.Printf("[ERROR] Initial time slot check failed: %v", err)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := sm.checkAndEnforceTimeSlot(); err != nil {
					sm.logger.Printf("[ERROR] Time slot check failed: %v", err)
				}
			}
		}
	}()
}

// GetScheduleDescription gets the description of a schedule
func (sm *ScheduleManager) GetScheduleDescription(ctx context.Context, className string) (*client.ScheduleDescription, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	handle, exists := sm.schedules[className]
	if !exists {
		return nil, fmt.Errorf("no schedule found for class %s", className)
	}

	return handle.Describe(ctx)
}

// DeleteSchedule deletes a schedule
func (sm *ScheduleManager) DeleteSchedule(ctx context.Context, className string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	handle, exists := sm.schedules[className]
	if !exists {
		return fmt.Errorf("no schedule found for class %s", className)
	}

	if err := handle.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete schedule for %s: %w", className, err)
	}

	delete(sm.schedules, className)
	sm.logger.Printf("[INFO] Deleted schedule for class %s", className)
	return nil
}

// filterConfigForClass returns a config filtered for a specific class
func filterConfigForClass(cfg *workflows.Config, className string) *workflows.Config {
	filteredConfig := *cfg
	filteredConfig.Specs = []workflows.ClassSpec{}

	for _, spec := range cfg.Specs {
		if spec.ClassName == className {
			filteredConfig.Specs = append(filteredConfig.Specs, spec)
		}
	}

	return &filteredConfig
}

// GetCurrentState returns the current state of the scheduler, usefull for monitoring and health
func (sm *ScheduleManager) GetCurrentState() (map[string]bool, *TimeSlot) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	activeSchedules := make(map[string]bool)
	for className := range sm.schedules {
		activeSchedules[className] = true
	}

	return activeSchedules, sm.activeSlot
}

// CreateOrGetClassSchedule creates a new schedule if it doesn't exist, or returns existing one
func (sm *ScheduleManager) CreateOrGetClassSchedule(ctx context.Context, className string, cfg *workflows.Config, opts *ScheduleOptions) error {
	sm.logger.Printf("[INFO] CreateOrGetClassSchedule called for class: %s", className)

	scheduleID := fmt.Sprintf("warcraft-logs-%s-schedule", className)

	// Check if schedule already exists
	scheduleHandle := sm.client.ScheduleClient().GetHandle(ctx, scheduleID)
	_, err := scheduleHandle.Describe(ctx)

	if err == nil {
		sm.logger.Printf("[INFO] Schedule already exists for class %s", className)
		sm.schedules[className] = scheduleHandle
		return nil
	}

	// If schedule doesn't exist, create it
	return sm.CreateClassSchedule(ctx, className, cfg, opts)
}

// TriggerSchedule triggers an immediate execution of the schedule
func (sm *ScheduleManager) TriggerSchedule(ctx context.Context, className string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	handle, exists := sm.schedules[className]
	if !exists {
		return fmt.Errorf("no schedule found for class %s", className)
	}

	// Verify the class is in the current time slot
	slot := FindTimeSlotForClass(className)
	if slot == nil {
		return fmt.Errorf("class %s not configured for any time slot", className)
	}

	currentSlot := GetCurrentTimeSlot()
	if currentSlot == nil || currentSlot.ID != slot.ID {
		return fmt.Errorf("cannot trigger schedule for %s outside its time slot", className)
	}

	err := handle.Trigger(ctx, client.ScheduleTriggerOptions{})
	if err != nil {
		return fmt.Errorf("failed to trigger schedule for %s: %w", className, err)
	}

	sm.logger.Printf("[INFO] Triggered immediate execution for class %s", className)
	return nil
}
