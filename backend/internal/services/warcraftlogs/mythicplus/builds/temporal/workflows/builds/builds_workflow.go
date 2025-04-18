// builds/temporal/workflows/builds/builds_workflow.go
package warcraftlogsBuildsTemporalWorkflowsBuilds

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// BuildsWorkflow implements the builds extraction workflow
type BuildsWorkflow struct{}

// NewBuildsWorkflow creates a new instance of the builds workflow
func NewBuildsWorkflow() definitions.BuildsWorkflow {
	return &BuildsWorkflow{}
}

// Execute runs the builds workflow
func (w *BuildsWorkflow) Execute(ctx workflow.Context, params models.BuildsWorkflowParams) (*models.BuildsWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting builds workflow",
		"batchSizeParam", params.BatchSize,
		"numWorkersParam", params.NumWorkers,
		"reportBatchSizeParam", params.ReportBatchSize)

	// Generate a unique batchID if not provided
	if params.BatchID == "" {
		params.BatchID = fmt.Sprintf("builds-workflow-%s", uuid.New().String())
		logger.Info("Generated workflow batch ID for tracking", "batchID", params.BatchID)
	}

	// Initialize the final result
	finalResult := &models.BuildsWorkflowResult{
		StartedAt:         workflow.Now(ctx),
		BatchID:           params.BatchID,
		BuildsByClassSpec: make(map[string]int32),
	}

	// WorkflowState Tracking
	workflowInfo := workflow.GetInfo(ctx)
	workflowStateID := fmt.Sprintf("builds-%s", workflowInfo.WorkflowExecution.ID)
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

	// Create the initial workflow state - Using string for the activity name
	err := workflow.ExecuteActivity(stateCtx, definitions.CreateWorkflowStateActivity, &warcraftlogsBuilds.WorkflowState{
		ID:           workflowStateID,
		WorkflowType: "builds",
		StartedAt:    finalResult.StartedAt,
		Status:       "running",
		CreatedAt:    finalResult.StartedAt,
		UpdatedAt:    finalResult.StartedAt,
	}).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to create workflow state", "error", err)
	}

	// Configuration des Activités Principales
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 12,
		HeartbeatTimeout:    time.Minute * 15,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 15,
			MaximumAttempts:    3,
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

	// Récupération des rapports à traiter
	var reportsToProcess []*warcraftlogsBuilds.Report
	err = workflow.ExecuteActivity(activityCtx,
		definitions.GetReportsNeedingBuildExtractionActivity,
		params.BatchSize,
	).Get(ctx, &reportsToProcess)

	if err != nil {
		logger.Error("Failed to get reports needing build extraction", "error", err)
		finalWorkflowState := &warcraftlogsBuilds.WorkflowState{ID: workflowStateID, Status: "failed", ErrorMessage: fmt.Sprintf("GetReportsNeeding... failed: %v", err), UpdatedAt: workflow.Now(ctx)}
		_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, finalWorkflowState).Get(ctx, nil)
		finalResult.CompletedAt = workflow.Now(ctx)
		return finalResult, err
	}

	logger.Info("Retrieved reports for build extraction", "count", len(reportsToProcess))

	if len(reportsToProcess) == 0 {
		logger.Info("No reports need build extraction.")
		finalWorkflowState := &warcraftlogsBuilds.WorkflowState{ID: workflowStateID, Status: "completed", CompletedAt: workflow.Now(ctx), ItemsProcessed: 0, UpdatedAt: workflow.Now(ctx)}
		_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, finalWorkflowState).Get(ctx, nil)
		finalResult.CompletedAt = workflow.Now(ctx)
		return finalResult, nil
	}

	// Traitement des rapports par lots
	batchProcessingSize := 5 // Default value
	if params.ReportBatchSize > 0 {
		batchProcessingSize = int(params.ReportBatchSize)
	}

	totalBuildsProcessed := int32(0)
	totalReportsProcessedSuccess := int32(0)
	totalReportsFailedProcessing := int32(0)
	buildsByClassSpecAggregated := make(map[string]int32)
	numReportsToProcess := len(reportsToProcess)

	for i := 0; i < numReportsToProcess; i += batchProcessingSize {
		end := i + batchProcessingSize
		if end > numReportsToProcess {
			end = numReportsToProcess
		}
		batch := reportsToProcess[i:end]

		batchNum := (i / batchProcessingSize) + 1
		logger.Info("Processing builds extraction batch")

		// Update the progress state
		workflowStateUpdate := &warcraftlogsBuilds.WorkflowState{
			ID:              workflowStateID,
			LastProcessedID: fmt.Sprintf("report_batch_start_index_%d", i),
			ItemsProcessed:  int(totalReportsProcessedSuccess),
			UpdatedAt:       workflow.Now(ctx),
			Status:          "running", // Ensure the status remains 'running' during processing
		}
		_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowStateUpdate)

		// Execute the activity that processes builds and marks reports
		var activityResult models.BuildsActivityResult
		activityName := definitions.ProcessBuildsActivity

		activityErr := workflow.ExecuteActivity(activityCtx,
			activityName,
			batch,
		).Get(ctx, &activityResult)

		if activityErr != nil {
			logger.Error("Build processing activity failed fatally for batch", "batchNum", batchNum, "error", activityErr)
			totalReportsFailedProcessing += int32(len(batch))
			continue
		}

		// The activity has finished, aggregate its results
		totalBuildsProcessed += activityResult.ProcessedBuildsCount
		totalReportsProcessedSuccess += activityResult.SuccessCount
		totalReportsFailedProcessing += activityResult.FailureCount

		for classSpec, count := range activityResult.BuildsByClassSpec {
			buildsByClassSpecAggregated[classSpec] += count
		}

		logger.Info("Completed processing builds for reports batch")

		workflow.Sleep(ctx, time.Second*1)
	}

	// Finalization
	finalResult.BuildsProcessed = totalBuildsProcessed
	finalResult.ReportsProcessed = totalReportsProcessedSuccess
	finalResult.BuildsByClassSpec = buildsByClassSpecAggregated
	finalResult.CompletedAt = workflow.Now(ctx)

	// Update the final workflow state
	finalStatus := "completed"
	finalErrorMsg := ""
	if totalReportsFailedProcessing > 0 {
		finalStatus = "completed_with_errors"
		finalErrorMsg = fmt.Sprintf("%d reports failed during build processing/marking", totalReportsFailedProcessing)
	}
	finalWorkflowState := &warcraftlogsBuilds.WorkflowState{
		ID:             workflowStateID,
		Status:         finalStatus,
		ErrorMessage:   finalErrorMsg,
		CompletedAt:    finalResult.CompletedAt,
		ItemsProcessed: int(totalReportsProcessedSuccess),
		UpdatedAt:      finalResult.CompletedAt,
	}
	err = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, finalWorkflowState).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to update final workflow state", "error", err)
	}

	logger.Info("Builds workflow completed")

	return finalResult, nil
}
