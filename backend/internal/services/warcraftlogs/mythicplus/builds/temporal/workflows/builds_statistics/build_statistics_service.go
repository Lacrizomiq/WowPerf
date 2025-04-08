// build_statistics_service.go
package warcraftlogsBuildsTemporalWorkflowsBuildsstatistics

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// ProcessBuildStatistics processes the equipment statistics
func ProcessBuildStatistics(
	ctx workflow.Context,
	config models.AnalysisWorkflowConfig,
	state *state.AnalysisWorkflowState,
) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting build statistics processing")

	// Activity configuration options
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 24,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    int32(config.RetryAttempts),
		},
		HeartbeatTimeout: time.Minute * 10,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Process each class/spec/dungeon combination
	for _, spec := range config.Specs {
		for _, dungeon := range config.Dungeons {
			// Identify the combination
			combinationKey := fmt.Sprintf("%s_%s_%d", spec.ClassName, spec.SpecName, dungeon.EncounterID)

			// Check if already processed
			if state.ProcessedCombinations[combinationKey] {
				logger.Info("Skipping already processed combination",
					"class", spec.ClassName,
					"spec", spec.SpecName,
					"dungeon", dungeon.Name)
				continue
			}

			// Update the current state
			state.CurrentSpec = &spec
			state.CurrentDungeon = &dungeon

			logger.Info("Processing build statistics",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name)

			// Execute the activity
			var activityResult models.BuildsAnalysisResult
			err := workflow.ExecuteActivity(
				ctx,
				definitions.ProcessBuildStatisticsActivity,
				spec.ClassName,
				spec.SpecName,
				uint(dungeon.EncounterID),
				config.BatchSize,
			).Get(ctx, &activityResult)

			if err != nil {
				if common.IsRateLimitError(err) {
					return err // Propagation of the rate limit error
				}
				return fmt.Errorf("failed to process build statistics: %w", err)
			}

			// Update the state
			state.ProcessedCombinations[combinationKey] = true
			state.Results.ItemsAnalyzed += activityResult.ItemsProcessed
			state.Results.TotalBuilds += activityResult.BuildsProcessed

			logger.Info("Successfully processed build statistics",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name,
				"buildsProcessed", activityResult.BuildsProcessed,
				"itemsProcessed", activityResult.ItemsProcessed)

		}

		// Increment the specs counter
		state.Results.SpecsProcessed++
	}

	state.Results.DungeonsProcessed = int32(len(config.Dungeons))

	logger.Info("Completed all build statistics processing")
	return nil
}
