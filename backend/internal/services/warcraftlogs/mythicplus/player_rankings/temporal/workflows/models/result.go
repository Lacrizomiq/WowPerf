// result.go
package warcraftlogsPlayerRankingsTemporalWorkflowsModels

import (
	"time"

	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"
)

// WorkflowResult contient les résultats et statistiques d'exécution du workflow
type WorkflowResult struct {
	// Métriques générales
	DungeonsProcessed int `json:"dungeons_processed"` // Nombre de donjons traités
	RankingsCount     int `json:"rankings_count"`     // Nombre total de classements récupérés
	SpecsProcessed    int `json:"specs_processed"`    // Nombre de spécialisations traitées

	// Métriques temporelles
	StartTime      time.Time     `json:"start_time"`      // Début de l'exécution
	EndTime        time.Time     `json:"end_time"`        // Fin de l'exécution
	TotalDuration  time.Duration `json:"total_duration"`  // Durée totale d'exécution
	FetchDuration  time.Duration `json:"fetch_duration"`  // Durée de récupération des classements
	StoreDuration  time.Duration `json:"store_duration"`  // Durée de stockage des classements
	MetricDuration time.Duration `json:"metric_duration"` // Durée de calcul des métriques

	// Statistiques par rôle
	TankCount   int `json:"tank_count"`   // Nombre de tanks
	HealerCount int `json:"healer_count"` // Nombre de healers
	DPSCount    int `json:"dps_count"`    // Nombre de DPS

	// Identifiant d'exécution et statut
	BatchID string `json:"batch_id"`        // ID de batch unique
	Error   string `json:"error,omitempty"` // Message d'erreur si échec

	// Données de récapitulation optionnelles
	GlobalRankings *playerRankingModels.GlobalRankings `json:"global_rankings,omitempty"`
}
