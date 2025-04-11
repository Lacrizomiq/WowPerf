package warcraftlogsBuildsTemporalScheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

// Constants for schedule IDs
const (
	// Existing schedules (keep them for now but will be removed soon)
	prodScheduleID    = "warcraft-logs-weekly-sync"
	testScheduleID    = "warcraft-logs-priest-test"
	analyzeScheduleID = "warcraft-logs-analyze"

	// New schedules for decoupled workflows
	rankingsScheduleID = "warcraft-logs-rankings"
	reportsScheduleID  = "warcraft-logs-reports"
	buildsScheduleID   = "warcraft-logs-builds"
)

// ScheduleManager manages Temporal schedules for WarcraftLogs workflows
type ScheduleManager struct {
	client client.Client
	logger *log.Logger

	// Existing handles (keep them for now but will be removed soon)
	prodSchedule    client.ScheduleHandle
	testSchedule    client.ScheduleHandle
	analyzeSchedule client.ScheduleHandle

	// New handles (decoupled workflows, will be used instead of the existing ones)
	rankingsSchedule client.ScheduleHandle
	reportsSchedule  client.ScheduleHandle
	buildsSchedule   client.ScheduleHandle
}

// NewScheduleManager creates a new ScheduleManager instance
func NewScheduleManager(temporalClient client.Client, logger *log.Logger) *ScheduleManager {
	return &ScheduleManager{
		client: temporalClient,
		logger: logger,
	}
}

// CreateRankingsSchedule creates the rankings schedule for the RankingsWorkflow
func (sm *ScheduleManager) CreateRankingsSchedule(ctx context.Context, params models.RankingsWorkflowParams, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := rankingsScheduleID

	// Generate a unique BatchID for this schedule if not provided
	if params.BatchID == "" {
		params.BatchID = fmt.Sprintf("rankings-%s", uuid.New().String())
	}

	workflowID := fmt.Sprintf("warcraft-logs-rankings-%s", time.Now().UTC().Format("2006-01-02"))

	// Create the schedule without automatic triggering (No CRON expression)
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		// No CronExpressions to avoid automatic triggering
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  "RankingsWorkflow",
			TaskQueue: DefaultScheduleConfig.TaskQueue,
			Args:      []interface{}{params},
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    opts.Retry.InitialInterval,
				BackoffCoefficient: opts.Retry.BackoffCoefficient,
				MaximumInterval:    opts.Retry.MaximumInterval,
				MaximumAttempts:    int32(opts.Retry.MaximumAttempts),
			},
			WorkflowRunTimeout: opts.Timeout,
		},
		Paused: opts.Paused, // Paused by default if specified in options
	}

	handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
	if err != nil {
		return fmt.Errorf("failed to create rankings schedule: %w", err)
	}

	sm.rankingsSchedule = handle
	sm.logger.Printf("[INFO] Created rankings workflow schedule: %s", scheduleID)
	return nil
}

// CreateReportsSchedule creates a schedule for the Reports workflow
func (sm *ScheduleManager) CreateReportsSchedule(ctx context.Context, params *models.ReportsWorkflowParams, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := reportsScheduleID

	// Generate a unique BatchID if not provided
	if params.BatchID == "" {
		params.BatchID = fmt.Sprintf("reports-%s", uuid.New().String())
	}

	workflowID := fmt.Sprintf("warcraft-logs-reports-%s", time.Now().UTC().Format("2006-01-02"))

	// Create schedule without automatic triggering (no CronExpressions)
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		// No CronExpressions to avoid automatic triggering
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  "ReportsWorkflow",
			TaskQueue: DefaultScheduleConfig.TaskQueue,
			Args:      []interface{}{params},
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    opts.Retry.InitialInterval,
				BackoffCoefficient: opts.Retry.BackoffCoefficient,
				MaximumInterval:    opts.Retry.MaximumInterval,
				MaximumAttempts:    int32(opts.Retry.MaximumAttempts),
			},
			WorkflowRunTimeout: opts.Timeout,
		},
		Paused: opts.Paused, // Paused by default if specified in options
	}

	handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
	if err != nil {
		return fmt.Errorf("failed to create reports schedule: %w", err)
	}

	sm.reportsSchedule = handle
	sm.logger.Printf("[INFO] Created reports workflow schedule: %s", scheduleID)
	return nil
}

// CreateBuildsSchedule creates a schedule for the Builds workflow
func (sm *ScheduleManager) CreateBuildsSchedule(ctx context.Context, params *models.BuildsWorkflowParams, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := buildsScheduleID

	// Generate a unique BatchID if not provided
	if params.BatchID == "" {
		params.BatchID = fmt.Sprintf("builds-%s", uuid.New().String())
	}

	workflowID := fmt.Sprintf("warcraft-logs-builds-%s", time.Now().UTC().Format("2006-01-02"))

	// Create schedule without automatic triggering (no CronExpressions)
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		// No CronExpressions to avoid automatic triggering
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  "BuildsWorkflow",
			TaskQueue: DefaultScheduleConfig.TaskQueue,
			Args:      []interface{}{params},
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    opts.Retry.InitialInterval,
				BackoffCoefficient: opts.Retry.BackoffCoefficient,
				MaximumInterval:    opts.Retry.MaximumInterval,
				MaximumAttempts:    int32(opts.Retry.MaximumAttempts),
			},
			WorkflowRunTimeout: opts.Timeout,
		},
		Paused: opts.Paused, // Paused by default if specified in options
	}

	handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
	if err != nil {
		return fmt.Errorf("failed to create builds schedule: %w", err)
	}

	sm.buildsSchedule = handle
	sm.logger.Printf("[INFO] Created builds workflow schedule: %s", scheduleID)
	return nil
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

