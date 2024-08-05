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
	blizardhandler, err := blizzard.NewHandler()
	if err != nil {
		log.Fatalf("Failed to initialize blizzard client: %v", err)
	}

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(config))

	// Raider.io API
	r.GET("/characters", rioHandler.GetCharacterProfile)
	r.GET("/characters/mythic-plus-scores", rioHandler.GetCharacterMythicPlusScores)
	r.GET("/characters/raid-progression", rioHandler.GetCharacterRaidProgression)

	// Blizzard API
	r.GET("/blizzard/characters/:realmSlug/:characterName", blizardhandler.GetCharacterProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/mythic-keystone-profile", blizardhandler.GetCharacterMythicKeystoneProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/equipment", blizardhandler.GetCharacterEquipment)
	r.GET("/blizzard/characters/:realmSlug/:characterName/specializations", blizardhandler.GetCharacterSpecializations)

	log.Fatal(r.Run(":8080"))
}
