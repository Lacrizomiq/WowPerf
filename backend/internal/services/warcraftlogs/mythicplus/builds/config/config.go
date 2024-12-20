// config.go
package warcraftlogsBuildsConfig

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// DefaultConfig is the default configuration
var defaultConfig = &Config{
	Rankings: RankingsConfig{
		MaxRankingsPerSpec: 150,
		UpdateInterval:     7 * 24 * time.Hour,
		Batch: BatchConfig{
			Size:        100,
			MaxPages:    2,
			RetryDelay:  5 * time.Second,
			MaxAttempts: 3,
		},
	},
	Worker: WorkerConfig{
		NumWorkers:   3,
		RequestDelay: 500 * time.Millisecond,
	},
	Specs: []ClassSpec{
		{ClassName: "Priest", SpecName: "Discipline"},
	},
	Dungeons: []Dungeon{
		{ID: 12660, EncounterID: 12660, Name: "Ara-Kara", Slug: "arakara-city-of-echoes"},
		{ID: 12669, EncounterID: 12669, Name: "City of Threads", Slug: "city-of-threads"},
		{ID: 60670, EncounterID: 60670, Name: "Grim Batol", Slug: "grim-batol"},
		{ID: 62290, EncounterID: 62290, Name: "Mists of Tirna Scithe", Slug: "mists-of-tirna-scithe"},
		{ID: 61822, EncounterID: 61822, Name: "Siege of Boralus", Slug: "siege-of-boralus"},
		{ID: 12662, EncounterID: 12662, Name: "The Dawnbreaker", Slug: "the-dawnbreaker"},
		{ID: 62286, EncounterID: 62286, Name: "The Necrotic Wake", Slug: "the-necrotic-wake"},
		{ID: 12652, EncounterID: 12652, Name: "The Stonevault", Slug: "the-stonevault"},
	},
}

// Load loads the configuration from the given file or uses the default configuration
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		return defaultConfig, nil
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := defaultConfig
	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	if config.Worker.NumWorkers < 1 {
		return fmt.Errorf("number of workers must be at least 1")
	}
	if config.Rankings.Batch.MaxAttempts < 0 {
		return fmt.Errorf("retry attempts cannot be negative")
	}
	if config.Rankings.MaxRankingsPerSpec < 100 {
		return fmt.Errorf("max rankings per spec must be at least 100")
	}
	if len(config.Specs) == 0 {
		return fmt.Errorf("at least one spec must be configured")
	}
	if len(config.Dungeons) == 0 {
		return fmt.Errorf("at least one dungeon must be configured")
	}
	if config.Rankings.Batch.Size < 1 {
		return fmt.Errorf("batch size must be at least 1")
	}
	if config.Rankings.Batch.MaxPages < 1 {
		return fmt.Errorf("max pages must be at least 1")
	}
	return nil
}
