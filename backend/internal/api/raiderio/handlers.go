package raiderio

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"wowperf/internal/services/raiderio"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Client *raiderio.Client
}

func NewHandler() *Handler {
	return &Handler{
		Client: raiderio.NewCLient(),
	}
}

func (h *Handler) GetCharacterProfile(c *gin.Context) {
	region := c.Query("region")
	realm := c.Query("realm")
	name := c.Query("name")
	fieldsQuery := c.Query("fields")

	if region == "" || realm == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	fields := parseFields(fieldsQuery)

	talentsRequested := containsField(fields, "talents") || containsField(fields, "talents:categorized")

	profile, err := h.Client.GetCharacterProfile(region, realm, name, fields)
	if err != nil {
		log.Printf("Error getting character profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character profile"})
		return
	}

	if !talentsRequested {
		profile.Talents = nil
		profile.TalentLoadout = nil
	}

	c.JSON(http.StatusOK, profile)
}

func containsField(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

func parseFields(fieldsQuery string) []string {
	if fieldsQuery == "" {
		return nil
	}

	fields := strings.Split(fieldsQuery, ",")
	parsedFields := make([]string, 0, len(fields))

	for _, field := range fields {
		if strings.HasPrefix(field, "mythic_plus_scores_by_season") {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) == 2 {
				seasons := strings.Split(parts[1], ":")
				for _, season := range seasons {
					parsedFields = append(parsedFields, fmt.Sprintf("%s:%s", parts[0], season))
				}
			} else {
				parsedFields = append(parsedFields, field)
			}
		} else {
			parsedFields = append(parsedFields, field)
		}
	}
	return parsedFields
}

func (h *Handler) GetCharacterMythicPlusScores(c *gin.Context) {
	region := c.Query("region")
	realm := c.Query("realm")
	name := c.Query("name")
	seasons := c.QueryArray("season")

	if region == "" || realm == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	if len(seasons) == 0 {
		seasons = []string{"current"}
	}

	fields := make([]string, len(seasons))
	for i, season := range seasons {
		fields[i] = fmt.Sprintf("mythic_plus_scores_by_season:%s", season)
	}

	profile, err := h.Client.GetCharacterProfile(region, realm, name, fields)
	if err != nil {
		log.Printf("Error getting mythic plus scores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mythic plus scores", "details": err.Error()})
		return
	}

	if len(profile.MythicPlusScoresBySeason) > 0 {
		c.JSON(http.StatusOK, profile.MythicPlusScoresBySeason)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Mythic+ scores found for the specified seasons"})
	}
}

func (h Handler) GetCharacterRaidProgression(c *gin.Context) {
	region := c.Query("region")
	realm := c.Query("realm")
	name := c.Query("name")

	if region == "" || realm == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters"})
		return
	}

	fields := []string{"raid_progression"}

	profile, err := h.Client.GetCharacterProfile(region, realm, name, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile.RaidProgression)
}
