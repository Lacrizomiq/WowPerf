package database

import (
	"fmt"
	"log"
	"path/filepath"
	"time"
	mythicplus "wowperf/internal/database/static/mythicplus"
	raids "wowperf/internal/database/static/raids"
	talents "wowperf/internal/database/static/talents"
	models "wowperf/internal/models/raiderio/mythicrundetails"

	"gorm.io/gorm"
)

const (
	staticMythicPlusPath = "./static/M+/"
	staticRaidsPath      = "./static/Raid/"
)

func SeedDatabase(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Initialize UpdateState
	var updateState models.UpdateState
	if err := db.FirstOrCreate(&updateState, models.UpdateState{LastUpdateTime: time.Now().Add(-25 * time.Hour)}).Error; err != nil {
		return fmt.Errorf("error initializing UpdateState: %v", err)
	}

	// Seed Mythic+ data
	if err := seedMythicPlusData(db); err != nil {
		return fmt.Errorf("error seeding Mythic+ data: %v", err)
	}

	// Seed Talents data
	if err := talents.SeedTalents(db); err != nil {
		return fmt.Errorf("error seeding Talents data: %v", err)
	}

	// Seed Raids data
	if err := seedRaidsData(db); err != nil {
		return fmt.Errorf("error seeding Raids data: %v", err)
	}

	log.Println("Database seeding completed successfully.")
	return nil
}

func seedMythicPlusData(db *gorm.DB) error {
	absPath, err := filepath.Abs(staticMythicPlusPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	if err := mythicplus.SeedSeasons(db, filepath.Join(absPath, "DF", "MMDF.json")); err != nil {
		return err
	}
	if err := mythicplus.SeedSeasons(db, filepath.Join(absPath, "TWW", "S1MMTWW.json")); err != nil {
		return err
	}
	if err := mythicplus.SeedAffixes(db, filepath.Join(absPath, "affix.json")); err != nil {
		return err
	}
	return nil
}

func seedRaidsData(db *gorm.DB) error {
	absPath, err := filepath.Abs(staticRaidsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	if err := raids.SeedRaids(db, filepath.Join(absPath, "TWW", "raids.json")); err != nil {
		return err
	}

	if err := raids.SeedRaids(db, filepath.Join(absPath, "DF", "raids.json")); err != nil {
		return err
	}

	return nil
}
