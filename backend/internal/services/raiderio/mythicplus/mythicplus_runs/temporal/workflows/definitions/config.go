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

	// Valider les éléments essentiels
	if len(config.Seasons) == 0 {
		return nil, fmt.Errorf("at least one season must be configured")
	}
	if len(config.Regions) == 0 {
		return nil, fmt.Errorf("at least one region must be configured")
	}
	if len(config.Dungeons) == 0 {
		return nil, fmt.Errorf("at least one dungeon must be configured")
	}

	// Valeurs par défaut pour les paramètres d'exécution
	pagesPerDungeon := config.Execution.PagesPerDungeon
	if pagesPerDungeon == 0 {
		pagesPerDungeon = 5 // Par défaut : 5 pages par donjon
	}

	maxConcurrency := config.Execution.MaxConcurrency
	if maxConcurrency == 0 {
		maxConcurrency = 3 // Par défaut : 3 workers parallèles
	}

	retryAttempts := config.Execution.RetryAttempts
	if retryAttempts == 0 {
		retryAttempts = 3 // Par défaut : 3 tentatives
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
