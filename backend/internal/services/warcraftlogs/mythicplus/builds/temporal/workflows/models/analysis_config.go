package warcraftlogsBuildsTemporalWorkflowsModels

import "time"

// == Decoupled workflows ==

// EquipmentAnalysisWorkflowParams contains the parameters for the equipment analysis workflow
// It defines the input configuration for analyzing player equipment.
type EquipmentAnalysisWorkflowParams struct {
	Spec          []ClassSpec   `json:"spec"`           // ClassSpec is a struct that contains the class name and spec name
	Dungeon       []Dungeon     `json:"dungeon"`        // Dungeon is a struct that contains the dungeon name and encounter ID
	BatchSize     int32         `json:"batch_size"`     // Batch size for processing
	Concurrency   int32         `json:"concurrency"`    // Number of parallel processes
	RetryAttempts int32         `json:"retry_attempts"` // Number of retries in case of failure
	RetryDelay    time.Duration `json:"retry_delay"`    // Delay between retries
	BatchID       string        `json:"batch_id"`       // Batch ID for the workflow
}

// TalentAnalysisWorkflowParams contains the parameters for the talent analysis workflow
// It defines the input configuration for analyzing player talent builds.

type TalentAnalysisWorkflowParams struct {
	Spec          []ClassSpec   `json:"spec"`           // ClassSpec is a struct that contains the class name and spec name
	Dungeon       []Dungeon     `json:"dungeon"`        // Dungeon is a struct that contains the dungeon name and encounter ID
	BatchSize     int32         `json:"batch_size"`     // Batch size for processing
	Concurrency   int32         `json:"concurrency"`    // Number of parallel processes
	RetryAttempts int32         `json:"retry_attempts"` // Number of retries in case of failure
	RetryDelay    time.Duration `json:"retry_delay"`    // Delay between retries
	BatchID       string        `json:"batch_id"`       // Batch ID for the workflow
}

// StatAnalysisWorkflowParams contains the parameters for the stats analysis workflow
// It defines the input configuration for analyzing player stats distribution.
type StatAnalysisWorkflowParams struct {
	Spec          []ClassSpec   `json:"spec"`           // ClassSpec is a struct that contains the class name and spec name
	Dungeon       []Dungeon     `json:"dungeon"`        // Dungeon is a struct that contains the dungeon name and encounter ID
	BatchSize     int32         `json:"batch_size"`     // Batch size for processing
	Concurrency   int32         `json:"concurrency"`    // Number of parallel processes
	RetryAttempts int32         `json:"retry_attempts"` // Number of retries in case of failure
	RetryDelay    time.Duration `json:"retry_delay"`    // Delay between retries
	BatchID       string        `json:"batch_id"`       // Batch ID for the workflow
}

// == Legacy workflows ==

// AnalysisWorkflowConfig contains the specific parameters for the analysis workflow
// This is the main configuration for the analysis workflow
// Note: This will not be used anymore.
// It is kept here for reference and in case we need to use it again.
type AnalysisWorkflowConfig struct {
	Specs         []ClassSpec   `json:"specs"`          // Reuse the ClassSpec structure
	Dungeons      []Dungeon     `json:"dungeons"`       // Reuse the Dungeon structure
	BatchSize     int32         `json:"batch_size"`     // Batch size for processing
	Concurrency   int32         `json:"concurrency"`    // Number of parallel processes
	RetryAttempts int32         `json:"retry_attempts"` // Number of retries in case of failure
	RetryDelay    time.Duration `json:"retry_delay"`    // Delay between retries
}
