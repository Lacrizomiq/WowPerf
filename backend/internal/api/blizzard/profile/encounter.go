package profile

import (
	"net/http"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/profile"
	wrapper "wowperf/internal/wrapper/blizzard"

	"github.com/gin-gonic/gin"
)

type EncounterSummaryHandler struct {
	Service *blizzard.Service
}

type EncounterDungeonHandler struct {
	Service *blizzard.Service
}

type EncounterRaidHandler struct {
	Service *blizzard.Service
}

func NewEncounterSummaryHandler(service *blizzard.Service) *EncounterSummaryHandler {
	return &EncounterSummaryHandler{
		Service: service,
	}
}

func NewEncounterDungeonHandler(service *blizzard.Service) *EncounterDungeonHandler {
	return &EncounterDungeonHandler{
		Service: service,
	}
}

func NewEncounterRaidHandler(service *blizzard.Service) *EncounterRaidHandler {
	return &EncounterRaidHandler{
		Service: service,
	}
}

// GetCharacterEncounterSummary retrieves a character's encounter summary.
func (h *EncounterSummaryHandler) GetCharacterEncounterSummary(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	encounterSummary, err := profile.GetCharacterEncounterSummary(h.Service.Profile, region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve encounter summary"})
		return
	}

	c.JSON(http.StatusOK, encounterSummary)
}

// GetCharacterEncounterDungeon retrieves a character's dungeon encounters.
func (h *EncounterDungeonHandler) GetCharacterEncounterDungeon(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	encounterDungeon, err := profile.GetCharacterDungeonEncounters(h.Service.Profile, region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dungeon encounters"})
		return
	}

	c.JSON(http.StatusOK, encounterDungeon)
}

// GetCharacterEncounterRaid retrieves a character's raid encounters.
func (h *EncounterRaidHandler) GetCharacterEncounterRaid(c *gin.Context) {
	region := c.Query("region")
	realmSlug := c.Param("realmSlug")
	characterName := c.Param("characterName")
	namespace := c.Query("namespace")
	locale := c.Query("locale")

	if region == "" || realmSlug == "" || characterName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	encounterRaid, err := profile.GetCharacterRaidEncounters(h.Service.Profile, region, realmSlug, characterName, namespace, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve raid encounters"})
		return
	}

	if encounterRaid == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No raid encounters found"})
		return
	}

	transformedData, err := wrapper.TransformRaidData(encounterRaid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transform raid data"})
		return
	}

	c.JSON(http.StatusOK, transformedData)
}
