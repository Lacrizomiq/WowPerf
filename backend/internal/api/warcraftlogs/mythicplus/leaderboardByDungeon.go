package warcraftlogs

import (
	"log"
	"net/http"
	"strconv"
	service "wowperf/internal/services/warcraftlogs"
	dungeons "wowperf/internal/services/warcraftlogs/dungeons"

	"github.com/gin-gonic/gin"
)

type DungeonLeaderboardHandler struct {
	dungeonService *service.WarcraftLogsClientService
}

func NewDungeonLeaderboardHandler(dungeonService *service.WarcraftLogsClientService) *DungeonLeaderboardHandler {
	return &DungeonLeaderboardHandler{dungeonService: dungeonService}
}

// GetDungeonLeaderboardByPlayer returns the dungeon leaderboard for a given encounter and page
func (h *DungeonLeaderboardHandler) GetDungeonLeaderboardByPlayer(c *gin.Context) {

	// Mandatory params
	encounterID, err := strconv.Atoi(c.Query("encounterID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid encounter ID"})
		return
	}

	// Get the page number from the query parameters
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Create the params
	params := dungeons.LeaderboardParams{
		EncounterID: encounterID,
		Page:        page,
		// Optional params
		ServerRegion: c.Query("serverRegion"),
		ServerSlug:   c.Query("serverSlug"),
		ClassName:    c.Query("className"),
		SpecName:     c.Query("specName"),
	}

	// Log the params for debugging
	log.Printf("Fetching dungeon leaderboard for params: %+v", params)

	// Get the dungeon leaderboard
	leaderboard, err := dungeons.GetDungeonLeaderboardByPlayer(h.dungeonService, params)
	if err != nil {
		log.Printf("Error getting dungeon leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dungeon leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// GetDungeonLeaderboardByTeam returns the dungeon leaderboard for a given encounter and page
func (h *DungeonLeaderboardHandler) GetDungeonLeaderboardByTeam(c *gin.Context) {
	// Get the dungeon ID from the query parameters
	encounterID, err := strconv.Atoi(c.Query("encounterID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid encounter ID"})
		return
	}

	// Get the page number from the query parameters
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Get the dungeon leaderboard
	leaderboard, err := dungeons.GetDungeonLeaderboardByTeam(h.dungeonService, encounterID, page)
	if err != nil {
		log.Printf("Error getting dungeon leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dungeon leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}
