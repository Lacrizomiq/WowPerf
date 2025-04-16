package warcraftlogsBuildsTemporalScheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

// Constants for schedule IDs
const (

	// New schedules for decoupled workflows
	rankingsScheduleID          = "warcraft-logs-rankings"
	reportsScheduleID           = "warcraft-logs-reports"
	buildsScheduleID            = "warcraft-logs-builds"
	equipmentAnalysisScheduleID = "warcraft-logs-equipment-analysis"
	talentAnalysisScheduleID    = "warcraft-logs-talent-analysis"
	statAnalysisScheduleID      = "warcraft-logs-stat-analysis"
)

// ScheduleManager manages Temporal schedules for WarcraftLogs workflows
type ScheduleManager struct {
	client client.Client
	logger *log.Logger

	// New handles (decoupled workflows, will be used instead of the existing ones)
	rankingsSchedule          client.ScheduleHandle
	reportsSchedule           client.ScheduleHandle
	buildsSchedule            client.ScheduleHandle
	equipmentAnalysisSchedule client.ScheduleHandle
	talentAnalysisSchedule    client.ScheduleHandle
	statAnalysisSchedule      client.ScheduleHandle
}

// NewScheduleManager creates a new ScheduleManager instance
func NewScheduleManager(temporalClient client.Client, logger *log.Logger) *ScheduleManager {
	return &ScheduleManager{
		client: temporalClient,
		logger: logger,
	}
}

// == Creation of schedules ==

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
			Workflow:  definitions.RankingsWorkflowName,
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
			Workflow:  definitions.ReportsWorkflowName,
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
			Workflow:  definitions.BuildsWorkflowName,
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

// CreateEquipmentAnalysisSchedule creates the equipment analysis workflow schedule
func (sm *ScheduleManager) CreateEquipmentAnalysisSchedule(ctx context.Context, params *models.EquipmentAnalysisWorkflowParams, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := equipmentAnalysisScheduleID
	workflowID := fmt.Sprintf("warcraft-logs-equipment-analysis-%s", time.Now().UTC().Format("2006-01-02"))

	// Create the schedule without automatic triggering (No CRON expressions)
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		// No CronExpressions to avoid automatic triggering
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  definitions.AnalyzeBuildsWorkflowName,
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
		return fmt.Errorf("failed to create equipment analysis schedule: %w", err)
	}

	sm.equipmentAnalysisSchedule = handle
	sm.logger.Printf("[INFO] Created equipment analysis workflow schedule: %s", scheduleID)
	return nil
}

// CreateTalentAnalysisSchedule creates the talent analysis workflow schedule
func (sm *ScheduleManager) CreateTalentAnalysisSchedule(ctx context.Context, params *models.TalentAnalysisWorkflowParams, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := talentAnalysisScheduleID
	workflowID := fmt.Sprintf("warcraft-logs-talent-analysis-%s", time.Now().UTC().Format("2006-01-02"))

	// Create the schedule without automatic triggering (No CRON expressions)
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		// No CronExpressions to avoid automatic triggering
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  definitions.AnalyzeTalentsWorkflowName,
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
		return fmt.Errorf("failed to create talent analysis schedule: %w", err)
	}

	sm.talentAnalysisSchedule = handle
	sm.logger.Printf("[INFO] Created talent analysis workflow schedule: %s", scheduleID)
	return nil
}

// CreateStatAnalysisSchedule creates the stat analysis workflow schedule
func (sm *ScheduleManager) CreateStatAnalysisSchedule(ctx context.Context, params *models.StatAnalysisWorkflowParams, opts *ScheduleOptions) error {
	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	scheduleID := statAnalysisScheduleID
	workflowID := fmt.Sprintf("warcraft-logs-stat-analysis-%s", time.Now().UTC().Format("2006-01-02"))

	// Create the schedule without automatic triggering (No CRON expressions)
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		// No CronExpressions to avoid automatic triggering
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  definitions.AnalyzeStatStatisticsWorkflowName,
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
		return fmt.Errorf("failed to create stat analysis schedule: %w", err)
	}

	sm.statAnalysisSchedule = handle
	sm.logger.Printf("[INFO] Created stat analysis workflow schedule: %s", scheduleID)
	return nil
}

