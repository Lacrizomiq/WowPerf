// config.go
package warcraftlogsBuildsTemporalWorkflowsDefinitions

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

/* This file :

- Handle the load of the config
- Give a default config
- Include config validation

*/

// LoadConfig loads configuration from file or returns default values
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
