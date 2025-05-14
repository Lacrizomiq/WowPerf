package warcraftlogsBuildMetrics

import (
	"encoding/json"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

// WorkflowMetricsSnapshot represents a snapshot of the metrics of a workflow
// This structure will be serialized to JSON and stored in WorkflowState.PerformanceMetrics
type WorkflowMetricsSnapshot struct {
	// Common metrics
	WorkflowType      string  `json:"workflow_type"`
	ClassName         string  `json:"class_name,omitempty"`
	BatchID           string  `json:"batch_id,omitempty"`
	ItemsProcessed    int     `json:"items_processed"`
	TotalToProcess    int     `json:"total_to_process,omitempty"`
	ProgressPercent   float64 `json:"progress_percent"`
	ApiRequestsCount  int     `json:"api_requests_count,omitempty"`
	ContinuationCount int     `json:"continuation_count,omitempty"`

	// Durations
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time,omitempty"`
	TotalDuration int64     `json:"total_duration_ms,omitempty"` // in milliseconds

	// Operations and performance metrics
	Operations map[string]OperationMetrics `json:"operations,omitempty"`

	// Specific metrics by workflow type (to be filled according to the type)
	ReportsMetrics *ReportsMetrics `json:"reports_metrics,omitempty"`
	BuildsMetrics  *BuildsMetrics  `json:"builds_metrics,omitempty"`

	// Errors
	Errors map[string]int `json:"errors,omitempty"`
}

// OperationMetrics contains the metrics for a specific operation
type OperationMetrics struct {
	Count       int     `json:"count"`
	TotalTimeMs int64   `json:"total_time_ms"`
	AvgTimeMs   float64 `json:"avg_time_ms"`
	MaxTimeMs   int64   `json:"max_time_ms"`
	MinTimeMs   int64   `json:"min_time_ms"`
}

// ReportsMetrics contains specific metrics for the Reports workflow
type ReportsMetrics struct {
	RankingsProcessed int `json:"rankings_processed"`
	ReportsProcessed  int `json:"reports_processed"`
	FailedReports     int `json:"failed_reports"`
	RateLimitHits     int `json:"rate_limit_hits"`
}

// BuildsMetrics contains specific metrics for the Builds workflow
type BuildsMetrics struct {
	BuildsProcessed   int            `json:"builds_processed"`
	ChildrenLaunched  int            `json:"children_launched"`
	ChildrenCompleted int            `json:"children_completed"`
	ChildrenFailed    int            `json:"children_failed"`
	BuildsByClassSpec map[string]int `json:"builds_by_class_spec,omitempty"`
}

// MetricsCollector is a structure to collect metrics
type MetricsCollector struct {
	workflowType    string
	className       string
	batchID         string
	startTime       time.Time
	snapshot        WorkflowMetricsSnapshot
	operationTimers map[string]time.Time // To measure the durations of ongoing operations
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(workflowType, className, batchID string) *MetricsCollector {
	now := time.Now()
	return &MetricsCollector{
		workflowType: workflowType,
		className:    className,
		batchID:      batchID,
		startTime:    now,
		snapshot: WorkflowMetricsSnapshot{
			WorkflowType: workflowType,
			ClassName:    className,
			BatchID:      batchID,
			StartTime:    now,
			Operations:   make(map[string]OperationMetrics),
			Errors:       make(map[string]int),
		},
		operationTimers: make(map[string]time.Time),
	}
}

// StartOperation starts a timer for an operation
func (m *MetricsCollector) StartOperation(operation string) {
	m.operationTimers[operation] = time.Now()
}

// EndOperation ends a timer for an operation and records its duration
func (m *MetricsCollector) EndOperation(operation string) time.Duration {
	if startTime, ok := m.operationTimers[operation]; ok {
		duration := time.Since(startTime)
		delete(m.operationTimers, operation)

		// Update our internal metrics
		durationMs := duration.Milliseconds()

		if _, ok := m.snapshot.Operations[operation]; !ok {
			m.snapshot.Operations[operation] = OperationMetrics{
				Count:       0,
				TotalTimeMs: 0,
				AvgTimeMs:   0,
				MaxTimeMs:   0,
				MinTimeMs:   durationMs,
			}
		}

		metrics := m.snapshot.Operations[operation]
		metrics.Count++
		metrics.TotalTimeMs += durationMs
		metrics.AvgTimeMs = float64(metrics.TotalTimeMs) / float64(metrics.Count)

		if durationMs > metrics.MaxTimeMs {
			metrics.MaxTimeMs = durationMs
		}
		if durationMs < metrics.MinTimeMs {
			metrics.MinTimeMs = durationMs
		}

		m.snapshot.Operations[operation] = metrics

		return duration
	}

	return 0
}

// RecordItemsProcessed records processed items
func (m *MetricsCollector) RecordItemsProcessed(count int, status string) {
	m.snapshot.ItemsProcessed += count

	// Calculate progress percentage if possible
	if m.snapshot.TotalToProcess > 0 {
		m.snapshot.ProgressPercent = float64(m.snapshot.ItemsProcessed) / float64(m.snapshot.TotalToProcess) * 100
	}
}

// RecordAPIRequest records an API request
func (m *MetricsCollector) RecordAPIRequest() {
	m.snapshot.ApiRequestsCount++

}

// RecordError records an error
func (m *MetricsCollector) RecordError(errorType string) {
	if _, ok := m.snapshot.Errors[errorType]; !ok {
		m.snapshot.Errors[errorType] = 0
	}
	m.snapshot.Errors[errorType]++

}

// RecordContinuation records a workflow continuation
func (m *MetricsCollector) RecordContinuation() {
	m.snapshot.ContinuationCount++

}

// RecordChildWorkflow records the launch of a child workflow
func (m *MetricsCollector) RecordChildWorkflow(childType string) {
	if m.snapshot.BuildsMetrics == nil {
		m.snapshot.BuildsMetrics = &BuildsMetrics{
			ChildrenLaunched:  0,
			ChildrenCompleted: 0,
			ChildrenFailed:    0,
			BuildsByClassSpec: make(map[string]int),
		}
	}

	m.snapshot.BuildsMetrics.ChildrenLaunched++

}

// RecordChildWorkflowCompletion records the completion of a child workflow
func (m *MetricsCollector) RecordChildWorkflowCompletion(success bool) {
	if m.snapshot.BuildsMetrics == nil {
		m.snapshot.BuildsMetrics = &BuildsMetrics{
			ChildrenLaunched:  0,
			ChildrenCompleted: 0,
			ChildrenFailed:    0,
			BuildsByClassSpec: make(map[string]int),
		}
	}

	if success {
		m.snapshot.BuildsMetrics.ChildrenCompleted++
	} else {
		m.snapshot.BuildsMetrics.ChildrenFailed++
	}
}

// RecordBuildByClassSpec records a build by class/spec
func (m *MetricsCollector) RecordBuildByClassSpec(class, spec string, count int) {
	if m.snapshot.BuildsMetrics == nil {
		m.snapshot.BuildsMetrics = &BuildsMetrics{
			BuildsByClassSpec: make(map[string]int),
		}
	}

	key := class + "-" + spec
	m.snapshot.BuildsMetrics.BuildsByClassSpec[key] = m.snapshot.BuildsMetrics.BuildsByClassSpec[key] + count

}

// RecordRateLimitHit records a rate limit hit
func (m *MetricsCollector) RecordRateLimitHit() {
	if m.snapshot.ReportsMetrics == nil {
		m.snapshot.ReportsMetrics = &ReportsMetrics{
			RateLimitHits: 0,
		}
	}

	m.snapshot.ReportsMetrics.RateLimitHits++

	// Also record as an error
	m.RecordError("rate_limit")
}

// SetTotalToProcess sets the total number of items to process
func (m *MetricsCollector) SetTotalToProcess(total int) {
	m.snapshot.TotalToProcess = total

	// Recalculate the progress percentage
	if total > 0 {
		m.snapshot.ProgressPercent = float64(m.snapshot.ItemsProcessed) / float64(total) * 100
	}
}

// SetReportsMetrics initializes or updates the Reports metrics
func (m *MetricsCollector) SetReportsMetrics(rankingsProcessed, reportsProcessed, failedReports int) {
	m.snapshot.ReportsMetrics = &ReportsMetrics{
		RankingsProcessed: rankingsProcessed,
		ReportsProcessed:  reportsProcessed,
		FailedReports:     failedReports,
	}
}

// SetBuildsMetrics initializes or updates the Builds metrics
func (m *MetricsCollector) SetBuildsMetrics(buildsProcessed int) {
	if m.snapshot.BuildsMetrics == nil {
		m.snapshot.BuildsMetrics = &BuildsMetrics{
			BuildsByClassSpec: make(map[string]int),
		}
	}

	m.snapshot.BuildsMetrics.BuildsProcessed = buildsProcessed
}

// Finish finishes the metrics collection and returns the final snapshot
func (m *MetricsCollector) Finish(status string) WorkflowMetricsSnapshot {
	endTime := time.Now()
	m.snapshot.EndTime = endTime
	m.snapshot.TotalDuration = endTime.Sub(m.startTime).Milliseconds()

	return m.snapshot
}

// GetSnapshot returns the current snapshot of the metrics
func (m *MetricsCollector) GetSnapshot() WorkflowMetricsSnapshot {
	return m.snapshot
}

// ToJSON converts the snapshot to JSON for storage in WorkflowState
func (m *MetricsCollector) ToJSON() ([]byte, error) {
	return json.Marshal(m.snapshot)
}

// UpdateWorkflowState updates the workflow state with the current metrics
func (m *MetricsCollector) UpdateWorkflowState(state *warcraftlogsBuilds.WorkflowState) error {
	metricsJSON, err := m.ToJSON()
	if err != nil {
		return err
	}

	state.PerformanceMetrics = metricsJSON
	state.ItemsProcessed = m.snapshot.ItemsProcessed
	state.TotalItemsToProcess = m.snapshot.TotalToProcess
	state.ProgressPercentage = m.snapshot.ProgressPercent
	state.ApiRequestsCount = m.snapshot.ApiRequestsCount
	state.ContinuationCount = m.snapshot.ContinuationCount

	// Calculate the estimated completion date
	if m.snapshot.ProgressPercent > 0 {
		elapsed := time.Since(m.startTime)
		totalEstimated := time.Duration(float64(elapsed) / m.snapshot.ProgressPercent * 100)
		state.EstimatedCompletion = m.startTime.Add(totalEstimated)
	}

	return nil
}
