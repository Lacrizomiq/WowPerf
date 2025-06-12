// progress.go
package warcraftlogsBuildsTemporalWorkflowsModels

import "time"

// Phase represents different workflow phases
type Phase string

const (
	PhaseRankings          Phase = "rankings"
	PhaseReports           Phase = "reports"
	PhaseBuilds            Phase = "builds"
	PhaseEquipmentAnalysis Phase = "equipment_analysis"
	PhaseTalentAnalysis    Phase = "talent_analysis"
	PhaseStatAnalysis      Phase = "stat_analysis"
)

// WorkflowProgress tracks the overall progress of a workflow
type WorkflowProgress struct {
	CurrentPhase    Phase                 `json:"current_phase"`
	CompletedPhases map[Phase]bool        `json:"completed_phases"`
	PhaseProgress   map[Phase]PhaseStatus `json:"phase_progress"`
	StartedAt       time.Time             `json:"started_at"`
	LastUpdateAt    time.Time             `json:"last_update_at"`
}

// PhaseStatus represents the status of a specific workflow phase
type PhaseStatus struct {
	ProcessedItems int32     `json:"processed_items"`
	TotalItems     int32     `json:"total_items"`
	StartedAt      time.Time `json:"started_at"`
	CompletedAt    time.Time `json:"completed_at"`
	CurrentSpec    ClassSpec `json:"current_spec"`
}

// NewWorkflowProgress creates a new WorkflowProgress instance
func NewWorkflowProgress() *WorkflowProgress {
	return &WorkflowProgress{
		CompletedPhases: make(map[Phase]bool),
		PhaseProgress:   make(map[Phase]PhaseStatus),
		StartedAt:       time.Now(),
		LastUpdateAt:    time.Now(),
	}
}
