package warcraftlogsPlayerRankingsTemporalWorkflowsModels

import "time"

// PlayerRankingWorkflowParams contient les paramètres du workflow player_rankings
type PlayerRankingWorkflowParams struct {
	// Paramètres des donjons et classes
	Dungeons []Dungeon   `json:"dungeons"`
	Specs    []ClassSpec `json:"specs"`

	// Paramètres de l'exécution
	PagesPerDungeon int `json:"pages_per_dungeon"`
	MaxConcurrency  int `json:"max_concurrency"`

	// Paramètres de retry
	RetryAttempts int           `json:"retry_attempts"`
	RetryDelay    time.Duration `json:"retry_delay"`

	// Identifiant unique pour cette exécution
	BatchID string `json:"batch_id"`
}

// ClassSpec représente une spécialisation d'une classe WoW
type ClassSpec struct {
	ClassName string `json:"class_name" yaml:"class_name"`
	SpecName  string `json:"spec_name" yaml:"spec_name"`
}

// Dungeon représente un donjon Mythic+
type Dungeon struct {
	ID          int    `json:"id" yaml:"id"`
	EncounterID int    `json:"encounter_id" yaml:"encounter_id"`
	Name        string `json:"name" yaml:"name"`
	Slug        string `json:"slug" yaml:"slug"`
}
