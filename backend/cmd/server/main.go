package main

import (
	"log"
	apiBlizzard "wowperf/internal/api/blizzard"
	"wowperf/internal/api/raiderio"

	"wowperf/internal/database"

	mythicplus "wowperf/internal/models/mythicplus"
	serviceBlizzard "wowperf/internal/services/blizzard"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	err = db.AutoMigrate(&mythicplus.Season{}, &mythicplus.Dungeon{}, &mythicplus.Affix{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	if err := database.SeedDatabase(db); err != nil {
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
