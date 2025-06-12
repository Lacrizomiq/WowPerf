// models/config.go
package raiderioMythicPlusRunsModels

type MythicRunsWorkflowParams struct {
	// Configuration des données
	Seasons  []Season  `json:"seasons" yaml:"seasons"`
	Regions  []string  `json:"regions" yaml:"regions"`
	Dungeons []Dungeon `json:"dungeons" yaml:"dungeons"`

	// Paramètres d'exécution (du YAML ou defaults)
	PagesPerDungeon int `json:"pages_per_dungeon" yaml:"pages_per_dungeon"`
	MaxConcurrency  int `json:"max_concurrency" yaml:"max_concurrency"`
	RetryAttempts   int `json:"retry_attempts" yaml:"retry_attempts"`

	// Runtime
	BatchID string `json:"batch_id"`
}

type Season struct {
	ID   string `json:"season_id" yaml:"season_id"`
	Name string `json:"name" yaml:"name"`
}

type Dungeon struct {
	Slug string `json:"slug" yaml:"slug"`
	Name string `json:"name" yaml:"name"`
}
