package warcraftlogsBuildsTemporalWorkflowsRankings

import (
	"time"

	"go.temporal.io/sdk/workflow"

	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// RankingsWorkflow implements the definitions.RankingsWorkflow interface
type RankingsWorkflow struct {
	stateManager *state.Manager
	processor    *Processor
}

// NewRankingsWorkflow creates a new rankings workflow
func NewRankingsWorkflow() definitions.RankingsWorkflow {
	return &RankingsWorkflow{
		stateManager: state.NewManager(),
		processor:    NewProcessor(),
	}
}

// Execute runs the rankings workflow
func (w *RankingsWorkflow) Execute(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting rankings workflow")

	// Load or initialize state
	if err := w.stateManager.LoadCheckpoint(ctx); err != nil {
		return nil, err
	}

	// Get class from workflow ID
	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	className := common.ExtractClassFromWorkflowID(workflowID)
	if className == "" {
		return nil, &common.WorkflowError{
			Type:      common.ErrorTypeConfiguration,
			Message:   "invalid workflow ID format",
			Retryable: false,
		}
	}

	// Filter specs for this class
	specs := common.FilterSpecsForClass(params.Specs, className)
	if len(specs) == 0 {
		return nil, &common.WorkflowError{
			Type:      common.ErrorTypeConfiguration,
			Message:   "no specs found for class",
			Retryable: false,
		}
	}

	state := w.stateManager.GetState()

	// Process each spec
	for _, spec := range specs {
		specKey := common.GenerateSpecKey(spec)
		if state.ProcessedSpecs[specKey] {
			continue
		}

		// Process rankings for each dungeon
		for _, dungeon := range params.Dungeons {
			dungeonKey := common.GenerateDungeonKey(spec, dungeon)
			if state.ProcessedDungeons[dungeonKey] {
				continue
			}

			// Use activity name from definitions and properly handle the result
			var batchResult models.BatchResult
			err := workflow.ExecuteActivity(ctx,
				definitions.FetchRankingsActivity,
				spec, dungeon, params.Rankings.Batch,
			).Get(ctx, &batchResult)

			if err != nil {
				if common.IsRateLimitError(err) {
					// Save state and continue as new workflow
					w.stateManager.SaveCheckpoint(ctx)
					return nil, workflow.NewContinueAsNewError(ctx, workflowID, params)
				}
				logger.Error("Failed to process rankings",
					"class", spec.ClassName,
					"spec", spec.SpecName,
					"dungeon", dungeon.Name,
					"error", err)
				continue
			}

			// Update state
			state.CurrentSpec = &spec
			state.CurrentDungeon = &dungeon
			state.ProcessedDungeons[dungeonKey] = true
			state.PartialResults.RankingsProcessed += batchResult.ProcessedItems

			// Update progress
			w.stateManager.UpdateProgress(models.PhaseRankings, state.PartialResults.RankingsProcessed)

			// Small delay between dungeons
			workflow.Sleep(ctx, time.Second*2)
		}

		// Mark spec as processed
		state.ProcessedSpecs[specKey] = true
	}

	return &models.WorkflowResult{
		RankingsProcessed: state.PartialResults.RankingsProcessed,
		StartedAt:         state.StartedAt,
		CompletedAt:       workflow.Now(ctx),
	}, nil
}
