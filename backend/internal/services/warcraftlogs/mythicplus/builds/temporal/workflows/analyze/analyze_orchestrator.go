// analyze_orchestrator.go
package warcraftlogsBuildsTemporalWorkflowsAnalyze

import (
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// Orchestrator manages the execution flow of different analysis phases
type AnalyzeOrchestrator struct {
	completedPhases map[models.Phase]bool
}

// NewOrchestrator creates a new instance of the analysis orchestrator
func NewAnalyzeOrchestrator() *AnalyzeOrchestrator {
	return &AnalyzeOrchestrator{
		completedPhases: make(map[models.Phase]bool),
	}
}

// IsPhaseCompleted checks if a specific phase has been completed
func (o *AnalyzeOrchestrator) IsPhaseCompleted(phase models.Phase) bool {
	return o.completedPhases[phase]
}

// MarkPhaseCompleted marks a phase as completed
func (o *AnalyzeOrchestrator) MarkPhaseCompleted(phase models.Phase) {
	o.completedPhases[phase] = true
}

// GetCompletedPhases returns all completed phases
func (o *AnalyzeOrchestrator) GetCompletedPhases() map[models.Phase]bool {
	// Return a copy to prevent external modification
	copy := make(map[models.Phase]bool)
	for k, v := range o.completedPhases {
		copy[k] = v
	}
	return copy
}

// ResetPhases resets all phases to incomplete
func (o *AnalyzeOrchestrator) ResetPhases() {
	o.completedPhases = make(map[models.Phase]bool)
}
