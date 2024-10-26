// services/warcraftlogs/rankingsUpdater.go
package warcraftlogs

import (
	"context"
	"log"
	"time"
	rankingsModels "wowperf/internal/models/warcraftlogs/mythicplus"
	"wowperf/pkg/cache"

	"gorm.io/gorm"
)

const updateInterval = 24 * time.Hour

const (
	DungeonAraKara       = 12660
	DungeonCityOfThreads = 12669
	DungeonGrimBatol     = 60670
	DungeonMists         = 62290
	DungeonSiege         = 61822
	DungeonDawnbreaker   = 12662
	DungeonNecroticWake  = 62286
	DungeonStonevault    = 12652
)

type RankingsUpdater struct {
	db             *gorm.DB
	rankingService *RankingsService
}

func NewRankingsUpdater(db *gorm.DB, rankingService *RankingsService) *RankingsUpdater {
	return &RankingsUpdater{
		db:             db,
		rankingService: rankingService,
	}
}

// InvalidateCache invalidates the cache when the rankings are updated
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

// Update rankings
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
	pagesPerDungeon := 5

	rankings, err := u.rankingService.GetGlobalRankings(ctx, dungeonIDs, pagesPerDungeon)
	if err != nil {
		log.Printf("Error getting global rankings: %v", err)
		return err
	}

	// Debug logs to see the data retrieved
	log.Printf("Tanks count: %d", len(rankings.Tanks.Players))
	log.Printf("Healers count: %d", len(rankings.Healers.Players))
	log.Printf("DPS count: %d", len(rankings.DPS.Players))

	err = u.db.Transaction(func(tx *gorm.DB) error {
		// Hard delete all existing rankings
		if err := tx.Exec("DELETE FROM player_rankings").Error; err != nil {
			log.Printf("Error deleting existing rankings: %v", err)
			return err
		}

		// Count new rankings to be inserted for logging
		var totalNewRankings int

		// Collect and insert rankings for each role
		for _, role := range []struct {
			name    string
			players []PlayerScore
		}{
			{"tank", rankings.Tanks.Players},
			{"healer", rankings.Healers.Players},
			{"dps", rankings.DPS.Players},
		} {
			log.Printf("Processing %s - %d players", role.name, len(role.players))

			var newRankings []rankingsModels.PlayerRanking
			for _, player := range role.players {
				log.Printf("Processing player %s (%s %s) with %d runs",
					player.Name, player.Class, player.Spec, len(player.Runs))

				for _, run := range player.Runs {
					newRankings = append(newRankings, rankingsModels.PlayerRanking{
						DungeonID: run.DungeonID,
						PlayerID:  player.ID,
						Name:      player.Name,
						Class:     player.Class,
						Spec:      player.Spec,
						Role:      role.name,
						Score:     run.Score,
						UpdatedAt: time.Now(),
					})
				}
			}

			// Insert rankings for this role
			if len(newRankings) > 0 {
				log.Printf("Attempting to insert %d %s rankings", len(newRankings), role.name)
				if err := tx.Create(&newRankings).Error; err != nil {
					log.Printf("Error inserting %s rankings: %v", role.name, err)
					return err
				}
				totalNewRankings += len(newRankings)
				log.Printf("Successfully inserted %d %s rankings", len(newRankings), role.name)
			} else {
				log.Printf("No rankings to insert for %s", role.name)
			}
		}

		log.Printf("Total rankings inserted: %d", totalNewRankings)

		// Update the last update time
		updateResult := tx.Model(&rankingsModels.RankingsUpdateState{}).
			Where("1 = 1").
			Updates(map[string]interface{}{
				"last_update_time": time.Now(),
			})

		if updateResult.Error != nil {
			log.Printf("Error updating rankings state: %v", updateResult.Error)
			return updateResult.Error
		}

		if updateResult.RowsAffected == 0 {
			log.Println("Creating new rankings update state")
			if err := tx.Create(&rankingsModels.RankingsUpdateState{
				LastUpdateTime: time.Now(),
			}).Error; err != nil {
				log.Printf("Error creating rankings state: %v", err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error during rankings update: %v", err)
		return err
	}

	// Verify the insert
	var count int64
	if err := u.db.Model(&rankingsModels.PlayerRanking{}).Count(&count).Error; err != nil {
		log.Printf("Error counting rankings after update: %v", err)
	} else {
		log.Printf("Total rankings in database after update: %d", count)
	}

	// Invalidate the cache after a successful update
	u.InvalidateCache(ctx)
	log.Println("Cache invalidated after successful update")

	log.Println("Rankings update completed successfully")
	return nil
}

// Start periodic update
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

// Check and update rankings if the interval has been exceeded or if the table is empty
func (u *RankingsUpdater) CheckAndUpdate() {
	// First check if we have any rankings
	var rankingsCount int64
	if err := u.db.Model(&rankingsModels.PlayerRanking{}).Count(&rankingsCount).Error; err != nil {
		log.Printf("Error checking rankings count: %v", err)
		return
	}

	var state rankingsModels.RankingsUpdateState
	result := u.db.First(&state)

	if result.Error == gorm.ErrRecordNotFound || rankingsCount == 0 {
		log.Println("No rankings update state found or rankings table is empty. Creating initial state and forcing update.")
		state = rankingsModels.RankingsUpdateState{LastUpdateTime: time.Now().Add(-updateInterval)}
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
	log.Printf("Current update state: Last update was %v ago", timeSinceLastUpdate)
	log.Printf("Current rankings count: %d", rankingsCount)

	if timeSinceLastUpdate >= updateInterval || rankingsCount == 0 {
		log.Printf("Update needed: Time since last update: %v, Rankings count: %d", timeSinceLastUpdate, rankingsCount)
		if err := u.UpdateRankings(context.Background()); err != nil {
			log.Printf("Error during scheduled rankings update: %v", err)
			return
		}
	} else {
		log.Printf("Rankings are up to date (count: %d)", rankingsCount)
	}
}
