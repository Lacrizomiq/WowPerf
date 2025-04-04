package warcraftlogsBuildsTemporalWorkflowsReports

import (
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

	// Process each spec
	for _, spec := range params.Specs {
		// Process each dungeon for this spec
		for _, dungeon := range params.Dungeons {
			logger.Info("Processing reports for spec and dungeon",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name)

			// Get rankings to process
			var rankings []*warcraftlogsBuilds.ClassRanking
			err := workflow.ExecuteActivity(activityCtx,
				definitions.GetStoredRankingsActivity,
				spec.ClassName,
				spec.SpecName,
				dungeon.EncounterID).Get(ctx, &rankings)

			if err != nil {
				logger.Error("Failed to get stored rankings",
					"spec", spec.SpecName,
					"dungeon", dungeon.Name,
					"error", err)
				continue
			}

			if len(rankings) == 0 {
				logger.Info("No rankings found to process",
					"spec", spec.SpecName,
					"dungeon", dungeon.Name)
				continue
			}

			// Update state information for potential resume
			state.CurrentSpec = &spec
			state.CurrentDungeon = &dungeon

			logger.Info("Processing reports for rankings",
				"count", len(rankings),
				"spec", spec.SpecName,
				"dungeon", dungeon.Name)

			// Process reports based on the rankings
			var batchResult models.BatchResult
			err = workflow.ExecuteActivity(activityCtx,
				definitions.ProcessReportsActivity,
				rankings).Get(ctx, &batchResult)

			if err != nil {
				if common.IsRateLimitError(err) {
					// No need to update state as we already did above
					return err
				}
				logger.Error("Failed to process reports",
					"spec", spec.SpecName,
					"dungeon", dungeon.Name,
					"error", err)
				continue
			}

			// Update state with results
			state.PartialResults.ReportsProcessed += batchResult.ProcessedItems

			logger.Info("Completed processing reports for spec/dungeon",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name,
				"processedCount", batchResult.ProcessedItems,
				"totalProcessed", state.PartialResults.ReportsProcessed)
		}
	}

	logger.Info("All reports processed successfully",
		"totalProcessed", state.PartialResults.ReportsProcessed)
	return nil
}
