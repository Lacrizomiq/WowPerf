package warcraftlogsBuildsTemporalWorkflowsSync

import (
	"go.temporal.io/sdk/workflow"

	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
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

// handleWorkflowError processes workflow errors and determines next steps
func handleWorkflowError(ctx workflow.Context, err error, state *state.WorkflowState) (*models.WorkflowResult, error) {
	if common.IsRateLimitError(err) {
		return nil, workflow.NewContinueAsNewError(ctx,
			definitions.SyncWorkflowName, // Utilisation de la constante de definitions
			state.PartialResults)
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
