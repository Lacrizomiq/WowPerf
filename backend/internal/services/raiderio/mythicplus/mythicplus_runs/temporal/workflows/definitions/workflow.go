package raiderioMythicPlusRunsDefinitions

import (
	models "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/models"

	"go.temporal.io/sdk/workflow"
)

// MythicPlusRunsWorkflow définit l'interface pour le workflow de gestion des runs M+
// Il orchestre la récupération des runs depuis l'API Raider.IO, leur traitement et stockage en base de données.
type MythicPlusRunsWorkflow interface {
	// Execute exécute le workflow avec les paramètres spécifiés et retourne le résultat
	// Traite toutes les combinaisons (Season, Region, Dungeon) avec retry automatique
	Execute(ctx workflow.Context, params models.MythicRunsWorkflowParams) (*models.MythicRunsWorkflowResult, error)
}
