package database

import (
	"fmt"
	"log"
	"path/filepath"
	"time"
	staticMythicPlus "wowperf/internal/database/static/mythicplus"
	staticRaids "wowperf/internal/database/static/raids"
	staticTalents "wowperf/internal/database/static/talents"

	serviceRaiderio "wowperf/internal/services/raiderio"
	mythicplusUpdate "wowperf/internal/services/raiderio/mythicplus"

	mythicplus "wowperf/internal/models/mythicplus"
	raiderioMythicPlus "wowperf/internal/models/raiderio/mythicrundetails"
	raids "wowperf/internal/models/raids"
	talents "wowperf/internal/models/talents"

	rankingsModels "wowperf/internal/models/warcraftlogs/mythicplus"

	migrations "wowperf/internal/database/migrations"

	"gorm.io/gorm"
)

const (
	staticMythicPlusPath = "./data/static/M+"
	staticRaidsPath      = "./data/static/Raid"
)

// InitializeDatabase sets up the database with all required data
func InitializeDatabase(db *gorm.DB) error {
	log.Println("Initializing database...")

	if err := migrations.RunMigrations(db); err != nil {
		return fmt.Errorf("error performing migrations: %v", err)
	}

	if err := ensureUpdateState(db); err != nil {
		return fmt.Errorf("error ensuring update state: %v", err)
	}

	if err := ensureInitialData(db); err != nil {
		return fmt.Errorf("error ensuring initial data: %v", err)
	}

	log.Println("Database initialization completed successfully.")
	return nil
}

// ensureUpdateState initializes update state records if they don't exist
func ensureUpdateState(db *gorm.DB) error {
	// Mythic+ Update State
	var updateState raiderioMythicPlus.UpdateState
	if err := db.FirstOrCreate(&updateState, raiderioMythicPlus.UpdateState{LastUpdateTime: time.Now().Add(-25 * time.Hour)}).Error; err != nil {
		return fmt.Errorf("error initializing UpdateState: %v", err)
	}

	// WarcraftLogs Rankings Update State
	if err := rankingsModels.InitializeRankingsUpdateState(db); err != nil {
		return fmt.Errorf("error initializing RankingsUpdateState: %v", err)
	}
	return nil
}

// ensureInitialData ensures all required data is present in the database
func ensureInitialData(db *gorm.DB) error {
	if err := ensureMythicPlusData(db); err != nil {
		return err
	}

	if err := ensureTalentsData(db); err != nil {
		return err
	}

	if err := ensureRaidsData(db); err != nil {
		return err
	}

	if err := ensureDungeonStats(db); err != nil {
		return err
	}

	if err := ensureRankingsData(db); err != nil {
		return err
	}

	return nil
}

