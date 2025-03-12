// progress_state.go
package warcraftlogsBuildsTemporalWorkflowsState

import (
	"time"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// ProgressState tracks detailed progress of each phase
type ProgressState struct {
	CurrentPhase     models.Phase
	RankingsProgress *PhaseProgress
	ReportsProgress  *PhaseProgress
	BuildsProgress   *PhaseProgress
}

// PhaseProgress tracks progress of a specific phase
type PhaseProgress struct {
	ProcessedItems int32
	TotalItems     int32
	StartTime      time.Time
	LastUpdateTime time.Time
	Errors         []string
}

// NewProgressState creates a new progress state
func NewProgressState() *ProgressState {
	return &ProgressState{
		RankingsProgress: &PhaseProgress{StartTime: time.Now()},
		ReportsProgress:  &PhaseProgress{StartTime: time.Now()},
		BuildsProgress:   &PhaseProgress{StartTime: time.Now()},
	}
}

// UpdatePhaseProgress updates progress for the current phase
func (p *ProgressState) UpdatePhaseProgress(phase models.Phase, processed, total int32) {
	var progress *PhaseProgress
	switch phase {
	case models.PhaseRankings:
		progress = p.RankingsProgress
	case models.PhaseReports:
		progress = p.ReportsProgress
	case models.PhaseBuilds:
		progress = p.BuildsProgress
	}

	if progress != nil {
		progress.ProcessedItems = processed
		progress.TotalItems = total
		progress.LastUpdateTime = time.Now()
	}
}
