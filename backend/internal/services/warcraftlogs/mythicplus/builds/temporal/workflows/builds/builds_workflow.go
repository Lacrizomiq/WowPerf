package warcraftlogsBuildsTemporalWorkflowsBuilds

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	warcraftlogsBuildMetrics "wowperf/internal/services/warcraftlogs/mythicplus/builds/metrics"
	warcraftlogsBuildsTemporalWorkflowsCommon "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// BuildsWorkflow implements the builds extraction parent workflow
type BuildsWorkflow struct{}

// NewBuildsWorkflow creates a new instance of the builds workflow
func NewBuildsWorkflow() definitions.BuildsWorkflow {
	return &BuildsWorkflow{}
}

// Execute runs the builds workflow
func (w *BuildsWorkflow) Execute(ctx workflow.Context, params models.BuildsWorkflowParams) (*models.BuildsWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)

	// Initialize metrics collector
	metrics := warcraftlogsBuildMetrics.NewMetricsCollector(
		"builds_workflow",
		"",
		params.BatchID,
	)

	logger.Info("Starting builds parent workflow",
		"batchID", params.BatchID,
		"batchSize", params.BatchSize,
		"numWorkers", params.NumWorkers,
		"reportBatchSize", params.ReportBatchSize,
		"offset", params.Offset,
		"continuationCount", params.ContinuationCount)

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
	workflowID := workflowInfo.WorkflowExecution.ID
	runID := workflowInfo.WorkflowExecution.RunID
	workflowStateID := fmt.Sprintf("builds-%s-%s", workflowID, runID)

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

	// Create initial workflow state or update if continuation
	var initialState *warcraftlogsBuilds.WorkflowState

	if params.ContinuationCount == 0 {
		// New execution - create a new state
		initialState = &warcraftlogsBuilds.WorkflowState{
			ID:                workflowStateID,
			WorkflowType:      "builds",
			BatchID:           params.BatchID,
			StartedAt:         finalResult.StartedAt,
			Status:            "running",
			CreatedAt:         finalResult.StartedAt,
			UpdatedAt:         finalResult.StartedAt,
			ContinuationCount: 0,
		}

		// Set initial metrics on workflow state
		metrics.UpdateWorkflowState(initialState)

		err := workflow.ExecuteActivity(stateCtx, definitions.CreateWorkflowStateActivity, initialState).Get(ctx, nil)
		if err != nil {
			logger.Error("Failed to create workflow state", "error", err)
			metrics.RecordError("workflow_state_creation")
		}
	} else {
		// Continuation - update existing state
		stateUpdate := &warcraftlogsBuilds.WorkflowState{
			ID:                workflowStateID,
			UpdatedAt:         workflow.Now(ctx),
			ContinuationCount: int(params.ContinuationCount),
			ItemsProcessed:    int(params.AlreadyProcessed),
		}

		// Update metrics
		metrics.UpdateWorkflowState(stateUpdate)
		metrics.RecordItemsProcessed(int(params.AlreadyProcessed), "already_processed")

		err := workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, stateUpdate).Get(ctx, nil)
		if err != nil {
			logger.Error("Failed to update workflow state after continuation", "error", err)
		}
	}

	// Configuration of Main Activities
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

	// Determine total reports to process if this is the first execution
	var totalToProcess int32 = params.TotalToProcess
	if params.ContinuationCount == 0 {
		// Count reports to process
		metrics.StartOperation("count_reports")
		var countResult int64
		err := workflow.ExecuteActivity(activityCtx,
			definitions.CountReportsNeedingBuildExtractionActivity,
			10*24*time.Hour,
		).Get(ctx, &countResult)
		metrics.EndOperation("count_reports")

		if err != nil {
			logger.Error("Failed to count reports needing build extraction", "error", err)
			metrics.RecordError("count_reports_activity")

			// Update workflow state with error
			workflowState := &warcraftlogsBuilds.WorkflowState{
				ID:           workflowStateID,
				Status:       "failed",
				ErrorMessage: fmt.Sprintf("Failed to count reports: %v", err),
				UpdatedAt:    workflow.Now(ctx),
			}
			metrics.UpdateWorkflowState(workflowState)
			_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

			finalResult.CompletedAt = workflow.Now(ctx)
			return finalResult, err
		}

		if countResult == 0 {
			logger.Info("No reports need build extraction.")

			// Update workflow state to completed
			workflowState := &warcraftlogsBuilds.WorkflowState{
				ID:          workflowStateID,
				Status:      "completed",
				CompletedAt: workflow.Now(ctx),
				UpdatedAt:   workflow.Now(ctx),
			}
			metrics.UpdateWorkflowState(workflowState)
			_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

			finalResult.CompletedAt = workflow.Now(ctx)
			metrics.Finish("completed_empty")
			return finalResult, nil
		}

		totalToProcess = int32(countResult)
		logger.Info("Found reports to process", "totalCount", totalToProcess)

		// Update state with total count
		workflowState := &warcraftlogsBuilds.WorkflowState{
			ID:                  workflowStateID,
			TotalItemsToProcess: int(totalToProcess),
			UpdatedAt:           workflow.Now(ctx),
		}
		metrics.UpdateWorkflowState(workflowState)
		_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)
	}

	// Determine page size
	pageSize := params.PageSize // Reasonable default
	if params.PageSize > 0 {
		pageSize = params.PageSize
	}

	// Get a page of reports to process
	metrics.StartOperation("get_reports_batch")
	var reportsToProcess []*warcraftlogsBuilds.Report
	err := workflow.ExecuteActivity(activityCtx,
		definitions.GetReportsNeedingBuildExtractionActivity,
		pageSize, params.Offset, 10*24*time.Hour,
	).Get(ctx, &reportsToProcess)
	metrics.EndOperation("get_reports_batch")

	if err != nil {
		logger.Error("Failed to get reports needing build extraction", "error", err)
		metrics.RecordError("get_reports_activity")

		// Update workflow state with error
		workflowState := &warcraftlogsBuilds.WorkflowState{
			ID:           workflowStateID,
			Status:       "failed",
			ErrorMessage: fmt.Sprintf("Failed to get reports: %v", err),
			UpdatedAt:    workflow.Now(ctx),
		}
		metrics.UpdateWorkflowState(workflowState)
		_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

		finalResult.CompletedAt = workflow.Now(ctx)
		return finalResult, err
	}

	if len(reportsToProcess) == 0 {
		logger.Info("No more reports to process in this page.")

		// If we've already processed some reports, consider it completed
		if params.AlreadyProcessed > 0 {
			logger.Info("All reports have been processed across all pages.")

			// Update workflow state to completed
			workflowState := &warcraftlogsBuilds.WorkflowState{
				ID:                 workflowStateID,
				Status:             "completed",
				CompletedAt:        workflow.Now(ctx),
				UpdatedAt:          workflow.Now(ctx),
				ItemsProcessed:     int(params.AlreadyProcessed),
				ProgressPercentage: 100.0,
			}
			metrics.UpdateWorkflowState(workflowState)
			_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

			finalResult.CompletedAt = workflow.Now(ctx)
			metrics.Finish("completed")
			return finalResult, nil
		}
	}

	logger.Info("Retrieved reports for build extraction",
		"count", len(reportsToProcess),
		"offset", params.Offset)

	// Set total to process for metrics
	metrics.SetTotalToProcess(int(totalToProcess))

	// Setup sliding window for child workflows
	maxConcurrentChildren := 5 // Default value
	if params.NumWorkers > 0 {
		maxConcurrentChildren = int(params.NumWorkers)
	}

	// Define batch size for child workflows
	batchSize := 10 // Default value
	if params.ReportBatchSize > 0 {
		batchSize = int(params.ReportBatchSize)
	}

	// Calculate total number of batches for this page
	totalBatches := (len(reportsToProcess) + batchSize - 1) / batchSize

	logger.Info("Preparing child workflows",
		"totalReports", len(reportsToProcess),
		"batchSize", batchSize,
		"totalBatches", totalBatches,
		"maxConcurrentChildren", maxConcurrentChildren)

	// Create slices to track progress and results
	results := make([]*models.BuildsBatchResult, totalBatches)
	activeFutures := make(map[int]workflow.Future)
	nextBatchIndex := 0
	completedBatches := 0

	// Child workflow options with defaults
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		TaskQueue: func() string {
			if params.TaskQueue == "" {
				return models.DefaultTaskQueue // Use the default value defined in models/config.go
			}
			return params.TaskQueue
		}(),
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    3,
		},
		WorkflowExecutionTimeout: time.Hour * 4, // Timeout for child workflows
	}

	// Process batches with sliding window
	for completedBatches < totalBatches {
		// Launch new child workflows up to max concurrent limit
		for len(activeFutures) < maxConcurrentChildren && nextBatchIndex < totalBatches {
			startIdx := nextBatchIndex * batchSize
			endIdx := min(startIdx+batchSize, len(reportsToProcess))
			batchReports := reportsToProcess[startIdx:endIdx]

			// Create child workflow params
			batchParams := models.BuildsBatchParams{
				Reports:               batchReports,
				BatchID:               fmt.Sprintf("%s-batch-%d-page-%d", params.BatchID, nextBatchIndex, params.ContinuationCount),
				ParentWorkflowStateID: workflowStateID,
			}

			// Generate unique child workflow ID
			childWorkflowID := fmt.Sprintf("builds-batch-%d-page-%d-%s", nextBatchIndex, params.ContinuationCount, uuid.New().String())

			// Configure child workflow options
			childOptions := childWorkflowOptions
			childOptions.WorkflowID = childWorkflowID
			childCtx := workflow.WithChildOptions(ctx, childOptions)

			// Launch child workflow
			metrics.StartOperation("launch_child_workflow")
			logger.Info("Launching child workflow", "batchIndex", nextBatchIndex, "childID", childWorkflowID)

			future := workflow.ExecuteChildWorkflow(childCtx,
				definitions.ProcessBuildsBatchWorkflow,
				batchParams)

			activeFutures[nextBatchIndex] = future
			metrics.RecordChildWorkflow("builds_batch")
			metrics.EndOperation("launch_child_workflow")

			nextBatchIndex++

			// Small delay to avoid overwhelming Temporal
			workflow.Sleep(ctx, time.Millisecond*100)
		}

		// Check if we have active futures to monitor
		if len(activeFutures) > 0 {
			// Create selector to wait for any child workflow to complete
			selector := workflow.NewSelector(ctx)

			// Add each future to the selector
			completedIndexes := make([]int, 0)

			for idx, future := range activeFutures {
				idx := idx // Capture variable for closure
				future := future

				selector.AddFuture(future, func(f workflow.Future) {
					completedIndexes = append(completedIndexes, idx)
				})
			}

			// Wait for at least one child workflow to complete
			selector.Select(ctx)

			// Process completed child workflows
			for _, idx := range completedIndexes {
				future := activeFutures[idx]
				var result models.BuildsBatchResult
				err := future.Get(ctx, &result)

				if err != nil {
					logger.Error("Child workflow failed", "batchIndex", idx, "error", err)
					metrics.RecordError("child_workflow_error")
					metrics.RecordChildWorkflowCompletion(false)

					// Set a minimal result with error info
					results[idx] = &models.BuildsBatchResult{
						Status: "failed",
						Error:  err.Error(),
					}
				} else {
					logger.Info("Child workflow completed successfully",
						"batchIndex", idx,
						"status", result.Status,
						"reportsProcessed", result.ReportsProcessed,
						"buildsProcessed", result.BuildsProcessed)

					metrics.RecordChildWorkflowCompletion(true)
					metrics.RecordItemsProcessed(int(result.ReportsProcessed), "success")

					// Process builds by class/spec
					for classSpec, count := range result.BuildsByClassSpec {
						parts := strings.Split(classSpec, "-")
						if len(parts) == 2 {
							metrics.RecordBuildByClassSpec(parts[0], parts[1], int(count))
						}
					}

					results[idx] = &result
				}

				// Remove this future from active futures
				delete(activeFutures, idx)
				completedBatches++

				// Update workflow state with progress
				progressState := &warcraftlogsBuilds.WorkflowState{
					ID:              workflowStateID,
					LastProcessedID: fmt.Sprintf("batch-%d-page-%d", idx, params.ContinuationCount),
					ItemsProcessed:  int(params.AlreadyProcessed) + completedBatches*batchSize,
					UpdatedAt:       workflow.Now(ctx),
				}
				metrics.UpdateWorkflowState(progressState)
				_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, progressState).Get(ctx, nil)
			}

			// Small delay before next iteration
			workflow.Sleep(ctx, time.Millisecond*50)
		} else if nextBatchIndex >= totalBatches {
			// All batches have been launched but not yet completed
			// Just wait for a moment before checking again
			workflow.Sleep(ctx, time.Second)
		}
	}

	// All child workflows have completed, aggregate results
	logger.Info("All child workflows for this page have completed, aggregating results")

	// Aggregate results from all child workflows for this page
	pageProcessedReports := int32(0)
	pageProcessedBuilds := int32(0)
	pageClassSpecBuilds := make(map[string]int32)

	for _, batchResult := range results {
		if batchResult != nil {
			pageProcessedBuilds += batchResult.BuildsProcessed
			pageProcessedReports += batchResult.ReportsProcessed

			// Merge builds by class/spec maps
			for classSpec, count := range batchResult.BuildsByClassSpec {
				pageClassSpecBuilds[classSpec] += count
			}
		}
	}

	// Update cumulative totals
	totalProcessedReports := params.AlreadyProcessed + pageProcessedReports

	// Calculate new offset and decide if we need to continue with a new page
	newOffset := params.Offset + int32(len(reportsToProcess))
	newContinuationCount := params.ContinuationCount + 1

	logger.Info("Page processing complete",
		"pageProcessedReports", pageProcessedReports,
		"pageProcessedBuilds", pageProcessedBuilds,
		"totalProcessedReports", totalProcessedReports,
		"totalToProcess", totalToProcess,
		"newOffset", newOffset,
		"continuation", newContinuationCount)

	// If we've processed less than the total number of reports, continue with a new page
	if totalProcessedReports < totalToProcess {
		// Prepare parameters for ContinueAsNew
		continueParams := params
		continueParams.Offset = newOffset
		continueParams.AlreadyProcessed = totalProcessedReports
		continueParams.TotalToProcess = totalToProcess
		continueParams.ContinuationCount = newContinuationCount

		// Update workflow state before ContinueAsNew
		progressState := &warcraftlogsBuilds.WorkflowState{
			ID:                  workflowStateID,
			Status:              "continuing",
			LastProcessedID:     fmt.Sprintf("page-%d", params.ContinuationCount),
			ItemsProcessed:      int(totalProcessedReports),
			TotalItemsToProcess: int(totalToProcess),
			ProgressPercentage:  float64(totalProcessedReports) / float64(totalToProcess) * 100.0,
			ContinuationCount:   int(newContinuationCount),
			UpdatedAt:           workflow.Now(ctx),
		}
		metrics.UpdateWorkflowState(progressState)
		_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, progressState).Get(ctx, nil)

		logger.Info("Continue-As-New to process next page",
			"continuationCount", newContinuationCount,
			"newOffset", newOffset,
			"totalProcessed", totalProcessedReports,
			"totalToProcess", totalToProcess)

		// Synchronize state before continue-as-new
		err := warcraftlogsBuildsTemporalWorkflowsCommon.SyncStateBeforeContinueAsNew(
			ctx,
			workflowStateID,
			"builds",
			metrics,
			"continuing",
		)
		if err != nil {
			logger.Warn("Failed to sync state before continue-as-new", "error", err)
		}

		return nil, workflow.NewContinueAsNewError(ctx, definitions.BuildsWorkflowName, continueParams)
	}

	// If we've finished processing all reports, finalize the workflow
	finalResult.BuildsProcessed = pageProcessedBuilds
	finalResult.ReportsProcessed = pageProcessedReports
	finalResult.BuildsByClassSpec = pageClassSpecBuilds
	finalResult.CompletedAt = workflow.Now(ctx)

	// Update final workflow state
	finalState := &warcraftlogsBuilds.WorkflowState{
		ID:                  workflowStateID,
		Status:              "completed",
		CompletedAt:         finalResult.CompletedAt,
		BatchID:             params.BatchID,
		ItemsProcessed:      int(totalProcessedReports),
		TotalItemsToProcess: int(totalToProcess),
		ProgressPercentage:  100.0,
		UpdatedAt:           finalResult.CompletedAt,
	}

	metrics.UpdateWorkflowState(finalState)
	_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, finalState).Get(ctx, nil)

	logger.Info("Builds workflow completed",
		"status", "completed",
		"totalReportsProcessed", totalProcessedReports,
		"totalBuildsProcessed", pageProcessedBuilds,
		"duration", finalResult.CompletedAt.Sub(finalResult.StartedAt))

	metrics.Finish("completed")
	return finalResult, nil
}

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
