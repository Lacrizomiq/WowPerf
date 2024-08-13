package profile

import (
	"net/http"
	"wowperf/internal/services/blizzard"
	wrapper "wowperf/internal/wrapper/blizzard"

	"github.com/gin-gonic/gin"
)

type CharacterProfileHandler struct {
	Service *blizzard.Service
}

func NewCharacterProfileHandler(service *blizzard.Service) *CharacterProfileHandler {
	return &CharacterProfileHandler{
		Service: service,
	}
}

// GetCharacterProfile retrieves a character's profile information, including name, realm, and class.
func (h *CharacterProfileHandler) GetCharacterProfile(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	characterData, err := h.Service.Client.GetCharacterProfile(region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character profile"})
		return
	}

	mediaData, err := h.Service.Client.GetCharacterMedia(region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character media"})
		return
	}

	profile, err := wrapper.TransformCharacterInfo(characterData, mediaData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transform character profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}
