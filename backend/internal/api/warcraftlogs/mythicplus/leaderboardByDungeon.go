package warcraftlogs

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	dungeonService "wowperf/internal/services/warcraftlogs/dungeons"
	"wowperf/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

	cacheKey := fmt.Sprintf("warcraftlogs:dungeon:%d:page:%d", encounterID, page)

	var leaderboard interface{}
	err = cache.Get(cacheKey, &leaderboard)
	if err == nil {
		c.JSON(http.StatusOK, leaderboard)
		return
	}

	if err != redis.Nil {
		log.Printf("Error getting dungeon leaderboard from cache: %v", err)
	}

	// Get the dungeon leaderboard
	leaderboard, err = h.dungeonService.GetDungeonLeaderboard(encounterID, page)
	if err != nil {
		log.Printf("Error getting dungeon leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dungeon leaderboard"})
		return
	}

	if err := cache.Set(cacheKey, leaderboard, 8*time.Hour); err != nil {
		log.Printf("Error setting dungeon leaderboard in cache: %v", err)
	}

	c.JSON(http.StatusOK, leaderboard)
}
