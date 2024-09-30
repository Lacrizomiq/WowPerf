package raiderio

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"wowperf/internal/services/raiderio"
	"wowperf/pkg/cache"

	raiderioRaid "wowperf/internal/services/raiderio/raids"

	"github.com/gin-gonic/gin"
)

type RaidLeaderboardHandler struct {
	Service *raiderio.RaiderIOService
}

func NewRaidLeaderboardHandler(service *raiderio.RaiderIOService) *RaidLeaderboardHandler {
	return &RaidLeaderboardHandler{
		Service: service,
	}
}

func (h *RaidLeaderboardHandler) GetRaidLeaderboard(c *gin.Context) {
	raid := c.Query("raid")
	difficulty := c.Query("difficulty")
	region := c.Query("region")
	limitStr := c.Query("limit")
	pageStr := c.Query("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	if raid == "" || difficulty == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	validDifficulties := map[string]bool{"normal": true, "heroic": true, "mythic": true}
	if !validDifficulties[difficulty] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid difficulty parameter"})
		return
	}

	validRegions := map[string]bool{"world": true, "us": true, "eu": true, "tw": true, "kr": true, "cn": true}
	if !validRegions[region] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region parameter"})
		return
	}

	cacheKey := fmt.Sprintf("raid_leaderboard_%s_%s_%s_%d_%d", raid, difficulty, region, limit, page)

	var leaderboard map[string]interface{}
	err = cache.Get(cacheKey, &leaderboard)
	if err == nil {
		c.JSON(http.StatusOK, leaderboard)
		return
	}

	leaderboard, err = raiderioRaid.GetRaidLeaderboard(h.Service, raid, difficulty, region, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cache.Set(cacheKey, leaderboard, 1*time.Hour)

	c.JSON(http.StatusOK, leaderboard)
}
