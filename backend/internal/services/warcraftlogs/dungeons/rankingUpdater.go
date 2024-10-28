package warcraftlogs

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"
	"wowperf/pkg/cache"

	"github.com/lib/pq"

	"gorm.io/gorm"
)

const (
	updateInterval = 24 * time.Hour
	batchSize      = 1000 // Maximum batch size for insertion
)

// InvalidateCache invalidates the Redis cache
func (u *RankingsUpdater) InvalidateCache(ctx context.Context) {
	patterns := []string{
		"warcraftlogs:global:*",
		"warcraftlogs:role:*",
		"warcraftlogs:class:*",
		"warcraftlogs:spec:*",
		"warcraftlogs:dungeon:*",
	}

	for _, pattern := range patterns {
		keys, err := cache.GetRedisClient().Keys(ctx, pattern).Result()
		if err != nil {
			log.Printf("Warning: Error getting cache keys for pattern %s: %v", pattern, err)
			continue
		}

		if len(keys) > 0 {
			if err := cache.GetRedisClient().Del(ctx, keys...).Err(); err != nil {
				log.Printf("Warning: Error deleting cache keys for pattern %s: %v", pattern, err)
			} else {
				log.Printf("Successfully invalidated %d cache keys for pattern %s", len(keys), pattern)
			}
		}
	}
}

// insertRankingsInBatches inserts rankings in batches
func (u *RankingsUpdater) insertRankingsInBatches(tx *gorm.DB, rankings []playerRankingModels.PlayerRanking) error {
	// Use prepared SQL values for affixes
	baseSQL := `
	INSERT INTO player_rankings (
			created_at, updated_at, dungeon_id, name, class, spec, role, 
			amount, hard_mode_level, duration, start_time, report_code, 
			report_fight_id, report_start_time, guild_id, guild_name, 
			guild_faction, server_id, server_name, server_region, 
			bracket_data, faction, affixes, medal, score, leaderboard
	) VALUES `

	batchSize := 100
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
		valueArgs := make([]interface{}, 0, len(batch)*25)

		for j, ranking := range batch {
			valueStrings[j] = fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				j*26+1, j*26+2, j*26+3, j*26+4, j*26+5, j*26+6, j*26+7, j*26+8, j*26+9, j*26+10,
				j*26+11, j*26+12, j*26+13, j*26+14, j*26+15, j*26+16, j*26+17, j*26+18, j*26+19, j*26+20,
				j*26+21, j*26+22, j*26+23, j*26+24, j*26+25, j*26+26,
			)

			now := time.Now()
			valueArgs = append(valueArgs,
				now, now, ranking.DungeonID, ranking.Name, ranking.Class,
				ranking.Spec, ranking.Role, ranking.Amount, ranking.HardModeLevel,
				ranking.Duration, ranking.StartTime, ranking.ReportCode,
				ranking.ReportFightID, ranking.ReportStartTime, ranking.GuildID,
				ranking.GuildName, ranking.GuildFaction, ranking.ServerID,
				ranking.ServerName, ranking.ServerRegion, ranking.BracketData,
				ranking.Faction, pq.Array(ranking.Affixes), ranking.Medal,
				ranking.Score, ranking.Leaderboard,
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

// UpdateRankings updates the rankings in the database
func (u *RankingsUpdater) UpdateRankings(ctx context.Context) error {
	log.Println("Starting rankings update...")

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
	pagesPerDungeon := 10

	rankings, err := u.rankingService.GetGlobalRankings(ctx, dungeonIDs, pagesPerDungeon)
	if err != nil {
		log.Printf("Error getting global rankings: %v", err)
		return err
	}

	err = u.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing rankings
		if err := tx.Exec("DELETE FROM player_rankings").Error; err != nil {
			return fmt.Errorf("failed to delete existing rankings: %w", err)
		}

		var newRankings []playerRankingModels.PlayerRanking

		// Processing function for each role
		processRoleRankings := func(players []PlayerScore, role string) {
			for _, player := range players {
				for _, run := range player.Runs {
					newRankings = append(newRankings, playerRankingModels.PlayerRanking{
						DungeonID:     run.DungeonID,
						Name:          player.Name,
						Class:         player.Class,
						Spec:          player.Spec,
						Role:          player.Role,
						Amount:        player.Amount,
						HardModeLevel: run.HardModeLevel,
						Duration:      run.Duration,
						StartTime:     run.StartTime,

						// Report information
						ReportCode:      run.Report.Code,
						ReportFightID:   run.Report.FightID,
						ReportStartTime: run.Report.StartTime,

						// Guild information
						GuildID:      player.Guild.ID,
						GuildName:    player.Guild.Name,
						GuildFaction: player.Guild.Faction,

						// Server information
						ServerID:     player.Server.ID,
						ServerName:   player.Server.Name,
						ServerRegion: player.Server.Region,

						// Other information
						BracketData: run.BracketData,
						Faction:     player.Faction,
						Affixes:     run.Affixes,
						Medal:       run.Medal,
						Score:       run.Score,
						Leaderboard: 0, // or the appropriate value if available
						UpdatedAt:   time.Now(),
					})
				}
			}
		}

		// Processing rankings for each role
		processRoleRankings(rankings.Tanks.Players, "Tank")
		processRoleRankings(rankings.Healers.Players, "Healer")
		processRoleRankings(rankings.DPS.Players, "DPS")

		// Insert new rankings in batches
		if len(newRankings) > 0 {
			log.Printf("Preparing to insert %d total rankings", len(newRankings))
			if err := u.insertRankingsInBatches(tx, newRankings); err != nil {
				return fmt.Errorf("failed to insert new rankings: %w", err)
			}
		} else {
			log.Println("No new rankings to insert")
		}

		// Update timestamp
		updateState := playerRankingModels.RankingsUpdateState{
			LastUpdateTime: time.Now(),
		}
		if err := tx.Save(&updateState).Error; err != nil {
			return fmt.Errorf("failed to update rankings state: %w", err)
		}

		log.Printf("Successfully processed %d rankings", len(newRankings))
		return nil
	})

	if err != nil {
		log.Printf("Error during rankings update: %v", err)
		return err
	}

	// Invalidate cache after a successful update
	u.InvalidateCache(ctx)
	log.Println("Rankings update completed successfully")
	return nil
}

