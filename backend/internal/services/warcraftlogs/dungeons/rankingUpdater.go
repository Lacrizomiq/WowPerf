// services/warcraftlogs/rankingsUpdater.go
package warcraftlogs

import (
	"context"
	"log"
	"time"
	rankingsModels "wowperf/internal/models/warcraftlogs/mythicplus"

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
	pagesPerDungeon := 3

	rankings, err := u.rankingService.GetGlobalRankings(ctx, dungeonIDs, pagesPerDungeon)
	if err != nil {
		return err
	}

	err = u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&rankingsModels.PlayerRanking{}).Error; err != nil {
			return err
		}

		for _, role := range []struct {
			name    string
			players []PlayerScore
		}{
			{"tank", rankings.Tanks.Players},
			{"healer", rankings.Healers.Players},
			{"dps", rankings.DPS.Players},
		} {
			for _, player := range role.players {
				for _, run := range player.Runs {
					ranking := rankingsModels.PlayerRanking{
						DungeonID: run.DungeonID,
						PlayerID:  player.ID,
						Name:      player.Name,
						Class:     player.Class,
						Spec:      player.Spec,
						Role:      role.name,
						Score:     run.Score,
						UpdatedAt: time.Now(),
					}
					if err := tx.Create(&ranking).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error during rankings update: %v", err)
		return err
	}

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

// Check and update rankings if the interval has been exceeded
func (u *RankingsUpdater) CheckAndUpdate() {
	var state rankingsModels.RankingsUpdateState
	result := u.db.First(&state)

	if result.Error == gorm.ErrRecordNotFound {
		log.Println("No rankings update state found. Creating initial state and forcing update.")
		u.db.Create(&rankingsModels.RankingsUpdateState{LastUpdateTime: time.Now().Add(-updateInterval)})
		u.UpdateRankings(context.Background())
		return
	}

	if time.Since(state.LastUpdateTime) >= updateInterval {
		log.Println("Rankings update interval exceeded. Performing update.")
		if err := u.UpdateRankings(context.Background()); err != nil {
			log.Printf("Error during scheduled rankings update: %v", err)
			return
		}
		u.db.Model(&state).Update("LastUpdateTime", time.Now())
	}
}
