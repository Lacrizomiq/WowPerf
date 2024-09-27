package raiderio

import (
	"net/http"
	"strconv"
	"wowperf/internal/services/raiderio"
	raiderioMythicPlus "wowperf/internal/services/raiderio/mythicplus"
	raiderioMythicPlusWrapper "wowperf/internal/wrapper/raiderio"

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

	// Get the raw data from the Raider.io Mythic Plus Runs Detail API
	rawRunDetails, err := raiderioMythicPlus.GetMythicPlusRunsDetails(h.Service, season, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transform the raw data into a more usable format with the wrapper
	transformedRunDetails, err := raiderioMythicPlusWrapper.TransformMythicPlusRun(rawRunDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transformedRunDetails)
}
