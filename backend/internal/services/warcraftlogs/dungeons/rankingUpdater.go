package warcraftlogs

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"
	service "wowperf/internal/services/warcraftlogs"
	middleware "wowperf/middleware/cache"
	cache "wowperf/pkg/cache"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	MinimumUpdateInterval = 12 * time.Hour
	DefaultUpdateInterval = 24 * time.Hour
	updateLockKey         = "warcraftlogs:rankings:update:lock"
	batchSize             = 100
)

// RankingsUpdater is responsible for updating the rankings in the database
type RankingsUpdater struct {
	db           *gorm.DB
	service      *service.WarcraftLogsClientService
	cache        cache.CacheService
	cacheManager *middleware.CacheManager
}

func NewRankingsUpdater(db *gorm.DB, service *service.WarcraftLogsClientService, cache cache.CacheService, cacheManager *middleware.CacheManager) *RankingsUpdater {
	return &RankingsUpdater{
		db:           db,
		service:      service,
		cache:        cache,
		cacheManager: cacheManager,
	}
}

// StartPeriodicUpdate starts the periodic updates
func (r *RankingsUpdater) StartPeriodicUpdate(ctx context.Context) {
	log.Println("Starting WarcraftLogs rankings periodic update...")

	// Performing an initial check immediately
	log.Println("Performing initial check...")
	if err := r.checkAndUpdate(ctx); err != nil {
		log.Printf("Initial check error: %v", err)
	}

	ticker := time.NewTicker(MinimumUpdateInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("Ticker triggered, checking for updates...")
				if err := r.checkAndUpdate(ctx); err != nil {
					log.Printf("Periodic check error: %v", err)
				}
			}
		}
	}()
}

// initializeUpdateState initializes the update state if necessary
func (r *RankingsUpdater) initializeUpdateState() error {
	var state playerRankingModels.RankingsUpdateState
	result := r.db.First(&state)

	if result.Error == gorm.ErrRecordNotFound {
		state = playerRankingModels.RankingsUpdateState{
			LastUpdateTime: time.Now(),
		}
		if err := r.db.Create(&state).Error; err != nil {
			return fmt.Errorf("failed to create initial state: %w", err)
		}
		log.Println("Created initial rankings update state")
	} else if result.Error != nil {
		return fmt.Errorf("failed to check update state: %w", result.Error)
	}

	return nil
}

func (r *RankingsUpdater) updateLastUpdateTime(ctx context.Context) error {
	return r.db.Exec(`
			INSERT INTO rankings_update_states (id, last_update_time, updated_at)
			VALUES (1, NOW(), NOW())
			ON CONFLICT (id) DO UPDATE 
			SET last_update_time = EXCLUDED.last_update_time,
					updated_at = EXCLUDED.updated_at
	`).Error
}

// checkAndUpdate checks and performs the update if necessary
func (r *RankingsUpdater) checkAndUpdate(ctx context.Context) error {
	// Checking the distributed lock
	locked, err := r.cache.SetNX(ctx, updateLockKey, time.Now().String(), 1*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to check update lock: %w", err)
	}
	if !locked {
		return fmt.Errorf("update already in progress")
	}
	defer r.cache.Delete(ctx, updateLockKey)

	// Getting the unique state
	state, err := playerRankingModels.GetOrCreateRankingsUpdateState(r.db)
	if err != nil {
		return fmt.Errorf("failed to get rankings state: %w", err)
	}

	timeSinceLastUpdate := time.Since(state.LastUpdateTime)
	log.Printf("Time since last rankings update: %v", timeSinceLastUpdate)

	if timeSinceLastUpdate >= MinimumUpdateInterval {
		log.Printf("Rankings update needed, last update was %v ago", timeSinceLastUpdate)

		if err := r.UpdateRankings(ctx); err != nil {
			return fmt.Errorf("failed to update rankings: %w", err)
		}

		if err := r.updateLastUpdateTime(ctx); err != nil {
			return fmt.Errorf("failed to update timestamp: %w", err)
		}

		log.Println("Rankings update completed successfully")
	} else {
		log.Printf("Skipping update: last update was %v ago", timeSinceLastUpdate)
	}

	return nil
}

