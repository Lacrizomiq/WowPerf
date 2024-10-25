// internal/api/warcraftlogs/mythicplus/ranking.go
package warcraftlogs

import (
	"net/http"
	rankingsModels "wowperf/internal/models/warcraftlogs/mythicplus"
	dungeons "wowperf/internal/services/warcraftlogs/dungeons"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RankingHandler struct {
	Service *dungeons.RankingsService
	DB      *gorm.DB
}

func NewRankingHandler(service *dungeons.RankingsService, db *gorm.DB) *RankingHandler {
	return &RankingHandler{
		Service: service,
		DB:      db,
	}
}

func (h *RankingHandler) GetRankings(c *gin.Context) {
	var rankings []rankingsModels.PlayerRanking
	if err := h.DB.Find(&rankings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rankings)
}
