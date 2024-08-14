package profile

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"wowperf/internal/services/blizzard"
	profileService "wowperf/internal/services/blizzard/profile"
	wrapper "wowperf/internal/wrapper/blizzard"

	"github.com/gin-gonic/gin"
)

type SpecializationsHandler struct {
	Service *blizzard.Service
}

func NewSpecializationsHandler(service *blizzard.Service) *SpecializationsHandler {
	return &SpecializationsHandler{
		Service: service,
	}
}

// GetCharacterSpecializations retrieves a character's specializations, including spec groups, specs, and spec tiers.
func (h *SpecializationsHandler) GetCharacterSpecializations(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	profileNamespace := c.Query("namespace")
	locale := c.Query("locale")

	staticNamespace := strings.Replace(profileNamespace, "profile", "static", 1)

	log.Printf("Fetching specializations for %s-%s (region: %s, profile namespace: %s, static namespace: %s, locale: %s)", realmSlug, characterName, region, profileNamespace, staticNamespace, locale)

	// first i retrieve the character profile
	characterData, err := profileService.GetCharacterProfile(h.Service.Profile, region, realmSlug, characterName, profileNamespace, locale)
	if err != nil {
		log.Printf("Error fetching character profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve character profile: %v", err)})
		return
	}

	// then transform the character profile into a struct
	characterProfile, err := wrapper.TransformCharacterInfo(characterData, nil)
	if err != nil {
		log.Printf("Error transforming character profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to transform character profile: %v", err)})
		return
	}
	// then i retrieve the specializations
	specializations, err := profileService.GetCharacterSpecializations(h.Service.Profile, region, realmSlug, characterName, profileNamespace, locale)
	if err != nil {
		log.Printf("Error fetching character specializations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve character specializations: %v", err)})
		return
	}

	// then i transform the specializations into a talent loadout
	talentLoadout, err := wrapper.TransformCharacterTalents(specializations, h.Service.GameDataClient, region, staticNamespace, locale, characterProfile.TreeID, characterProfile.SpecID)
	if err != nil {
		log.Printf("Error transforming character talents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to transform character talents: %v", err)})
		return
	}

	c.JSON(http.StatusOK, talentLoadout)
}
