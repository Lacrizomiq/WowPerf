package profile

import (
	"log"
	"net/http"
	"strconv"

	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/profile"
	wrapper "wowperf/internal/wrapper/blizzard"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MythicKeystoneProfileHandler struct {
	Service *blizzard.Service
}

type MythicKeystoneSeasonDetailsHandler struct {
	Service *blizzard.Service
	DB      *gorm.DB
}

func NewMythicKeystoneProfileHandler(service *blizzard.Service) *MythicKeystoneProfileHandler {
	return &MythicKeystoneProfileHandler{
		Service: service,
	}
}

func NewMythicKeystoneSeasonDetailsHandler(service *blizzard.Service, db *gorm.DB) *MythicKeystoneSeasonDetailsHandler {
	return &MythicKeystoneSeasonDetailsHandler{
		Service: service,
		DB:      db,
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
func (h *MythicKeystoneSeasonDetailsHandler) GetCharacterMythicKeystoneSeasonBestRuns(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	seasonIdStr := c.Param("seasonId")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || seasonIdStr == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	// Validate season ID
	seasonId, err := strconv.Atoi(seasonIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID format"})
		return
	}

	// Get season slug
	seasonSlug, exists := wrapper.SeasonSlugMapping[seasonId]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown season ID"})
		return
	}

	// Retrieve raw data from blizzard API
	rawData, err := profile.GetCharacterMythicKeystoneSeasonDetails(h.Service.Profile, region, realmSlug, characterName, seasonIdStr, namespace, locale)
	if err != nil {
		log.Printf("Error retrieving mythic keystone season details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic keystone season details"})
		return
	}

	log.Printf("Raw data retrieved successfully for %s-%s, season %s", realmSlug, characterName, seasonIdStr)

	// Transform the raw data into a more usable format from the wrapper
	transformedData, err := wrapper.TransformMythicPlusBestRuns(rawData, h.DB, seasonSlug)
	if err != nil {
		log.Printf("Error transforming mythic keystone season details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transform mythic keystone season details"})
		return
	}

	log.Printf("Data transformed successfully, returning %d runs", len(transformedData))
	c.JSON(http.StatusOK, transformedData)
}
