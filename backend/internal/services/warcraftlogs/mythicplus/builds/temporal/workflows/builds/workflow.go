package warcraftlogsBuildsTemporalWorkflowsBuilds

import (
	"go.temporal.io/sdk/workflow"

	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// BuildsWorkflow implements the definitions.ProcessBuildBatchWorkflow interface
type BuildsWorkflow struct {
	stateManager *state.Manager
	processor    *Processor
}

// NewBuildsWorkflow creates a new builds workflow
func NewBuildsWorkflow() definitions.ProcessBuildBatchWorkflow {
	return &BuildsWorkflow{
		stateManager: state.NewManager(),
		processor:    NewProcessor(),
	}
}

// Execute runs the builds workflow
func (w *BuildsWorkflow) Execute(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting builds workflow")

	// Load or initialize state
	if err := w.stateManager.LoadCheckpoint(ctx); err != nil {
		return nil, err
	}

	state := w.stateManager.GetState()

	// Get reports to process
	reports, err := w.processor.GetStoredReports(ctx)
	if err != nil {
		logger.Error("Failed to get stored reports", "error", err)
		return nil, err
	}

	if len(reports) > 0 {
		logger.Info("Processing builds from reports", "count", len(reports))

		// Process builds in batches
		batchResult, err := w.processor.ProcessBuilds(ctx, reports, params.Worker)
		if err != nil {
			if common.IsRateLimitError(err) {
				w.stateManager.SaveCheckpoint(ctx)
				return nil, workflow.NewContinueAsNewError(ctx, workflow.GetInfo(ctx).WorkflowExecution.ID, params)
			}
			logger.Error("Failed to process builds", "error", err)
			return nil, err
		}

		// Update state with results
		state.PartialResults.BuildsProcessed += batchResult.ProcessedItems

		// Update progress
		w.stateManager.UpdateProgress(models.PhaseBuilds, state.PartialResults.BuildsProcessed)
	}

	return &models.WorkflowResult{
		BuildsProcessed: state.PartialResults.BuildsProcessed,
		StartedAt:       state.StartedAt,
		CompletedAt:     workflow.Now(ctx),
	}, nil
}
