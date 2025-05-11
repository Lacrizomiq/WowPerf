// internal/api/warcraftlogs/handler.go
package warcraftlogs

import (
	"time"

	// Services
	service "wowperf/internal/services/warcraftlogs"
	leaderboard "wowperf/internal/services/warcraftlogs/dungeons"
	mythicplusanalytics "wowperf/internal/services/warcraftlogs/mythicplus/analytics"

	// API
	mythicplus "wowperf/internal/api/warcraftlogs/mythicplus"
	mythicplusbuildsAnalysis "wowperf/internal/api/warcraftlogs/mythicplus/builds"
	character "wowperf/internal/api/warcraftlogs/mythicplus/character"

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
		Dungeon       *mythicplus.DungeonLeaderboardHandler
		Global        *mythicplus.GlobalLeaderboardHandler
		Leaderboard   *mythicplus.DungeonLeaderboardHandler
		Analysis      *mythicplus.GlobalLeaderboardAnalysisHandler
		Builds        *mythicplusbuildsAnalysis.MythicPlusBuildsAnalysisHandler
		SpecEvolution *mythicplus.SpecEvolutionMetricsAnalysisHandler
	}
	cache        cache.CacheService
	cacheManager *middleware.CacheManager
}

func NewHandler(
	globalService *leaderboard.GlobalLeaderboardService,
	analysisService *leaderboard.GlobalLeaderboardAnalysisService,
	buildsAnalysisService *mythicplusanalytics.BuildAnalysisService,
	specEvolutionService *leaderboard.SpecEvolutionMetricsAnalysisService,
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
			Dungeon       *mythicplus.DungeonLeaderboardHandler
			Global        *mythicplus.GlobalLeaderboardHandler
			Leaderboard   *mythicplus.DungeonLeaderboardHandler
			Analysis      *mythicplus.GlobalLeaderboardAnalysisHandler
			Builds        *mythicplusbuildsAnalysis.MythicPlusBuildsAnalysisHandler
			SpecEvolution *mythicplus.SpecEvolutionMetricsAnalysisHandler
		}{
			Dungeon:       mythicplus.NewDungeonLeaderboardHandler(warcraftLogsService),
			Global:        mythicplus.NewGlobalLeaderboardHandler(globalService),
			Leaderboard:   mythicplus.NewDungeonLeaderboardHandler(warcraftLogsService),
			Analysis:      mythicplus.NewGlobalLeaderboardAnalysisHandler(analysisService),
			Builds:        mythicplusbuildsAnalysis.NewMythicPlusBuildsAnalysisHandler(buildsAnalysisService),
			SpecEvolution: mythicplus.NewSpecEvolutionMetricsAnalysisHandler(specEvolutionService),
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

			// Builds analysis for Mythic+
			builds := mythicplus.Group("/builds/analysis")
			{
				// Popular items by slot
				builds.GET("/items", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetPopularItemsBySlot)

				// Popular items by slot across all encounters
				builds.GET("/items/global", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetGlobalPopularItemsBySlot)

				// Enchant usage
				builds.GET("/enchants", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetEnchantUsage)

				// Gem usage
				builds.GET("/gems", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetGemUsage)

				// Talent builds
				builds.GET("/talents/top", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetTopTalentBuilds)
				builds.GET("/talents/dungeons", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetTalentBuildsByDungeon)

				// Stats
				builds.GET("/stats", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetStatPriorities)

				// Optimal build
				builds.GET("/optimal", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetOptimalBuild)

				// Spec comparison
				builds.GET("/specs/comparison", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetSpecComparison)

				// Class spec summary
				builds.GET("/summary", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Builds.GetClassSpecSummary)
			}

			// Evolution metrics routes
			evolution := mythicplus.Group("/evolution")
			{
				// Get spec evolution metrics
				evolution.GET("/spec", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.SpecEvolution.GetSpecEvolution)

				// Get current ranking
				evolution.GET("/ranking", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.SpecEvolution.GetCurrentRanking)

				// Get latest metrics date
				evolution.GET("/latest-date", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.SpecEvolution.GetLatestMetricsDate)

				// Get specs for a class
				evolution.GET("/class/specs", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.SpecEvolution.GetClassSpecs)

				// Get all available classes
				evolution.GET("/classes", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.SpecEvolution.GetAvailableClasses)

				// Get specs for a role
				evolution.GET("/role/specs", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.SpecEvolution.GetRoleSpecs)

				// Get historical data for a spec
				evolution.GET("/spec/history", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.SpecEvolution.GetSpecHistoricalData)

				// Get available dungeons
				evolution.GET("/dungeons", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.SpecEvolution.GetAvailableDungeons)

				// Get top specs for a dungeon
				evolution.GET("/dungeons/top-specs", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.SpecEvolution.GetTopSpecsForDungeon)
			}

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

			// Get dungeon media
			mythicplus.GET("/analysis/specs/dungeons/media", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlus.Analysis.GetDungeonMedia)

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
