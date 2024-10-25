package warcraftlogs

import (
	"log"
	"net/http"
	"strconv"
	dungeonService "wowperf/internal/services/warcraftlogs/dungeons"

	"github.com/gin-gonic/gin"
)

type DungeonLeaderboardHandler struct {
	dungeonService *dungeonService.DungeonService
}

func NewDungeonLeaderboardHandler(dungeonService *dungeonService.DungeonService) *DungeonLeaderboardHandler {
	return &DungeonLeaderboardHandler{dungeonService: dungeonService}
}

func (h *DungeonLeaderboardHandler) GetDungeonLeaderboard(c *gin.Context) {
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
	leaderboard, err := h.dungeonService.GetDungeonLeaderboard(encounterID, page)
	if err != nil {
		log.Printf("Error getting dungeon leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dungeon leaderboard"})
		return
	}

	// Log the dungeon leaderboard for debugging
	log.Printf("Got leaderboard for dungeon %d, page %d: %v", encounterID, page, leaderboard)
	if leaderboard != nil {
		log.Printf("Leaderboard: %d rankings", len(leaderboard.Rankings))
		for i, ranking := range leaderboard.Rankings {
			log.Printf("Ranking %d: %v", i, ranking)
		}
	}

	c.JSON(http.StatusOK, leaderboard)
}
