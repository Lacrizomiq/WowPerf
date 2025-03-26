// workflow_state.go
package warcraftlogsBuildsTemporalWorkflowsState

import (
	"time"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// WorkflowState represents the complete state of a workflow
type WorkflowState struct {
	// D'après sync_workflow.go
	CurrentSpec       *models.ClassSpec
	CurrentDungeon    *models.Dungeon
	ProcessedSpecs    map[string]bool
	ProcessedDungeons map[string]bool
	LastCheckpoint    time.Time
	StartedAt         time.Time
	PartialResults    *models.WorkflowResult
}

// NewWorkflowState creates a new workflow state instance
func NewWorkflowState() *WorkflowState {
	return &WorkflowState{
		ProcessedSpecs:    make(map[string]bool),
		ProcessedDungeons: make(map[string]bool),
		StartedAt:         time.Now(),
		PartialResults:    &models.WorkflowResult{StartedAt: time.Now()},
	}
}
