package raiderioMythicPlus

import (
	"log"
	"sync"
	"time"
	models "wowperf/internal/models/raiderio/mythicrundetails"
	"wowperf/internal/services/raiderio"

	"gorm.io/gorm"
)

// UpdateDungeonStats updates the dungeon stats in the database
func UpdateDungeonStats(db *gorm.DB, rioService *raiderio.RaiderIOService) error {
	seasons := []string{"season-tww-1"}
	regions := []string{"world", "us", "eu", "tw", "kr", "cn"}
	dungeonSlugs := []string{"all", "arakara-city-of-echoes", "city-of-threads", "grim-batol", "mists-of-tirna-scithe", "siege-of-boralus", "the-dawnbreaker", "the-necrotic-wake", "the-stonevault"}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit the concurrency to 5 goroutines at a time

	for _, season := range seasons {
		for _, region := range regions {
			wg.Add(1)
			go func(season, region string) {
				defer wg.Done()
				semaphore <- struct{}{}        // Acquire a place in the semaphore
				defer func() { <-semaphore }() // Release the place at the end

				dungeonStats, err := GetAllDungeonStats(rioService, season, region, dungeonSlugs)
				if err != nil {
					log.Printf("Error getting stats for %s %s: %v", season, region, err)
					return
				}

				err = db.Transaction(func(tx *gorm.DB) error {
					for _, stats := range dungeonStats {
						var dbStats models.DungeonStats
						result := tx.Where("season = ? AND region = ? AND dungeon_slug = ?", season, region, stats.DungeonName).First(&dbStats)
						if result.Error == gorm.ErrRecordNotFound {
							dbStats = models.DungeonStats{
								Season:      season,
								Region:      region,
								DungeonSlug: stats.DungeonName,
							}
						}

						dbStats.RoleStats = stats.RoleStats
						dbStats.UpdatedAt = time.Now()

						if result.Error == gorm.ErrRecordNotFound {
							if err := tx.Create(&dbStats).Error; err != nil {
								return err
							}
						} else {
							if err := tx.Save(&dbStats).Error; err != nil {
								return err
							}
						}
					}
					return nil
				})

				if err != nil {
					log.Printf("Error updating stats for %s %s: %v", season, region, err)
				}
			}(season, region)
		}
	}

	wg.Wait()
	return nil
}

// StartWeeklyDungeonStatsUpdate starts a ticker that updates the dungeon stats once a week
func StartWeeklyDungeonStatsUpdate(db *gorm.DB, rioService *raiderio.RaiderIOService) {
	ticker := time.NewTicker(7 * 24 * time.Hour) // Once a week
	go func() {
		for range ticker.C {
			if err := UpdateDungeonStats(db, rioService); err != nil {
				log.Printf("Error updating dungeon stats: %v", err)
			}
		}
	}()
}
