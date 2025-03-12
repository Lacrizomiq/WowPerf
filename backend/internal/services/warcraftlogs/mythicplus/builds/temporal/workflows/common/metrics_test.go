package warcraftlogsBuildsTemporalWorkflowsCommon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMetricsCollector tests the creation of a new metrics collector.
// It verifies that:
// - Collector is properly initialized
// - Start time is set
// - Metrics map is initialized
func TestNewMetricsCollector(t *testing.T) {
	t.Log("Starting metrics collector initialization test")

	collector := NewMetricsCollector()
	require.NotNil(t, collector, "Metrics collector should not be nil")

	t.Log("Validating collector initialization...")
	assert.NotNil(t, collector.metrics, "Metrics map should be initialized")
	assert.False(t, collector.startTime.IsZero(), "Start time should be set")

	t.Logf("Collector initialized at: %v", collector.startTime)
	t.Log("Metrics collector initialization validated successfully")
}

// TestRecordDuration tests the duration recording functionality.
// It verifies that:
// - Durations are properly recorded
// - Multiple durations for different operations are stored correctly
// - Durations are stored with correct keys
func TestRecordDuration(t *testing.T) {
	t.Log("Starting duration recording tests")

	collector := NewMetricsCollector()

	testCases := []struct {
		name      string
		operation string
		duration  time.Duration
	}{
		{
			name:      "Short operation",
			operation: "quick_task",
			duration:  100 * time.Millisecond,
		},
		{
			name:      "Long operation",
			operation: "heavy_task",
			duration:  5 * time.Second,
		},
		{
			name:      "Zero duration",
			operation: "instant_task",
			duration:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Recording duration for operation: %s", tc.operation)
			t.Logf("Duration value: %v", tc.duration)

			collector.RecordDuration(tc.operation, tc.duration)

			metricKey := tc.operation + "_duration"
			recordedDuration, exists := collector.metrics[metricKey]

			assert.True(t, exists, "Duration metric should be recorded")
			assert.Equal(t, tc.duration, recordedDuration)

			t.Logf("Duration recorded with key: %s", metricKey)
			t.Log("Duration recording validated successfully")
		})
	}
}

// TestRecordCount tests the count recording functionality.
// It verifies that:
// - Count values are properly recorded
// - Multiple counts for different metrics are stored correctly
// - Count values are stored with correct keys
func TestRecordCount(t *testing.T) {
	t.Log("Starting count recording tests")

	collector := NewMetricsCollector()

	testCases := []struct {
		name   string
		metric string
		value  int32
	}{
		{
			name:   "Positive count",
			metric: "processed_items",
			value:  42,
		},
		{
			name:   "Zero count",
			metric: "errors",
			value:  0,
		},
		{
			name:   "Large count",
			metric: "total_requests",
			value:  99999,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Recording count for metric: %s", tc.metric)
			t.Logf("Count value: %d", tc.value)

			collector.RecordCount(tc.metric, tc.value)

			recordedCount, exists := collector.metrics[tc.metric]

			assert.True(t, exists, "Count metric should be recorded")
			assert.Equal(t, tc.value, recordedCount)

			t.Logf("Count recorded with key: %s", tc.metric)
			t.Log("Count recording validated successfully")
		})
	}
}

// TestGetMetrics tests the retrieval of all metrics.
// It verifies that:
// - All recorded metrics are returned
// - Total duration is calculated and included
// - Metrics are not modified during retrieval
func TestGetMetrics(t *testing.T) {
	t.Log("Starting metrics retrieval test")

	collector := NewMetricsCollector()

	// Record some test metrics
	collector.RecordDuration("operation1", 1*time.Second)
	collector.RecordCount("count1", 10)

	t.Log("Recorded test metrics, waiting briefly before retrieval...")
	time.Sleep(100 * time.Millisecond) // Ensure some time passes

	t.Log("Retrieving all metrics")
	metrics := collector.GetMetrics()

	t.Log("Validating retrieved metrics...")

	// Verify recorded metrics exist
	assert.Equal(t, time.Second, metrics["operation1_duration"])
	assert.Equal(t, int32(10), metrics["count1"])

	// Verify total duration is present and reasonable
	totalDuration, exists := metrics["total_duration"]
	assert.True(t, exists, "Total duration should be included")
	duration, ok := totalDuration.(time.Duration)
	assert.True(t, ok, "Total duration should be a time.Duration")
	assert.True(t, duration >= 100*time.Millisecond, "Total duration should be at least the sleep time")

	t.Logf("Retrieved %d metrics", len(metrics))
	for key, value := range metrics {
		t.Logf("Metric - %s: %v", key, value)
	}
	t.Log("Metrics retrieval validated successfully")
}
