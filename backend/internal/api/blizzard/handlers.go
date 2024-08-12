package blizzard

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"wowperf/internal/services/blizzard"
	wrapper "wowperf/internal/wrapper/blizzard"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Client         *blizzard.Client
	GameDataClient *blizzard.GameDataClient
}

func NewHandler() (*Handler, error) {
	client, err := blizzard.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Blizzard client: %w", err)
	}

	gameDataClient, err := blizzard.NewGameDataClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Game Data client: %w", err)
	}

	return &Handler{
		Client:         client,
		GameDataClient: gameDataClient,
	}, nil
}

// GetCharacterProfile retrieves a character's profile information, including name, realm, and class.
func (h *Handler) GetCharacterProfile(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	characterData, err := h.Client.GetCharacterProfile(region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character profile"})
		return
	}

	mediaData, err := h.Client.GetCharacterMedia(region, realmSlug, characterName, namespace, locale)
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

func (h *Handler) GetCharacterMythicKeystoneProfile(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	details, err := h.Client.GetCharacterMythicKeystoneProfile(region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone season details"})
		return
	}

	c.JSON(http.StatusOK, details)
}

func (h *Handler) GetCharacterMythicKeystoneSeasonDetails(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	seasonId := c.Param("seasonId")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	details, err := h.Client.GetCharacterMythicKeystoneSeasonDetails(region, realmSlug, characterName, seasonId, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone season details"})
		return
	}

	c.JSON(http.StatusOK, details)
}

// GetCharacterEquipment retrieves a character's equipment, including items, gems, and enchantments.
func (h *Handler) GetCharacterEquipment(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	equipmentData, err := h.Client.GetCharacterEquipment(region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character equipment"})
		return
	}

	transformedGear, err := wrapper.TransformCharacterGear(equipmentData, h.GameDataClient, region, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transform character equipment"})
		return
	}

	c.JSON(http.StatusOK, transformedGear)
}

// GetCharacterSpecializations retrieves a character's specializations, including spec groups, specs, and spec tiers.
func (h *Handler) GetCharacterSpecializations(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	log.Printf("Fetching specializations for %s-%s (region: %s, namespace: %s, locale: %s)", realmSlug, characterName, region, namespace, locale)

	specializations, err := h.Client.GetCharacterSpecializations(region, realmSlug, characterName, namespace, locale)
	if err != nil {
		log.Printf("Error fetching character specializations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve character specializations: %v", err)})
		return
	}

	talentLoadout, err := wrapper.TransformCharacterTalents(specializations, h.GameDataClient, region, namespace, locale)
	if err != nil {
		log.Printf("Error transforming character talents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to transform character talents: %v", err)})
		return
	}

	c.JSON(http.StatusOK, talentLoadout)
}

// GetCharacterMedia retrieves a character's media assets, including avatar, inset avatar, main raw, and character media.
func (h *Handler) GetCharacterMedia(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	media, err := h.Client.GetCharacterMedia(region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character media"})
		return
	}

	c.JSON(http.StatusOK, media)
}

// GetItemMedia retrieves the media assets for an item.
func (h *Handler) GetItemMedia(c *gin.Context) {
	if h.GameDataClient == nil {
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

	mediaData, err := h.GameDataClient.GetItemMedia(id, region, namespace, locale)
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

// GetPlayableSpecializationIndex retrieves an index of playable specializations
func (h *Handler) GetPlayableSpecializationIndex(c *gin.Context) {
	if h.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	log.Printf("Requesting playable specialization index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := h.GameDataClient.GetPlayableSpecializationIndex(region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving playable specialization index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve playable specialization index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetPlayableSpecialization retrieves a playable specialization
func (h *Handler) GetPlayableSpecialization(c *gin.Context) {
	if h.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	specID := c.Param("specId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting playable specialization for SpecID: %s, Region: %s, Namespace: %s, Locale: %s", specID, region, namespace, locale)

	if specID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(specID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
		return
	}

	data, err := h.GameDataClient.GetPlayableSpecialization(id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving playable specialization: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playable specialization not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve playable specialization"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetPlayableSpecializationMedia retrieves the media assets for a playable specialization
func (h *Handler) GetPlayableSpecializationMedia(c *gin.Context) {
	if h.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	specID := c.Param("specId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting playable specialization media for SpecID: %s, Region: %s, Namespace: %s, Locale: %s", specID, region, namespace, locale)

	if specID == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	id, err := strconv.Atoi(specID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
		return
	}

	data, err := h.GameDataClient.GetPlayableSpecializationMedia(id, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving playable specialization media: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playable specialization media not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve playable specialization media"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetTalentTreeIndex retrieves an index of talent trees
func (h *Handler) GetTalentTreeIndex(c *gin.Context) {
	if h.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	log.Printf("Requesting talent tree index for Region: %s, Namespace: %s, Locale: %s", region, namespace, locale)

	index, err := h.GameDataClient.GetTalentTreeIndex(region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving talent tree index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve talent tree index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, index)
}

// GetTalentTree retrieves a talent tree by spec ID
func (h *Handler) GetTalentTree(c *gin.Context) {
	if h.GameDataClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Game Data Client not initialized"})
		return
	}

	talentTreeID := c.Param("talentTreeId")
	specID := c.Param("specId")
	region := c.Query("region")
	namespace := c.DefaultQuery("namespace", fmt.Sprintf("static-%s", region))
	locale := c.DefaultQuery("locale", "en_US")

	if namespace == "" {
		namespace = fmt.Sprintf("static-%s", region)
	}

	log.Printf("Requesting talent tree for TalentTreeID: %s, SpecID: %s, Region: %s, Namespace: %s, Locale: %s", talentTreeID, specID, region, namespace, locale)

	if talentTreeID == "" || specID == "" || region == "" {
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

	data, err := h.GameDataClient.GetTalentTree(treeID, specIDInt, region, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving talent tree: %v", err)
		if strings.Contains(err.Error(), "404") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent tree not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve talent tree"})
		}
		return
	}

	c.JSON(http.StatusOK, data)
}
