package profile

import (
	"net/http"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/profile"

	"github.com/gin-gonic/gin"
)

type CharacterMediaHandler struct {
	Service *blizzard.Service
}

func NewCharacterMediaHandler(service *blizzard.Service) *CharacterMediaHandler {
	return &CharacterMediaHandler{
		Service: service,
	}
}

// GetCharacterMedia retrieves a character's media assets, including avatar, inset avatar, main raw, and character media.
func (h *CharacterMediaHandler) GetCharacterMedia(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	media, err := profile.GetCharacterMedia(h.Service.Profile, region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character media"})
		return
	}

	c.JSON(http.StatusOK, media)
}
