package warcraftlogsBuildsTemporal

import (
	"time"
	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

// WorkflowParams contains the parameters for the rankings workflow
type WorkflowParams struct {
	// Config of specs and dungeons
	Specs    []ClassSpec `json:"specs"`
	Dungeons []Dungeon   `json:"dungeons"`

	// Batch processing config
	BatchConfig BatchConfig `json:"batch"`

	// Rankings config
	Rankings struct {
		MaxRankingsPerSpec int           `json:"max_rankings_per_spec"`
		UpdateInterval     time.Duration `json:"update_interval"`
	} `json:"rankings"`
}

// ClassSpec and Dungeon correspond to my YAML configuration
type ClassSpec struct {
	ClassName string `json:"class_name"`
	SpecName  string `json:"spec_name"`
}

type Dungeon struct {
	ID          uint   `json:"id"`
	EncounterID uint   `json:"encounter_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
}

// BatchResult represents the result of a batch of rankings
type BatchResult struct {
	ClassName   string
	SpecName    string
	EncounterID uint
	Rankings    []*warcraftlogsBuilds.ClassRanking
	HasMore     bool
	ProcessedAt time.Time
}

// BatchConfig represents the batch processing configuration
type BatchConfig struct {
	Size        uint          `json:"size"`
	MaxPages    uint          `json:"max_pages"`
	RetryDelay  time.Duration `json:"retry_delay"`
	MaxAttempts int           `json:"max_attempts"`
}

// ReportProcessingResult represents the result of the report processing
type ReportProcessingResult struct {
	ProcessedReports int
	SuccessCount     int
	FailureCount     int
	ProcessedAt      time.Time
}

// BuildsProcessingResult represents the result of the builds processing
type BuildsProcessingResult struct {
	ProcessedBuilds int
	SuccessCount    int
	FailureCount    int
	ProcessedAt     time.Time
}

// WorkflowResult represents the result of the workflow
type WorkflowResult struct {
	RankingsProcessed int
	ReportsProcessed  int
	BuildsProcessed   int
	StartedAt         time.Time
	CompletedAt       time.Time
}

// Activity names
const (
	FetchRankingsActivityName   = "fetch-rankings"
	ProcessReportsActivityName  = "process-reports"
	ProcessPlayerBuildsActivity = "process-player-builds"
)

// RankingsSyncWorkflow is the workflow for synchronizing rankings
type RankingsSyncWorkflow interface {
	Execute(params WorkflowParams) (WorkflowResult, error)
}

// Interface of the activities
type RankingsActivity interface {
	FetchAndStore(spec ClassSpec, dungeon Dungeon, batchConfig BatchConfig) (*BatchResult, error)
}

type ReportsActivity interface {
	ProcessReports(rankings []*warcraftlogsBuilds.ClassRanking) (*ReportProcessingResult, error)
}

type PlayerBuildsActivity interface {
	ProcessBuilds(reportCodes []string) (*BuildsProcessingResult, error)
}