// ensureMythicPlusData ensures Mythic+ data is up to date
// Modified to use UpdateSeasons and UpdateAffixes for keeping data in sync with JSON files
func ensureMythicPlusData(db *gorm.DB) error {
	log.Println("Checking and updating Mythic+ data...")
	absPath, err := filepath.Abs(staticMythicPlusPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Check if tables are empty - use original seed functions for initial population
	var count int64
	db.Model(&mythicplus.Season{}).Count(&count)
	if count == 0 {
		log.Println("Seeding initial Mythic+ data...")
		if err := staticMythicPlus.SeedSeasons(db, filepath.Join(absPath, "DF", "MMDF.json")); err != nil {
			return err
		}
		if err := staticMythicPlus.SeedSeasons(db, filepath.Join(absPath, "TWW", "S1MMTWW.json")); err != nil {
			return err
		}
		if err := staticMythicPlus.SeedAffixes(db, filepath.Join(absPath, "affix.json")); err != nil {
			return err
		}
	} else {
		// Tables have data, use the new update functions to check for changes
		log.Println("Updating existing Mythic+ data...")
		if err := staticMythicPlus.UpdateSeasons(db, filepath.Join(absPath, "DF", "MMDF.json")); err != nil {
			return err
		}
		if err := staticMythicPlus.UpdateSeasons(db, filepath.Join(absPath, "TWW", "S1MMTWW.json")); err != nil {
			return err
		}
		if err := staticMythicPlus.UpdateAffixes(db, filepath.Join(absPath, "affix.json")); err != nil {
			return err
		}
	}
	return nil
}

// ensureTalentsData ensures Talents data is present
func ensureTalentsData(db *gorm.DB) error {
	var count int64
	db.Model(&talents.TalentTree{}).Count(&count)
	if count == 0 {
		log.Println("Seeding Talents data...")
		if err := staticTalents.SeedTalents(db); err != nil {
			return fmt.Errorf("error seeding Talents data: %v", err)
		}
	}
	return nil
}

// ensureRaidsData ensures Raids data is up to date
// Modified to use UpdateRaids for keeping data in sync with JSON files
func ensureRaidsData(db *gorm.DB) error {
	log.Println("Checking and updating Raids data...")
	absPath, err := filepath.Abs(staticRaidsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Check if tables are empty - use original seed functions for initial population
	var count int64
	db.Model(&raids.Raid{}).Count(&count)
	if count == 0 {
		log.Println("Seeding initial Raids data...")
		if err := staticRaids.SeedRaids(db, filepath.Join(absPath, "TWW", "raids.json")); err != nil {
			return err
		}
		if err := staticRaids.SeedRaids(db, filepath.Join(absPath, "DF", "raids.json")); err != nil {
			return err
		}
	} else {
		// Tables have data, use the new update function to check for changes
		log.Println("Updating existing Raids data...")
		if err := staticRaids.UpdateRaids(db, filepath.Join(absPath, "TWW", "raids.json")); err != nil {
			return err
		}
		if err := staticRaids.UpdateRaids(db, filepath.Join(absPath, "DF", "raids.json")); err != nil {
			return err
		}
	}
	return nil
}

// ensureDungeonStats ensures DungeonStats data is present and up to date
func ensureDungeonStats(db *gorm.DB) error {
	if mythicplusUpdate.IsDungeonStatsEmpty(db) {
		log.Println("DungeonStats table is effectively empty. Resetting update state and initiating update...")
		if err := mythicplusUpdate.RemoveInitialDungeonStats(db); err != nil {
			log.Printf("Error removing initial DungeonStats: %v", err)
		}
		mythicplusUpdate.ResetUpdateState(db)

		rioService, err := serviceRaiderio.NewRaiderIOService()
		if err != nil {
			return fmt.Errorf("failed to initialize raiderio service: %v", err)
		}

		const maxRetries = 3
		var updateErr error
		for i := 0; i < maxRetries; i++ {
			updateErr = mythicplusUpdate.UpdateDungeonStats(db, rioService)
			if updateErr == nil {
				break
			}
			log.Printf("Attempt %d failed to update DungeonStats: %v. Retrying...", i+1, updateErr)
			time.Sleep(time.Second * 5)
		}

		if updateErr != nil {
			return fmt.Errorf("error performing DungeonStats update after %d attempts: %v", maxRetries, updateErr)
		}

		log.Println("DungeonStats update completed successfully")
	} else {
		log.Println("DungeonStats are not empty, checking if update is needed...")
		if mythicplusUpdate.CheckAndSetUpdateLock(db) {
			rioService, err := serviceRaiderio.NewRaiderIOService()
			if err != nil {
				return fmt.Errorf("failed to initialize raiderio service: %v", err)
			}
			if err := mythicplusUpdate.UpdateDungeonStats(db, rioService); err != nil {
				return fmt.Errorf("error updating DungeonStats: %v", err)
			}
			log.Println("DungeonStats update completed successfully")
		} else {
			log.Println("DungeonStats are up to date")
		}
	}

	return nil
}

// ensureRankingsData ensures WarcraftLogs Rankings data is present
func ensureRankingsData(db *gorm.DB) error {
	// Check if the rankings table is empty
	if err := db.AutoMigrate(&rankingsModels.PlayerRanking{}, &rankingsModels.RankingsUpdateState{}); err != nil {
		return fmt.Errorf("error migrating rankings models: %v", err)
	}

	var count int64
	db.Model(&rankingsModels.PlayerRanking{}).Count(&count)
	if count == 0 {
		log.Println("Initial rankings data is empty, will be populated on first update...")
	}
	return nil
}

// UpdateGameData manually triggers an update of all game data from JSON files
// This function can be called from an API endpoint or CLI command
func UpdateGameData(db *gorm.DB) error {
	log.Println("Manually updating all game data...")

	// Update Mythic+ data
	if err := ensureMythicPlusData(db); err != nil {
		return fmt.Errorf("error updating Mythic+ data: %v", err)
	}

	// Update Raids data
	if err := ensureRaidsData(db); err != nil {
		return fmt.Errorf("error updating Raids data: %v", err)
	}

	log.Println("Game data update completed successfully.")
	return nil
}
