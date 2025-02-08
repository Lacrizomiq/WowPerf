package warcraftlogsBuildsTemporalWorkflows

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/workflow"
	"gopkg.in/yaml.v2"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

// ClassSpec represents a WoW class specialization
type ClassSpec struct {
	ClassName string `json:"class_name" yaml:"class_name"`
	SpecName  string `json:"spec_name" yaml:"spec_name"`
}

// Dungeon represents a Mythic+ dungeon
type Dungeon struct {
	ID          uint   `json:"id" yaml:"id"`
	EncounterID uint   `json:"encounter_id" yaml:"encounter_id"`
	Name        string `json:"name" yaml:"name"`
	Slug        string `json:"slug" yaml:"slug"`
}

// BatchConfig defines parameters for batch processing
type BatchConfig struct {
	Size        uint          `json:"size" yaml:"size"`
	RetryDelay  time.Duration `json:"retry_delay" yaml:"retry_delay"`
	MaxAttempts int           `json:"max_attempts" yaml:"max_attempts"`
}

// RankingsConfig contains settings for rankings processing
type RankingsConfig struct {
	MaxRankingsPerSpec int           `json:"max_rankings_per_spec" yaml:"max_rankings_per_spec"`
	UpdateInterval     time.Duration `json:"update_interval" yaml:"update_interval"`
	Batch              BatchConfig   `json:"batch" yaml:"batch"`
}

// WorkerConfig defines worker pool settings
type WorkerConfig struct {
	NumWorkers   int           `json:"num_workers" yaml:"num_workers"`
	RequestDelay time.Duration `json:"request_delay" yaml:"request_delay"`
}

// Config represents the root configuration structure
type Config struct {
	Rankings RankingsConfig `json:"rankings" yaml:"rankings"`
	Worker   WorkerConfig   `json:"worker" yaml:"worker"`
	Specs    []ClassSpec    `json:"specs" yaml:"specs"`
	Dungeons []Dungeon      `json:"dungeons" yaml:"dungeons"`
}

// WorkflowParams contains parameters needed for workflow execution
type WorkflowParams struct {
	// Config contains the complete configuration
	Config *Config `json:"config"`
	// Progress contains the progress of the workflow
	Progress *WorkflowProgress
	// WorkflowID is the ID of the workflow
	WorkflowID string `json:"workflow_id"`
}

// BatchResult represents the outcome of a rankings batch processing
type BatchResult struct {
	ClassName   string
	SpecName    string
	EncounterID uint
	Rankings    []*warcraftlogsBuilds.ClassRanking
	ProcessedAt time.Time
}

// ReportProcessingResult contains statistics about report processing
type ReportProcessingResult struct {
	ProcessedReports []*warcraftlogsBuilds.Report
	ProcessedCount   int
	SuccessCount     int
	FailureCount     int
	ProcessedAt      time.Time
}

// BuildsProcessingResult contains statistics about build processing
type BuildsProcessingResult struct {
	ProcessedBuilds int
	SuccessCount    int
	FailureCount    int
	ProcessedAt     time.Time
}

// WorkflowResult represents the final outcome of a workflow execution
type WorkflowResult struct {
	RankingsProcessed int
	ReportsProcessed  int
	BuildsProcessed   int
	StartedAt         time.Time
	CompletedAt       time.Time
}

// WorkflowProgress tracks the progress of spec and dungeon processing
type WorkflowProgress struct {
	CompletedSpecs      map[string]bool // Map of completed specs
	CompletedDungeons   map[string]bool // Map of completed dungeons
	CurrentSpecIndex    int             // Index of the current spec being processed
	CurrentDungeonIndex int             // Index of the current dungeon being processed
	PartialResults      *WorkflowResult // Accumulated results
	Stats               *ProgressStats
}

// ProgressStats contains more stats about the progress
type ProgressStats struct {
	TotalSpecs        int
	TotalDungeons     int
	ProcessedSpecs    int
	ProcessedDungeons int
	CurrentSpec       *ClassSpec
	CurrentDungeon    *Dungeon
	StartedAt         time.Time
	LastSpecUpdate    time.Time
	LastDungeonUpdate time.Time
}

// QuotaExceededError is an error type for quota exceeded errors
type QuotaExceededError struct {
	Message string
	ResetIn time.Duration
}