// == Triggering of schedules ==

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

// TriggerEquipmentAnalysisNow triggers the immediate execution of the equipment analysis schedule
func (sm *ScheduleManager) TriggerEquipmentAnalysisNow(ctx context.Context) error {
	if sm.equipmentAnalysisSchedule == nil {
		return fmt.Errorf("no equipment analysis schedule has been created")
	}
	return sm.equipmentAnalysisSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// TriggerTalentAnalysisNow triggers the immediate execution of the talent analysis schedule
func (sm *ScheduleManager) TriggerTalentAnalysisNow(ctx context.Context) error {
	if sm.talentAnalysisSchedule == nil {
		return fmt.Errorf("no talent analysis schedule has been created")
	}
	return sm.talentAnalysisSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// TriggerStatAnalysisNow triggers the immediate execution of the stat analysis schedule
func (sm *ScheduleManager) TriggerStatAnalysisNow(ctx context.Context) error {
	if sm.statAnalysisSchedule == nil {
		return fmt.Errorf("no stat analysis schedule has been created")
	}
	return sm.statAnalysisSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// == Pausing and unpausing of schedules ==

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

// PauseEquipmentAnalysisSchedule pauses the equipment analysis schedule
func (sm *ScheduleManager) PauseEquipmentAnalysisSchedule(ctx context.Context) error {
	if sm.equipmentAnalysisSchedule == nil {
		return fmt.Errorf("no equipment analysis schedule has been created")
	}
	return sm.equipmentAnalysisSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// PauseTalentAnalysisSchedule pauses the talent analysis schedule
func (sm *ScheduleManager) PauseTalentAnalysisSchedule(ctx context.Context) error {
	if sm.talentAnalysisSchedule == nil {
		return fmt.Errorf("no talent analysis schedule has been created")
	}
	return sm.talentAnalysisSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// PauseStatAnalysisSchedule pauses the stat analysis schedule
func (sm *ScheduleManager) PauseStatAnalysisSchedule(ctx context.Context) error {
	if sm.statAnalysisSchedule == nil {
		return fmt.Errorf("no stat analysis schedule has been created")
	}
	return sm.statAnalysisSchedule.Pause(ctx, client.SchedulePauseOptions{})
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

// UnpauseEquipmentAnalysisSchedule reactivates the equipment analysis schedule
func (sm *ScheduleManager) UnpauseEquipmentAnalysisSchedule(ctx context.Context) error {
	if sm.equipmentAnalysisSchedule == nil {
		return fmt.Errorf("no equipment analysis schedule has been created")
	}
	return sm.equipmentAnalysisSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// UnpauseTalentAnalysisSchedule reactivates the talent analysis schedule
func (sm *ScheduleManager) UnpauseTalentAnalysisSchedule(ctx context.Context) error {
	if sm.talentAnalysisSchedule == nil {
		return fmt.Errorf("no talent analysis schedule has been created")
	}
	return sm.talentAnalysisSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// UnpauseStatAnalysisSchedule reactivates the stat analysis schedule
func (sm *ScheduleManager) UnpauseStatAnalysisSchedule(ctx context.Context) error {
	if sm.statAnalysisSchedule == nil {
		return fmt.Errorf("no stat analysis schedule has been created")
	}
	return sm.statAnalysisSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// DeleteSchedule deletes a schedule by its ID
func (sm *ScheduleManager) DeleteSchedule(ctx context.Context, scheduleID string) error {
	handle := sm.client.ScheduleClient().GetHandle(ctx, scheduleID)
	return handle.Delete(ctx)
}

// CleanupDecoupledSchedules cleans up the decoupled schedules
func (sm *ScheduleManager) CleanupDecoupledSchedules(ctx context.Context) error {
	// List and delete decoupled schedules
	schedules := []string{rankingsScheduleID, reportsScheduleID, buildsScheduleID, equipmentAnalysisScheduleID, talentAnalysisScheduleID, statAnalysisScheduleID}
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
	sm.equipmentAnalysisSchedule = nil
	sm.talentAnalysisSchedule = nil
	sm.statAnalysisSchedule = nil

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

	if err := sm.CleanupOldWorkflows(ctx); err != nil {
		return fmt.Errorf("failed to cleanup workflows: %w", err)
	}

	return nil
}
