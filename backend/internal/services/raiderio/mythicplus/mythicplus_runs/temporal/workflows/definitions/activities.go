package raiderioMythicPlusRunsDefinitions

import (
	"context"

	models "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/models"
)

// Noms constants des activities - correspondant exactement aux méthodes d'activity
const (
	// Activity principale de récupération et traitement des runs M+
	FetchAndProcessMythicPlusRunsActivity = "FetchAndProcessMythicPlusRunsActivity" // Récupère et traite les runs pour toutes les combinaisons

	// Nom du workflow principal
	MythicPlusRunsWorkflowName = "MythicPlusRunsWorkflow" // Workflow de gestion des runs M+
)

// MythicPlusRunsActivity définit l'interface pour les activities liées aux runs M+
type MythicPlusRunsActivity interface {
	// Récupère et traite les runs M+ pour toutes les combinaisons (Season, Region, Dungeon)
	// Orchestre Query + Repository et retourne les statistiques globales
	FetchAndProcessMythicPlusRunsActivity(
		ctx context.Context,
		params models.MythicRunsWorkflowParams,
	) (*models.RunsProcessingStats, error)
}
