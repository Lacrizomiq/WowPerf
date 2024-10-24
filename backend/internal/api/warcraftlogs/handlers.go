package warcraftlogs

import (
	rankingsHandler "wowperf/internal/api/warcraftlogs/mythicplus"
	rankingsService "wowperf/internal/services/warcraftlogs/dungeons"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	Rankings *rankingsHandler.RankingsHandler
}

func NewHandler(rankingService *rankingsService.RankingsService, db *gorm.DB) *Handler {
	return &Handler{
		Rankings: rankingsHandler.NewRankingsHandler(rankingService, db),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	warcraftlogs := router.Group("/warcraftlogs")
	{
		mythicplus := warcraftlogs.Group("/mythicplus")
		{
			mythicplus.GET("/rankings", h.Rankings.GetRankings)
		}
	}
}