// UpdateRankings updates the rankings in the database
func (r *RankingsUpdater) UpdateRankings(ctx context.Context) error {
	dungeonIDs := []int{
		DungeonAraKara,
		DungeonCityOfThreads,
		DungeonGrimBatol,
		DungeonMists,
		DungeonSiege,
		DungeonDawnbreaker,
		DungeonNecroticWake,
		DungeonStonevault,
	}
	pagesPerDungeon := 1

	rankings, err := GetGlobalRankings(r.service, ctx, dungeonIDs, pagesPerDungeon)
	if err != nil {
		return fmt.Errorf("failed to get global rankings: %w", err)
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Deleting the old data
		if err := tx.Exec("DELETE FROM player_rankings").Error; err != nil {
			return fmt.Errorf("failed to delete existing rankings: %w", err)
		}

		// Preparing the new data
		var newRankings []playerRankingModels.PlayerRanking
		processRoleRankings := func(players []PlayerScore) {
			for _, player := range players {
				for _, run := range player.Runs {
					newRankings = append(newRankings, playerRankingModels.PlayerRanking{
						DungeonID:       run.DungeonID,
						Name:            player.Name,
						Class:           player.Class,
						Spec:            player.Spec,
						Role:            player.Role,
						Amount:          player.Amount,
						HardModeLevel:   run.HardModeLevel,
						Duration:        run.Duration,
						StartTime:       run.StartTime,
						ReportCode:      run.Report.Code,
						ReportFightID:   run.Report.FightID,
						ReportStartTime: run.Report.StartTime,
						GuildID:         player.Guild.ID,
						GuildName:       player.Guild.Name,
						GuildFaction:    player.Guild.Faction,
						ServerID:        player.Server.ID,
						ServerName:      player.Server.Name,
						ServerRegion:    player.Server.Region,
						BracketData:     run.BracketData,
						Faction:         player.Faction,
						Affixes:         run.Affixes,
						Medal:           run.Medal,
						Score:           run.Score,
						Leaderboard:     0,
					})
				}
			}
		}

		processRoleRankings(rankings.Tanks.Players)
		processRoleRankings(rankings.Healers.Players)
		processRoleRankings(rankings.DPS.Players)

		// Inserting the new data by batches
		return r.insertRankingsInBatches(tx, newRankings)
	})
}

// insertRankingsInBatches inserts the rankings by batches
func (r *RankingsUpdater) insertRankingsInBatches(tx *gorm.DB, rankings []playerRankingModels.PlayerRanking) error {
	baseSQL := `
    INSERT INTO player_rankings (
        created_at, updated_at, dungeon_id, name, class, spec, role, 
        amount, hard_mode_level, duration, start_time, report_code, 
        report_fight_id, report_start_time, guild_id, guild_name, 
        guild_faction, server_id, server_name, server_region, 
        bracket_data, faction, affixes, medal, score, leaderboard
    ) VALUES `

	totalRankings := len(rankings)
	batches := int(math.Ceil(float64(totalRankings) / float64(batchSize)))

	for i := 0; i < batches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalRankings {
			end = totalRankings
		}

		batch := rankings[start:end]
		valueStrings := make([]string, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*26)

		for j, ranking := range batch {
			valueStrings[j] = fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				j*26+1, j*26+2, j*26+3, j*26+4, j*26+5, j*26+6, j*26+7, j*26+8, j*26+9, j*26+10,
				j*26+11, j*26+12, j*26+13, j*26+14, j*26+15, j*26+16, j*26+17, j*26+18, j*26+19, j*26+20,
				j*26+21, j*26+22, j*26+23, j*26+24, j*26+25, j*26+26,
			)

			now := time.Now()
			valueArgs = append(valueArgs,
				now, // created_at
				now, // updated_at
				ranking.DungeonID,
				ranking.Name,
				ranking.Class,
				ranking.Spec,
				ranking.Role,
				ranking.Amount,
				ranking.HardModeLevel,
				ranking.Duration,
				ranking.StartTime,
				ranking.ReportCode,
				ranking.ReportFightID,
				ranking.ReportStartTime,
				ranking.GuildID,
				ranking.GuildName,
				ranking.GuildFaction,
				ranking.ServerID,
				ranking.ServerName,
				ranking.ServerRegion,
				ranking.BracketData,
				ranking.Faction,
				pq.Array(ranking.Affixes),
				ranking.Medal,
				ranking.Score,
				ranking.Leaderboard,
			)
		}

		query := baseSQL + strings.Join(valueStrings, ",")
		if err := tx.Exec(query, valueArgs...).Error; err != nil {
			return fmt.Errorf("failed to insert batch %d: %w", i+1, err)
		}

		log.Printf("Inserted batch %d/%d (%d rankings)", i+1, batches, len(batch))
	}

	return nil
}

// invalidateCache invalidates the caches related to rankings
func (r *RankingsUpdater) invalidateCache(ctx context.Context) error {
	tags := []string{
		"rankings",
		"leaderboard",
		"global-rankings",
		"class-rankings",
		"spec-rankings",
	}
	return r.cacheManager.InvalidateByTags(ctx, tags)
}
