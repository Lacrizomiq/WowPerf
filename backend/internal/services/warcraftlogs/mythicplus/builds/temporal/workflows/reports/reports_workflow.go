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
	err := workflow.ExecuteActivity(stateCtx, "CreateWorkflowState", &warcraftlogsBuilds.WorkflowState{
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
		_ = workflow.ExecuteActivity(stateCtx, "UpdateWorkflowState", workflowState).Get(ctx, nil)

		return result, err
	}

	logger.Info("Retrieved rankings needing report processing", "count", len(rankingsToProcess))

	if len(rankingsToProcess) == 0 {
		logger.Info("No rankings need report processing")

		// Complete workflow state with no work done
		workflowState := &warcraftlogsBuilds.WorkflowState{
			ID:             workflowStateID,
			Status:         "completed",
			CompletedAt:    workflow.Now(ctx),
			ItemsProcessed: 0,
			UpdatedAt:      workflow.Now(ctx),
		}
		_ = workflow.ExecuteActivity(stateCtx, "UpdateWorkflowState", workflowState).Get(ctx, nil)

		result.CompletedAt = workflow.Now(ctx)
		return result, nil
	}

	// Process reports in batches
	const batchSize = 10
	totalReportsProcessed := 0
	totalRankingsProcessed := len(rankingsToProcess)
	apiRequestsCount := 0
	failedReports := 0

	for i := 0; i < len(rankingsToProcess); i += batchSize {
		end := i + batchSize
		if end > len(rankingsToProcess) {
			end = len(rankingsToProcess)
		}

		batch := rankingsToProcess[i:end]
		logger.Info("Processing rankings batch",
			"batchSize", len(batch),
			"progress", fmt.Sprintf("%d/%d", i+len(batch), len(rankingsToProcess)))

		// Update workflow state with current batch
		workflowState := &warcraftlogsBuilds.WorkflowState{
			ID:              workflowStateID,
			LastProcessedID: fmt.Sprintf("batch-%d", i/batchSize),
			ItemsProcessed:  i,
			UpdatedAt:       workflow.Now(ctx),
		}
		_ = workflow.ExecuteActivity(stateCtx, "UpdateWorkflowState", workflowState).Get(ctx, nil)

		// Process the batch through the API to fetch and store reports
		var batchResult models.BatchResult
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
				_ = workflow.ExecuteActivity(stateCtx, "UpdateWorkflowState", workflowState).Get(ctx, nil)

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
		var reportCodes []string
		for _, report := range batchResult.ProcessedReports {
			reportCodes = append(reportCodes, report.Code)
		}

		if len(reportCodes) > 0 {
			// Mark reports for build processing
			err = workflow.ExecuteActivity(activityCtx,
				definitions.MarkReportsForBuildProcessingActivity,
				reportCodes, params.BatchID).Get(ctx, nil)

			if err != nil {
				logger.Error("Failed to mark reports for build processing", "error", err)
				// Continue even if marking fails
			}
		}

		totalReportsProcessed += len(batchResult.ProcessedReports)
		apiRequestsCount += len(batch) * 2 // Approximately 2 API calls per ranking

		logger.Info("Completed processing reports batch",
			"batchProcessed", len(batchResult.ProcessedReports),
			"totalProcessed", totalReportsProcessed)

		// Small delay between batches
		workflow.Sleep(ctx, time.Second*2)
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
	_ = workflow.ExecuteActivity(stateCtx, "UpdateWorkflowState", workflowState).Get(ctx, nil)

	logger.Info("Reports workflow completed",
		"rankingsProcessed", totalRankingsProcessed,
		"reportsProcessed", totalReportsProcessed,
		"failedReports", failedReports,
		"apiRequests", apiRequestsCount,
		"duration", result.CompletedAt.Sub(result.StartedAt))

	return result, nil
}
