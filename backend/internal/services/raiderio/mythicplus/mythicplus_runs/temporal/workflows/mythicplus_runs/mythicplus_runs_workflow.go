// mythicplus_runs_workflow.go
package raiderioMythicPlusRunsWorkflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	definitions "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/definitions"
	models "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/models"
)

// MythicPlusRunsWorkflow implémente le workflow principal pour les runs M+
type MythicPlusRunsWorkflow struct{}

// NewMythicPlusRunsWorkflow crée une nouvelle instance du workflow
func NewMythicPlusRunsWorkflow() *MythicPlusRunsWorkflow {
	return &MythicPlusRunsWorkflow{}
}

// Execute exécute le workflow avec les paramètres fournis
func (w *MythicPlusRunsWorkflow) Execute(
	ctx workflow.Context,
	params models.MythicRunsWorkflowParams,
) (*models.MythicRunsWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting mythic+ runs workflow",
		"seasonsCount", len(params.Seasons),
		"regionsCount", len(params.Regions),
		"dungeonsCount", len(params.Dungeons),
		"pagesPerDungeon", params.PagesPerDungeon,
		"maxConcurrency", params.MaxConcurrency,
		"batchID", params.BatchID)

	// Préparer le résultat
	result := &models.MythicRunsWorkflowResult{
		StartTime:         workflow.Now(ctx),
		BatchID:           params.BatchID,
		Success:           false, // Sera mis à true à la fin
		RegionsProcessed:  len(params.Regions),
		DungeonsProcessed: len(params.Dungeons),
		RegionStats:       make(map[string]models.RegionStats),
	}

	// Calcul du nombre total de combinaisons
	totalCombinations := len(params.Seasons) * len(params.Regions) * len(params.Dungeons)
	logger.Info("Workflow will process combinations", "totalCombinations", totalCombinations)

	// Configurer les options des activities avec retry
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 6 * time.Hour, // Long timeout pour traiter tous les donjons
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    30 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Minute,
			MaximumAttempts:    int32(params.RetryAttempts),
		},
		HeartbeatTimeout: 5 * time.Minute, // Heartbeat fréquent pour monitoring
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// 1. Exécuter l'activity principale de fetch et processing
	logger.Info("Starting mythic+ runs fetch and processing")
	fetchStart := workflow.Now(ctx)

	var processingStats *models.RunsProcessingStats
	err := workflow.ExecuteActivity(
		ctx,
		definitions.FetchAndProcessMythicPlusRunsActivity,
		params,
	).Get(ctx, &processingStats)

	if err != nil {
		logger.Error("Failed to fetch and process mythic+ runs", "error", err)
		result.Error = err.Error()
		result.EndTime = workflow.Now(ctx)
		result.TotalDuration = result.EndTime.Sub(result.StartTime)
		return result, err
	}

	result.FetchDuration = workflow.Now(ctx).Sub(fetchStart)
	logger.Info("Fetch and processing completed",
		"duration", result.FetchDuration,
		"totalFetched", processingStats.TotalFetched,
		"totalStored", processingStats.TotalStored,
		"totalSkipped", processingStats.TotalSkipped)

	// 2. Mapper les résultats vers le résultat final
	result.TotalRunsFetched = processingStats.TotalFetched
	result.TotalRunsStored = processingStats.TotalStored
	result.TotalRunsUpdated = processingStats.TotalUpdated
	result.TotalRunsSkipped = processingStats.TotalSkipped
	result.RegionStats = processingStats.RegionStats

	// 3. Calcul des compositions d'équipe (optionnel, peut être fait plus tard via SQL)
	logger.Info("Starting team compositions build")
	compositionsStart := workflow.Now(ctx)

	// Pour l'instant, on simule juste le calcul - tu feras ça en SQL
	// Tu peux ajouter une activity ici plus tard si besoin
	result.BuildCompositionsDuration = workflow.Now(ctx).Sub(compositionsStart)
	result.TotalTeamCompositions = 0 // À calculer via SQL

	// 4. Build roster (optionnel, peut être fait plus tard via SQL)
	logger.Info("Starting roster build")
	rosterStart := workflow.Now(ctx)

	// Pour l'instant, on simule juste le calcul - tu feras ça en SQL
	result.BuildRosterDuration = workflow.Now(ctx).Sub(rosterStart)
	result.TotalRosterEntries = result.TotalRunsStored * 5 // Estimation : 5 membres par run

	// 5. Finaliser le résultat
	result.EndTime = workflow.Now(ctx)
	result.TotalDuration = result.EndTime.Sub(result.StartTime)
	result.Success = true

	logger.Info("Mythic+ runs workflow completed successfully",
		"totalDuration", result.TotalDuration,
		"totalRunsStored", result.TotalRunsStored,
		"totalRunsSkipped", result.TotalRunsSkipped,
		"regionsProcessed", result.RegionsProcessed,
		"dungeonsProcessed", result.DungeonsProcessed)

	return result, nil
}

// GetWorkflowName retourne le nom du workflow pour l'enregistrement
func (w *MythicPlusRunsWorkflow) GetWorkflowName() string {
	return definitions.MythicPlusRunsWorkflowName
}
