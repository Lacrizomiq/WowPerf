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
	Global     *mythicplus.GlobalLeaderboardHandler
}

func NewHandler(rankingService *dungeons.RankingsService, dungeonService *dungeons.DungeonService, db *gorm.DB) *Handler {
	return &Handler{
		MythicPlus: mythicplus.NewRankingHandler(rankingService, db),
		Dungeon:    mythicplus.NewDungeonLeaderboardHandler(dungeonService),
		Global:     mythicplus.NewGlobalLeaderboardHandler(rankingService),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	warcraftlogs := router.Group("/warcraftlogs")
	{
		mythicplus := warcraftlogs.Group("/mythicplus")
		{
			// Get all the rankings for all dungeons to seed the database
			mythicplus.GET("/rankings", h.MythicPlus.GetRankings)

			// Get the leaderboard for a specific dungeon
			mythicplus.GET("/rankings/dungeon", h.Dungeon.GetDungeonLeaderboard)

			// Get the global leaderboard
			mythicplus.GET("/global/leaderboard", h.Global.GetGlobalLeaderboard)

			// Get the global leaderboard by role
			mythicplus.GET("/global/leaderboard/role", h.Global.GetRoleLeaderboard)

			// Get the global leaderboard by class
			mythicplus.GET("/global/leaderboard/class", h.Global.GetClassLeaderboard)

			// Get the global leaderboard by spec
			mythicplus.GET("/global/leaderboard/spec", h.Global.GetSpecLeaderboard)
		}
	}
}

/*
	For /global/leaderboard

	/global/leaderboard?limit=100                    // Top 100 of all players
	/global/leaderboard?limit=10                     // Top 10 of all players
*/
