package warcraftlogsBuildsTemporalWorkflowsReports

import (
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// ReportsWorkflow implements the definitions.ReportsWorkflow interface
type ReportsWorkflow struct {
	stateManager *state.Manager
	processor    *Processor
}

// NewReportsWorkflow creates a new reports workflow
func NewReportsWorkflow() definitions.ProcessBuildBatchWorkflow {
	return &ReportsWorkflow{
		stateManager: state.NewManager(),
		processor:    NewProcessor(),
	}
}

// Execute runs the reports workflow
func (w *ReportsWorkflow) Execute(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting reports workflow")

	// Load or initialize state
	if err := w.stateManager.LoadCheckpoint(ctx); err != nil {
		return nil, err
	}

	state := w.stateManager.GetState()

	// Get rankings to process using definitions constant
	var rankings []*warcraftlogsBuilds.ClassRanking
	err := workflow.ExecuteActivity(ctx,
		definitions.GetStoredRankingsActivity,
		state.CurrentSpec.ClassName,
		state.CurrentSpec.SpecName,
		state.CurrentDungeon.EncounterID,
	).Get(ctx, &rankings)

	if err != nil {
		logger.Error("Failed to get stored rankings", "error", err)
		return nil, err
	}

	if len(rankings) > 0 {
		logger.Info("Processing reports for rankings", "count", len(rankings))

		// Use definitions constant for processing reports
		var batchResult models.BatchResult
		err := workflow.ExecuteActivity(ctx,
			definitions.ProcessReportsActivity,
			rankings,
			params.Worker,
		).Get(ctx, &batchResult)

		if err != nil {
			if common.IsRateLimitError(err) {
				w.stateManager.SaveCheckpoint(ctx)
				return nil, workflow.NewContinueAsNewError(ctx, workflow.GetInfo(ctx).WorkflowExecution.ID, params)
			}
			logger.Error("Failed to process reports",
				"spec", state.CurrentSpec.SpecName,
				"error", err)
			return nil, err
		}

		// Update state with results
		state.PartialResults.ReportsProcessed += batchResult.ProcessedItems

		// Update progress
		w.stateManager.UpdateProgress(models.PhaseReports, state.PartialResults.ReportsProcessed)
	}

	return &models.WorkflowResult{
		ReportsProcessed: state.PartialResults.ReportsProcessed,
		StartedAt:        state.StartedAt,
		CompletedAt:      workflow.Now(ctx),
	}, nil
}
