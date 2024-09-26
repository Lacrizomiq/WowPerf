package raiderio

import (
	raiderioMythicPlus "wowperf/internal/api/raiderio/mythicplus"
	"wowperf/internal/services/raiderio"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	MythicPlusBestRun *raiderioMythicPlus.MythicPlusBestRunHandler
}

func NewHandler(service *raiderio.RaiderIOService) *Handler {
	return &Handler{
		MythicPlusBestRun: raiderioMythicPlus.NewMythicPlusBestRunHandler(service),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/mythic-plus/best-runs", h.MythicPlusBestRun.GetMythicPlusBestRuns)
}
