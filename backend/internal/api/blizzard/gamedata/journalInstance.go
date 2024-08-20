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

type JournalInstanceIndexHandler struct {
	Service *blizzard.Service
}

type JournalInstanceByIDHandler struct {
	Service *blizzard.Service
}

type JournalInstanceMediaHandler struct {
	Service *blizzard.Service
}

func NewJournalInstanceIndexHandler(service *blizzard.Service) *JournalInstanceIndexHandler {
	return &JournalInstanceIndexHandler{
		Service: service,
	}
}

func NewJournalInstanceByIDHandler(service *blizzard.Service) *JournalInstanceByIDHandler {
	return &JournalInstanceByIDHandler{
		Service: service,
	}
}

func NewJournalInstanceMediaHandler(service *blizzard.Service) *JournalInstanceMediaHandler {
	return &JournalInstanceMediaHandler{
		Service: service,
	}
}

// GetJournalInstanceIndex retrieves an index of journal instances
func (h *JournalInstanceIndexHandler) GetJournalInstanceIndex(c *gin.Context) {
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

	log.Printf("Requesting journal instance index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := gamedataService.GetJournalInstancesIndex(h.Service.GameData, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving journal instance index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve journal instance index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetJournalInstanceByID retrieves a journal instance by ID
func (h *JournalInstanceByIDHandler) GetJournalInstanceByID(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	instanceID := c.Param("instanceId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting journal instance for JournalInstanceID: %s, Region: %s, Namespace: %s, Locale: %s", instanceID, region, namespace, locale)

	if instanceID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(instanceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid journal instance ID"})
		return
	}

	data, err := gamedataService.GetJournalInstanceByID(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving journal instance: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Journal instance not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve journal instance"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetJournalInstanceMedia retrieves the media assets for a journal instance
func (h *JournalInstanceMediaHandler) GetJournalInstanceMedia(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	instanceID := c.Param("instanceId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting journal instance media for JournalInstanceID: %s, Region: %s, Namespace: %s, Locale: %s", instanceID, region, namespace, locale)

	if instanceID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(instanceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid journal instance ID"})
		return
	}

	mediaData, err := gamedataService.GetJournalInstanceMedia(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving journal instance media: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Journal instance media not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve journal instance media"})
		}
		return
	}

	c.JSON(http.StatusOK, mediaData)
}
