// config.go
package warcraftlogsBuildsTemporalWorkflowsModels

import "time"

// TemporalConstants holds all Temporal-related constants
const (
	// DefaultNamespace is the Temporal namespace used by all components
	DefaultNamespace = "default"

	// DefaultTaskQueue is the main task queue for workflow execution
	DefaultTaskQueue = "warcraft-logs-sync"
)

// RankingsWorkflowParams contains only the parameters for the rankings workflow
type RankingsWorkflowParams struct {
	Specs              []ClassSpec   `json:"specs"`
	Dungeons           []Dungeon     `json:"dungeons"`
	MaxRankingsPerSpec int32         `json:"max_rankings_per_spec"`
	BatchSize          int32         `json:"batch_size"`
	RetryDelay         time.Duration `json:"retry_delay"`
	MaxAttempts        int32         `json:"max_attempts"`
	BatchID            string        `json:"batch_id"`
}

// ReportsWorkflowParams contains only the parameters for the reports workflow
type ReportsWorkflowParams struct {
	ClassName        string        `json:"class_name"`
	BatchSize        int32         `json:"batch_size"`
	NumWorkers       int32         `json:"num_workers"`
	RequestDelay     time.Duration `json:"request_delay"`
	ProcessingWindow time.Duration `json:"processing_window"`
	BatchID          string        `json:"batch_id"`

	// New fields for ContinueAsNew
	TotalProcessedRankings int32 `json:"total_processed_rankings"`
	TotalProcessedReports  int32 `json:"total_processed_reports"`
	TotalAPIRequests       int32 `json:"total_api_requests"`
	TotalFailedReports     int32 `json:"total_failed_reports"`
	ContinuationCount      int32 `json:"continuation_count"`

	// for workflow continuation
	ParentWorkflowID string `json:"parent_workflow_id"`
}

// BuildsWorkflowParams contains only the parameters for the builds workflow
type BuildsWorkflowParams struct {
	BatchSize       int32  `json:"batch_size"`
	NumWorkers      int32  `json:"num_workers"`
	ReportBatchSize int32  `json:"report_batch_size"`
	BatchID         string `json:"batch_id"`
	TaskQueue       string `json:"task_queue"`

	// For workflow continuation
	ParentWorkflowStateID string `json:"parent_workflow_state_id"`

	// New fields for pagination
	Offset            int32 `json:"offset"`             // Current position in pagination
	PageSize          int32 `json:"page_size"`          // Number of reports to process per page
	TotalToProcess    int32 `json:"total_to_process"`   // Total number of reports to process
	AlreadyProcessed  int32 `json:"already_processed"`  // Number of reports already processed
	ContinuationCount int32 `json:"continuation_count"` // Number of times the workflow has been continued
}

// WorkflowConfig represents the root configuration structure
// Legacy config
// TODO: Remove this when the new workflow struct is fully implemented
type WorkflowConfig struct {
	Rankings RankingsConfig `json:"rankings" yaml:"rankings"`
	Worker   WorkerConfig   `json:"worker" yaml:"worker"`
	Specs    []ClassSpec    `json:"specs" yaml:"specs"`
	Dungeons []Dungeon      `json:"dungeons" yaml:"dungeons"`
}

// RankingsConfig contains settings for rankings processing
type RankingsConfig struct {
	MaxRankingsPerSpec int32         `json:"max_rankings_per_spec" yaml:"max_rankings_per_spec"`
	UpdateInterval     time.Duration `json:"update_interval" yaml:"update_interval"`
	Batch              BatchConfig   `json:"batch" yaml:"batch"`
}

// WorkerConfig defines worker pool settings
type WorkerConfig struct {
	NumWorkers   int32         `json:"num_workers" yaml:"num_workers"`
	RequestDelay time.Duration `json:"request_delay" yaml:"request_delay"`
}

// ClassSpec represents a WoW class specialization
type ClassSpec struct {
	ClassName string `json:"class_name" yaml:"class_name"`
	SpecName  string `json:"spec_name" yaml:"spec_name"`
}

// Dungeon represents a Mythic+ dungeon
type Dungeon struct {
	ID          uint32 `json:"id" yaml:"id"`
	EncounterID uint32 `json:"encounter_id" yaml:"encounter_id"`
	Name        string `json:"name" yaml:"name"`
	Slug        string `json:"slug" yaml:"slug"`
}
