package main

import (
	"log"
	"wowperf/internal/api/blizzard"
	"wowperf/internal/api/raiderio"

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

	rioHandler := raiderio.NewHandler()
	blizzardhandler, err := blizzard.NewHandler()
	if err != nil {
		log.Fatalf("Failed to initialize blizzard client: %v", err)
	}

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

	// Blizzard Profile API
	r.GET("/blizzard/characters/:realmSlug/:characterName", blizzardhandler.GetCharacterProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/mythic-keystone-profile", blizzardhandler.GetCharacterMythicKeystoneProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/equipment", blizzardhandler.GetCharacterEquipment)
	r.GET("/blizzard/characters/:realmSlug/:characterName/specializations", blizzardhandler.GetCharacterSpecializations)
	r.GET("/blizzard/characters/:realmSlug/:characterName/mythic-keystone-profile/season/:seasonId", blizzardhandler.GetCharacterMythicKeystoneSeasonDetails)
	r.GET("/blizzard/characters/:realmSlug/:characterName/character-media", blizzardhandler.GetCharacterMedia)

	// Blizzard Game Data API
	r.GET("/blizzard/data/item/:itemId/media", blizzardhandler.GetItemMedia)

	log.Fatal(r.Run(":8080"))
}
