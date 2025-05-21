package warcraftlogsPlayerRankingsActivities

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"

	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"
	warcraftlogs "wowperf/internal/services/warcraftlogs"
	playerRankingsQueries "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/queries"
	playerRankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/repository"
	models "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/models"
)

// PlayerRankingsActivity handles all activities related to player rankings
// Implémente l'interface definitions.PlayerRankingsActivity
type PlayerRankingsActivity struct {
	client     *warcraftlogs.WarcraftLogsClientService
	repository *playerRankingsRepository.PlayerRankingsRepository
}

// NewPlayerRankingsActivity creates a new player rankings activity handler
func NewPlayerRankingsActivity(
	client *warcraftlogs.WarcraftLogsClientService,
	repository *playerRankingsRepository.PlayerRankingsRepository,
) *PlayerRankingsActivity {
	return &PlayerRankingsActivity{
		client:     client,
		repository: repository,
	}
}

// FetchAllDungeonRankings récupère les classements de plusieurs donjons en parallèle
// Elle stocke directement les résultats en base de données et retourne uniquement des statistiques
func (a *PlayerRankingsActivity) FetchAllDungeonRankings(
	ctx context.Context,
	dungeonIDs []int,
	pagesPerDungeon int,
	maxConcurrency int,
) (*models.RankingsStats, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting rankings fetch for multiple dungeons",
		"dungeonCount", len(dungeonIDs),
		"pagesPerDungeon", pagesPerDungeon,
		"maxConcurrency", maxConcurrency)

	startTime := time.Now()
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)

	var mu sync.Mutex
	var allRankings []playerRankingModels.PlayerRanking
	errorsChan := make(chan error, len(dungeonIDs))

	// Compteurs pour les statistiques de rôles
	var tankCount, healerCount, dpsCount int
	var statsMutex sync.Mutex

	for _, dungeonID := range dungeonIDs {
		wg.Add(1)
		go func(dID int) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			// Record heartbeat
			activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing dungeon %d", dID))

			// Process each page
			var dungeonRankings []playerRankingModels.PlayerRanking

			for page := 1; page <= pagesPerDungeon; page++ {
				// Prepare query parameters
				params := playerRankingsQueries.LeaderboardParams{
					EncounterID: dID,
					Page:        page,
				}

				// Fetch rankings
				dungeonData, err := playerRankingsQueries.GetDungeonLeaderboardByPlayer(a.client, params)
				if err != nil {
					errorsChan <- fmt.Errorf("failed to fetch dungeon %d page %d: %w", dID, page, err)
					return
				}

				// Compter les rôles pour chaque classement
				localTankCount := 0
				localHealerCount := 0
				localDpsCount := 0

				// Convert to PlayerRanking objects
				for _, ranking := range dungeonData.Rankings {
					// Determine role based on class and spec
					role := determineRole(ranking.Class, ranking.Spec)

					// Track role counts
					switch role {
					case "Tank":
						localTankCount++
					case "Healer":
						localHealerCount++
					case "DPS":
						localDpsCount++
					}

					// Create PlayerRanking object
					playerRanking := playerRankingModels.PlayerRanking{
						DungeonID:       dID,
						Name:            ranking.Name,
						Class:           ranking.Class,
						Spec:            ranking.Spec,
						Role:            role,
						Amount:          ranking.Amount,
						HardModeLevel:   ranking.HardModeLevel,
						Duration:        ranking.Duration,
						StartTime:       ranking.StartTime,
						ReportCode:      ranking.Report.Code,
						ReportFightID:   ranking.Report.FightID,
						ReportStartTime: ranking.Report.StartTime,
						GuildID:         ranking.Guild.ID,
						GuildName:       ranking.Guild.Name,
						GuildFaction:    ranking.Guild.Faction,
						ServerID:        ranking.Server.ID,
						ServerName:      ranking.Server.Name,
						ServerRegion:    ranking.Server.Region,
						BracketData:     ranking.BracketData,
						Faction:         ranking.Faction,
						Affixes:         ranking.Affixes,
						Medal:           ranking.Medal,
						Score:           ranking.Score,
						Leaderboard:     0,
					}

					dungeonRankings = append(dungeonRankings, playerRanking)
				}

				// Mettre à jour les compteurs globaux
				statsMutex.Lock()
				tankCount += localTankCount
				healerCount += localHealerCount
				dpsCount += localDpsCount
				statsMutex.Unlock()

				// If there are no more pages, break the loop
				if !dungeonData.HasMorePages {
					break
				}
			}

			// Add results to global array with mutex protection
			mu.Lock()
			allRankings = append(allRankings, dungeonRankings...)
			mu.Unlock()

			logger.Info("Completed dungeon rankings fetch",
				"dungeonID", dID,
				"rankingsCount", len(dungeonRankings))
		}(dungeonID)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorsChan)

	// Process errors
	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return nil, temporal.NewApplicationError(
			fmt.Sprintf("Failed to fetch rankings for some dungeons: %v", errors),
			"FETCH_ERROR",
		)
	}

	logger.Info("Successfully fetched rankings for all dungeons",
		"totalRankingsCount", len(allRankings))

	// Stockage direct en base de données
	if len(allRankings) > 0 {
		// D'abord supprimer les classements existants
		if err := a.repository.DeleteExistingRankings(ctx); err != nil {
			logger.Error("Failed to delete existing rankings", "error", err)
			return nil, temporal.NewApplicationError(
				fmt.Sprintf("Failed to delete existing rankings: %v", err),
				"DB_ERROR",
			)
		}

		// Stocker les nouveaux classements
		if err := a.repository.StoreRankingsByBatches(ctx, allRankings); err != nil {
			logger.Error("Failed to store rankings", "error", err)
			return nil, temporal.NewApplicationError(
				fmt.Sprintf("Failed to store rankings: %v", err),
				"DB_ERROR",
			)
		}

		logger.Info("Successfully stored rankings directly in database", "count", len(allRankings))
	}

	// Créer et retourner les statistiques
	stats := &models.RankingsStats{
		TotalCount:        len(allRankings),
		DungeonsProcessed: len(dungeonIDs),
		ProcessingTime:    time.Since(startTime),
		TankCount:         tankCount,
		HealerCount:       healerCount,
		DPSCount:          dpsCount,
	}

	return stats, nil
}

