package raiderioMythicPlusRunsModels

import (
	"time"
)

// MythicRunsWorkflowResult contient les résultats et métriques d'exécution du workflow
type MythicRunsWorkflowResult struct {
	// Métriques temporelles
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	TotalDuration time.Duration `json:"total_duration"`

	// Métriques par étape
	FetchDuration             time.Duration `json:"fetch_duration"`              // Durée récupération API
	BuildCompositionsDuration time.Duration `json:"build_compositions_duration"` // Durée build compositions
	BuildRosterDuration       time.Duration `json:"build_roster_duration"`       // Durée build roster

	// Métriques de processing
	TotalRunsFetched      int `json:"total_runs_fetched"`      // Runs récupérés de l'API
	TotalRunsStored       int `json:"total_runs_stored"`       // Runs stockés en DB (nouveaux)
	TotalRunsUpdated      int `json:"total_runs_updated"`      // Runs mis à jour (meilleur score)
	TotalRunsSkipped      int `json:"total_runs_skipped"`      // Runs ignorés (score inférieur)
	TotalTeamCompositions int `json:"total_team_compositions"` // Compositions créées/utilisées
	TotalRosterEntries    int `json:"total_roster_entries"`    // Entrées dans run_roster

	RegionsProcessed  int `json:"regions_processed"`  // Nombre de régions traitées
	DungeonsProcessed int `json:"dungeons_processed"` // Nombre de donjons traités

	// Détail par région (optionnel pour debug)
	RegionStats map[string]RegionStats `json:"region_stats,omitempty"`

	// Identifiant et statut
	BatchID string `json:"batch_id"`        // ID unique du batch
	Success bool   `json:"success"`         // Workflow réussi ou non
	Error   string `json:"error,omitempty"` // Message d'erreur si échec
}

// RegionStats détaille les métriques par région
type RegionStats struct {
	RunsFetched int           `json:"runs_fetched"`
	RunsStored  int           `json:"runs_stored"`
	Duration    time.Duration `json:"duration"`
}

// RunsProcessingStats contient uniquement les stats de processing (pour transfer entre activities)
type RunsProcessingStats struct {
	TotalFetched int                    `json:"total_fetched"`
	TotalStored  int                    `json:"total_stored"`
	TotalUpdated int                    `json:"total_updated"`
	TotalSkipped int                    `json:"total_skipped"`
	Duration     time.Duration          `json:"duration"`
	RegionStats  map[string]RegionStats `json:"region_stats,omitempty"`
}

// CompositionBuildStats contient les stats de construction des compositions
type CompositionBuildStats struct {
	NewCompositions      int           `json:"new_compositions"`      // Nouvelles compositions créées
	ExistingCompositions int           `json:"existing_compositions"` // Compositions existantes réutilisées
	Duration             time.Duration `json:"duration"`
}

// RosterBuildStats contient les stats de construction du roster
type RosterBuildStats struct {
	RosterEntries int           `json:"roster_entries"` // Entrées créées dans run_roster
	Duration      time.Duration `json:"duration"`
}
