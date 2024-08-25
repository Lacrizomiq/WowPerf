package database

import (
	"fmt"
	"log"
	"path/filepath"
	mythicplus "wowperf/internal/database/static/mythicplus"
	talents "wowperf/internal/database/static/talents"

	"gorm.io/gorm"
)

const (
	staticMythicPlusPath = "./static/M+/"
)

func SeedDatabase(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Seed Mythic+ data
	if err := seedMythicPlusData(db); err != nil {
		return fmt.Errorf("error seeding Mythic+ data: %v", err)
	}

	// Seed Talents data
	if err := talents.SeedTalents(db); err != nil {
		return fmt.Errorf("error seeding Talents data: %v", err)
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
