package gamedata

import (
	"fmt"
	"net/http"
	"wowperf/internal/services/blizzard"
	gamedataService "wowperf/internal/services/blizzard/gamedata"

	"github.com/gin-gonic/gin"
)

type RealmsIndexHandler struct {
	Service *blizzard.Service
}

func NewRealmsIndexHandler(service *blizzard.Service) *RealmsIndexHandler {
	return &RealmsIndexHandler{
		Service: service,
	}
}

func (h *RealmsIndexHandler) GetRealmsIndex(c *gin.Context) {
	region := c.Query("region")
	if region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "region parameter is required"})
		return
	}

	namespace := c.DefaultQuery("namespace", fmt.Sprintf("dynamic-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	realms, err := gamedataService.GetRealmsIndex(h.Service.GameData, region, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, realms)
}
