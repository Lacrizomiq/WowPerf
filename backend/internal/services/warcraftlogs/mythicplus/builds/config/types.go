// types.go
package warcraftlogsBuildsConfig

import "time"

// ClassSpec represents a class and spec
type ClassSpec struct {
	ClassName string `json:"class_name" yaml:"class_name"`
	SpecName  string `json:"spec_name" yaml:"spec_name"`
}

// Dungeon represents a dungeon
type Dungeon struct {
	ID          uint   `json:"id" yaml:"id"`
	EncounterID uint   `json:"encounter_id" yaml:"encounter_id"`
	Name        string `json:"name" yaml:"name"`
	Slug        string `json:"slug" yaml:"slug"`
}

// BatchConfig represents the batch processing configuration
type BatchConfig struct {
	Size        uint          `json:"size" yaml:"size"`
	MaxPages    uint          `json:"max_pages" yaml:"max_pages"`
	RetryDelay  time.Duration `json:"retry_delay" yaml:"retry_delay"`
	MaxAttempts int           `json:"max_attempts" yaml:"max_attempts"`
}

// RankingsConfig represents the rankings configuration
type RankingsConfig struct {
	MaxRankingsPerSpec int           `json:"max_rankings_per_spec" yaml:"max_rankings_per_spec"`
	UpdateInterval     time.Duration `json:"update_interval" yaml:"update_interval"`
	Batch              BatchConfig   `json:"batch" yaml:"batch"`
}

// WorkerConfig represents the worker pool configuration
type WorkerConfig struct {
	NumWorkers   int           `json:"num_workers" yaml:"num_workers"`
	RequestDelay time.Duration `json:"request_delay" yaml:"request_delay"`
}

// Config represents the overall configuration
type Config struct {
	Rankings RankingsConfig `json:"rankings" yaml:"rankings"`
	Worker   WorkerConfig   `json:"worker" yaml:"worker"`
	Specs    []ClassSpec    `json:"specs" yaml:"specs"`
	Dungeons []Dungeon      `json:"dungeons" yaml:"dungeons"`
}
