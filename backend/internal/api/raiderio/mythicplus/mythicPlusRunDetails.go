package raiderio

import (
	"net/http"
	"strconv"
	"wowperf/internal/services/raiderio"
	raiderioMythicPlus "wowperf/internal/services/raiderio/mythicplus"

	"github.com/gin-gonic/gin"
)

type MythicPlusRunDetailsHandler struct {
	Service *raiderio.RaiderIOService
}

func NewMythicPlusRunDetailsHandler(service *raiderio.RaiderIOService) *MythicPlusRunDetailsHandler {
	return &MythicPlusRunDetailsHandler{
		Service: service,
	}
}

func (h *MythicPlusRunDetailsHandler) GetMythicPlusRunDetails(c *gin.Context) {
	season := c.Query("season")
	idStr := c.Query("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}

	if season == "" || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	runDetails, err := raiderioMythicPlus.GetMythicPlusRunsDetails(h.Service, season, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, runDetails)
}
