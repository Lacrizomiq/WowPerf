// player_rankings_workflow.go
package warcraftlogsPlayerRankingsWorkflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/models"
)

// PlayerRankingsWorkflow implémente le workflow principal pour les classements des joueurs
type PlayerRankingsWorkflow struct{}

// NewPlayerRankingsWorkflow crée une nouvelle instance du workflow
func NewPlayerRankingsWorkflow() *PlayerRankingsWorkflow {
	return &PlayerRankingsWorkflow{}
}

// Execute exécute le workflow avec les paramètres fournis
func (w *PlayerRankingsWorkflow) Execute(ctx workflow.Context, params models.PlayerRankingWorkflowParams) (*models.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting player rankings workflow",
		"dungeonCount", len(params.Dungeons),
		"specCount", len(params.Specs),
		"maxConcurrency", params.MaxConcurrency)

	// Préparer le résultat
	result := &models.WorkflowResult{
		StartTime:         workflow.Now(ctx),
		BatchID:           params.BatchID,
		DungeonsProcessed: len(params.Dungeons),
		SpecsProcessed:    len(params.Specs),
	}

	// Extraire les IDs des donjons
	dungeonIDs := make([]int, 0, len(params.Dungeons))
	for _, dungeon := range params.Dungeons {
		dungeonIDs = append(dungeonIDs, dungeon.ID)
	}

	// Configurer les options des activities avec retry
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * time.Duration(params.RetryDelay),
			BackoffCoefficient: 2.0,
			MaximumInterval:    5 * time.Minute,
			MaximumAttempts:    int32(params.RetryAttempts),
		},
		HeartbeatTimeout: 2 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// 1. Récupérer les classements pour tous les donjons
	logger.Info("Starting rankings fetch")
	fetchStart := workflow.Now(ctx)

	var rankings []playerRankingModels.PlayerRanking
	err := workflow.ExecuteActivity(
		ctx,
		definitions.FetchAllDungeonRankingsActivity,
		dungeonIDs,
		params.PagesPerDungeon,
		params.MaxConcurrency,
	).Get(ctx, &rankings)

	if err != nil {
		logger.Error("Failed to fetch rankings", "error", err)
		result.Error = err.Error()
		result.EndTime = workflow.Now(ctx)
		result.TotalDuration = result.EndTime.Sub(result.StartTime)
		return result, err
	}

	result.FetchDuration = workflow.Now(ctx).Sub(fetchStart)
	result.RankingsCount = len(rankings)
	logger.Info("Fetch completed", "rankingsCount", result.RankingsCount, "duration", result.FetchDuration)

	// 2. Stocker les classements en base de données
	logger.Info("Starting rankings storage")
	storeStart := workflow.Now(ctx)

	err = workflow.ExecuteActivity(
		ctx,
		definitions.StoreRankingsActivity,
		rankings,
	).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to store rankings", "error", err)
		result.Error = err.Error()
		result.EndTime = workflow.Now(ctx)
		result.TotalDuration = result.EndTime.Sub(result.StartTime)
		return result, err
	}

	result.StoreDuration = workflow.Now(ctx).Sub(storeStart)
	logger.Info("Storage completed", "duration", result.StoreDuration)

	// 3. Calculer les métriques quotidiennes
	logger.Info("Starting daily metrics calculation")
	metricsStart := workflow.Now(ctx)

	err = workflow.ExecuteActivity(
		ctx,
		definitions.CalculateDailyMetricsActivity,
	).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to calculate daily metrics", "error", err)
		result.Error = err.Error()
		result.EndTime = workflow.Now(ctx)
		result.TotalDuration = result.EndTime.Sub(result.StartTime)
		return result, err
	}

	result.MetricDuration = workflow.Now(ctx).Sub(metricsStart)
	logger.Info("Metrics calculation completed", "duration", result.MetricDuration)

	// 4. Tenter de récupérer les statistiques globales (optionnel, ne bloque pas le workflow)
	var globalRankings *playerRankingModels.GlobalRankings
	err = workflow.ExecuteActivity(
		ctx,
		"GetGlobalRankings", // Activité supplémentaire à implémenter si nécessaire
	).Get(ctx, &globalRankings)

	if err == nil && globalRankings != nil {
		result.GlobalRankings = globalRankings
		result.TankCount = globalRankings.Tanks.Count
		result.HealerCount = globalRankings.Healers.Count
		result.DPSCount = globalRankings.DPS.Count
		logger.Info("Retrieved global rankings",
			"tanks", result.TankCount,
			"healers", result.HealerCount,
			"dps", result.DPSCount)
	}

	// Finaliser le résultat
	result.EndTime = workflow.Now(ctx)
	result.TotalDuration = result.EndTime.Sub(result.StartTime)

	logger.Info("Player rankings workflow completed successfully",
		"totalDuration", result.TotalDuration,
		"rankingsCount", result.RankingsCount)

	return result, nil
}
