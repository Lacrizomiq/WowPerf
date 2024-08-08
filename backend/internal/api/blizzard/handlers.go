package blizzard

import (
	"net/http"

	"wowperf/internal/services/blizzard"
	wrapper "wowperf/internal/wrapper/blizzard"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Client *blizzard.Client
}

func NewHandler() (*Handler, error) {
	client, err := blizzard.NewClient()
	if err != nil {
		return nil, err
	}

	return &Handler{Client: client}, nil
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

func (h *Handler) GetCharacterEquipment(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	equipment, err := h.Client.GetCharacterEquipment(region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character equipment"})
		return
	}

	c.JSON(http.StatusOK, equipment)
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