// LoadConfig loads configuration from file or returns default values
func LoadConfig(configPath string) (*Config, error) {
	logger := log.New(os.Stdout, "[CONFIG] ", log.LstdFlags)

	// Default configuration values
	defaultConfig := &Config{
		Rankings: RankingsConfig{
			MaxRankingsPerSpec: 150,
			UpdateInterval:     7 * 24 * time.Hour,
			Batch: BatchConfig{
				Size:        5,
				RetryDelay:  5 * time.Second,
				MaxAttempts: 3,
			},
		},
		Worker: WorkerConfig{
			NumWorkers:   3,
			RequestDelay: 500 * time.Millisecond,
		},
	}

	if configPath == "" {
		logger.Printf("No config path provided, using default configuration")
		return defaultConfig, nil
	}

	logger.Printf("Loading configuration from: %s", configPath)
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := defaultConfig
	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Log loaded configuration details
	logger.Printf("Configuration loaded successfully:")
	logger.Printf("- Specs loaded: %d", len(config.Specs))
	for _, spec := range config.Specs {
		logger.Printf("  * %s-%s", spec.ClassName, spec.SpecName)
	}
	logger.Printf("- Dungeons loaded: %d", len(config.Dungeons))
	for _, dungeon := range config.Dungeons {
		logger.Printf("  * %s (ID: %d)", dungeon.Name, dungeon.ID)
	}
	logger.Printf("- Rankings config: max=%d, interval=%v",
		config.Rankings.MaxRankingsPerSpec,
		config.Rankings.UpdateInterval)
	logger.Printf("- Worker config: workers=%d, delay=%v",
		config.Worker.NumWorkers,
		config.Worker.RequestDelay)

	// Detailed validation
	if err := validateConfigDetails(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// validateConfig checks if the configuration is valid
func validateConfig(config *Config) error {
	if config.Worker.NumWorkers < 1 {
		return fmt.Errorf("number of workers must be at least 1")
	}
	if config.Rankings.Batch.MaxAttempts < 0 {
		return fmt.Errorf("retry attempts cannot be negative")
	}
	if config.Rankings.MaxRankingsPerSpec <= 0 {
		return fmt.Errorf("max rankings per spec must be greater than 0")
	}
	if len(config.Specs) == 0 {
		return fmt.Errorf("at least one spec must be configured")
	}
	if len(config.Dungeons) == 0 {
		return fmt.Errorf("at least one dungeon must be configured")
	}
	return nil
}

func validateConfigDetails(config *Config) error {
	// Base validation
	if err := validateConfig(config); err != nil {
		return err
	}

	// Validate specs uniqueness and format
	specMap := make(map[string]bool)
	for _, spec := range config.Specs {
		// Validate spec format
		if spec.ClassName == "" {
			return fmt.Errorf("empty class name found in specs")
		}
		if spec.SpecName == "" {
			return fmt.Errorf("empty spec name found for class %s", spec.ClassName)
		}

		// Check for duplicates
		key := fmt.Sprintf("%s-%s", spec.ClassName, spec.SpecName)
		if specMap[key] {
			return fmt.Errorf("duplicate spec found: %s", key)
		}
		specMap[key] = true
	}

	// Validate dungeons uniqueness and format
	dungeonMap := make(map[uint]bool)
	for _, dungeon := range config.Dungeons {
		// Validate dungeon format
		if dungeon.Name == "" {
			return fmt.Errorf("empty dungeon name found for ID %d", dungeon.ID)
		}
		if dungeon.ID == 0 {
			return fmt.Errorf("invalid dungeon ID found for %s", dungeon.Name)
		}
		if dungeon.Slug == "" {
			return fmt.Errorf("empty slug found for dungeon %s", dungeon.Name)
		}

		// Check for duplicates
		if dungeonMap[dungeon.ID] {
			return fmt.Errorf("duplicate dungeon ID found: %d (%s)", dungeon.ID, dungeon.Name)
		}
		dungeonMap[dungeon.ID] = true
	}

	// Validate worker configuration
	if config.Worker.RequestDelay < 100*time.Millisecond {
		return fmt.Errorf("request delay too low (minimum 100ms)")
	}
	if config.Worker.NumWorkers < 1 {
		return fmt.Errorf("at least one worker required")
	}
	if config.Worker.NumWorkers > 10 {
		return fmt.Errorf("too many workers (maximum 10)")
	}

	return nil
}

// Activity name constants - matching exactly with the activity methods
const (
	// Rankings activities
	FetchRankingsActivity     = "FetchAndStore"
	GetStoredRankingsActivity = "GetStoredRankings"

	// Reports activities
	ProcessReportsActivity      = "ProcessReports"
	GetProcessedReportsActivity = "GetProcessedReports"
	GetReportsBatchActivity     = "GetReportsBatch"
	CountAllReportsActivity     = "CountAllReports"

	// Player builds activities
	ProcessBuildsActivity     = "ProcessAllBuilds"
	CountPlayerBuildsActivity = "CountPlayerBuilds"

	// Rate limit activities
	ReserveRateLimitPointsActivity = "ReservePoints"
	ReleaseRateLimitPointsActivity = "ReleasePoints"
	CheckRemainingPointsActivity   = "CheckRemainingPoints"

	// Sub-workflow names
	ProcessBuildBatchWorkflowName = "ProcessBuildBatch"
	SyncWorkflowName              = "SyncWorkflow"
)

// RankingsSyncWorkflow defines the interface for rankings synchronization workflow
type RankingsSyncWorkflow interface {
	Execute(ctx workflow.Context, params WorkflowParams) (*WorkflowResult, error)
}

// ProcessBuildBatchWorkflow defines the interface for build batch processing workflow
type ProcessBuildBatchWorkflow interface {
	Execute(ctx workflow.Context, params BuildBatchParams) (*BuildBatchResult, error)
}

// RankingsActivity defines the interface for rankings-related activities
type RankingsActivity interface {
	FetchAndStore(ctx context.Context, spec ClassSpec, dungeon Dungeon, batchConfig BatchConfig) (*BatchResult, error)
	GetStoredRankings(ctx context.Context, className, specName string, encounterID uint) ([]*warcraftlogsBuilds.ClassRanking, error)
}

// ReportsActivity defines the interface for report-related activities
type ReportsActivity interface {
	ProcessReports(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking) (*ReportProcessingResult, error)
	GetReportsBatch(ctx context.Context, batchSize int, offset int) ([]*warcraftlogsBuilds.Report, error)
	CountAllReports(ctx context.Context) (int64, error)
}

// PlayerBuildsActivity defines the interface for player build-related activities
type PlayerBuildsActivity interface {
	ProcessAllBuilds(ctx context.Context, reports []*warcraftlogsBuilds.Report) (*BuildsProcessingResult, error)
	CountPlayerBuilds(ctx context.Context) (int64, error)
}

// RateLimitActivity defines the interface for rate limiting operations
type RateLimitActivity interface {
	ReservePoints(ctx context.Context, params WorkflowParams) error
	ReleasePoints(ctx context.Context, params WorkflowParams) error
	CheckRemainingPoints(ctx context.Context, params WorkflowParams) (float64, error)
}
