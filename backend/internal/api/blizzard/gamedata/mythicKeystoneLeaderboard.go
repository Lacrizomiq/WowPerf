package gamedata

import (
	"fmt"
	"net/http"
	"strconv"
	"wowperf/internal/services/blizzard"
	gamedataService "wowperf/internal/services/blizzard/gamedata"

	"github.com/gin-gonic/gin"
)

type MythicKeystoneLeaderboardHandler struct {
	Service *blizzard.Service
}

func NewMythicKeystoneLeaderboardHandler(service *blizzard.Service) *MythicKeystoneLeaderboardHandler {
	return &MythicKeystoneLeaderboardHandler{
		Service: service,
	}
}

// GetMythicKeystoneLeaderboardIndex retrieves an index of mythic keystone leaderboards for a connected realm.
func (h *MythicKeystoneLeaderboardHandler) GetMythicKeystoneLeaderboardIndex(c *gin.Context) {
	connectedRealmID, err := strconv.Atoi(c.Param("connectedRealmId")) // Changed from "connectedRealmID" to "connectedRealmId"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connected realm ID"})
		return
	}
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("dynamic-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if connectedRealmID == 0 || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	leaderboardData, err := gamedataService.GetMythicKeystoneLeaderboardIndex(h.Service.GameData, connectedRealmID, region, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone leaderboard index"})
		return
	}

	c.JSON(http.StatusOK, leaderboardData)
}
