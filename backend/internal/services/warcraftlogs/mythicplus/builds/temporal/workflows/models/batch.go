// batch.go
package warcraftlogsBuildsTemporalWorkflowsModels

import "time"

// BatchConfig defines parameters for batch processing
type BatchConfig struct {
	Size        int32         `json:"size" yaml:"size"`
	RetryDelay  time.Duration `json:"retry_delay" yaml:"retry_delay"`
	MaxAttempts int32         `json:"max_attempts" yaml:"max_attempts"`
}

// BatchResult represents the outcome of a batch processing operation
type BatchResult struct {
	// D'après rankings_workflow.go
	ClassName      string    `json:"class_name"`
	SpecName       string    `json:"spec_name"`
	EncounterID    uint32    `json:"encounter_id"`
	ProcessedItems int32     `json:"processed_items"`
	RankingsCount  int32     `json:"rankings_count"`
	ProcessedAt    time.Time `json:"processed_at"`
}
