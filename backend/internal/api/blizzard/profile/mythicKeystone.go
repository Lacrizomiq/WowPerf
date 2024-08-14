package profile

import (
	"net/http"

	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/profile"

	"github.com/gin-gonic/gin"
)

type MythicKeystoneProfileHandler struct {
	Service *blizzard.Service
}

type MythicKeystoneSeasonDetailsHandler struct {
	Service *blizzard.Service
}

func NewMythicKeystoneProfileHandler(service *blizzard.Service) *MythicKeystoneProfileHandler {
	return &MythicKeystoneProfileHandler{
		Service: service,
	}
}

func NewMythicKeystoneSeasonDetailsHandler(service *blizzard.Service) *MythicKeystoneSeasonDetailsHandler {
	return &MythicKeystoneSeasonDetailsHandler{
		Service: service,
	}
}

// GetCharacterMythicKeystoneProfile retrieves a character's mythic keystone profile information, including seasons, tiers, and keystone upgrades.
func (h *MythicKeystoneProfileHandler) GetCharacterMythicKeystoneProfile(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	details, err := profile.GetCharacterMythicKeystoneProfile(h.Service.Profile, region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone season details"})
		return
	}

	c.JSON(http.StatusOK, details)
}

// GetCharacterMythicKeystoneSeasonDetails retrieves a character's mythic keystone season details, including seasons, tiers, and keystone upgrades.
func (h *MythicKeystoneSeasonDetailsHandler) GetCharacterMythicKeystoneSeasonDetails(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	seasonId := c.Param("seasonId")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || seasonId == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	details, err := profile.GetCharacterMythicKeystoneSeasonDetails(h.Service.Profile, region, realmSlug, characterName, seasonId, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone season details"})
		return
	}

	c.JSON(http.StatusOK, details)
}
