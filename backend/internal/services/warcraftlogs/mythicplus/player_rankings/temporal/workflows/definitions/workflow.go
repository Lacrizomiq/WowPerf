// workflow.go
package warcraftlogsPlayerRankingsTemporalWorkflowsDefinitions

import (
	models "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/models"

	"go.temporal.io/sdk/workflow"
)

// PlayerRankingsWorkflow définit l'interface pour le workflow de gestion des classements des joueurs
// Il orchestre la récupération des classements, leur stockage et le calcul des métriques.
type PlayerRankingsWorkflow interface {
	// Execute exécute le workflow avec les paramètres spécifiés et retourne le résultat
	Execute(ctx workflow.Context, params models.PlayerRankingWorkflowParams) (*models.WorkflowResult, error)
}
