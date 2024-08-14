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

type MythicKeystoneAffixIndexHandler struct {
	Service *blizzard.Service
}

type MythicKeystoneAffixByIDHandler struct {
	Service *blizzard.Service
}

type MythicKeystoneAffixMediaHandler struct {
	Service *blizzard.Service
}

func NewMythicKeystoneAffixIndexHandler(service *blizzard.Service) *MythicKeystoneAffixIndexHandler {
	return &MythicKeystoneAffixIndexHandler{
		Service: service,
	}
}

func NewMythicKeystoneAffixByIDHandler(service *blizzard.Service) *MythicKeystoneAffixByIDHandler {
	return &MythicKeystoneAffixByIDHandler{
		Service: service,
	}
}

func NewMythicKeystoneAffixMediaHandler(service *blizzard.Service) *MythicKeystoneAffixMediaHandler {
	return &MythicKeystoneAffixMediaHandler{
		Service: service,
	}
}

// GetMythicKeystoneAffixIndex retrieves an index of mythic keystone affixes
func (h *MythicKeystoneAffixIndexHandler) GetMythicKeystoneAffixIndex(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting mythic keystone affix index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := gamedataService.GetMythicKeystoneAffixIndex(h.Service.GameData, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone affix index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve mythic keystone affix index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetMythicKeystoneAffixByID retrieves a mythic keystone affix by ID
func (h *MythicKeystoneAffixByIDHandler) GetMythicKeystoneAffixByID(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	affixID := c.Param("affixId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting mythic keystone affix for MythicKeystoneAffixID: %s, Region: %s, Namespace: %s, Locale: %s", affixID, region, namespace, locale)

	if affixID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(affixID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mythic keystone affix ID"})
		return
	}

	data, err := gamedataService.GetMythicKeystoneAffixByID(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone affix: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mythic keystone affix not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone affix"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetMythicKeystoneAffixMedia retrieves the media assets for a mythic keystone affix
func (h *MythicKeystoneAffixMediaHandler) GetMythicKeystoneAffixMedia(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	affixID := c.Param("affixId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting mythic keystone affix media for MythicKeystoneAffixID: %s, Region: %s, Namespace: %s, Locale: %s", affixID, region, namespace, locale)

	if affixID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(affixID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mythic keystone affix ID"})
		return
	}

	mediaData, err := gamedataService.GetMythicKeystoneAffixMedia(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone affix media: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mythic keystone affix media not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone affix media"})
		}
		return
	}

	c.JSON(http.StatusOK, mediaData)
}
