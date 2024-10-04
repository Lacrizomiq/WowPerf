package raiderio

import (
	"net/http"
	models "wowperf/internal/models/raiderio/mythicrundetails"
	"wowperf/internal/services/raiderio"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DungeonStatsHandler struct {
	Service *raiderio.RaiderIOService
	DB      *gorm.DB
}

func NewDungeonStatsHandler(service *raiderio.RaiderIOService, db *gorm.DB) *DungeonStatsHandler {
	return &DungeonStatsHandler{
		Service: service,
		DB:      db,
	}
}

func (h *DungeonStatsHandler) GetDungeonStats(c *gin.Context) {
	season := c.Query("season")
	region := c.Query("region")

	if season == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	var dungeonStats []models.DungeonStats
	if err := h.DB.Where("season = ? AND region = ?", season, region).Find(&dungeonStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dungeon stats"})
		return
	}

	c.JSON(http.StatusOK, dungeonStats)
}
