package raiderioMythicPlusRunsActivities

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.temporal.io/sdk/activity"

	raiderioService "wowperf/internal/services/raiderio"
	queries "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/queries"
	repository "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/repository"
	models "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/models"
)

// workUnit represente une unité de travail pour le traitement (Season, Region, Dungeon)
type workUnit struct {
	season  models.Season
	region  string
	dungeon models.Dungeon
}

// MythicPlusRunsActivity gère toutes les activités liées aux runs M+
type MythicPlusRunsActivity struct {
	client     *raiderioService.RaiderIOService
	repository *repository.MythicPlusRunsRepository
}

// NewMythicPlusRunsActivity crée un nouveau gestionnaire d'activités pour les runs M+
func NewMythicPlusRunsActivity(
	client *raiderioService.RaiderIOService,
	repo *repository.MythicPlusRunsRepository,
) *MythicPlusRunsActivity {
	return &MythicPlusRunsActivity{
		client:     client,
		repository: repo,
	}
}

// FetchAndProcessMythicPlusRunsActivity
func (a *MythicPlusRunsActivity) FetchAndProcessMythicPlusRunsActivity(
	ctx context.Context,
	params models.MythicRunsWorkflowParams,
) (*models.RunsProcessingStats, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting mythic+ runs fetch and processing",
		"seasonsCount", len(params.Seasons),
		"regionsCount", len(params.Regions),
		"dungeonsCount", len(params.Dungeons),
		"pagesPerDungeon", params.PagesPerDungeon,
		"maxConcurrency", params.MaxConcurrency,
		"batchID", params.BatchID)

	startTime := time.Now()

	// Crée toutes les combinaisons (Season, Region, Dungeon)
	var workUnits []workUnit
	for _, season := range params.Seasons {
		for _, region := range params.Regions {
			for _, dungeon := range params.Dungeons {
				workUnits = append(workUnits, workUnit{
					season:  season,
					region:  region,
					dungeon: dungeon,
				})
			}
		}
	}

	logger.Info("Created work units", "totalUnits", len(workUnits))

	// Worker pool pour paralléliser le traitement
	var wg sync.WaitGroup
	sem := make(chan struct{}, params.MaxConcurrency)

	// Stats globales thread-safe
	var mu sync.Mutex
	globalStats := &models.RunsProcessingStats{
		RegionStats: make(map[string]models.RegionStats),
	}

	// Process chaque combinaison - SANS gestion d'erreur globale
	for i, unit := range workUnits {
		wg.Add(1)
		go func(workIndex int, work workUnit) {
			defer wg.Done()

			sem <- struct{}{}        // Acquérir le sémaphore
			defer func() { <-sem }() // Libérer le sémaphore

			// Enregistrer le heartbeat avec la progression
			activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing %s/%s/%s (%d/%d)",
				work.season.ID, work.region, work.dungeon.Slug, workIndex+1, len(workUnits)))

			// Process cette combinaison - SI ERREUR, Temporal va retry automatiquement
			unitStats, err := a.processWorkUnit(ctx, work, params)
			if err != nil {
				// Log error seulement - pas de collection d'erreurs
				logger.Warn("Failed to process work unit - Temporal will retry",
					"season", work.season.ID,
					"region", work.region,
					"dungeon", work.dungeon.Slug,
					"error", err)
				return // Exit cette goroutine, Temporal va retry l'activity entière
			}

			// Fusionner les stats thread-safe
			mu.Lock()
			globalStats.TotalFetched += unitStats.TotalFetched
			globalStats.TotalStored += unitStats.TotalStored
			globalStats.TotalUpdated += unitStats.TotalUpdated
			globalStats.TotalSkipped += unitStats.TotalSkipped

			// Aggreger les stats de la région
			if regionStat, exists := globalStats.RegionStats[work.region]; exists {
				regionStat.RunsFetched += unitStats.TotalFetched
				regionStat.RunsStored += unitStats.TotalStored
				regionStat.Duration += unitStats.Duration
				globalStats.RegionStats[work.region] = regionStat
			} else {
				globalStats.RegionStats[work.region] = models.RegionStats{
					RunsFetched: unitStats.TotalFetched,
					RunsStored:  unitStats.TotalStored,
					Duration:    unitStats.Duration,
				}
			}
			mu.Unlock()

			logger.Info("Completed work unit",
				"season", work.season.ID,
				"region", work.region,
				"dungeon", work.dungeon.Slug,
				"fetched", unitStats.TotalFetched,
				"stored", unitStats.TotalStored,
				"skipped", unitStats.TotalSkipped)

		}(i, unit)
	}

	// Attendre que tous les travailleurs soient terminés
	wg.Wait()

	globalStats.Duration = time.Since(startTime)

	// Log final stats - SANS vérification d'erreur
	logger.Info("Mythic+ runs processing completed",
		"totalDuration", globalStats.Duration,
		"totalFetched", globalStats.TotalFetched,
		"totalStored", globalStats.TotalStored,
		"totalSkipped", globalStats.TotalSkipped)

	return globalStats, nil
}

// processWorkUnit - Si erreur ici, Temporal retry automatiquement
func (a *MythicPlusRunsActivity) processWorkUnit(
	ctx context.Context,
	unit workUnit,
	params models.MythicRunsWorkflowParams,
) (*models.RunsProcessingStats, error) {
	startTime := time.Now()
	stats := &models.RunsProcessingStats{}

	// Traite toutes les pages pour cette combinaison
	for page := 1; page <= params.PagesPerDungeon; page++ {
		queryParams := queries.MythicPlusRunsParams{
			Season:  unit.season.ID,
			Region:  unit.region,
			Dungeon: unit.dungeon.Slug,
			Page:    page,
		}

		// 1. Fetch data via Query - SI ERREUR, Temporal retry
		runs, err := queries.GetMythicPlusRuns(a.client, queryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch runs (page %d): %w", page, err)
		}

		// Check context cancellation between API calls
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if len(runs) == 0 {
			// Plus de données sur cette page, stop la pagination
			break
		}

		stats.TotalFetched += len(runs)

		// 2. Process via Repository - SI ERREUR, Temporal retry
		repoStats, err := a.repository.ProcessRuns(runs, params.BatchID)
		if err != nil {
			return nil, fmt.Errorf("failed to process runs (page %d): %w", page, err)
		}

		// 3. Aggregate stats
		stats.TotalStored += repoStats.NewRuns
		stats.TotalUpdated += 0 // Pas d'updates dans notre logique
		stats.TotalSkipped += repoStats.SkippedRuns
	}

	stats.Duration = time.Since(startTime)
	return stats, nil
}
