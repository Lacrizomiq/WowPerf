package raiderio

import (
	raiderioMythicPlus "wowperf/internal/api/raiderio/mythicplus"
	raiderioRaid "wowperf/internal/api/raiderio/raids"
	"wowperf/internal/services/raiderio"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	MythicPlusBestRun    *raiderioMythicPlus.MythicPlusBestRunHandler
	MythicPlusRunDetails *raiderioMythicPlus.MythicPlusRunDetailsHandler
	RaidLeaderboard      *raiderioRaid.RaidLeaderboardHandler
}

func NewHandler(service *raiderio.RaiderIOService) *Handler {
	return &Handler{
		MythicPlusBestRun:    raiderioMythicPlus.NewMythicPlusBestRunHandler(service),
		MythicPlusRunDetails: raiderioMythicPlus.NewMythicPlusRunDetailsHandler(service),
		RaidLeaderboard:      raiderioRaid.NewRaidLeaderboardHandler(service),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/mythic-plus/best-runs", h.MythicPlusBestRun.GetMythicPlusBestRuns)
	router.GET("/mythic-plus/run-details", h.MythicPlusRunDetails.GetMythicPlusRunDetails)
	router.GET("/raids/leaderboard", h.RaidLeaderboard.GetRaidLeaderboard)
}
