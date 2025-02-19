// config.go
package warcraftlogsBuildsTemporalWorkflowsModels

import "time"

// WorkflowConfig represents the root configuration structure
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
