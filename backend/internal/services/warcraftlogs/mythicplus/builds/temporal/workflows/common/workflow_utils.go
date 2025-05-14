package warcraftlogsBuildsTemporalWorkflowsCommon

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	warcraftlogsBuildMetrics "wowperf/internal/services/warcraftlogs/mythicplus/builds/metrics"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// SyncStateBeforeContinueAsNew synchronizes the workflow state in the database
// before calling ContinueAsNew
func SyncStateBeforeContinueAsNew(
	ctx workflow.Context,
	workflowStateID string,
	workflowType string,
	metrics *warcraftlogsBuildMetrics.MetricsCollector,
	status string,
) error {
	logger := workflow.GetLogger(ctx)

	// Preparation of the workflowState to update
	workflowState := &warcraftlogsBuilds.WorkflowState{
		ID:        workflowStateID,
		Status:    status,
		UpdatedAt: workflow.Now(ctx),
	}

	// Update metrics
	metrics.RecordContinuation()

	// Add metrics to the workflowState
	if err := metrics.UpdateWorkflowState(workflowState); err != nil {
		logger.Error("Failed to update metrics in workflow state", "error", err)
		return err
	}

	// Update the database via the activity
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.5,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

	err := workflow.ExecuteActivity(activityCtx,
		definitions.UpdateWorkflowStateActivity,
		workflowState).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to update workflow state before ContinueAsNew", "error", err)
		// Continue even if there is an error because the Continue-As-New must be made
		return err
	}

	logger.Info("Workflow state synchronized before ContinueAsNew",
		"workflowStateID", workflowStateID,
		"continuationCount", metrics.GetSnapshot().ContinuationCount,
		"itemsProcessed", metrics.GetSnapshot().ItemsProcessed)

	return nil
}

// SyncReportsWorkflowParams synchronizes the Reports workflow parameters with the metrics
// to prepare the ContinueAsNew
func SyncReportsWorkflowParams(
	params models.ReportsWorkflowParams,
	metrics *warcraftlogsBuildMetrics.MetricsCollector,
	localRankingsProcessed int,
	localReportsProcessed int,
	localAPIRequestsCount int,
	localFailedReports int,
) models.ReportsWorkflowParams {
	// Update the parameters for the ContinueAsNew
	updatedParams := params
	updatedParams.TotalProcessedRankings += int32(localRankingsProcessed)
	updatedParams.TotalProcessedReports += int32(localReportsProcessed)
	updatedParams.TotalAPIRequests += int32(localAPIRequestsCount)
	updatedParams.TotalFailedReports += int32(localFailedReports)
	updatedParams.ContinuationCount++

	// Add fields if necessary for metrics tracking
	metricsSnapshot := metrics.GetSnapshot()
	if metricsSnapshot.ReportsMetrics != nil {
		//  Can add other metrics here if necessary
	}

	return updatedParams
}
