package gamedata

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"wowperf/internal/services/blizzard"
	gamedataService "wowperf/internal/services/blizzard/gamedata"
	wrapper "wowperf/internal/wrapper/blizzard"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TalentTreeIndexHandler struct {
	Service *blizzard.Service
}

type TalentTreeHandler struct {
	DB *gorm.DB
}

type TalentTreeNodesHandler struct {
	Service *blizzard.Service
}

type TalentIndexHandler struct {
	Service *blizzard.Service
}

type TalentByIDHandler struct {
	Service *blizzard.Service
}

func NewTalentTreeIndexHandler(service *blizzard.Service) *TalentTreeIndexHandler {
	return &TalentTreeIndexHandler{
		Service: service,
	}
}

func NewTalentTreeHandler(db *gorm.DB) *TalentTreeHandler {
	return &TalentTreeHandler{
		DB: db,
	}
}

func NewTalentTreeNodesHandler(service *blizzard.Service) *TalentTreeNodesHandler {
	return &TalentTreeNodesHandler{
		Service: service,
	}
}

func NewTalentIndexHandler(service *blizzard.Service) *TalentIndexHandler {
	return &TalentIndexHandler{
		Service: service,
	}
}

func NewTalentByIDHandler(service *blizzard.Service) *TalentByIDHandler {
	return &TalentByIDHandler{
		Service: service,
	}
}

// GetTalentTreeIndex retrieves an index of talent trees
func (h *TalentTreeIndexHandler) GetTalentTreeIndex(c *gin.Context) {
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

	log.Printf("Requesting talent tree index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := gamedataService.GetTalentTreeIndex(h.Service.GameData, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving talent tree index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve talent tree index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetTalentTree retrieves a talent tree by spec ID from the database
func (h *TalentTreeHandler) GetTalentTree(c *gin.Context) {
	talentTreeID := c.Param("talentTreeId")
	specID := c.Param("specId")

	log.Printf("Requesting talent tree for TalentTreeID: %s, SpecID: %s", talentTreeID, specID)

	if talentTreeID == "" || specID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	treeID, err := strconv.Atoi(talentTreeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talent tree ID"})
		return
	}

	specIDInt, err := strconv.Atoi(specID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
		return
	}

	talentTree, err := wrapper.GetFullTalentTree(h.DB, treeID, specIDInt)
	if err != nil {
		log.Printf("Error retrieving talent tree: %v", err)
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent tree not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve talent tree"})
		}
		return
	}

	c.JSON(http.StatusOK, talentTree)
}

// GetTalentTreeNodes retrieves the nodes of a talent tree as well as links to associated playable specializations given a talent tree id
func (h *TalentTreeNodesHandler) GetTalentTreeNodes(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	talentTreeID := c.Param("talentTreeId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting talent tree nodes for TalentTreeID: %s, Region: %s, Namespace: %s, Locale: %s", talentTreeID, region, namespace, locale)

	if talentTreeID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	treeID, err := strconv.Atoi(talentTreeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talent tree ID"})
		return
	}

	data, err := gamedataService.GetTalentTreeNodes(h.Service.GameData, treeID, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving talent tree nodes: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent tree not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve talent tree nodes"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetTalentIndex retrieves an index of talents
func (h *TalentIndexHandler) GetTalentIndex(c *gin.Context) {
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

	log.Printf("Requesting talent index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := gamedataService.GetTalentIndex(h.Service.GameData, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving talent index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve talent index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetTalentByID retrieves a talent by ID
func (h *TalentByIDHandler) GetTalentByID(c *gin.Context) {
	if h.Service.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	talentID := c.Param("talentId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting talent for TalentID: %s, Region: %s, Namespace: %s, Locale: %s", talentID, region, namespace, locale)

	if talentID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(talentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talent ID"})
		return
	}

	data, err := gamedataService.GetTalentByID(h.Service.GameData, id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving talent: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve talent"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}
