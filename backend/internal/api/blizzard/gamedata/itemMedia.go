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

type ItemMediaHandler struct {
	Service *blizzard.Service
}

func NewItemMediaHandler(service *blizzard.Service) *ItemMediaHandler {
	return &ItemMediaHandler{
		Service: service,
	}
}

// GetItemMedia retrieves the media assets for an item.
func (h *ItemMediaHandler) GetItemMedia(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	itemID := c.Param("itemId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting item media for ItemID: %s, Region: %s, Namespace: %s, Locale: %s", itemID, region, namespace, locale)

	if itemID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	mediaData, err := gamedataService.GetItemMedia(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving item media: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item media not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve item media"})
		}
		return
	}

	c.JSON(http.StatusOK, mediaData)
}