// TriggerRankingsNow triggers the immediate execution of the rankings schedule
func (sm *ScheduleManager) TriggerRankingsNow(ctx context.Context) error {
	if sm.rankingsSchedule == nil {
		return fmt.Errorf("no rankings schedule has been created")
	}
	return sm.rankingsSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// TriggerReportsNow triggers the immediate execution of the reports schedule
func (sm *ScheduleManager) TriggerReportsNow(ctx context.Context) error {
	if sm.reportsSchedule == nil {
		return fmt.Errorf("no reports schedule has been created")
	}
	return sm.reportsSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// TriggerBuildsNow triggers the immediate execution of the builds schedule
func (sm *ScheduleManager) TriggerBuildsNow(ctx context.Context) error {
	if sm.buildsSchedule == nil {
		return fmt.Errorf("no builds schedule has been created")
	}
	return sm.buildsSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// TriggerSyncNow triggers the immediate execution of the production schedule
// Legacy, will be removed soon
// Not used anymore
func (sm *ScheduleManager) TriggerSyncNow(ctx context.Context) error {
	if sm.prodSchedule == nil {
		return fmt.Errorf("no production schedule has been created")
	}
	return sm.prodSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// TriggerTestNow triggers the immediate execution of the test schedule
// Legacy, will be removed soon
// Not used anymore
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

// PauseRankingsSchedule pauses the rankings schedule
func (sm *ScheduleManager) PauseRankingsSchedule(ctx context.Context) error {
	if sm.rankingsSchedule == nil {
		return fmt.Errorf("no rankings schedule has been created")
	}
	return sm.rankingsSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// PauseReportsSchedule pauses the reports schedule
func (sm *ScheduleManager) PauseReportsSchedule(ctx context.Context) error {
	if sm.reportsSchedule == nil {
		return fmt.Errorf("no reports schedule has been created")
	}
	return sm.reportsSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// PauseBuildsSchedule pauses the builds schedule
func (sm *ScheduleManager) PauseBuildsSchedule(ctx context.Context) error {
	if sm.buildsSchedule == nil {
		return fmt.Errorf("no builds schedule has been created")
	}
	return sm.buildsSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// PauseSchedule pauses the production schedule
// Legacy, will be removed soon
// Not used anymore
func (sm *ScheduleManager) PauseSchedule(ctx context.Context) error {
	if sm.prodSchedule == nil {
		return fmt.Errorf("no production schedule has been created")
	}
	return sm.prodSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// PauseTestSchedule pauses the test schedule
// Legacy, will be removed soon
// Not used anymore
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

// UnpauseRankingsSchedule reactivates the rankings schedule
func (sm *ScheduleManager) UnpauseRankingsSchedule(ctx context.Context) error {
	if sm.rankingsSchedule == nil {
		return fmt.Errorf("no rankings schedule has been created")
	}
	return sm.rankingsSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// UnpauseReportsSchedule reactivates the reports schedule
func (sm *ScheduleManager) UnpauseReportsSchedule(ctx context.Context) error {
	if sm.reportsSchedule == nil {
		return fmt.Errorf("no reports schedule has been created")
	}
	return sm.reportsSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// UnpauseBuildsSchedule reactivates the builds schedule
func (sm *ScheduleManager) UnpauseBuildsSchedule(ctx context.Context) error {
	if sm.buildsSchedule == nil {
		return fmt.Errorf("no builds schedule has been created")
	}
	return sm.buildsSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// UnpauseSchedule reactivates the production schedule
// Legacy, will be removed soon
// Not used anymore
func (sm *ScheduleManager) UnpauseSchedule(ctx context.Context) error {
	if sm.prodSchedule == nil {
		return fmt.Errorf("no production schedule has been created")
	}
	return sm.prodSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// UnpauseTestSchedule reactivates the test schedule
// Legacy, will be removed soon
// Not used anymore
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

// CleanupDecoupledSchedules cleans up the decoupled schedules
func (sm *ScheduleManager) CleanupDecoupledSchedules(ctx context.Context) error {
	// List and delete decoupled schedules
	schedules := []string{rankingsScheduleID, reportsScheduleID, buildsScheduleID}
	for _, id := range schedules {
		handle := sm.client.ScheduleClient().GetHandle(ctx, id)
		if err := handle.Delete(ctx); err != nil {
			sm.logger.Printf("[WARN] Failed to delete schedule %s: %v", id, err)
			continue
		}
		sm.logger.Printf("[INFO] Deleted schedule: %s", id)
	}

	// Reset references
	sm.rankingsSchedule = nil
	sm.reportsSchedule = nil
	sm.buildsSchedule = nil

	return nil
}

// CleanupExistingSchedules deletes existing schedules
// Legacy, will be removed soon
// Not used anymore
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
// Legacy, will be removed soon
// Not used anymore
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
