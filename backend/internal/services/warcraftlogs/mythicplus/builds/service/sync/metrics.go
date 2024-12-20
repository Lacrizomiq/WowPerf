package warcraftlogsBuildsSync

import (
	"sync"
	"time"
)

// SyncMetrics track all synchronization operations  and their metrics
type SyncMetrics struct {
	StartTime        time.Time
	EndTime          time.Time
	BatchesTotal     int
	BatchesProcessed int

	// Rankings metrics
	Rankings struct {
		Total             int
		New               int
		Updated           int
		Deleted           int
		Unchanged         int
		ProcessesSpecs    map[string]int // Count per spec
		ProcessesDungeons map[uint]int   // Count per dungeon
	}

	// Reports metrics
	Reports struct {
		Total   int
		New     int
		Updated int
		Deleted int
		Skipped int
	}

	// PlayerBuilds metrics
	PlayerBuilds struct {
		Total   int
		New     int
		Updated int
		Deleted int
		Skipped int
	}

	// Error tracking
	Errors     []string
	ErrorCount int

	mu sync.Mutex
}

// NewSyncMetrics creates a new SyncMetrics instance
func NewSyncMetrics() *SyncMetrics {
	return &SyncMetrics{
		StartTime: time.Now(),
		Rankings: struct {
			Total             int
			New               int
			Updated           int
			Deleted           int
			Unchanged         int
			ProcessesSpecs    map[string]int
			ProcessesDungeons map[uint]int
		}{
			ProcessesSpecs:    make(map[string]int),
			ProcessesDungeons: make(map[uint]int),
		},
	}
}

// RecordBatchProcessed increments the batch processed counter
func (m *SyncMetrics) RecordBatchProcessed(spec string, dungeonID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.BatchesProcessed++
	m.Rankings.ProcessesSpecs[spec]++
	m.Rankings.ProcessesDungeons[dungeonID]++
}

// RecordRankingChanges updates the rankings metrics
func (m *SyncMetrics) RecordRankingChanges(total int, new int, updated int, deleted int, unchanged int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Rankings.New += new
	m.Rankings.Updated += updated
	m.Rankings.Deleted += deleted
	m.Rankings.Unchanged += unchanged
	m.Rankings.Total = new + updated + unchanged
}

// RecordReportChanges updates the reports metrics
func (m *SyncMetrics) RecordReportChanges(new, updated, deleted, skipped int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Reports.New += new
	m.Reports.Updated += updated
	m.Reports.Deleted += deleted
	m.Reports.Skipped += skipped
	m.Reports.Total = new + updated
}

// RecordPlayerBuildChanges updates the player builds metrics
func (m *SyncMetrics) RecordPlayerBuildChanges(new, updated, deleted, skipped int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.PlayerBuilds.New += new
	m.PlayerBuilds.Updated += updated
	m.PlayerBuilds.Deleted += deleted
	m.PlayerBuilds.Skipped += skipped
	m.PlayerBuilds.Total = new + updated
}

// RecordError adds an error message to the metrics
func (m *SyncMetrics) RecordError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Errors = append(m.Errors, err.Error())
	m.ErrorCount++
}

// Complete marks the sync operation as completed
func (m *SyncMetrics) Complete() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.EndTime = time.Now()
}

// GetDuration returns the total duration of the sync operation
func (m *SyncMetrics) GetDuration() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.EndTime.IsZero() {
		return time.Since(m.StartTime)
	}

	return m.EndTime.Sub(m.StartTime)
}

// GetSummary returns a summary of the metrics
func (m *SyncMetrics) GetSummary() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	return map[string]interface{}{
		"duration":          m.GetDuration().String(),
		"batches_total":     m.BatchesTotal,
		"batches_processed": m.BatchesProcessed,
		"rankings": map[string]int{
			"total":     m.Rankings.Total,
			"new":       m.Rankings.New,
			"updated":   m.Rankings.Updated,
			"deleted":   m.Rankings.Deleted,
			"unchanged": m.Rankings.Unchanged,
		},
		"reports": map[string]int{
			"total":   m.Reports.Total,
			"new":     m.Reports.New,
			"updated": m.Reports.Updated,
			"deleted": m.Reports.Deleted,
			"skipped": m.Reports.Skipped,
		},
		"player_builds": map[string]int{
			"total":   m.PlayerBuilds.Total,
			"new":     m.PlayerBuilds.New,
			"updated": m.PlayerBuilds.Updated,
			"deleted": m.PlayerBuilds.Deleted,
			"skipped": m.PlayerBuilds.Skipped,
		},
		"error_count": m.ErrorCount,
	}
}