// StoreRankings stores the rankings in the database
// It first deletes existing rankings and then stores new ones in batches
// Corresponds to definitions.StoreRankingsActivity
func (a *PlayerRankingsActivity) StoreRankings(
	ctx context.Context,
	rankings []playerRankingModels.PlayerRanking,
) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting rankings storage", "rankingsCount", len(rankings))

	// Delete existing rankings
	if err := a.repository.DeleteExistingRankings(ctx); err != nil {
		logger.Error("Failed to delete existing rankings", "error", err)
		return temporal.NewApplicationError(
			fmt.Sprintf("Failed to delete existing rankings: %v", err),
			"DB_ERROR",
		)
	}

	// Store new rankings
	if err := a.repository.StoreRankingsByBatches(ctx, rankings); err != nil {
		logger.Error("Failed to store rankings", "error", err)
		return temporal.NewApplicationError(
			fmt.Sprintf("Failed to store rankings: %v", err),
			"DB_ERROR",
		)
	}

	logger.Info("Successfully stored rankings", "count", len(rankings))
	return nil
}

// CalculateDailyMetrics calculates daily metrics for specializations
// Corresponds to definitions.CalculateDailyMetricsActivity
func (a *PlayerRankingsActivity) CalculateDailyMetrics(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting daily metrics calculation")

	startTime := time.Now()

	// Calculate metrics
	if err := a.repository.CalculateDailySpecMetrics(ctx); err != nil {
		logger.Error("Failed to calculate daily metrics", "error", err)
		return temporal.NewApplicationError(
			fmt.Sprintf("Failed to calculate daily metrics: %v", err),
			"METRICS_ERROR",
		)
	}

	duration := time.Since(startTime)
	logger.Info("Successfully calculated daily metrics", "duration", duration)
	return nil
}

// determineRole determines the role (Tank, Healer, DPS) based on class and spec
func determineRole(class, spec string) string {
	tanks := map[string][]string{
		"Warrior":     {"Protection"},
		"Paladin":     {"Protection"},
		"DeathKnight": {"Blood"},
		"DemonHunter": {"Vengeance"},
		"Druid":       {"Guardian"},
		"Monk":        {"Brewmaster"},
	}

	healers := map[string][]string{
		"Priest":  {"Holy", "Discipline"},
		"Paladin": {"Holy"},
		"Druid":   {"Restoration"},
		"Shaman":  {"Restoration"},
		"Monk":    {"Mistweaver"},
		"Evoker":  {"Preservation"},
	}

	if specs, ok := tanks[class]; ok {
		for _, tankSpec := range specs {
			if spec == tankSpec {
				return "Tank"
			}
		}
	}

	if specs, ok := healers[class]; ok {
		for _, healSpec := range specs {
			if spec == healSpec {
				return "Healer"
			}
		}
	}

	return "DPS"
}
