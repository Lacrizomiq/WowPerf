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
	prodScheduleID    = "warcraft-logs-weekly-sync"
	testScheduleID    = "warcraft-logs-priest-test"
	analyzeScheduleID = "warcraft-logs-analyze"
)

// ScheduleManager handles Temporal schedule for WarcraftLogs sync
type ScheduleManager struct {
	client          client.Client
	prodSchedule    client.ScheduleHandle
	testSchedule    client.ScheduleHandle
	analyzeSchedule client.ScheduleHandle
	logger          *log.Logger
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

// CreateAnalyzeSchedule creates the schedule for the analyze workflow
func (sm *ScheduleManager) CreateAnalyzeSchedule(ctx context.Context, cfg *models.WorkflowConfig, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := analyzeScheduleID
	workflowID := fmt.Sprintf("warcraft-logs-analyze-%s", time.Now().UTC().Format("2006-01-02"))

	analyzeDay := 3 // Wednesday

	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{fmt.Sprintf("0 %d * * %d", DefaultScheduleConfig.Hour, analyzeDay)},
			Jitter:          time.Minute,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  "AnalyzeBuildsWorkflow",
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
		return fmt.Errorf("failed to create analyze schedule: %w", err)
	}

	sm.analyzeSchedule = handle
	sm.logger.Printf("[INFO] Created weekly analysis schedule %s for Wednesday %d AM UTC", scheduleID, DefaultScheduleConfig.Hour)
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

// TriggerAnalyzeNow triggers the immediate execution of the analyze schedule
func (sm *ScheduleManager) TriggerAnalyzeNow(ctx context.Context) error {
	if sm.analyzeSchedule == nil {
		return fmt.Errorf("no analyze schedule has been created")
	}
	return sm.analyzeSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
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

// PauseAnalyzeSchedule pauses the analyze schedule
func (sm *ScheduleManager) PauseAnalyzeSchedule(ctx context.Context) error {
	if sm.analyzeSchedule == nil {
		return fmt.Errorf("no analyze schedule has been created")
	}
	return sm.analyzeSchedule.Pause(ctx, client.SchedulePauseOptions{})
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

// UnpauseAnalyzeSchedule reactivates the analyze schedule
func (sm *ScheduleManager) UnpauseAnalyzeSchedule(ctx context.Context) error {
	if sm.analyzeSchedule == nil {
		return fmt.Errorf("no analyze schedule has been created")
	}
	return sm.analyzeSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// DeleteSchedule deletes a schedule by its ID
func (sm *ScheduleManager) DeleteSchedule(ctx context.Context, scheduleID string) error {
	handle := sm.client.ScheduleClient().GetHandle(ctx, scheduleID)
	return handle.Delete(ctx)
}

// CleanupExistingSchedules deletes existing schedules
func (sm *ScheduleManager) CleanupExistingSchedules(ctx context.Context) error {
	// List and delete existing schedules
	schedules := []string{prodScheduleID, testScheduleID, analyzeScheduleID}
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
	sm.analyzeSchedule = nil

	return nil
}

// CleanupOldWorkflows terminates running workflows
func (sm *ScheduleManager) CleanupOldWorkflows(ctx context.Context) error {
	var terminatedCount int

	// Define the workflow types to clean
	workflowTypes := []string{
		"SyncWorkflow",
		"ProcessBuildBatchWorkflow",
		"AnalyzeBuildsWorkflow",
	}

	// Process each workflow type separately
	for _, workflowType := range workflowTypes {
		// Build a query for a single workflow type
		query := fmt.Sprintf("WorkflowType='%s'", workflowType)

		sm.logger.Printf("[INFO] Listing workflows of type: %s", workflowType)

		// Retrieve the workflows of this type
		resp, err := sm.client.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
			Namespace: models.DefaultNamespace,
			Query:     query,
		})

		if err != nil {
			sm.logger.Printf("[WARN] Failed to list workflows of type %s: %v", workflowType, err)
			continue
		}

		// Process each retrieved workflow
		for _, execution := range resp.Executions {
			workflowID := execution.Execution.WorkflowId
			runID := execution.Execution.RunId

			// Only terminate running workflows
			if execution.Status != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
				sm.logger.Printf("[INFO] Skipping non-running workflow: %s (type: %s, status: %s)",
					workflowID, workflowType, execution.Status.String())
				continue
			}

			// Terminate the workflow
			err := sm.client.TerminateWorkflow(ctx, workflowID, runID, "Cleanup of old workflows")
			if err != nil {
				sm.logger.Printf("[WARN] Failed to terminate workflow %s (type: %s): %v",
					workflowID, workflowType, err)
				continue
			}

			sm.logger.Printf("[INFO] Terminated workflow: %s (type: %s)", workflowID, workflowType)
			terminatedCount++
		}
	}

	sm.logger.Printf("[INFO] Cleanup completed - terminated %d workflows", terminatedCount)
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
