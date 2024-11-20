package raiderio

import (
	"time"
	raiderioMythicPlus "wowperf/internal/api/raiderio/mythicplus"
	raiderioRaid "wowperf/internal/api/raiderio/raids"
	"wowperf/internal/services/raiderio"

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
	cache                       cache.CacheService
	cacheManager                *middleware.CacheManager
}

func NewHandler(service *raiderio.RaiderIOService, db *gorm.DB, cache cache.CacheService, cacheManager *middleware.CacheManager) *Handler {
	// Cache configuration
	cacheConfig := middleware.CacheConfig{
		Cache:      cache,
		Expiration: 24 * time.Hour,
		KeyPrefix:  "raiderio",
		Metrics:    true,
	}

	return &Handler{
		MythicPlusBestRun:           raiderioMythicPlus.NewMythicPlusBestRunHandler(service),
		MythicPlusRunDetails:        raiderioMythicPlus.NewMythicPlusRunDetailsHandler(service),
		RaidLeaderboard:             raiderioRaid.NewRaidLeaderboardHandler(service),
		DungeonStats:                raiderioMythicPlus.NewDungeonStatsHandler(service, db),
		CharacterMythicPlusBestRuns: raiderioMythicPlus.NewCharacterMythicPlusBestRunsHandler(service),
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
		}

		// Raids API for raider.io Raids API Data
		raids := raiderio.Group("/raids")
		{
			raids.GET("/leaderboard", h.cacheManager.CacheMiddleware(routeConfig), h.RaidLeaderboard.GetRaidLeaderboard)
		}
	}
}
