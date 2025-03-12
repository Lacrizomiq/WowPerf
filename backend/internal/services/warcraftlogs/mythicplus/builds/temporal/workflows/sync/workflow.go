package warcraftlogsBuildsTemporalWorkflowsSync

import (
	"go.temporal.io/sdk/workflow"

	builds "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/builds"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	rankings "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/rankings"
	reports "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/reports"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// SyncWorkflow implements the definitions.SyncWorkflow interface
type SyncWorkflow struct {
	stateManager *state.Manager
	orchestrator *Orchestrator
}

// NewSyncWorkflow creates a new instance of the sync workflow
func NewSyncWorkflow() definitions.SyncWorkflow {
	return &SyncWorkflow{
		stateManager: state.NewManager(),
		orchestrator: NewOrchestrator(),
	}
}

// Execute runs the main synchronization workflow
func (w *SyncWorkflow) Execute(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting sync workflow orchestration")

	if err := w.stateManager.LoadCheckpoint(ctx); err != nil {
		return nil, err
	}

	state := w.stateManager.GetState()

	// Process Rankings phase using constant from definitions
	if !w.orchestrator.IsPhaseCompleted(models.PhaseRankings) {
		// Ensure we're using the interface from definitions
		var rankingsWorkflow definitions.RankingsWorkflow
		rankingsWorkflow = rankings.NewRankingsWorkflow()

		rankingsResult, err := rankingsWorkflow.Execute(ctx, params)
		if err != nil {
			return handleWorkflowError(ctx, err, state)
		}
		state.PartialResults.RankingsProcessed = rankingsResult.RankingsProcessed
		w.orchestrator.MarkPhaseCompleted(models.PhaseRankings)
		w.stateManager.UpdateProgress(models.PhaseRankings, rankingsResult.RankingsProcessed)
	}

	// Process Reports phase
	if !w.orchestrator.IsPhaseCompleted(models.PhaseReports) {
		var reportsWorkflow definitions.ProcessBuildBatchWorkflow
		reportsWorkflow = reports.NewReportsWorkflow()

		reportsResult, err := reportsWorkflow.Execute(ctx, params)
		if err != nil {
			return handleWorkflowError(ctx, err, state)
		}
		state.PartialResults.ReportsProcessed = reportsResult.ReportsProcessed
		w.orchestrator.MarkPhaseCompleted(models.PhaseReports)
		w.stateManager.UpdateProgress(models.PhaseReports, reportsResult.ReportsProcessed)
	}

	// Process Builds phase
	if !w.orchestrator.IsPhaseCompleted(models.PhaseBuilds) {
		var buildsWorkflow definitions.ProcessBuildBatchWorkflow
		buildsWorkflow = builds.NewBuildsWorkflow()

		buildsResult, err := buildsWorkflow.Execute(ctx, params)
		if err != nil {
			return handleWorkflowError(ctx, err, state)
		}
		state.PartialResults.BuildsProcessed = buildsResult.BuildsProcessed
		w.orchestrator.MarkPhaseCompleted(models.PhaseBuilds)
		w.stateManager.UpdateProgress(models.PhaseBuilds, buildsResult.BuildsProcessed)
	}

	return &models.WorkflowResult{
		RankingsProcessed: state.PartialResults.RankingsProcessed,
		ReportsProcessed:  state.PartialResults.ReportsProcessed,
		BuildsProcessed:   state.PartialResults.BuildsProcessed,
		StartedAt:         state.StartedAt,
		CompletedAt:       workflow.Now(ctx),
	}, nil
}
