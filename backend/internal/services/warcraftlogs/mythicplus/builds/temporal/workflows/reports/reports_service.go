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
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// ProcessAllReports processes all reports for the rankings stored in the state
// It updates the state with the processing results
func ProcessAllReports(
	ctx workflow.Context,
	params models.WorkflowConfig,
	state *state.WorkflowState,
) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting reports processing")

	// Configure activity options
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 24,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    3,
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

	// Retrieve all unique report references
	var reportRefs []*warcraftlogsBuilds.ClassRanking
	err := workflow.ExecuteActivity(activityCtx,
		definitions.GetUniqueReportReferencesActivity).Get(ctx, &reportRefs)

	if err != nil {
		logger.Error("Failed to get unique report references", "error", err)
		return err
	}

	logger.Info("Retrieved unique report references", "count", len(reportRefs))

	if len(reportRefs) == 0 {
		logger.Info("No report references found to process")
		return nil
	}

	// Process reports in batches to avoid memory overload
	const batchSize = 10
	totalProcessed := 0

	for i := 0; i < len(reportRefs); i += batchSize {
		end := i + batchSize
		if end > len(reportRefs) {
			end = len(reportRefs)
		}

		batch := reportRefs[i:end]
		logger.Info("Processing report batch",
			"batchSize", len(batch),
			"progress", fmt.Sprintf("%d/%d", i+len(batch), len(reportRefs)))

		var batchResult models.BatchResult
		err := workflow.ExecuteActivity(activityCtx,
			definitions.ProcessReportsActivity,
			batch).Get(ctx, &batchResult)

		if err != nil {
			if common.IsRateLimitError(err) {
				// Save the current state to resume later
				return err
			}
			logger.Error("Failed to process report batch", "error", err)
			continue
		}

		totalProcessed += int(batchResult.ProcessedItems)
		state.PartialResults.ReportsProcessed += batchResult.ProcessedItems

		logger.Info("Completed processing report batch",
			"batchProcessed", batchResult.ProcessedItems,
			"totalProcessed", totalProcessed)
	}

	logger.Info("All reports processed successfully",
		"totalProcessed", totalProcessed)
	return nil
}
