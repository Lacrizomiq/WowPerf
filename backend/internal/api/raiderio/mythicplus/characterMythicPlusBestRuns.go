package raiderio

import (
	"net/http"
	"wowperf/internal/services/raiderio"
	raiderioMythicPlus "wowperf/internal/services/raiderio/mythicplus"

	"github.com/gin-gonic/gin"
)

type CharacterMythicPlusBestRunsHandler struct {
	Service *raiderio.RaiderIOService
}

func NewCharacterMythicPlusBestRunsHandler(service *raiderio.RaiderIOService) *CharacterMythicPlusBestRunsHandler {
	return &CharacterMythicPlusBestRunsHandler{
		Service: service,
	}
}

func (h *CharacterMythicPlusBestRunsHandler) GetCharacterMythicPlusBestRuns(c *gin.Context) {
	region := c.Query("region")
	realm := c.Query("realm")
	name := c.Query("name")
	fields := c.Query("fields")

	runInfos, err := raiderioMythicPlus.GetCharacterMythicPlusBestRuns(h.Service, region, realm, name, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, runInfos)
}
