package warcraftlogsBuildsTemporalWorkflowsReports

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
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
		"batchID", params.BatchID,
		"batchSize", params.BatchSize,
		"processingWindow", params.ProcessingWindow)

	// Initialize result with start time and batch ID
	result := &models.ReportsWorkflowResult{
		StartedAt: workflow.Now(ctx),
		BatchID:   params.BatchID,
	}

	// WorkflowState Tracking
	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	workflowStateID := fmt.Sprintf("reports-%s", workflowID)

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

	// Create workflow state
	err := workflow.ExecuteActivity(stateCtx, definitions.CreateWorkflowStateActivity, &warcraftlogsBuilds.WorkflowState{
		ID:              workflowStateID,
		WorkflowType:    "reports",
		StartedAt:       workflow.Now(ctx),
		Status:          "running",
		ItemsProcessed:  0,
		LastProcessedID: "",
		CreatedAt:       workflow.Now(ctx),
		UpdatedAt:       workflow.Now(ctx),
	}).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to create workflow state", "error", err)
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
	totalReportsProcessed := 0
	totalRankingsProcessed := 0
	apiRequestsCount := 0
	failedReports := 0
	batchNumber := 0

	// Add a log to indicate the start of batch processing
	logger.Info("Starting to process all rankings in batches",
		"batchSize", params.BatchSize)

	// External loop to retrieve all rankings until exhaustion
	for {
		batchNumber++

		// Get rankings needing report processing
		var rankingsToProcess []*warcraftlogsBuilds.ClassRanking
		err = workflow.ExecuteActivity(activityCtx,
			definitions.GetRankingsNeedingReportProcessingActivity,
			params.BatchSize, params.ProcessingWindow).Get(ctx, &rankingsToProcess)

		if err != nil {
			logger.Error("Failed to get rankings needing report processing", "error", err)

			// Update workflow state with error
			workflowState := &warcraftlogsBuilds.WorkflowState{
				ID:           workflowStateID,
				Status:       "failed",
				ErrorMessage: fmt.Sprintf("Failed to get rankings: %v", err),
				UpdatedAt:    workflow.Now(ctx),
			}
			_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

			return result, err
		}

		if len(rankingsToProcess) == 0 {
			logger.Info("No more rankings need report processing")
			break
		}

		// Improvement of the general progress log
		logger.Info("Processing batch of rankings",
			"batchNumber", batchNumber,
			"batchSize", len(rankingsToProcess),
			"totalProcessedSoFar", totalRankingsProcessed,
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

			// Update workflow state with current batch
			workflowState := &warcraftlogsBuilds.WorkflowState{
				ID:              workflowStateID,
				LastProcessedID: fmt.Sprintf("batch-%d-%d", batchNumber, i/batchSize),
				ItemsProcessed:  totalRankingsProcessed + i,
				UpdatedAt:       workflow.Now(ctx),
			}
			_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

			// Process the batch through the API to fetch and store reports
			var batchResult models.ReportProcessingResult
			err := workflow.ExecuteActivity(activityCtx,
				definitions.ProcessReportsActivity,
				batch).Get(ctx, &batchResult)

			if err != nil {
				if common.IsRateLimitError(err) {
					// Update workflow state for rate limit
					workflowState := &warcraftlogsBuilds.WorkflowState{
						ID:           workflowStateID,
						Status:       "rate_limited",
						ErrorMessage: fmt.Sprintf("Rate limit reached: %v", err),
						UpdatedAt:    workflow.Now(ctx),
					}
					_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

					logger.Info("Rate limit reached during reports processing")

					result.CompletedAt = workflow.Now(ctx)
					result.ReportsProcessed = int32(totalReportsProcessed)
					result.RankingsProcessed = int32(totalRankingsProcessed)
					result.FailedReports = int32(failedReports)
					result.APIRequestsCount = int32(apiRequestsCount)

					return result, err
				}

				logger.Error("Failed to process reports batch", "error", err)
				failedReports += len(batch)
				continue
			}

			// Extract report codes for marking
			if len(batchResult.ProcessedReports) > 0 {
				// Mark reports for build processing
				err = workflow.ExecuteActivity(activityCtx,
					definitions.MarkReportsForBuildProcessingActivity,
					batchResult.ProcessedReports, params.BatchID).Get(ctx, nil)

				if err != nil {
					logger.Error("Failed to mark reports for build processing", "error", err)
					// Continue even if marking fails
				}
			}

			totalReportsProcessed += int(batchResult.ProcessedCount)
			apiRequestsCount += int(len(batch) * 2) // Approximately 2 API calls per ranking

			// Log only for the last sub-batch to reduce verbosity
			if i+batchSize >= len(rankingsToProcess) {
				logger.Info("Completed processing reports batch",
					"batchProcessed", len(batchResult.ProcessedReports),
					"totalProcessed", totalReportsProcessed)
			}

			// Small delay between batches
			workflow.Sleep(ctx, time.Second*2)
		}

		// Update the total count of processed rankings after each main batch
		totalRankingsProcessed += batchRankingsProcessed

		logger.Info("Batch complete",
			"batchNumber", batchNumber,
			"processedInBatch", batchRankingsProcessed,
			"totalProcessed", totalRankingsProcessed,
			"progress", fmt.Sprintf("Processed %d rankings so far", totalRankingsProcessed))

		// Small delay between large data retrievals
		workflow.Sleep(ctx, time.Second*5)
	}

	// Set final results
	result.ReportsProcessed = int32(totalReportsProcessed)
	result.RankingsProcessed = int32(totalRankingsProcessed)
	result.FailedReports = int32(failedReports)
	result.APIRequestsCount = int32(apiRequestsCount)
	result.CompletedAt = workflow.Now(ctx)

	// Complete workflow state
	workflowState := &warcraftlogsBuilds.WorkflowState{
		ID:             workflowStateID,
		Status:         "completed",
		CompletedAt:    workflow.Now(ctx),
		ItemsProcessed: totalReportsProcessed,
		UpdatedAt:      workflow.Now(ctx),
	}
	_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

	// Keep this final log exactly as you requested
	logger.Info("Reports workflow completed",
		"rankingsProcessed", totalRankingsProcessed,
		"reportsProcessed", totalReportsProcessed,
		"failedReports", failedReports,
		"apiRequests", apiRequestsCount,
		"duration", result.CompletedAt.Sub(result.StartedAt))

	return result, nil
}
