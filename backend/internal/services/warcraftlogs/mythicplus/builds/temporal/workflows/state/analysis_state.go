// analysis_state.go
package warcraftlogsBuildsTemporalWorkflowsState

import (
	"time"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// AnalysisWorkflowState follows the progress of the analysis
type AnalysisWorkflowState struct {
	CurrentSpec           *models.ClassSpec
	CurrentDungeon        *models.Dungeon
	ProcessedCombinations map[string]bool // Class/spec/dungeon combinations processed
	LastCheckpoint        time.Time
	StartedAt             time.Time
	Results               *models.AnalysisWorkflowResult
}

// NewAnalysisWorkflowState creates a new instance of the state
func NewAnalysisWorkflowState() *AnalysisWorkflowState {
	return &AnalysisWorkflowState{
		ProcessedCombinations: make(map[string]bool),
		StartedAt:             time.Now(),
		Results:               &models.AnalysisWorkflowResult{StartedAt: time.Now()},
	}
}
