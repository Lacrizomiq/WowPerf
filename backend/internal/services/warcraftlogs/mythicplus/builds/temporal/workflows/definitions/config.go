// config.go
package warcraftlogsBuildsTemporalWorkflowsDefinitions

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

/* This file :

- Handle the load of the config
- Give a default config
- Include config validation

*/

// === NEW FUNCTIONS FOR DECOUPLED WORKFLOWS ===

// LoadRankingsParams loads the parameters for the rankings workflow
func LoadRankingsParams(configPath string) (*models.RankingsWorkflowParams, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &models.RankingsWorkflowParams{
		Specs:              config.Specs,
		Dungeons:           config.Dungeons,
		MaxRankingsPerSpec: config.Rankings.MaxRankingsPerSpec,
		BatchSize:          config.Rankings.Batch.Size,
		RetryDelay:         config.Rankings.Batch.RetryDelay,
		MaxAttempts:        config.Rankings.Batch.MaxAttempts,
		BatchID:            fmt.Sprintf("rankings-%s", uuid.New().String()),
	}, nil
}

// LoadReportsParams loads the parameters for the reports workflow
func LoadReportsParams(configPath string) (*models.ReportsWorkflowParams, error) {
	// Check if the file exists
	_, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Fully independent configuration
	return &models.ReportsWorkflowParams{
		BatchSize:        10,                     // Optimized value for reports
		NumWorkers:       2,                      // Optimized value for reports
		RequestDelay:     500 * time.Millisecond, // Optimized value for reports
		ProcessingWindow: 7 * 24 * time.Hour,     // 7 days
		BatchID:          fmt.Sprintf("reports-%s", uuid.New().String()),
	}, nil
}

// LoadBuildsParams loads the parameters for the builds workflow
func LoadBuildsParams(configPath string) (*models.BuildsWorkflowParams, error) {
	// Check if the file exists
	_, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Fully independent configuration with optimized hardcoded values
	return &models.BuildsWorkflowParams{
		BatchSize:       5,  // Optimized value specific to builds
		NumWorkers:      4,  // Optimized value specific to builds
		ReportBatchSize: 10, // Optimized value specific to builds
		BatchID:         fmt.Sprintf("builds-%s", uuid.New().String()),
	}, nil
}

// LoadEquipmentAnalysisParams loads the parameters for the equipment analysis workflow
func LoadEquipmentAnalysisParams(configPath string) (*models.EquipmentAnalysisWorkflowParams, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &models.EquipmentAnalysisWorkflowParams{
		Spec:          config.Specs,    // Specs to analyze
		Dungeon:       config.Dungeons, // Dungeons to analyze
		BatchSize:     10,              // Batch size for the analysis
		Concurrency:   4,               // Number of concurrent workers
		RetryAttempts: 3,               // Number of retry attempts
		RetryDelay:    5 * time.Second, // Retry delay
	}, nil
}

// LoadStatAnalysisParams loads the parameters for the stat analysis workflow
func LoadStatAnalysisParams(configPath string) (*models.StatAnalysisWorkflowParams, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &models.StatAnalysisWorkflowParams{
		Spec:          config.Specs,    // Specs to analyze
		Dungeon:       config.Dungeons, // Dungeons to analyze
		BatchSize:     10,              // Batch size for the analysis
		Concurrency:   4,               // Number of concurrent workers
		RetryAttempts: 3,               // Number of retry attempts
		RetryDelay:    5 * time.Second, // Retry delay
	}, nil
}

// LoadTalentAnalysisParams loads the parameters for the talent analysis workflow
func LoadTalentAnalysisParams(configPath string) (*models.TalentAnalysisWorkflowParams, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &models.TalentAnalysisWorkflowParams{
		Spec:          config.Specs,    // Specs to analyze
		Dungeon:       config.Dungeons, // Dungeons to analyze
		BatchSize:     10,              // Batch size for the analysis
		Concurrency:   4,               // Number of concurrent workers
		RetryAttempts: 3,               // Number of retry attempts
		RetryDelay:    5 * time.Second, // Retry delay
	}, nil
}

// === LEGACY FUNCTIONS ===

// LoadConfig loads configuration from file or returns default values
// TODO: Remove this when the new workflow struct is fully implemented
func LoadConfig(configPath string) (*models.WorkflowConfig, error) {
	// If no config path provided, use default config
	if configPath == "" {
		return GetDefaultConfig(), nil
	}

	// Read configuration file
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &models.WorkflowConfig{}
	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Validate the loaded configuration
	if err := ValidateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetDefaultConfig returns default configuration values
// TODO: Remove this when the new workflow struct is fully implemented
func GetDefaultConfig() *models.WorkflowConfig {
	return &models.WorkflowConfig{
		Rankings: models.RankingsConfig{
			MaxRankingsPerSpec: 150,
			Batch: models.BatchConfig{
				Size:        5,
				RetryDelay:  5,
				MaxAttempts: 3,
			},
		},
		Worker: models.WorkerConfig{
			NumWorkers:   3,
			RequestDelay: 500,
		},
	}
}

// ValidateConfig performs validation on the configuration
// TODO: Remove this when the new workflow struct is fully implemented
func ValidateConfig(config *models.WorkflowConfig) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if config.Rankings.MaxRankingsPerSpec <= 0 {
		return fmt.Errorf("max rankings per spec must be greater than 0")
	}

	if config.Worker.NumWorkers <= 0 {
		return fmt.Errorf("number of workers must be at least 1")
	}

	if len(config.Specs) == 0 {
		return fmt.Errorf("at least one spec must be configured")
	}

	if len(config.Dungeons) == 0 {
		return fmt.Errorf("at least one dungeon must be configured")
	}

	return nil
}
