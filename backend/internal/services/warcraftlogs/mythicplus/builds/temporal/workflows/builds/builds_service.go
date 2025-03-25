package warcraftlogsBuildsTemporalWorkflowsBuilds

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

// ProcessAllBuilds processes all builds from reports
// It updates the state with the processing results
func ProcessAllBuilds(
	ctx workflow.Context,
	params models.WorkflowConfig,
	state *state.WorkflowState,
) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting builds processing")

	// Activity options configuration
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 24,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	// Count the number of reports to process (direct activity call)
	var totalCount int64
	err := workflow.ExecuteActivity(ctx,
		definitions.CountAllReportsActivity).Get(ctx, &totalCount)

	if err != nil {
		logger.Error("Failed to count stored reports", "error", err)
		return err
	}

	logger.Info("Total reports to process", "count", totalCount)

	// If no reports, exit
	if totalCount == 0 {
		logger.Info("No reports found to process builds")
		return nil
	}

	// Pagination configuration
	batchSize := int32(5) // batch size
	offset := int32(0)    // starting position

	// Process reports by batches
	for {
		// Get a batch of reports (direct activity call)
		var reports []*warcraftlogsBuilds.Report
		err := workflow.ExecuteActivity(ctx,
			definitions.GetReportsBatchActivity,
			batchSize, offset).Get(ctx, &reports)

		if err != nil {
			logger.Error("Failed to get stored reports", "error", err)
			return err
		}

		// If no more reports, exit the loop
		if len(reports) == 0 {
			logger.Info("No more reports found to process")
			break
		}

		logger.Info("Processing builds from reports batch",
			"count", len(reports),
			"batchNumber", offset/batchSize+1,
			"progress", fmt.Sprintf("%.1f%%", float64(offset)/float64(totalCount)*100))

		// Process the batch (direct activity call)
		var batchResult models.BatchResult
		err = workflow.ExecuteActivity(ctx,
			definitions.ProcessBuildsActivity,
			reports).Get(ctx, &batchResult)

		if err != nil {
			if common.IsRateLimitError(err) {
				logger.Info("Rate limit reached during builds processing")
				return err
			}
			logger.Error("Failed to process builds", "error", err)
			return err
		}

		// Update the state
		state.PartialResults.BuildsProcessed += batchResult.ProcessedItems

		logger.Info("Completed processing builds",
			"processedCount", batchResult.ProcessedItems,
			"totalProcessed", state.PartialResults.BuildsProcessed,
			"batch", offset/batchSize)

		// Pass to the next batch
		offset += batchSize
	}

	logger.Info("All builds processed successfully", "totalProcessed", state.PartialResults.BuildsProcessed)

	return nil
}
