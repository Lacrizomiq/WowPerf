package warcraftlogsBuildsTemporalWorkflowsRankings

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// ProcessAllRankings processes all rankings for the given specs and dungeons
// It updates the state as it processes each spec/dungeon combination
func ProcessAllRankings(
	ctx workflow.Context,
	params models.WorkflowConfig,
	state *state.WorkflowState,
) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting rankings processing")

	specsToProcess := params.Specs

	if len(specsToProcess) == 0 {
		return &common.WorkflowError{
			Type:      common.ErrorTypeConfiguration,
			Message:   "no specs found in configuration",
			Retryable: false,
		}
	}

	logger.Info("Processing specs", "count", len(specsToProcess))

	// Process each spec
	for _, spec := range specsToProcess {
		specKey := common.GenerateSpecKey(spec)
		if state.ProcessedSpecs[specKey] {
			logger.Info("Skipping already processed spec", "class", spec.ClassName, "spec", spec.SpecName)
			continue
		}

		logger.Info("Processing spec", "class", spec.ClassName, "spec", spec.SpecName)

		// Process rankings for each dungeon
		for _, dungeon := range params.Dungeons {
			dungeonKey := common.GenerateDungeonKey(spec, dungeon)
			if state.ProcessedDungeons[dungeonKey] {
				logger.Info("Skipping already processed dungeon", "dungeon", dungeon.Name, "for spec", spec.SpecName)
				continue
			}

			logger.Info("Processing dungeon", "dungeon", dungeon.Name, "for spec", spec.SpecName)

			// Configure activity options
			activityOpts := workflow.ActivityOptions{
				StartToCloseTimeout: time.Hour * 24,
				HeartbeatTimeout:    time.Minute * 10,
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:    time.Second * 5,
					BackoffCoefficient: 2.0,
					MaximumInterval:    time.Minute * 10,
					MaximumAttempts:    3,
				},
			}
			activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

			// Execute the rankings activity directly
			var batchResult models.BatchResult
			err := workflow.ExecuteActivity(activityCtx,
				definitions.FetchRankingsActivity,
				spec, dungeon, params.Rankings.Batch).Get(ctx, &batchResult)

			if err != nil {
				if common.IsRateLimitError(err) {
					// Rate limit reached, bail out and let main workflow handle continuation
					logger.Info("Rate limit reached during rankings processing")
					return err
				}
				logger.Error("Failed to process rankings",
					"class", spec.ClassName,
					"spec", spec.SpecName,
					"dungeon", dungeon.Name,
					"error", err)
				continue
			}

			// Update state with results
			state.CurrentSpec = &spec
			state.CurrentDungeon = &dungeon
			state.ProcessedDungeons[dungeonKey] = true
			state.PartialResults.RankingsProcessed += batchResult.ProcessedItems

			logger.Info("Successfully processed rankings",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name,
				"itemsProcessed", batchResult.ProcessedItems,
				"totalProcessed", state.PartialResults.RankingsProcessed)

			// Small delay between dungeons to avoid overwhelming the system
			workflow.Sleep(ctx, time.Second*2)
		}

		// Mark spec as processed
		state.ProcessedSpecs[specKey] = true
		logger.Info("Completed processing for spec", "class", spec.ClassName, "spec", spec.SpecName)
	}

	logger.Info("Rankings processing completed", "totalProcessed", state.PartialResults.RankingsProcessed)
	return nil
}
