package profile

import (
	"net/http"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/profile"

	"github.com/gin-gonic/gin"
)

type CharacterStatsHandler struct {
	Service *blizzard.Service
}

func NewCharacterStatsHandler(service *blizzard.Service) *CharacterStatsHandler {
	return &CharacterStatsHandler{
		Service: service,
	}
}

// GetCharacterStats retrieves the stats of a character.
func (h *CharacterStatsHandler) GetCharacterStats(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	characterData, err := profile.GetCharacterStats(h.Service.Profile, region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character stats"})
		return
	}

	c.JSON(http.StatusOK, characterData)
}
