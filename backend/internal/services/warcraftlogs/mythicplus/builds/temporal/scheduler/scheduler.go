package warcraftlogsBuildsTemporalScheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

// Constants for schedule IDs
const (
	prodScheduleID = "warcraft-logs-weekly-sync"
	testScheduleID = "warcraft-logs-priest-test"
)

// ScheduleManager handles Temporal schedule for WarcraftLogs sync
type ScheduleManager struct {
	client   client.Client
	schedule client.ScheduleHandle
	logger   *log.Logger
}

// NewScheduleManager creates a new ScheduleManager instance
func NewScheduleManager(temporalClient client.Client, logger *log.Logger) *ScheduleManager {
	return &ScheduleManager{
		client: temporalClient,
		logger: logger,
	}
}

// CreateSchedule creates the main weekly synchronization schedule
func (sm *ScheduleManager) CreateSchedule(ctx context.Context, cfg *workflows.Config, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := prodScheduleID
	workflowID := fmt.Sprintf("warcraft-logs-workflow-%s", time.Now().UTC().Format("2006-01-02"))

	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{fmt.Sprintf("0 %d * * %d", DefaultScheduleConfig.Hour, DefaultScheduleConfig.Day)},
			Jitter:          time.Minute,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  workflows.SyncWorkflow,
			TaskQueue: DefaultScheduleConfig.TaskQueue,
			Args:      []interface{}{workflows.WorkflowParams{Config: cfg, WorkflowID: workflowID}},
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    opts.Retry.InitialInterval,
				BackoffCoefficient: opts.Retry.BackoffCoefficient,
				MaximumInterval:    opts.Retry.MaximumInterval,
				MaximumAttempts:    int32(opts.Retry.MaximumAttempts),
			},
			WorkflowRunTimeout: opts.Timeout,
		},
	}

	handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
	if err != nil {
		return fmt.Errorf("failed to create production schedule: %w", err)
	}

	sm.schedule = handle
	sm.logger.Printf("[INFO] Created weekly sync schedule %s for Tuesday 2 AM UTC", scheduleID)
	return nil
}

// CreateTestSchedule creates a separate test schedule for Priest class only
func (sm *ScheduleManager) CreateTestSchedule(ctx context.Context, cfg *workflows.Config, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := testScheduleID
	workflowID := fmt.Sprintf("warcraft-logs-priest-test-%s", time.Now().UTC().Format("2006-01-02"))

	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{fmt.Sprintf("0 %d * * %d", DefaultScheduleConfig.Hour, DefaultScheduleConfig.Day)},
			Jitter:          time.Minute,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  workflows.SyncWorkflow,
			TaskQueue: DefaultScheduleConfig.TaskQueue,
			Args:      []interface{}{workflows.WorkflowParams{Config: cfg, WorkflowID: workflowID}},
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    opts.Retry.InitialInterval,
				BackoffCoefficient: opts.Retry.BackoffCoefficient,
				MaximumInterval:    opts.Retry.MaximumInterval,
				MaximumAttempts:    int32(opts.Retry.MaximumAttempts),
			},
			WorkflowRunTimeout: opts.Timeout,
		},
	}

	handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
	if err != nil {
		return fmt.Errorf("failed to create test schedule: %w", err)
	}

	sm.schedule = handle
	sm.logger.Printf("[TEST] Created Priest test schedule %s for immediate execution", scheduleID)
	return nil
}

// TriggerSyncNow triggers an immediate execution of the current schedule
func (sm *ScheduleManager) TriggerSyncNow(ctx context.Context) error {
	if sm.schedule == nil {
		return fmt.Errorf("no schedule has been created")
	}
	return sm.schedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// PauseSchedule pauses the current schedule
func (sm *ScheduleManager) PauseSchedule(ctx context.Context) error {
	if sm.schedule == nil {
		return fmt.Errorf("no schedule has been created")
	}
	return sm.schedule.Pause(ctx, client.SchedulePauseOptions{})
}

// UnpauseSchedule unpauses the current schedule
func (sm *ScheduleManager) UnpauseSchedule(ctx context.Context) error {
	if sm.schedule == nil {
		return fmt.Errorf("no schedule has been created")
	}
	return sm.schedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

func (sm *ScheduleManager) DeleteSchedule(ctx context.Context, scheduleID string) error {
	handle := sm.client.ScheduleClient().GetHandle(ctx, scheduleID)
	return handle.Delete(ctx)
}

// Add cleanup functionality
func (sm *ScheduleManager) CleanupExistingSchedules(ctx context.Context) error {
	// List and delete existing schedules
	schedules := []string{prodScheduleID, testScheduleID}
	for _, id := range schedules {
		if err := sm.DeleteSchedule(ctx, id); err != nil {
			sm.logger.Printf("[WARN] Failed to delete schedule %s: %v", id, err)
			continue
		}
		sm.logger.Printf("[INFO] Deleted schedule: %s", id)
	}
	return nil
}
