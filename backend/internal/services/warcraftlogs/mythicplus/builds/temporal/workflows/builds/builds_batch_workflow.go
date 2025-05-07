package warcraftlogsBuildsTemporalWorkflowsBuilds

import (
	"fmt" // Ajout de l'import nécessaire
	"strings"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds" // Ajout de l'import pour les modèles
	warcraftlogsBuildMetrics "wowperf/internal/services/warcraftlogs/mythicplus/builds/metrics"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// BuildsBatchWorkflow processes a batch of reports for builds extraction
type BuildsBatchWorkflow struct{}

// NewBuildsBatchWorkflow creates a new instance of BuildsBatchWorkflow
func NewBuildsBatchWorkflow() definitions.BuildsBatchWorkflow {
	return &BuildsBatchWorkflow{}
}

// Execute runs the workflow
func (w *BuildsBatchWorkflow) Execute(ctx workflow.Context, params models.BuildsBatchParams) (*models.BuildsBatchResult, error) {
	logger := workflow.GetLogger(ctx)

	// Initialize metrics collector
	metrics := warcraftlogsBuildMetrics.NewMetricsCollector(
		"builds_batch_workflow",
		"",
		params.BatchID,
	)

	logger.Info("Starting builds batch workflow",
		"batchID", params.BatchID,
		"reportsCount", len(params.Reports),
		"parentWorkflowStateID", params.ParentWorkflowStateID) // Ajout du log pour le parentWorkflowStateID

	// Initialize result
	result := &models.BuildsBatchResult{
		StartedAt:         workflow.Now(ctx),
		BatchID:           params.BatchID,
		BuildsByClassSpec: make(map[string]int32),
	}

	// Skip if no reports to process
	if len(params.Reports) == 0 {
		logger.Info("No reports to process in this batch")
		result.CompletedAt = workflow.Now(ctx)
		metrics.Finish("completed_empty")
		return result, nil
	}

	// Ajout : WorkflowState Tracking pour le workflow enfant
	workflowInfo := workflow.GetInfo(ctx)
	workflowID := workflowInfo.WorkflowExecution.ID
	runID := workflowInfo.WorkflowExecution.RunID
	workflowStateID := fmt.Sprintf("builds-batch-%s-%s", workflowID, runID)

	// Activity options for workflow state
	stateOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.5,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	stateCtx := workflow.WithActivityOptions(ctx, stateOpts)

	// Create initial workflow state for the child workflow
	initialState := &warcraftlogsBuilds.WorkflowState{
		ID:               workflowStateID,
		WorkflowType:     "builds_batch",
		ParentWorkflowID: params.ParentWorkflowStateID, // Référence au workflow state parent
		BatchID:          params.BatchID,
		StartedAt:        workflow.Now(ctx),
		Status:           "running",
		CreatedAt:        workflow.Now(ctx),
		UpdatedAt:        workflow.Now(ctx),
	}

	// Set initial metrics on workflow state
	metrics.UpdateWorkflowState(initialState)

	err := workflow.ExecuteActivity(stateCtx, definitions.CreateWorkflowStateActivity, initialState).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to create workflow state for batch", "error", err)
		metrics.RecordError("workflow_state_creation")
		// Continue execution even if state tracking fails
	}

	// Configure activity options
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 2, // Shorter timeout for child workflow
		HeartbeatTimeout:    time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    3,
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

	// Process all reports in a single activity call
	// This activity will handle the extraction and storage of builds
	metrics.StartOperation("process_builds")
	var activityResult models.BuildsActivityResult
	err = workflow.ExecuteActivity(activityCtx,
		definitions.ProcessBuildsActivity,
		params.Reports,
	).Get(ctx, &activityResult)
	metrics.EndOperation("process_builds")

	if err != nil {
		logger.Error("Failed to process builds", "error", err)
		metrics.RecordError("process_builds_activity")
		metrics.Finish("failed")

		// Update workflow state with error
		workflowState := &warcraftlogsBuilds.WorkflowState{
			ID:           workflowStateID,
			Status:       "failed",
			ErrorMessage: fmt.Sprintf("Failed to process builds: %v", err),
			UpdatedAt:    workflow.Now(ctx),
		}
		metrics.UpdateWorkflowState(workflowState)
		_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

		// Provide a minimal result with error information
		result.Status = "failed"
		result.Error = err.Error()
		result.CompletedAt = workflow.Now(ctx)
		return result, err
	}

	// Update metrics with results
	metrics.RecordItemsProcessed(int(activityResult.ProcessedBuildsCount), "success")
	metrics.SetBuildsMetrics(int(activityResult.ProcessedBuildsCount))

	// Record builds by class/spec
	for classSpec, count := range activityResult.BuildsByClassSpec {
		metrics.RecordBuildByClassSpec(
			classSpec[:strings.Index(classSpec, "-")],   // Class
			classSpec[strings.Index(classSpec, "-")+1:], // Spec
			int(count),
		)
	}

	// Update result
	result.BuildsProcessed = activityResult.ProcessedBuildsCount
	result.ReportsProcessed = activityResult.SuccessCount
	result.BuildsByClassSpec = activityResult.BuildsByClassSpec
	result.Status = "completed"
	result.CompletedAt = workflow.Now(ctx)

	// Update final workflow state
	finalState := &warcraftlogsBuilds.WorkflowState{
		ID:             workflowStateID,
		Status:         "completed",
		CompletedAt:    result.CompletedAt,
		ItemsProcessed: int(result.ReportsProcessed),
		UpdatedAt:      result.CompletedAt,
	}
	metrics.UpdateWorkflowState(finalState)
	_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, finalState).Get(ctx, nil)

	logger.Info("Builds batch workflow completed",
		"batchID", params.BatchID,
		"reportsProcessed", result.ReportsProcessed,
		"buildsProcessed", result.BuildsProcessed,
		"duration", result.CompletedAt.Sub(result.StartedAt))

	metrics.Finish("completed")
	return result, nil
}
