package raiderio

import (
	"time"
	raiderioMythicPlus "wowperf/internal/api/raiderio/mythicplus"
	raiderioMythicPlusAnalysis "wowperf/internal/api/raiderio/mythicplus_runs"
	raiderioRaid "wowperf/internal/api/raiderio/raids"
	"wowperf/internal/services/raiderio"
	analyticsService "wowperf/internal/services/raiderio/mythicplus/analytics"

	middleware "wowperf/middleware/cache"
	"wowperf/pkg/cache"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	MythicPlusBestRun           *raiderioMythicPlus.MythicPlusBestRunHandler
	MythicPlusRunDetails        *raiderioMythicPlus.MythicPlusRunDetailsHandler
	RaidLeaderboard             *raiderioRaid.RaidLeaderboardHandler
	DungeonStats                *raiderioMythicPlus.DungeonStatsHandler
	CharacterMythicPlusBestRuns *raiderioMythicPlus.CharacterMythicPlusBestRunsHandler
	MythicPlusAnalytics         *raiderioMythicPlusAnalysis.MythicPlusRunsAnalysisHandler

	cache        cache.CacheService
	cacheManager *middleware.CacheManager
}

func NewHandler(service *raiderio.RaiderIOService, db *gorm.DB, cache cache.CacheService, cacheManager *middleware.CacheManager) *Handler {
	// Cache configuration
	cacheConfig := middleware.CacheConfig{
		Cache:      cache,
		Expiration: 24 * time.Hour,
		KeyPrefix:  "raiderio",
		Metrics:    true,
	}

	analytics := analyticsService.NewMythicPlusRunsAnalysisService(db)

	return &Handler{
		MythicPlusBestRun:           raiderioMythicPlus.NewMythicPlusBestRunHandler(service),
		MythicPlusRunDetails:        raiderioMythicPlus.NewMythicPlusRunDetailsHandler(service),
		RaidLeaderboard:             raiderioRaid.NewRaidLeaderboardHandler(service),
		DungeonStats:                raiderioMythicPlus.NewDungeonStatsHandler(service, db),
		CharacterMythicPlusBestRuns: raiderioMythicPlus.NewCharacterMythicPlusBestRunsHandler(service),
		MythicPlusAnalytics:         raiderioMythicPlusAnalysis.NewMythicPlusRunsAnalysisHandler(analytics),
		cache:                       cache,
		cacheManager:                middleware.NewCacheManager(cacheConfig),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	routeConfig := middleware.RouteConfig{
		Enabled:    true,
		Expiration: 24 * time.Hour,
	}
	raiderio := router.Group("/raiderio")
	{
		// Mythic Plus API for raider.io Mythic Plus API Data
		mythicplus := raiderio.Group("/mythicplus")
		{
			mythicplus.GET("/best-runs", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusBestRun.GetMythicPlusBestRuns)
			mythicplus.GET("/run-details", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusRunDetails.GetMythicPlusRunDetails)
			mythicplus.GET("/dungeon-stats", h.cacheManager.CacheMiddleware(routeConfig), h.DungeonStats.GetDungeonStats)
			mythicplus.GET("/character-best-runs", h.cacheManager.CacheMiddleware(routeConfig), h.CharacterMythicPlusBestRuns.GetCharacterMythicPlusBestRuns)

			// Analytics endpoints avec le même cache (24h)
			analytics := mythicplus.Group("/analytics")
			{
				// Spécialisations globales
				specs := analytics.Group("/specs")
				{
					specs.GET("/tank", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetTankSpecializations)
					specs.GET("/healer", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetHealerSpecializations)
					specs.GET("/dps", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetDPSSpecializations)
					specs.GET("/:role", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetSpecializationsByRole)
				}

				// Compositions globales
				analytics.GET("/compositions", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetTopCompositions)

				// Analyses par donjon
				dungeons := analytics.Group("/dungeons")
				{
					// Spécialisations par donjon
					dungeonSpecs := dungeons.Group("/specs")
					{
						dungeonSpecs.GET("/tank", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetTankSpecsByDungeon)
						dungeonSpecs.GET("/healer", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetHealerSpecsByDungeon)
						dungeonSpecs.GET("/dps", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetDPSSpecsByDungeon)
					}

					// Compositions par donjon
					dungeons.GET("/compositions", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetCompositionsByDungeon)

					// Spécialisations pour un donjon et rôle spécifique
					dungeons.GET("/:dungeon_slug/specs/:role", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetSpecsByDungeonAndRole)
				}

				// Analyses avancées
				analytics.GET("/key-levels", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetSpecsByKeyLevel)
				analytics.GET("/regions", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetSpecsByRegion)

				// Statistiques utilitaires
				stats := analytics.Group("/stats")
				{
					stats.GET("/overall", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetOverallStats)
					stats.GET("/key-levels", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetKeyLevelDistribution)
					stats.GET("/dungeons", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetDungeonDistribution)
					stats.GET("/regions", h.cacheManager.CacheMiddleware(routeConfig), h.MythicPlusAnalytics.GetRegionDistribution)
				}
			}
		}

		// Raids API for raider.io Raids API Data
		raids := raiderio.Group("/raids")
		{
			raids.GET("/leaderboard", h.cacheManager.CacheMiddleware(routeConfig), h.RaidLeaderboard.GetRaidLeaderboard)
		}
	}
}
