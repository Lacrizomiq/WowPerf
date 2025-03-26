// state/manager.go
package warcraftlogsBuildsTemporalWorkflowsState

import (
	"time"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"

	"go.temporal.io/sdk/workflow"
)

// Manager handles state management for workflows
type Manager struct {
	state          *WorkflowState
	progress       *models.WorkflowProgress
	lastCheckpoint time.Time
}

// NewManager creates a new state manager
func NewManager() *Manager {
	return &Manager{
		state:    NewWorkflowState(),
		progress: models.NewWorkflowProgress(),
	}
}

// GetState returns the current workflow state
func (m *Manager) GetState() *WorkflowState {
	return m.state
}

// SaveCheckpoint saves the current state
func (m *Manager) SaveCheckpoint(ctx workflow.Context) error {
	m.lastCheckpoint = workflow.Now(ctx)
	m.state.LastCheckpoint = m.lastCheckpoint
	return nil
}

// LoadCheckpoint loads the last checkpoint
func (m *Manager) LoadCheckpoint(ctx workflow.Context) error {
	if workflow.GetInfo(ctx).Attempt > 1 {
		return workflow.GetLastCompletionResult(ctx, &m.state)
	}
	return nil
}

// UpdateProgress updates progress for the current phase
func (m *Manager) UpdateProgress(phase models.Phase, processedItems int32) {
	m.progress.CurrentPhase = phase
	if status, exists := m.progress.PhaseProgress[phase]; exists {
		status.ProcessedItems = processedItems
		m.progress.PhaseProgress[phase] = status
	}
	m.progress.LastUpdateAt = time.Now()
}
