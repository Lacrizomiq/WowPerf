package warcraftlogs

import (
	"log"
	"net/http"
	"strconv"
	service "wowperf/internal/services/warcraftlogs"
	character "wowperf/internal/services/warcraftlogs/character"

	"github.com/gin-gonic/gin"
)

type CharacterRankingHandler struct {
	characterRankingService *service.WarcraftLogsClientService
}

func NewCharacterRankingHandler(characterRankingService *service.WarcraftLogsClientService) *CharacterRankingHandler {
	return &CharacterRankingHandler{characterRankingService: characterRankingService}
}

// GetCharacterRanking returns the character ranking for a given character name, server slug, server region and zone ID
func (h *CharacterRankingHandler) GetCharacterRanking(c *gin.Context) {
	characterName := c.Query("characterName")
	serverSlug := c.Query("serverSlug")
	serverRegion := c.Query("serverRegion")

	// validate the character name
	if characterName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Character name is required"})
		return
	}

	// validate the server slug
	if serverSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server slug is required"})
		return
	}

	// validate the server region
	if serverRegion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server region is required"})
		return
	}

	// get the zone ID
	zoneID, err := strconv.Atoi(c.Query("zoneID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid zone ID"})
		return
	}

	// decode the character name
	decodedCharacterName, err := character.DecodeCharacterName(characterName)
	if err != nil {
		log.Printf("Failed to decode character name: %s %v", characterName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode character name"})
		return
	}

	// get the character ranking
	characterRanking, err := character.GetCharacterRanking(h.characterRankingService, decodedCharacterName, serverSlug, serverRegion, zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get character ranking"})
		return
	}

	// return the character ranking
	c.JSON(http.StatusOK, characterRanking)
}
