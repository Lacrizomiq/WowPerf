package warcraftlogsBuildsTemporalWorkflowsModels

import "time"

// AnalysisWorkflowConfig contains the specific parameters for the analysis workflow
type AnalysisWorkflowConfig struct {
	Specs         []ClassSpec   `json:"specs"`          // Reuse the ClassSpec structure
	Dungeons      []Dungeon     `json:"dungeons"`       // Reuse the Dungeon structure
	BatchSize     int32         `json:"batch_size"`     // Batch size for processing
	Concurrency   int32         `json:"concurrency"`    // Number of parallel processes
	RetryAttempts int32         `json:"retry_attempts"` // Number of retries in case of failure
	RetryDelay    time.Duration `json:"retry_delay"`    // Delay between retries
}
