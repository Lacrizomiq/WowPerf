package main

import (
	"log"
	"wowperf/internal/api/blizzard"
	"wowperf/internal/api/raiderio"

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
	blizardhandler, err := blizzard.NewHandler()
	if err != nil {
		log.Fatalf("Failed to initialize blizzard client: %v", err)
	}

	r := gin.Default()

	// Raider.io API
	r.GET("/characters", rioHandler.GetCharacterProfile)
	r.GET("/characters/mythic-plus-scores", rioHandler.GetCharacterMythicPlusScores)
	r.GET("/characters/raid-progression", rioHandler.GetCharacterRaidProgression)

	// Blizzard API
	r.GET("/blizzard/characters/:realmSlug/:characterName", blizardhandler.GetCharacterProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/mythic-keystone-profile", blizardhandler.GetCharacterMythicKeystoneProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/equipment", blizardhandler.GetCharacterEquipment)

	log.Fatal(r.Run(":8080"))
}
