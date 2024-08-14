package profile

import (
	"net/http"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/profile"
	wrapper "wowperf/internal/wrapper/blizzard"

	"github.com/gin-gonic/gin"
)

type EquipmentHandler struct {
	Service *blizzard.Service
}

func NewEquipmentHandler(service *blizzard.Service) *EquipmentHandler {
	return &EquipmentHandler{
		Service: service,
	}
}

// GetCharacterEquipment retrieves a character's equipment, including items, gems, and enchantments.
func (h *EquipmentHandler) GetCharacterEquipment(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	equipmentData, err := profile.GetCharacterEquipment(h.Service.Profile, region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character equipment"})
		return
	}

	transformedGear, err := wrapper.TransformCharacterGear(equipmentData, h.Service.GameDataClient, region, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transform character equipment"})
		return
	}

	c.JSON(http.StatusOK, transformedGear)
}
