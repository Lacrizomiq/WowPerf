// internal/api/warcraftlogs/handler.go
package warcraftlogs

import (
	"time"

	// Services
	service "wowperf/internal/services/warcraftlogs"

	// API
	mythicplus "wowperf/internal/api/warcraftlogs/mythicplus"
	character "wowperf/internal/api/warcraftlogs/mythicplus/character"
	leaderboard "wowperf/internal/services/warcraftlogs/dungeons"

	middleware "wowperf/middleware/cache"
	"wowperf/pkg/cache"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	// Regrouper les handlers logiquement
	Character struct {
		Ranking *character.CharacterRankingHandler
	}
	MythicPlus struct {
		Dungeon     *mythicplus.DungeonLeaderboardHandler
		Global      *mythicplus.GlobalLeaderboardHandler
		Leaderboard *mythicplus.DungeonLeaderboardHandler
		Analysis    *mythicplus.GlobalLeaderboardAnalysisHandler
	}
	cache        cache.CacheService
	cacheManager *middleware.CacheManager
}

func NewHandler(
	globalService *leaderboard.GlobalLeaderboardService,
	analysisService *leaderboard.GlobalLeaderboardAnalysisService,
	warcraftLogsService *service.WarcraftLogsClientService,
	db *gorm.DB,
	cache cache.CacheService,
	cacheManager *middleware.CacheManager,
) *Handler {
	return &Handler{
		Character: struct {
			Ranking *character.CharacterRankingHandler
		}{
			Ranking: character.NewCharacterRankingHandler(warcraftLogsService),
		},
		MythicPlus: struct {
			Dungeon     *mythicplus.DungeonLeaderboardHandler
			Global      *mythicplus.GlobalLeaderboardHandler
			Leaderboard *mythicplus.DungeonLeaderboardHandler
			Analysis    *mythicplus.GlobalLeaderboardAnalysisHandler
		}{
			Dungeon:     mythicplus.NewDungeonLeaderboardHandler(warcraftLogsService),
			Global:      mythicplus.NewGlobalLeaderboardHandler(globalService),
			Leaderboard: mythicplus.NewDungeonLeaderboardHandler(warcraftLogsService),
			Analysis:    mythicplus.NewGlobalLeaderboardAnalysisHandler(analysisService),
		},
		cache:        cache,
		cacheManager: cacheManager,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {

	routeConfig := middleware.RouteConfig{
		Enabled:    true,
		Expiration: 2 * time.Hour,
	}

	warcraftlogs := router.Group("/warcraftlogs")
	{
		// Character routes
		character := warcraftlogs.Group("/character")
		{
			// Get the character ranking for a given character name, server slug, server region and zone ID
			character.GET("/ranking/player", h.cacheManager.CacheMiddleware(routeConfig), h.Character.Ranking.GetCharacterRanking)
		}

		// Mythic+ routes
		mythicplus := warcraftlogs.Group("/mythicplus")
		{
			// Get all the rankings for all dungeons to seed the database
			// mythicplus.GET("/rankings", h.MythicPlus.GetRankings)

			// Get the leaderboard for a specific dungeon by team
			mythicplus.GET("/rankings/dungeon/team", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Dungeon.GetDungeonLeaderboardByTeam)

			// Get the leaderboard for a specific dungeon by player
			mythicplus.GET("/rankings/dungeon/player", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Leaderboard.GetDungeonLeaderboardByPlayer)

			// Get the global leaderboard
			mythicplus.GET("/global/leaderboard", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Global.GetGlobalLeaderboard)

			// Get the global leaderboard by role
			mythicplus.GET("/global/leaderboard/role", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Global.GetRoleLeaderboard)

			// Get the global leaderboard by class
			mythicplus.GET("/global/leaderboard/class", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Global.GetClassLeaderboard)

			// Get the global leaderboard by spec
			mythicplus.GET("/global/leaderboard/spec", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Global.GetSpecLeaderboard)

			// Analysis routes (moved under /mythicplus for consistency with M+ ecosystem)
			// Get average global scores per spec
			mythicplus.GET("/analysis/specs/avg-scores", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Analysis.GetSpecGlobalScores)

			// Get average global scores per class
			mythicplus.GET("/analysis/classes/avg-scores", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Analysis.GetClassGlobalScores)

			// Get max key levels per spec and dungeon
			mythicplus.GET("/analysis/specs/dungeons/max-levels-key", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Analysis.GetSpecDungeonMaxKeyLevels)

			// Get average key levels per dungeon
			mythicplus.GET("/analysis/dungeons/avg-levels-key", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Analysis.GetDungeonAvgKeyLevels)

			// Get top 10 players per spec
			mythicplus.GET("/analysis/players/top-specs", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Analysis.GetTop10PlayersPerSpec)

			// Get top 5 players per role
			mythicplus.GET("/analysis/players/top-roles", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Analysis.GetTop5PlayersPerRole)
		}
	}
}

/*
	For /global/leaderboard

	/global/leaderboard?limit=100                    // Top 100 of all players
	/global/leaderboard?limit=10                     // Top 10 of all players
*/
