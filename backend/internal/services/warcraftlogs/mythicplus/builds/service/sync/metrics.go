package warcraftlogsBuildsSync

import (
	"sync"
	"time"

	"wowperf/internal/services/warcraftlogs"
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

	// Rate limit metrics
	RateLimits struct {
		Hits              int           // Number of times the rate limit was hit
		TotalWaitTime     time.Duration // Total wait time for all rate limit hits
		MinPointsObserved float64       // Minimum points observed during rate limit hits
		MaxPointsObserved float64       // Maximum points observed during rate limit hits
		LastResetTime     time.Time     // Time when the rate limit was last reset
	}

	// Error tracking
	Errors struct {
		Total      int
		ByType     map[warcraftlogs.ErrorType]int
		Message    []ErrorEntry
		Retries    int
		MaxRetries int
	}

	mu sync.Mutex
}

// ErrorEntry represents an error entry with timestamp, type, message, retryable status, and retry count
type ErrorEntry struct {
	TimeStamp time.Time
	Type      warcraftlogs.ErrorType
	Message   string
	Retryable bool
	Retried   bool
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
		Errors: struct {
			Total      int
			ByType     map[warcraftlogs.ErrorType]int
			Message    []ErrorEntry
			Retries    int
			MaxRetries int
		}{
			ByType: make(map[warcraftlogs.ErrorType]int),
		},
	}
}

// RecordRateLimit records rate limit information
func (m *SyncMetrics) RecordRateLimit(info *warcraftlogs.RateLimitInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.RateLimits.Hits++
	m.RateLimits.TotalWaitTime += info.ResetIn

	if info.RemainingPoints < m.RateLimits.MinPointsObserved || m.RateLimits.MinPointsObserved == 0 {
		m.RateLimits.MinPointsObserved = info.RemainingPoints
	}
	if info.RemainingPoints > m.RateLimits.MaxPointsObserved {
		m.RateLimits.MaxPointsObserved = info.RemainingPoints
	}

	if !info.NextRefresh.IsZero() {
		m.RateLimits.LastResetTime = info.NextRefresh
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

// RecordError records an error with additional context
func (m *SyncMetrics) RecordError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Errors.Total++

	if wlErr, ok := err.(*warcraftlogs.WarcraftLogsError); ok {
		m.Errors.ByType[wlErr.Type]++
		m.Errors.Message = append(m.Errors.Message, ErrorEntry{
			TimeStamp: time.Now(),
			Type:      wlErr.Type,
			Message:   wlErr.Message,
			Retryable: wlErr.Retryable,
		})
	} else {
		m.Errors.ByType[warcraftlogs.ErrorTypeAPI]++
		m.Errors.Message = append(m.Errors.Message, ErrorEntry{
			TimeStamp: time.Now(),
			Type:      warcraftlogs.ErrorTypeAPI,
			Message:   err.Error(),
		})
	}

	// Garder seulement les 100 dernières erreurs
	if len(m.Errors.Message) > 100 {
		m.Errors.Message = m.Errors.Message[1:]
	}
}

// RecordRetry records a retry attempt
func (m *SyncMetrics) RecordRetry(successful bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Errors.Retries++
	if successful {
		m.Errors.Message[len(m.Errors.Message)-1].Retried = true
	}
}

// GetRateLimitSummary returns a summary of rate limit metrics
func (m *SyncMetrics) GetRateLimitSummary() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	return map[string]interface{}{
		"total_hits":          m.RateLimits.Hits,
		"total_wait_time":     m.RateLimits.TotalWaitTime.String(),
		"min_points_observed": m.RateLimits.MinPointsObserved,
		"max_points_observed": m.RateLimits.MaxPointsObserved,
		"last_reset":          m.RateLimits.LastResetTime,
	}
}

// GetErrorSummary returns a summary of error metrics
func (m *SyncMetrics) GetErrorSummary() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	return map[string]interface{}{
		"total_errors":       m.Errors.Total,
		"errors_by_type":     m.Errors.ByType,
		"total_retries":      m.Errors.Retries,
		"retry_success_rate": float64(m.Errors.Retries) / float64(m.Errors.Total),
	}
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
		"rate_limits": m.GetRateLimitSummary(),
		"errors":      m.GetErrorSummary(),
	}
}

// Complete marks the sync operation as completed
func (m *SyncMetrics) Complete() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.EndTime = time.Now()
}
