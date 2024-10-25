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

	"gorm.io/gorm"
)

const (
	staticMythicPlusPath = "./static/M+/"
	staticRaidsPath      = "./static/Raid/"
)

func InitializeDatabase(db *gorm.DB) error {
	log.Println("Initializing database...")

	if err := Migrate(db); err != nil {
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

func ensureUpdateState(db *gorm.DB) error {

	// Mythic+ Update State
	var updateState raiderioMythicPlus.UpdateState
	if err := db.FirstOrCreate(&updateState, raiderioMythicPlus.UpdateState{LastUpdateTime: time.Now().Add(-25 * time.Hour)}).Error; err != nil {
		return fmt.Errorf("error initializing UpdateState: %v", err)
	}

	// WarcraftLogs Rankings Update State
	var rankingsUpdateState rankingsModels.RankingsUpdateState
	if err := db.FirstOrCreate(&rankingsUpdateState, rankingsModels.RankingsUpdateState{LastUpdateTime: time.Now().Add(-25 * time.Hour)}).Error; err != nil {
		return fmt.Errorf("error initializing RankingsUpdateState: %v", err)
	}
	return nil
}

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

func ensureMythicPlusData(db *gorm.DB) error {
	var count int64
	db.Model(&mythicplus.Season{}).Count(&count)
	if count == 0 {
		log.Println("Seeding Mythic+ data...")
		absPath, err := filepath.Abs(staticMythicPlusPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %v", err)
		}

		if err := staticMythicPlus.SeedSeasons(db, filepath.Join(absPath, "DF", "MMDF.json")); err != nil {
			return err
		}
		if err := staticMythicPlus.SeedSeasons(db, filepath.Join(absPath, "TWW", "S1MMTWW.json")); err != nil {
			return err
		}
		if err := staticMythicPlus.SeedAffixes(db, filepath.Join(absPath, "affix.json")); err != nil {
			return err
		}
	}
	return nil
}

// Ensure Talents data is present
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

// Ensure Raids data is present
func ensureRaidsData(db *gorm.DB) error {
	var count int64
	db.Model(&raids.Raid{}).Count(&count)
	if count == 0 {
		log.Println("Seeding Raids data...")
		absPath, err := filepath.Abs(staticRaidsPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %v", err)
		}

		if err := staticRaids.SeedRaids(db, filepath.Join(absPath, "TWW", "raids.json")); err != nil {
			return err
		}

		if err := staticRaids.SeedRaids(db, filepath.Join(absPath, "DF", "raids.json")); err != nil {
			return err
		}
	}
	return nil
}

// Ensure DungeonStats data is present
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

// Ensure WarcraftLogs Rankings data is present
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
