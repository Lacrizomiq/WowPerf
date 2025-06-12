package warcraftlogsBuildsTemporalWorkflowsReports

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	warcraftlogsBuildMetrics "wowperf/internal/services/warcraftlogs/mythicplus/builds/metrics"
	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// ReportsWorkflow implements the reports workflow
type ReportsWorkflow struct{}

// NewReportsWorkflow creates a new instance of the reports workflow
func NewReportsWorkflow() definitions.ReportsWorkflow {
	return &ReportsWorkflow{}
}

// Execute runs the reports workflow
func (w *ReportsWorkflow) Execute(ctx workflow.Context, params models.ReportsWorkflowParams) (*models.ReportsWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting reports workflow",
		"class", params.ClassName,
		"batchID", params.BatchID,
		"batchSize", params.BatchSize,
		"processingWindow", params.ProcessingWindow,
		"continuationCount", params.ContinuationCount,
		"parentWorkflowID", params.ParentWorkflowID)

	// Initialize metrics collector
	metrics := warcraftlogsBuildMetrics.NewMetricsCollector(
		"reports_workflow",
		params.ClassName,
		params.BatchID,
	)
	defer metrics.Finish("completed")

	// Initialize metrics from params if it's a continuation
	if params.ContinuationCount > 0 {
		metrics.SetReportsMetrics(
			int(params.TotalProcessedRankings),
			int(params.TotalProcessedReports),
			int(params.TotalFailedReports),
		)
		metrics.RecordItemsProcessed(int(params.TotalProcessedRankings), "carried_over")

		// Set continuation count from params
		for i := 0; i < int(params.ContinuationCount); i++ {
			metrics.RecordContinuation()
		}
	}

	// Initialize result with start time, batch ID and previous totals
	result := &models.ReportsWorkflowResult{
		StartedAt:         workflow.Now(ctx),
		BatchID:           params.BatchID,
		ReportsProcessed:  params.TotalProcessedReports,
		RankingsProcessed: params.TotalProcessedRankings,
		FailedReports:     params.TotalFailedReports,
		APIRequestsCount:  params.TotalAPIRequests,
	}

	// WorkflowState Tracking
	workflowInfo := workflow.GetInfo(ctx)
	workflowID := workflowInfo.WorkflowExecution.ID
	runID := workflowInfo.WorkflowExecution.RunID
	workflowStateID := fmt.Sprintf("reports-%s-%s-%s", params.ClassName, workflowID, runID)

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

	// Create initial workflow state
	initialState := &warcraftlogsBuilds.WorkflowState{
		ID:                workflowStateID,
		WorkflowType:      "reports",
		StartedAt:         workflow.Now(ctx),
		Status:            "running",
		BatchID:           params.BatchID,
		ClassName:         params.ClassName,
		ItemsProcessed:    int(params.TotalProcessedRankings),
		ContinuationCount: int(params.ContinuationCount),
		ApiRequestsCount:  int(params.TotalAPIRequests),
		CreatedAt:         workflow.Now(ctx),
		UpdatedAt:         workflow.Now(ctx),
	}

	// Define the parent_workflow_id if it's a continuation
	if params.ContinuationCount > 0 && params.ParentWorkflowID != "" {
		initialState.ParentWorkflowID = params.ParentWorkflowID
		logger.Info("Setting parent workflow ID for continuation",
			"parentWorkflowID", params.ParentWorkflowID,
			"continuationCount", params.ContinuationCount)
	}

	// Set initial metrics on workflow state
	metrics.UpdateWorkflowState(initialState)

	err := workflow.ExecuteActivity(stateCtx, definitions.CreateWorkflowStateActivity, initialState).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to create workflow state", "error", err)
		metrics.RecordError("workflow_state_creation")
		// Continue execution even if state tracking fails
	}

	// Configure activity options for reports processing
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 12,
		HeartbeatTimeout:    time.Minute * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    3,
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

	// Process reports in batches
	localReportsProcessed := 0
	localRankingsProcessed := 0
	localAPIRequestsCount := 0
	localFailedReports := 0
	batchNumber := 0

	// Define a limit of batches before ContinueAsNew
	const maxBatchesBeforeContinue = 25 // Adjust this value as needed

	// External loop to retrieve all rankings until exhaustion
	for {
		batchNumber++

		// Get rankings needing report processing
		metrics.StartOperation("get_rankings_batch")
		var rankingsToProcess []*warcraftlogsBuilds.ClassRanking
		err = workflow.ExecuteActivity(activityCtx,
			definitions.GetRankingsNeedingReportProcessingActivity,
			params.ClassName, params.BatchSize, params.ProcessingWindow).Get(ctx, &rankingsToProcess)
		metrics.EndOperation("get_rankings_batch")

		if err != nil {
			logger.Error("Failed to get rankings needing report processing", "error", err)
			metrics.RecordError("get_rankings_activity")

			// Update workflow state with error
			workflowState := &warcraftlogsBuilds.WorkflowState{
				ID:           workflowStateID,
				Status:       "failed",
				ErrorMessage: fmt.Sprintf("Failed to get rankings: %v", err),
				UpdatedAt:    workflow.Now(ctx),
			}
			metrics.UpdateWorkflowState(workflowState)
			_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

			return result, err
		}

		if len(rankingsToProcess) == 0 {
			logger.Info("No more rankings need report processing")
			break
		}

		// Set total to process to help with progress percentage
		metrics.SetTotalToProcess(len(rankingsToProcess) + int(params.TotalProcessedRankings) + localRankingsProcessed)

		// Progress log
		logger.Info("Processing batch of rankings",
			"batchNumber", batchNumber,
			"batchSize", len(rankingsToProcess),
			"totalProcessedSoFar", params.TotalProcessedRankings+int32(localRankingsProcessed),
			"remainingToProcess", len(rankingsToProcess))

		// Processing retrieved batches
		const batchSize = 10
		batchRankingsProcessed := len(rankingsToProcess)

		for i := 0; i < len(rankingsToProcess); i += batchSize {
			end := i + batchSize
			if end > len(rankingsToProcess) {
				end = len(rankingsToProcess)
			}

			batch := rankingsToProcess[i:end]

			// Reduce detailed logs for each sub-batch
			if i == 0 || i+batchSize >= len(rankingsToProcess) {
				logger.Info("Processing rankings sub-batch",
					"batchSize", len(batch),
					"progress", fmt.Sprintf("%d/%d", i+len(batch), len(rankingsToProcess)))
			}

			// Update workflow state with current batch progress
			workflowState := &warcraftlogsBuilds.WorkflowState{
				ID:              workflowStateID,
				LastProcessedID: fmt.Sprintf("batch-%d-%d", batchNumber, i/batchSize),
				ItemsProcessed:  int(params.TotalProcessedRankings) + localRankingsProcessed + i,
				UpdatedAt:       workflow.Now(ctx),
			}
			metrics.UpdateWorkflowState(workflowState)
			_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

			// Process the batch through the API
			metrics.StartOperation("process_reports_batch")
			var batchResult models.ReportProcessingResult
			err := workflow.ExecuteActivity(activityCtx,
				definitions.ProcessReportsActivity,
				batch).Get(ctx, &batchResult)
			metrics.EndOperation("process_reports_batch")

			if err != nil {
				if common.IsRateLimitError(err) {
					// Rate limit case
					metrics.RecordRateLimitHit()

					// Update workflow state for rate limit
					workflowState := &warcraftlogsBuilds.WorkflowState{
						ID:           workflowStateID,
						Status:       "rate_limited",
						ErrorMessage: fmt.Sprintf("Rate limit reached: %v", err),
						UpdatedAt:    workflow.Now(ctx),
					}
					metrics.UpdateWorkflowState(workflowState)
					_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

					logger.Info("Rate limit reached during reports processing")

					// Synchronize state before continuing as new
					common.SyncStateBeforeContinueAsNew(ctx, workflowStateID, "reports_workflow", metrics, "continuing")

					// Update params for ContinueAsNew
					updatedParams := common.SyncReportsWorkflowParams(
						params,
						metrics,
						localRankingsProcessed,
						localReportsProcessed,
						localAPIRequestsCount,
						localFailedReports,
					)

					// Define the parent_workflow_id if it's a continuation
					if params.ParentWorkflowID == "" {
						// If it's the first continuation, use the current ID as parent
						updatedParams.ParentWorkflowID = workflowStateID
						logger.Info("Setting workflow as parent for first continuation",
							"parentID", workflowStateID)
					} else {
						// Keep the same parent
						updatedParams.ParentWorkflowID = params.ParentWorkflowID
						logger.Info("Keeping existing parent for continuation",
							"parentID", params.ParentWorkflowID)
					}

					// Continue with a new workflow
					return nil, workflow.NewContinueAsNewError(ctx, definitions.ReportsWorkflowName, updatedParams)
				}

				logger.Error("Failed to process reports batch", "error", err)
				metrics.RecordError("process_reports_activity")
				localFailedReports += len(batch)
				continue
			}

			// Update metrics for this batch
			localReportsProcessed += int(batchResult.ProcessedCount)
			localAPIRequestsCount += int(len(batch) * 2) // Approximately 2 API calls per ranking
			metrics.RecordItemsProcessed(int(len(batch)), "success")

			for i := 0; i < len(batch)*2; i++ {
				metrics.RecordAPIRequest()
			}

			// Extract report codes for marking
			if len(batchResult.ProcessedReports) > 0 {
				// Mark reports for build processing
				metrics.StartOperation("mark_reports_for_processing")
				err = workflow.ExecuteActivity(activityCtx,
					definitions.MarkReportsForBuildProcessingActivity,
					batchResult.ProcessedReports, params.BatchID).Get(ctx, nil)
				metrics.EndOperation("mark_reports_for_processing")

				if err != nil {
					logger.Error("Failed to mark reports for build processing", "error", err)
					metrics.RecordError("mark_reports_activity")
					// Continue even if marking fails
				}
			}

			// Small delay between batches
			workflow.Sleep(ctx, time.Second*2)
		}

		// Update the total count of processed rankings after each main batch
		localRankingsProcessed += batchRankingsProcessed

		logger.Info("Batch complete",
			"batchNumber", batchNumber,
			"processedInBatch", batchRankingsProcessed,
			"totalProcessed", params.TotalProcessedRankings+int32(localRankingsProcessed),
			"progress", fmt.Sprintf("Processed %d rankings so far", params.TotalProcessedRankings+int32(localRankingsProcessed)))

		// Check if we should continue with a new workflow
		if batchNumber >= maxBatchesBeforeContinue {
			logger.Info("Reached batch limit, continuing as new workflow",
				"processedBatches", batchNumber,
				"localProcessed", localRankingsProcessed,
				"totalProcessed", params.TotalProcessedRankings+int32(localRankingsProcessed))

			// Synchronize state before continuing as new
			common.SyncStateBeforeContinueAsNew(ctx, workflowStateID, "reports_workflow", metrics, "continuing")

			// Update params for ContinueAsNew
			updatedParams := common.SyncReportsWorkflowParams(
				params,
				metrics,
				localRankingsProcessed,
				localReportsProcessed,
				localAPIRequestsCount,
				localFailedReports,
			)

			// Define the parent_workflow_id if it's a continuation
			if params.ParentWorkflowID == "" {
				// If it's the first continuation, use the current ID as parent
				updatedParams.ParentWorkflowID = workflowStateID
				logger.Info("Setting workflow as parent for first continuation",
					"parentID", workflowStateID)
			} else {
				// Keep the same parent
				updatedParams.ParentWorkflowID = params.ParentWorkflowID
				logger.Info("Keeping existing parent for continuation",
					"parentID", params.ParentWorkflowID)
			}

			// Continue with a new workflow (clean history)
			return nil, workflow.NewContinueAsNewError(ctx, definitions.ReportsWorkflowName, updatedParams)
		}

		// Small delay between large data retrievals
		workflow.Sleep(ctx, time.Second*5)
	}

	// Set final results
	result.ReportsProcessed = params.TotalProcessedReports + int32(localReportsProcessed)
	result.RankingsProcessed = params.TotalProcessedRankings + int32(localRankingsProcessed)
	result.FailedReports = params.TotalFailedReports + int32(localFailedReports)
	result.APIRequestsCount = params.TotalAPIRequests + int32(localAPIRequestsCount)
	result.CompletedAt = workflow.Now(ctx)

	// Complete workflow state with metrics
	workflowState := &warcraftlogsBuilds.WorkflowState{
		ID:             workflowStateID,
		Status:         "completed",
		CompletedAt:    workflow.Now(ctx),
		ItemsProcessed: int(params.TotalProcessedRankings) + localRankingsProcessed,
		UpdatedAt:      workflow.Now(ctx),
	}
	metrics.UpdateWorkflowState(workflowState)
	_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

	logger.Info("Reports workflow completed",
		"rankingsProcessed", result.RankingsProcessed,
		"reportsProcessed", result.ReportsProcessed,
		"failedReports", result.FailedReports,
		"apiRequests", result.APIRequestsCount,
		"duration", result.CompletedAt.Sub(result.StartedAt),
		"continuationCount", params.ContinuationCount)

	return result, nil
}
