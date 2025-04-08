// analyze_workflow.go
package warcraftlogsBuildsTemporalWorkflowsAnalyze

import (
	"fmt"

	"go.temporal.io/sdk/workflow"

	buildsstatistics "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/builds_statistics"
	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// AnalyzeWorkflow implements the workflow for analyzing player builds
type AnalyzeWorkflow struct {
	stateManager *state.AnalysisStateManager
	orchestrator *AnalyzeOrchestrator
}

// NewAnalyzeWorkflow creates a new instance of the analysis workflow
func NewAnalyzeWorkflow() definitions.AnalyzeWorkflow {
	return &AnalyzeWorkflow{
		stateManager: state.NewAnalysisStateManager(),
		orchestrator: NewAnalyzeOrchestrator(),
	}
}

// Execute runs the main analysis workflow
// It orchestrates the three main phases: Equipment Analysis, Talent Analysis, and Stat Analysis
func (w *AnalyzeWorkflow) Execute(ctx workflow.Context, config models.AnalysisWorkflowConfig) (*models.AnalysisWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting analyze workflow orchestration",
		"workflowID", workflow.GetInfo(ctx).WorkflowExecution.ID,
		"attempt", workflow.GetInfo(ctx).Attempt)

	// Validate input parameters
	if len(config.Specs) == 0 {
		return nil, fmt.Errorf("no specs found in config")
	}

	if len(config.Dungeons) == 0 {
		return nil, fmt.Errorf("no dungeons found in config")
	}

	// Log configuration details
	if len(config.Specs) > 0 {
		logger.Info("First spec to analyze",
			"className", config.Specs[0].ClassName,
			"specName", config.Specs[0].SpecName)
	}

	if len(config.Dungeons) > 0 {
		logger.Info("First dungeon to analyze",
			"name", config.Dungeons[0].Name,
			"ID", config.Dungeons[0].EncounterID)
	}

	// Load checkpoint state if this is a continuation
	if err := w.stateManager.LoadCheckpoint(ctx); err != nil {
		return nil, err
	}

	state := w.stateManager.GetState()

	// Process Equipment Analysis phase
	if !w.orchestrator.IsPhaseCompleted(models.PhaseEquipmentAnalysis) {
		logger.Info("Starting equipment analysis phase")

		if err := buildsstatistics.ProcessBuildStatistics(ctx, config, state); err != nil {
			if common.IsRateLimitError(err) {
				// Save state and continue as new if we hit rate limits
				if saveErr := w.stateManager.SaveCheckpoint(ctx); saveErr != nil {
					logger.Error("Failed to save checkpoint", "error", saveErr)
				}
				return nil, workflow.NewContinueAsNewError(ctx, definitions.AnalyzeBuildsWorkflowName, config)
			}
			return nil, err
		}

		// Mark phase as completed
		w.orchestrator.MarkPhaseCompleted(models.PhaseEquipmentAnalysis)
		logger.Info("Equipment analysis phase completed", "itemsAnalyzed", state.Results.ItemsAnalyzed)
	}

	// Process Talent Analysis phase
	if !w.orchestrator.IsPhaseCompleted(models.PhaseTalentAnalysis) {
		logger.Info("Starting talent analysis phase")

		if err := buildsstatistics.ProcessTalentStatistics(ctx, config, state); err != nil {
			if common.IsRateLimitError(err) {
				// Save state and continue as new if we hit rate limits
				if saveErr := w.stateManager.SaveCheckpoint(ctx); saveErr != nil {
					logger.Error("Failed to save checkpoint", "error", saveErr)
				}
				return nil, workflow.NewContinueAsNewError(ctx, definitions.AnalyzeBuildsWorkflowName, config)
			}
			return nil, err
		}

		// Mark phase as completed
		w.orchestrator.MarkPhaseCompleted(models.PhaseTalentAnalysis)
		logger.Info("Talent analysis phase completed", "talentsAnalyzed", state.Results.TalentsAnalyzed)
	}

	// Process Stat Analysis phase
	if !w.orchestrator.IsPhaseCompleted(models.PhaseStatAnalysis) {
		logger.Info("Starting stat analysis phase")

		if err := buildsstatistics.ProcessStatStatistics(ctx, config, state); err != nil {
			if common.IsRateLimitError(err) {
				// Save state and continue as new if we hit rate limits
				if saveErr := w.stateManager.SaveCheckpoint(ctx); saveErr != nil {
					logger.Error("Failed to save checkpoint", "error", saveErr)
				}
				return nil, workflow.NewContinueAsNewError(ctx, definitions.AnalyzeBuildsWorkflowName, config)
			}
			return nil, err
		}

		// Mark phase as completed
		w.orchestrator.MarkPhaseCompleted(models.PhaseStatAnalysis)
		logger.Info("Stat analysis phase completed", "statsAnalyzed", state.Results.StatsAnalyzed)
	}

	// Prepare final result
	state.Results.CompletedAt = workflow.Now(ctx)

	logger.Info("Analyze workflow completed successfully",
		"totalBuilds", state.Results.TotalBuilds,
		"itemsAnalyzed", state.Results.ItemsAnalyzed,
		"talentsAnalyzed", state.Results.TalentsAnalyzed,
		"statsAnalyzed", state.Results.StatsAnalyzed)

	return state.Results, nil
}
