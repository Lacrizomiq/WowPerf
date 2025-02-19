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

// CreateSchedule creates the weekly synchronization schedule
func (sm *ScheduleManager) CreateSchedule(ctx context.Context, cfg *workflows.Config, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := "warcraft-logs-weekly-sync"
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
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	sm.schedule = handle
	sm.logger.Printf("[INFO] Created weekly sync schedule for Tuesday 2 AM UTC")
	return nil
}

// TriggerSyncNow triggers an immediate execution
func (sm *ScheduleManager) TriggerSyncNow(ctx context.Context) error {
	if sm.schedule == nil {
		return fmt.Errorf("no schedule has been created")
	}
	return sm.schedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// PauseSchedule pauses the schedule
func (sm *ScheduleManager) PauseSchedule(ctx context.Context) error {
	if sm.schedule == nil {
		return fmt.Errorf("no schedule has been created")
	}
	return sm.schedule.Pause(ctx, client.SchedulePauseOptions{})
}

// UnpauseSchedule unpauses the schedule
func (sm *ScheduleManager) UnpauseSchedule(ctx context.Context) error {
	if sm.schedule == nil {
		return fmt.Errorf("no schedule has been created")
	}
	return sm.schedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}
