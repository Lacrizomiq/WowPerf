// stat_statistics_service.go
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

// ProcessStatStatistics processes the character stat statistics
func ProcessStatStatistics(
	ctx workflow.Context,
	config models.AnalysisWorkflowConfig,
	state *state.AnalysisWorkflowState,
) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting character stat statistics processing")

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
			combinationKey := fmt.Sprintf("stat_%s_%s_%d", spec.ClassName, spec.SpecName, dungeon.EncounterID)

			// Check if already processed
			if state.ProcessedCombinations[combinationKey] {
				logger.Info("Skipping already processed stat combination",
					"class", spec.ClassName,
					"spec", spec.SpecName,
					"dungeon", dungeon.Name)
				continue
			}

			// Update the current state
			state.CurrentSpec = &spec
			state.CurrentDungeon = &dungeon

			logger.Info("Processing character stat statistics",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name)

			// Execute the activity
			var activityResult models.BuildsAnalysisResult
			err := workflow.ExecuteActivity(
				ctx,
				definitions.ProcessStatStatisticsActivity,
				spec.ClassName,
				spec.SpecName,
				uint(dungeon.EncounterID),
				config.BatchSize,
			).Get(ctx, &activityResult)

			if err != nil {
				if common.IsRateLimitError(err) {
					return err
				}
				return fmt.Errorf("failed to process stat statistics: %w", err)
			}

			// Update the state
			state.ProcessedCombinations[combinationKey] = true
			state.Results.StatsAnalyzed += activityResult.ItemsProcessed
			state.Results.TotalBuilds += activityResult.BuildsProcessed

			logger.Info("Successfully processed stat statistics",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name,
				"buildsProcessed", activityResult.BuildsProcessed,
				"statsProcessed", activityResult.ItemsProcessed)

			// Save a checkpoint after each combination
			if saveErr := workflow.ExecuteLocalActivity(ctx, saveCheckpoint, state).Get(ctx, nil); saveErr != nil {
				logger.Warn("Failed to save checkpoint", "error", saveErr)
			}
		}
	}

	logger.Info("Completed all character stat statistics processing")
	return nil
}
