// batch.go
package warcraftlogsBuildsTemporalWorkflowsModels

import (
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

// BatchConfig defines parameters for batch processing
type BatchConfig struct {
	Size        int32         `json:"size" yaml:"size"`
	RetryDelay  time.Duration `json:"retry_delay" yaml:"retry_delay"`
	MaxAttempts int32         `json:"max_attempts" yaml:"max_attempts"`
}

// BatchResult represents the outcome of a batch processing operation
type BatchResult struct {
	ClassName         string                       `json:"class_name"`           // Class name
	SpecName          string                       `json:"spec_name"`            // Spec name
	EncounterID       uint32                       `json:"encounter_id"`         // Encounter ID
	ProcessedItems    int32                        `json:"processed_items"`      // Processed items
	RankingsCount     int32                        `json:"rankings_count"`       // Rankings count
	ProcessedAt       time.Time                    `json:"processed_at"`         // Processed at
	ProcessedReports  []*warcraftlogsBuilds.Report `json:"processed_reports"`    // Reports processed in this batch
	BuildsByClassSpec map[string]int32             `json:"builds_by_class_spec"` // Number of builds by class+spec
}
