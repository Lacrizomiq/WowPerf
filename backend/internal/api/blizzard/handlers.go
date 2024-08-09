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

func (h *Handler) GetCharacterSpecializations(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	specializations, err := h.Client.GetCharacterSpecializations(region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character specializations"})
		return
	}

	c.JSON(http.StatusOK, specializations)
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
