// internal/api/warcraftlogs/handler.go
package warcraftlogs

import (
	"time"
	mythicplus "wowperf/internal/api/warcraftlogs/mythicplus"
	dungeons "wowperf/internal/services/warcraftlogs/dungeons"
	middleware "wowperf/middleware/cache"
	"wowperf/pkg/cache"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	MythicPlus          *mythicplus.RankingHandler
	Dungeon             *mythicplus.DungeonLeaderboardHandler
	Global              *mythicplus.GlobalLeaderboardHandler
	LeaderboardByPlayer *mythicplus.DungeonLeaderboardHandler
	cache               cache.CacheService
	cacheManager        *middleware.CacheManager
}

func NewHandler(rankingService *dungeons.RankingsService, dungeonService *dungeons.DungeonService, db *gorm.DB, cache cache.CacheService, cacheManager *middleware.CacheManager) *Handler {
	// Cache configuration
	cacheConfig := middleware.CacheConfig{
		Cache:      cache,
		Expiration: 8 * time.Hour,
		KeyPrefix:  "warcraftlogs",
		Metrics:    true,
	}
	return &Handler{
		MythicPlus:          mythicplus.NewRankingHandler(rankingService, db),
		Dungeon:             mythicplus.NewDungeonLeaderboardHandler(dungeonService),
		Global:              mythicplus.NewGlobalLeaderboardHandler(rankingService),
		LeaderboardByPlayer: mythicplus.NewDungeonLeaderboardHandler(dungeonService),
		cache:               cache,
		cacheManager:        middleware.NewCacheManager(cacheConfig),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {

	routeConfig := middleware.RouteConfig{
		Enabled:    true,
		Expiration: 8 * time.Hour,
	}

	warcraftlogs := router.Group("/warcraftlogs")
	{
		mythicplus := warcraftlogs.Group("/mythicplus")
		{
			// Get all the rankings for all dungeons to seed the database
			// mythicplus.GET("/rankings", h.MythicPlus.GetRankings)

			// Get the leaderboard for a specific dungeon by team
			mythicplus.GET("/rankings/dungeon/team", h.cacheManager.CacheMiddleware(routeConfig), h.Dungeon.GetDungeonLeaderboardByTeam)

			// Get the leaderboard for a specific dungeon by player
			mythicplus.GET("/rankings/dungeon/player", h.cacheManager.CacheMiddleware(routeConfig), h.LeaderboardByPlayer.GetDungeonLeaderboardByPlayer)

			// Get the global leaderboard
			mythicplus.GET("/global/leaderboard", h.cacheManager.CacheMiddleware(routeConfig), h.Global.GetGlobalLeaderboard)

			// Get the global leaderboard by role
			mythicplus.GET("/global/leaderboard/role", h.cacheManager.CacheMiddleware(routeConfig), h.Global.GetRoleLeaderboard)

			// Get the global leaderboard by class
			mythicplus.GET("/global/leaderboard/class", h.cacheManager.CacheMiddleware(routeConfig), h.Global.GetClassLeaderboard)

			// Get the global leaderboard by spec
			mythicplus.GET("/global/leaderboard/spec", h.cacheManager.CacheMiddleware(routeConfig), h.Global.GetSpecLeaderboard)
		}
	}
}

/*
	For /global/leaderboard

	/global/leaderboard?limit=100                    // Top 100 of all players
	/global/leaderboard?limit=10                     // Top 10 of all players
*/
