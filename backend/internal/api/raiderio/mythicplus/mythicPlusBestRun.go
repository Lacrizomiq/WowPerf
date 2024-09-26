package raiderio

import (
	"net/http"
	"strconv"
	"wowperf/internal/services/raiderio"
	raiderioMythicPlus "wowperf/internal/services/raiderio/mythicplus"

	"github.com/gin-gonic/gin"
)

type MythicPlusBestRunHandler struct {
	Service *raiderio.RaiderIOService
}

func NewMythicPlusBestRunHandler(service *raiderio.RaiderIOService) *MythicPlusBestRunHandler {
	return &MythicPlusBestRunHandler{
		Service: service,
	}
}

func (h *MythicPlusBestRunHandler) GetMythicPlusBestRuns(c *gin.Context) {
	season := c.Query("season")
	region := c.Query("region")
	dungeon := c.Query("dungeon")
	pageStr := c.Query("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	if season == "" || region == "" || dungeon == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	bestRuns, err := raiderioMythicPlus.GetMythicPlusBestRuns(h.Service, season, region, dungeon, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bestRuns)
}
