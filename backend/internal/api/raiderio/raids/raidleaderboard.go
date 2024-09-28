package raiderio

import (
	"net/http"
	"wowperf/internal/services/raiderio"
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

	if raid == "" || difficulty == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	leaderboard, err := raiderioRaid.GetRaidLeaderboard(h.Service, raid, difficulty, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}
