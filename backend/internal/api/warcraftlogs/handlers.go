// internal/api/warcraftlogs/handler.go
package warcraftlogs

import (
	mythicplus "wowperf/internal/api/warcraftlogs/mythicplus"
	dungeons "wowperf/internal/services/warcraftlogs/dungeons"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	MythicPlus *mythicplus.RankingHandler
	Dungeon    *mythicplus.DungeonLeaderboardHandler
}

func NewHandler(rankingService *dungeons.RankingsService, dungeonService *dungeons.DungeonService, db *gorm.DB) *Handler {
	return &Handler{
		MythicPlus: mythicplus.NewRankingHandler(rankingService, db),
		Dungeon:    mythicplus.NewDungeonLeaderboardHandler(dungeonService),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	warcraftlogs := router.Group("/warcraftlogs")
	{
		mythicplus := warcraftlogs.Group("/mythicplus")
		{
			mythicplus.GET("/rankings", h.MythicPlus.GetRankings)
			mythicplus.GET("/rankings/dungeon", h.Dungeon.GetDungeonLeaderboard)
		}
	}
}
