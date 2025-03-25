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

	// Validate state contains the necessary information
	if state.CurrentSpec == nil || state.CurrentDungeon == nil {
		return &common.WorkflowError{
			Type:      common.ErrorTypeConfiguration,
			Message:   "current spec and dungeon must be set before processing reports",
			Retryable: false,
		}
	}

	logger.Info("Processing reports for spec and dungeon",
		"class", state.CurrentSpec.ClassName,
		"spec", state.CurrentSpec.SpecName,
		"dungeon", state.CurrentDungeon.Name)

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

	// Get rankings to process (appel direct à l'activité)
	var rankings []*warcraftlogsBuilds.ClassRanking
	err := workflow.ExecuteActivity(activityCtx,
		definitions.GetStoredRankingsActivity,
		state.CurrentSpec.ClassName,
		state.CurrentSpec.SpecName,
		state.CurrentDungeon.EncounterID).Get(ctx, &rankings)

	if err != nil {
		logger.Error("Failed to get stored rankings",
			"spec", state.CurrentSpec.SpecName,
			"dungeon", state.CurrentDungeon.Name,
			"error", err)
		return err
	}

	if len(rankings) == 0 {
		logger.Info("No rankings found to process",
			"spec", state.CurrentSpec.SpecName,
			"dungeon", state.CurrentDungeon.Name)
		return nil
	}

	logger.Info("Processing reports for rankings",
		"count", len(rankings),
		"spec", state.CurrentSpec.SpecName,
		"dungeon", state.CurrentDungeon.Name)

	// Process reports based on the rankings (appel direct à l'activité)
	var batchResult models.BatchResult
	err = workflow.ExecuteActivity(activityCtx,
		definitions.ProcessReportsActivity,
		rankings).Get(ctx, &batchResult)

	if err != nil {
		if common.IsRateLimitError(err) {
			// Rate limit reached, bail out and let main workflow handle continuation
			logger.Info("Rate limit reached during reports processing")
			return err
		}
		logger.Error("Failed to process reports",
			"spec", state.CurrentSpec.SpecName,
			"dungeon", state.CurrentDungeon.Name,
			"error", err)
		return err
	}

	// Update state with results
	state.PartialResults.ReportsProcessed += batchResult.ProcessedItems

	logger.Info("Completed processing reports",
		"processedCount", batchResult.ProcessedItems,
		"totalProcessed", state.PartialResults.ReportsProcessed)

	return nil
}