// StartPeriodicUpdate starts a periodic update
func (u *RankingsUpdater) StartPeriodicUpdate() {
	ticker := time.NewTicker(updateInterval)
	go func() {
		for range ticker.C {
			if err := u.UpdateRankings(context.Background()); err != nil {
				log.Printf("Error updating rankings: %v", err)
			}
		}
	}()
}

// CheckAndUpdate checks and updates the rankings if necessary
func (u *RankingsUpdater) CheckAndUpdate() {
	var rankingsCount int64
	if err := u.db.Model(&playerRankingModels.PlayerRanking{}).Count(&rankingsCount).Error; err != nil {
		log.Printf("Error checking rankings count: %v", err)
		return
	}

	var state playerRankingModels.RankingsUpdateState
	result := u.db.First(&state)

	if result.Error == gorm.ErrRecordNotFound || rankingsCount == 0 {
		log.Println("No rankings update state found or rankings table is empty. Creating initial state and forcing update.")
		state = playerRankingModels.RankingsUpdateState{LastUpdateTime: time.Now().Add(-updateInterval)}
		if err := u.db.Create(&state).Error; err != nil {
			log.Printf("Error creating initial rankings state: %v", err)
			return
		}
		if err := u.UpdateRankings(context.Background()); err != nil {
			log.Printf("Error during initial rankings update: %v", err)
		}
		return
	}

	timeSinceLastUpdate := time.Since(state.LastUpdateTime)
	if timeSinceLastUpdate >= updateInterval || rankingsCount == 0 {
		if err := u.UpdateRankings(context.Background()); err != nil {
			log.Printf("Error during scheduled rankings update: %v", err)
		}
	} else {
		log.Printf("Rankings are up to date (count: %d)", rankingsCount)
	}
}
