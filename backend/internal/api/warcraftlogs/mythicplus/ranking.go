package warcraftlogs

import (
	"net/http"
	rankingsModels "wowperf/internal/models/warcraftlogs/mythicplus"
	rankingsService "wowperf/internal/services/warcraftlogs/dungeons"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RankingsHandler struct {
	Service *rankingsService.RankingsService
	DB      *gorm.DB
}

func NewRankingsHandler(service *rankingsService.RankingsService, db *gorm.DB) *RankingsHandler {
	return &RankingsHandler{Service: service, DB: db}
}

// Get rankings
func (h *RankingsHandler) GetRankings(c *gin.Context) {
	var rankings []rankingsModels.PlayerRanking
	if err := h.DB.Find(&rankings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rankings)
}
