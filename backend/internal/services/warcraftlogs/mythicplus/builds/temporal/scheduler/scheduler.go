package warcraftlogsBuildsTemporalScheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
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
	client       client.Client
	prodSchedule client.ScheduleHandle
	testSchedule client.ScheduleHandle
	logger       *log.Logger
}

// NewScheduleManager creates a new ScheduleManager instance
func NewScheduleManager(temporalClient client.Client, logger *log.Logger) *ScheduleManager {
	return &ScheduleManager{
		client: temporalClient,
		logger: logger,
	}
}

// CreateSchedule creates the main weekly synchronization schedule
func (sm *ScheduleManager) CreateSchedule(ctx context.Context, cfg *models.WorkflowConfig, opts *ScheduleOptions) error {
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
			Workflow:  "SyncWorkflow",
			TaskQueue: DefaultScheduleConfig.TaskQueue,
			Args:      []interface{}{*cfg},
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

	sm.prodSchedule = handle
	sm.logger.Printf("[INFO] Created weekly sync schedule %s for Tuesday 2 AM UTC", scheduleID)
	return nil
}

// CreateTestSchedule creates a separate test schedule for Priest class only
func (sm *ScheduleManager) CreateTestSchedule(ctx context.Context, cfg *models.WorkflowConfig, opts *ScheduleOptions) error {
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
			Workflow:  "SyncWorkflow",
			TaskQueue: DefaultScheduleConfig.TaskQueue,
			Args:      []interface{}{*cfg},
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

	sm.testSchedule = handle
	sm.logger.Printf("[TEST] Created Priest test schedule %s", scheduleID)
	return nil
}

// TriggerSyncNow triggers the immediate execution of the production schedule
func (sm *ScheduleManager) TriggerSyncNow(ctx context.Context) error {
	if sm.prodSchedule == nil {
		return fmt.Errorf("no production schedule has been created")
	}
	return sm.prodSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// TriggerTestNow triggers the immediate execution of the test schedule
func (sm *ScheduleManager) TriggerTestNow(ctx context.Context) error {
	if sm.testSchedule == nil {
		return fmt.Errorf("no test schedule has been created")
	}
	return sm.testSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// PauseSchedule pauses the production schedule
func (sm *ScheduleManager) PauseSchedule(ctx context.Context) error {
	if sm.prodSchedule == nil {
		return fmt.Errorf("no production schedule has been created")
	}
	return sm.prodSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// PauseTestSchedule pauses the test schedule
func (sm *ScheduleManager) PauseTestSchedule(ctx context.Context) error {
	if sm.testSchedule == nil {
		return fmt.Errorf("no test schedule has been created")
	}
	return sm.testSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// UnpauseSchedule reactivates the production schedule
func (sm *ScheduleManager) UnpauseSchedule(ctx context.Context) error {
	if sm.prodSchedule == nil {
		return fmt.Errorf("no production schedule has been created")
	}
	return sm.prodSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// UnpauseTestSchedule reactivates the test schedule
func (sm *ScheduleManager) UnpauseTestSchedule(ctx context.Context) error {
	if sm.testSchedule == nil {
		return fmt.Errorf("no test schedule has been created")
	}
	return sm.testSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// DeleteSchedule deletes a schedule by its ID
func (sm *ScheduleManager) DeleteSchedule(ctx context.Context, scheduleID string) error {
	handle := sm.client.ScheduleClient().GetHandle(ctx, scheduleID)
	return handle.Delete(ctx)
}

// CleanupExistingSchedules deletes existing schedules
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

	// Reset references
	sm.prodSchedule = nil
	sm.testSchedule = nil

	return nil
}

// CleanupOldWorkflows terminates running workflows
func (sm *ScheduleManager) CleanupOldWorkflows(ctx context.Context) error {
	// Get all workflows with SyncWorkflow type
	resp, err := sm.client.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		Namespace: models.DefaultNamespace,
		Query:     "WorkflowType='SyncWorkflow' OR WorkflowType='ProcessBuildBatchWorkflow'",
	})

	if err != nil {
		return fmt.Errorf("failed to list workflows: %w", err)
	}

	for _, execution := range resp.Executions {
		workflowID := execution.Execution.WorkflowId
		runID := execution.Execution.RunId

		// Only terminate if it's running
		if execution.Status != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			sm.logger.Printf("[INFO] Skipping non-running workflow: %s (status: %s)",
				workflowID, execution.Status.String())
			continue
		}

		err := sm.client.TerminateWorkflow(ctx, workflowID, runID, "Cleanup of old workflows")
		if err != nil {
			sm.logger.Printf("[WARN] Failed to terminate workflow %s: %v", workflowID, err)
			continue
		}
		sm.logger.Printf("[INFO] Terminated workflow: %s", workflowID)
	}

	return nil
}

// CleanupAll do a complete cleanup
func (sm *ScheduleManager) CleanupAll(ctx context.Context) error {
	if err := sm.CleanupExistingSchedules(ctx); err != nil {
		return fmt.Errorf("failed to cleanup schedules: %w", err)
	}

	if err := sm.CleanupOldWorkflows(ctx); err != nil {
		return fmt.Errorf("failed to cleanup workflows: %w", err)
	}

	return nil
}
