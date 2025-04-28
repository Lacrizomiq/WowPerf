package raiderioMythicPlus

import (
	"log"
	"sync"
	"time"
	models "wowperf/internal/models/raiderio/mythicrundetails"
	"wowperf/internal/services/raiderio"

	"gorm.io/gorm"
)

const updateInterval = 7 * 24 * time.Hour // 7 days

func IsDungeonStatsEmpty(db *gorm.DB) bool {
	var count int64
	err := db.Model(&models.DungeonStats{}).
		Where("season != 'initial' OR region != 'initial' OR dungeon_slug != 'initial'").
		Count(&count).Error
	if err != nil {
		log.Printf("Error checking if dungeon stats are empty: %v", err)
		return true // Assume empty if there's an error
	}
	log.Printf("Found %d non-initial dungeon stats records", count)
	return count == 0
}

func RemoveInitialDungeonStats(db *gorm.DB) error {
	return db.Where("season = 'initial' AND region = 'initial' AND dungeon_slug = 'initial'").
		Delete(&models.DungeonStats{}).Error
}

func CheckAndSetUpdateLock(db *gorm.DB) bool {
	if IsDungeonStatsEmpty(db) {
		log.Println("DungeonStats table is empty. Forcing update.")
		ResetUpdateState(db)
		return true
	}

	var state models.UpdateState
	result := db.First(&state)

	if result.Error == gorm.ErrRecordNotFound {
		log.Println("No update state found. Creating initial state and forcing update.")
		newState := models.UpdateState{LastUpdateTime: time.Now().Add(-updateInterval)}
		if err := db.Create(&newState).Error; err != nil {
			log.Printf("Error creating initial update state: %v", err)
			return false
		}
		return true
	} else if result.Error != nil {
		log.Printf("Error fetching update state: %v", result.Error)
		return false
	}

	timeSinceLastUpdate := time.Since(state.LastUpdateTime)
	log.Printf("Time since last update: %v", timeSinceLastUpdate)

	if timeSinceLastUpdate >= updateInterval {
		log.Println("Update interval exceeded. Performing update.")
		if err := db.Model(&state).Update("LastUpdateTime", time.Now()).Error; err != nil {
			log.Printf("Error updating last update time: %v", err)
			return false
		}
		return true
	}

	log.Println("Update not needed at this time.")
	return false
}

func ResetUpdateState(db *gorm.DB) {
	if err := db.Exec("DELETE FROM update_states").Error; err != nil {
		log.Printf("Error resetting update state: %v", err)
	}
}

// UpdateDungeonStats updates the dungeon stats in the database
func UpdateDungeonStats(db *gorm.DB, rioService *raiderio.RaiderIOService) error {

	log.Println("Starting dungeon stats update...")

	seasons := []string{"season-tww-2"}
	regions := []string{"world", "us", "eu", "tw", "kr"}
	dungeonSlugs := []string{"all", "cinderbrew-meadery", "darkflame-cleft",
		"operation-mechagon-workshop", "operation-floodgate", "priory-of-the-sacred-flame",
		"the-motherlode", "the-rookery", "theater-of-pain"}

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
						dbStats.SpecStats = stats.SpecStats
						dbStats.LevelStats = stats.LevelStats
						dbStats.TeamComp = stats.TeamComp
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
	ticker := time.NewTicker(updateInterval)
	go func() {
		for range ticker.C {
			if err := UpdateDungeonStats(db, rioService); err != nil {
				log.Printf("Error updating dungeon stats: %v", err)
			} else {
				log.Println("Weekly dungeon stats update completed successfully")
			}
		}
	}()
}
