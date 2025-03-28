// analysis_manager.go
package warcraftlogsBuildsTemporalWorkflowsState

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// AnalysisStateManager manages the state of the analysis workflow
type AnalysisStateManager struct {
	state          *AnalysisWorkflowState
	lastCheckpoint time.Time
}

// NewAnalysisStateManager creates a new state manager
func NewAnalysisStateManager() *AnalysisStateManager {
	return &AnalysisStateManager{
		state: NewAnalysisWorkflowState(),
	}
}

// GetState returns the current state
func (m *AnalysisStateManager) GetState() *AnalysisWorkflowState {
	return m.state
}

// SaveCheckpoint saves the current state
func (m *AnalysisStateManager) SaveCheckpoint(ctx workflow.Context) error {
	m.lastCheckpoint = workflow.Now(ctx)
	m.state.LastCheckpoint = m.lastCheckpoint
	return nil
}

// LoadCheckpoint loads the last checkpoint
func (m *AnalysisStateManager) LoadCheckpoint(ctx workflow.Context) error {
	if workflow.GetInfo(ctx).Attempt > 1 {
		return workflow.GetLastCompletionResult(ctx, &m.state)
	}
	return nil
}
