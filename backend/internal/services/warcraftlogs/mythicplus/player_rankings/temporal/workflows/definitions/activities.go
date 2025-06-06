// activities.go
package warcraftlogsPlayerRankingsTemporalWorkflowsDefinitions

import (
	"context"

	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"
)

// Noms constants des activities - correspondant exactement aux méthodes d'activity
const (
	// Activities de récupération et stockage des classements
	FetchAllDungeonRankingsActivity = "FetchAllDungeonRankings" // Récupère les classements pour tous les donjons
	StoreRankingsActivity           = "StoreRankings"           // Stocke les classements en base de données
	CalculateDailyMetricsActivity   = "CalculateDailyMetrics"   // Calcule les métriques quotidiennes

	// Nom du workflow principal
	PlayerRankingsWorkflowName = "PlayerRankingsWorkflow" // Workflow de gestion des classements
)

// PlayerRankingsActivity définit l'interface pour les activities liées aux classements
type PlayerRankingsActivity interface {
	// Récupère les classements pour plusieurs donjons en parallèle
	FetchAllDungeonRankings(
		ctx context.Context,
		dungeonIDs []int,
		pagesPerDungeon int,
		maxConcurrency int,
	) ([]playerRankingModels.PlayerRanking, error)

	// Stocke les classements en base de données
	StoreRankings(
		ctx context.Context,
		rankings []playerRankingModels.PlayerRanking,
	) error

	// Calcule les métriques quotidiennes pour toutes les spécialisations
	CalculateDailyMetrics(
		ctx context.Context,
	) error
}
