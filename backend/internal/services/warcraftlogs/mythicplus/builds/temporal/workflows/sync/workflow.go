package warcraftlogsBuildsTemporalWorkflowsSync

import (
	"fmt"

	"go.temporal.io/sdk/workflow"

	builds "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/builds"
	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	rankings "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/rankings"
	reports "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/reports"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// SyncWorkflow implements the definitions.SyncWorkflow interface
// This is the main workflow that orchestrates the entire process
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
// It orchestrates the three main phases: Rankings, Reports, and Builds
func (w *SyncWorkflow) Execute(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting sync workflow orchestration",
		"workflowID", workflow.GetInfo(ctx).WorkflowExecution.ID,
		"attempt", workflow.GetInfo(ctx).Attempt)

	// Validate input parameters
	if len(params.Specs) == 0 {
		return nil, fmt.Errorf("no specs found in config")
	}

	if len(params.Dungeons) == 0 {
		return nil, fmt.Errorf("no dungeons found in config")
	}

	// Log configuration details for verification
	if len(params.Specs) > 0 {
		logger.Info("First spec",
			"className", params.Specs[0].ClassName,
			"specName", params.Specs[0].SpecName)
	}

	if len(params.Dungeons) > 0 {
		logger.Info("First dungeon",
			"name", params.Dungeons[0].Name,
			"ID", params.Dungeons[0].EncounterID)
	}

	// Load checkpoint state if this is a continuation
	if err := w.stateManager.LoadCheckpoint(ctx); err != nil {
		return nil, err
	}

	state := w.stateManager.GetState()

	// Process Rankings phase
	if !w.orchestrator.IsPhaseCompleted(models.PhaseRankings) {
		logger.Info("Starting rankings phase")

		// Use the rankings service instead of a child workflow
		if err := rankings.ProcessAllRankings(ctx, params, state); err != nil {
			if common.IsRateLimitError(err) {
				// Save state and continue as new if we hit rate limits
				if saveErr := w.stateManager.SaveCheckpoint(ctx); saveErr != nil {
					logger.Error("Failed to save checkpoint", "error", saveErr)
				}
				return nil, workflow.NewContinueAsNewError(ctx, definitions.SyncWorkflowName, params)
			}
			return nil, err
		}

		// Mark phase as completed
		w.orchestrator.MarkPhaseCompleted(models.PhaseRankings)
		w.stateManager.UpdateProgress(models.PhaseRankings, state.PartialResults.RankingsProcessed)
		logger.Info("Rankings phase completed", "processed", state.PartialResults.RankingsProcessed)
	}

	// Process Reports phase
	if !w.orchestrator.IsPhaseCompleted(models.PhaseReports) {
		logger.Info("Starting reports phase")

		// Use the reports service instead of a child workflow
		if err := reports.ProcessAllReports(ctx, params, state); err != nil {
			if common.IsRateLimitError(err) {
				// Save state and continue as new if we hit rate limits
				if saveErr := w.stateManager.SaveCheckpoint(ctx); saveErr != nil {
					logger.Error("Failed to save checkpoint", "error", saveErr)
				}
				return nil, workflow.NewContinueAsNewError(ctx, definitions.SyncWorkflowName, params)
			}
			return nil, err
		}

		// Mark phase as completed
		w.orchestrator.MarkPhaseCompleted(models.PhaseReports)
		w.stateManager.UpdateProgress(models.PhaseReports, state.PartialResults.ReportsProcessed)
		logger.Info("Reports phase completed", "processed", state.PartialResults.ReportsProcessed)
	}

	// Process Builds phase
	if !w.orchestrator.IsPhaseCompleted(models.PhaseBuilds) {
		logger.Info("Starting builds phase")

		// Use the builds service instead of a child workflow
		if err := builds.ProcessAllBuilds(ctx, params, state); err != nil {
			if common.IsRateLimitError(err) {
				// Save state and continue as new if we hit rate limits
				if saveErr := w.stateManager.SaveCheckpoint(ctx); saveErr != nil {
					logger.Error("Failed to save checkpoint", "error", saveErr)
				}
				return nil, workflow.NewContinueAsNewError(ctx, definitions.SyncWorkflowName, params)
			}
			return nil, err
		}

		// Mark phase as completed
		w.orchestrator.MarkPhaseCompleted(models.PhaseBuilds)
		w.stateManager.UpdateProgress(models.PhaseBuilds, state.PartialResults.BuildsProcessed)
		logger.Info("Builds phase completed", "processed", state.PartialResults.BuildsProcessed)
	}

	// Prepare final result
	return &models.WorkflowResult{
		RankingsProcessed: state.PartialResults.RankingsProcessed,
		ReportsProcessed:  state.PartialResults.ReportsProcessed,
		BuildsProcessed:   state.PartialResults.BuildsProcessed,
		StartedAt:         state.StartedAt,
		CompletedAt:       workflow.Now(ctx),
	}, nil
}
