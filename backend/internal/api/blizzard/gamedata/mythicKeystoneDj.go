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

type MythicKeystoneIndexHandler struct {
	Service *blizzard.Service
}

type MythicKeystoneDungeonsIndexHandler struct {
	Service *blizzard.Service
}

type MythicKeystoneByIDHandler struct {
	Service *blizzard.Service
}

type MythicKeystonePeriodsIndexHandler struct {
	Service *blizzard.Service
}

type MythicKeystonePeriodByIDHandler struct {
	Service *blizzard.Service
}

type MythicKeystoneSeasonsIndexHandler struct {
	Service *blizzard.Service
}

type MythicKeystoneSeasonByIDHandler struct {
	Service *blizzard.Service
}

func NewMythicKeystoneIndexHandler(service *blizzard.Service) *MythicKeystoneIndexHandler {
	return &MythicKeystoneIndexHandler{
		Service: service,
	}
}

func NewMythicKeystoneDungeonsIndexHandler(service *blizzard.Service) *MythicKeystoneDungeonsIndexHandler {
	return &MythicKeystoneDungeonsIndexHandler{
		Service: service,
	}
}

func NewMythicKeystoneByIDHandler(service *blizzard.Service) *MythicKeystoneByIDHandler {
	return &MythicKeystoneByIDHandler{
		Service: service,
	}
}

func NewMythicKeystonePeriodsIndexHandler(service *blizzard.Service) *MythicKeystonePeriodsIndexHandler {
	return &MythicKeystonePeriodsIndexHandler{
		Service: service,
	}
}

func NewMythicKeystonePeriodByIDHandler(service *blizzard.Service) *MythicKeystonePeriodByIDHandler {
	return &MythicKeystonePeriodByIDHandler{
		Service: service,
	}
}

func NewMythicKeystoneSeasonsIndexHandler(service *blizzard.Service) *MythicKeystoneSeasonsIndexHandler {
	return &MythicKeystoneSeasonsIndexHandler{
		Service: service,
	}
}

func NewMythicKeystoneSeasonByIDHandler(service *blizzard.Service) *MythicKeystoneSeasonByIDHandler {
	return &MythicKeystoneSeasonByIDHandler{
		Service: service,
	}
}

// GetMythicKeystoneIndex retrieves an index of mythic keystones
func (h *MythicKeystoneIndexHandler) GetMythicKeystoneIndex(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("dynamic-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("dynamic-%s", region)
	}

	log.Printf("Requesting mythic keystone index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := gamedataService.GetMythicKeystoneIndex(h.Service.GameData, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve mythic keystone index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetMythicKeystoneDungeonsIndex retrieves an index of mythic keystone dungeons
func (h *MythicKeystoneDungeonsIndexHandler) GetMythicKeystoneDungeonsIndex(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("dynamic-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("dynamic-%s", region)
	}

	log.Printf("Requesting mythic keystone dungeons index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := gamedataService.GetMythicKeystoneDungeonsIndex(h.Service.GameData, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone dungeons index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve mythic keystone dungeons index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetMythicKeystoneByID retrieves a mythic keystone by ID
func (h *MythicKeystoneByIDHandler) GetMythicKeystoneByID(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	mythicKeystoneID := c.Param("mythicKeystoneId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("dynamic-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("dynamic-%s", region)
	}

	log.Printf("Requesting mythic keystone for MythicKeystoneID: %s, Region: %s, Namespace: %s, Locale: %s", mythicKeystoneID, region, namespace, locale)

	if mythicKeystoneID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(mythicKeystoneID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mythic keystone ID"})
		return
	}

	data, err := gamedataService.GetMythicKeystoneByID(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mythic keystone not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetMythicKeystonePeriodsIndex retrieves an index of mythic keystone periods
func (h *MythicKeystonePeriodsIndexHandler) GetMythicKeystonePeriodsIndex(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("dynamic-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("dynamic-%s", region)
	}

	log.Printf("Requesting mythic keystone periods index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := gamedataService.GetMythicKeystonePeriodsIndex(h.Service.GameData, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone periods index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve mythic keystone periods index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetMythicKeystonePeriodByID retrieves a mythic keystone period by periodID
func (h *MythicKeystonePeriodByIDHandler) GetMythicKeystonePeriodByID(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	periodID := c.Param("periodId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("dynamic-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("dynamic-%s", region)
	}

	log.Printf("Requesting mythic keystone period for PeriodID: %s, Region: %s, Namespace: %s, Locale: %s", periodID, region, namespace, locale)

	if periodID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period ID"})
		return
	}

	data, err := gamedataService.GetMythicKeystonePeriodByID(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone period: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mythic keystone period not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone period"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetMythicKeystoneSeasonsIndex retrieves an index of mythic keystone seasons
func (h *MythicKeystoneSeasonsIndexHandler) GetMythicKeystoneSeasonsIndex(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("dynamic-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("dynamic-%s", region)
	}

	log.Printf("Requesting mythic keystone seasons index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := gamedataService.GetMythicKeystoneSeasonsIndex(h.Service.GameData, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone seasons index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve mythic keystone seasons index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetMythicKeystoneSeasonByID retrieves a mythic keystone season by seasonID
func (h *MythicKeystoneSeasonByIDHandler) GetMythicKeystoneSeasonByID(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	seasonID := c.Param("seasonId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("dynamic-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("dynamic-%s", region)
	}

	log.Printf("Requesting mythic keystone season for SeasonID: %s, Region: %s, Namespace: %s, Locale: %s", seasonID, region, namespace, locale)

	if seasonID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(seasonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}

	data, err := gamedataService.GetMythicKeystoneSeasonByID(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone season: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mythic keystone season not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone season"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}
