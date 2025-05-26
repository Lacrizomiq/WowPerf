// definitions/config.go
package raiderioMythicPlusRunsDefinitions

import (
	"fmt"
	"os"

	models "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/models"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

func LoadMythicRunsParams(configPath string) (*models.MythicRunsWorkflowParams, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config struct {
		Seasons   []models.Season  `yaml:"seasons"`
		Regions   []string         `yaml:"regions"`
		Dungeons  []models.Dungeon `yaml:"dungeons"`
		Execution struct {
			PagesPerDungeon int `yaml:"pages_per_dungeon"`
			MaxConcurrency  int `yaml:"max_concurrency"`
			RetryAttempts   int `yaml:"retry_attempts"`
		} `yaml:"execution"`
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Validation + defaults
	if len(config.Dungeons) == 0 {
		return nil, fmt.Errorf("at least one dungeon must be configured")
	}

	return &models.MythicRunsWorkflowParams{
		Seasons:         config.Seasons,
		Regions:         config.Regions,
		Dungeons:        config.Dungeons,
		PagesPerDungeon: config.Execution.PagesPerDungeon,
		MaxConcurrency:  config.Execution.MaxConcurrency,
		RetryAttempts:   config.Execution.RetryAttempts,
		BatchID:         fmt.Sprintf("mythic-runs-%s", uuid.New().String()),
	}, nil
}
