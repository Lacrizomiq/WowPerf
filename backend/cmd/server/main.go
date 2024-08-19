package main

import (
	"log"
	apiBlizzard "wowperf/internal/api/blizzard"
	"wowperf/internal/api/raiderio"
	"wowperf/internal/database"
	"wowperf/internal/database/migrations"
	seeds "wowperf/internal/database/seeds/mythicdj"
	models "wowperf/internal/models/mythicdj"
	serviceBlizzard "wowperf/internal/services/blizzard"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return
	}

	// DB init
	database.InitDB()

	// Auto migrate the database
	err = database.DB.AutoMigrate(&models.MythicDj{}, &models.KeystoneUpgrade{}, &models.Affix{}, &models.Season{}, &models.Period{}, &models.MythicDjSeason{})
	if err != nil {
		log.Fatalf("Failed to perform auto migration: %v", err)
	}

	// Check for duplicates before cleanup
	checkDuplicates(database.DB)

	// Cleanup duplicates
	if err := migrations.CleanupDuplicateKeystoneUpgrades(database.DB); err != nil {
		log.Fatalf("Failed to cleanup duplicate keystone upgrades: %v", err)
	}

	// Check for duplicates after cleanup
	checkDuplicates(database.DB)

	// Add unique constraint
	if err := migrations.AddUniqueConstraintToKeystoneUpgrades(database.DB); err != nil {
		log.Fatalf("Failed to add unique constraint to keystone_upgrades table: %v", err)
	}

	// Seed the database
	if err := seeds.SeedAll(database.DB); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	blizzardService, err := serviceBlizzard.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize blizzard service: %v", err)
	}

	rioHandler := raiderio.NewHandler()
	blizzardHandler := apiBlizzard.NewHandler(blizzardService)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(config))

	// Raider.io API
	r.GET("/characters", rioHandler.GetCharacterProfile)
	r.GET("/characters/mythic-plus-scores", rioHandler.GetCharacterMythicPlusScores)
	r.GET("/characters/raid-progression", rioHandler.GetCharacterRaidProgression)

	// Blizzard API
	blizzardHandler.RegisterRoutes(r)

	log.Fatal(r.Run(":8080"))
}

func checkDuplicates(db *gorm.DB) {
	var result []struct {
		MythicDjID   uint
		UpgradeLevel int
		Count        int
	}
	db.Table("keystone_upgrades").
		Select("mythic_dj_id, upgrade_level, count(*)").
		Group("mythic_dj_id, upgrade_level").
		Having("count(*) > 1").
		Scan(&result)

	if len(result) > 0 {
		log.Println("Duplicates found:")
		for _, r := range result {
			log.Printf("MythicDjID: %d, UpgradeLevel: %d, Count: %d\n", r.MythicDjID, r.UpgradeLevel, r.Count)
		}
	} else {
		log.Println("No duplicates found")
	}
}
