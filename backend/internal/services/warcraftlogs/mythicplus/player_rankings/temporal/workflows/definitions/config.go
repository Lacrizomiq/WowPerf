// config.go
package warcraftlogsPlayerRankingsTemporalWorkflowsDefinitions

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	models "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/models"
)

// LoadPlayerRankingsParams charge les paramètres du workflow depuis le fichier de configuration
// et les combine avec des valeurs par défaut pour les paramètres non spécifiés
func LoadPlayerRankingsParams(configPath string) (*models.PlayerRankingWorkflowParams, error) {
	// Lire le fichier de configuration
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Structure pour charger la configuration
	var config struct {
		Specs    []models.ClassSpec `yaml:"specs"`
		Dungeons []models.Dungeon   `yaml:"dungeons"`
	}

	// Parser le fichier YAML
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Valider les éléments essentiels
	if len(config.Specs) == 0 {
		return nil, fmt.Errorf("at least one spec must be configured")
	}
	if len(config.Dungeons) == 0 {
		return nil, fmt.Errorf("at least one dungeon must be configured")
	}

	// Créer et retourner les paramètres complets du workflow
	return &models.PlayerRankingWorkflowParams{
		Specs:           config.Specs,
		Dungeons:        config.Dungeons,
		PagesPerDungeon: 2,
		MaxConcurrency:  3,
		RetryAttempts:   3,
		RetryDelay:      5 * time.Second,
		BatchID:         fmt.Sprintf("player-rankings-%s", uuid.New().String()),
	}, nil
}
