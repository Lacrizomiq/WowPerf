// metrics.go
package warcraftlogsBuildsTemporalWorkflowsCommon

import (
	"fmt"
	"time"
)

// MetricsCollector handles workflow metrics
type MetricsCollector struct {
	startTime time.Time
	metrics   map[string]interface{}
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		startTime: time.Now(),
		metrics:   make(map[string]interface{}),
	}
}

// RecordDuration records the duration of an operation
func (m *MetricsCollector) RecordDuration(operation string, duration time.Duration) {
	m.metrics[fmt.Sprintf("%s_duration", operation)] = duration
}

// RecordCount records a count metric
func (m *MetricsCollector) RecordCount(metric string, value int32) {
	m.metrics[metric] = value
}

// GetMetrics returns all collected metrics
func (m *MetricsCollector) GetMetrics() map[string]interface{} {
	m.metrics["total_duration"] = time.Since(m.startTime)
	return m.metrics
}
