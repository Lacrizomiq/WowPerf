package main

import (
	"log"
	api "wowperf/internal/api/raiderio"
	"wowperf/internal/services/raiderio"

	"github.com/gin-gonic/gin"
)

func main() {
	rioClient := raiderio.NewCLient()
	handler := api.NewHandler(rioClient)

	r := gin.Default()

	r.GET("/characters", handler.GetCharacterProfile)
	r.GET("/characters/mythic-plus-scores", handler.GetCharacterMythicPlusScores)
	r.GET("/characters/raid-progression", handler.GetCharacterRaidProgression)

	log.Fatal(r.Run(":8080"))
}
