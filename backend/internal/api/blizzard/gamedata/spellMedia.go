package gamedata

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"wowperf/internal/services/blizzard"
	gamedataService "wowperf/internal/services/blizzard/gamedata"

	"github.com/gin-gonic/gin"
)

type SpellMediaHandler struct {
	Service *blizzard.Service
}

func NewSpellMediaHandler(service *blizzard.Service) *SpellMediaHandler {
	return &SpellMediaHandler{
		Service: service,
	}
}

// GetSpellMedia retrieves the media assets for a spell
func (h *SpellMediaHandler) GetSpellMedia(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	spellID := c.Param("spellId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting spell media for SpellID: %s, Region: %s, Namespace: %s, Locale: %s", spellID, region, namespace, locale)

	if spellID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(spellID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spell ID"})
		return
	}

	mediaData, err := gamedataService.GetSpellMedia(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving spell media: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Spell media not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve spell media"})
		}
		return
	}

	c.JSON(http.StatusOK, mediaData)
}
