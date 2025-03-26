package warcraftlogsBuildsTemporalWorkflowsSync

import (
	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"

	"go.temporal.io/sdk/workflow"
)

// Orchestrator handles the coordination and phase tracking
type Orchestrator struct {
	completedPhases map[models.Phase]bool
}

// NewOrchestrator creates a new orchestrator instance
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		completedPhases: make(map[models.Phase]bool),
	}
}

// Todo : Use this function to handle workflow errors instead of coding it in each part of the workflow sync
func handleWorkflowError(ctx workflow.Context, err error, state *state.WorkflowState, stateManager *state.Manager) (*models.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx) // Get logger from context

	if common.IsRateLimitError(err) {
		// Save state and continue as new if we hit rate limits
		if saveErr := stateManager.SaveCheckpoint(ctx); saveErr != nil {
			logger.Error("Failed to save checkpoint", "error", saveErr)
		}
		return nil, workflow.NewContinueAsNewError(ctx, definitions.SyncWorkflowName, state.PartialResults)
	}
	return nil, err
}

// IsPhaseCompleted checks if a specific phase is completed
func (o *Orchestrator) IsPhaseCompleted(phase models.Phase) bool {
	return o.completedPhases[phase]
}

// MarkPhaseCompleted marks a phase as completed
func (o *Orchestrator) MarkPhaseCompleted(phase models.Phase) {
	o.completedPhases[phase] = true
}
