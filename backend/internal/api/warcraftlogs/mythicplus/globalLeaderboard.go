package warcraftlogs

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	dungeons "wowperf/internal/services/warcraftlogs/dungeons"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gin-gonic/gin"
)

type GlobalLeaderboardHandler struct {
	rankingsService *dungeons.RankingsService
}

func NewGlobalLeaderboardHandler(rankingsService *dungeons.RankingsService) *GlobalLeaderboardHandler {
	return &GlobalLeaderboardHandler{rankingsService: rankingsService}
}

// Helper to format the string with the first letter in uppercase
func formatNameCase(s string) string {
	if s == "" {
		return s
	}
	caser := cases.Title(language.English)
	return caser.String(strings.ToLower(strings.TrimSpace(s)))
}

// Helper to format the role with the first letter in uppercase
func formatRole(role string) string {
	formatted := strings.ToUpper(strings.TrimSpace(role))
	if formatted == "DPS" {
		return formatted
	}
	return formatNameCase(role)
}

// Helper to handle sorting parameters
func getOrderParams(c *gin.Context) (string, dungeons.OrderDirection) {
	orderBy := c.DefaultQuery("orderBy", "score")
	direction := dungeons.OrderDirection(strings.ToUpper(c.DefaultQuery("direction", "DESC")))

	// Validate direction
	if direction != dungeons.ASC && direction != dungeons.DESC {
		direction = dungeons.DESC
	}

	return orderBy, direction
}

// GetGlobalLeaderboard returns the global leaderboard
func (h *GlobalLeaderboardHandler) GetGlobalLeaderboard(c *gin.Context) {

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	orderBy, direction := getOrderParams(c)

	leaderboard, err := h.rankingsService.GetGlobalLeaderboard(c.Request.Context(), limit, orderBy, direction)
	if err != nil {
		log.Printf("Error getting global leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get global leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// GetRoleLeaderboard returns the leaderboard for a specific role
func (h *GlobalLeaderboardHandler) GetRoleLeaderboard(c *gin.Context) {
	role := formatRole(c.Query("role"))
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role parameter is required"})
		return
	}

	// Validation du rôle
	validRoles := map[string]bool{
		"Tank":   true,
		"Healer": true,
		"DPS":    true,
	}

	if !validRoles[role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid role. Must be one of: Tank, Healer, DPS. Got: %s", role),
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	orderBy, direction := getOrderParams(c)

	leaderboard, err := h.rankingsService.GetGlobalLeaderboardByRole(c.Request.Context(), role, limit, orderBy, direction)
	if err != nil {
		log.Printf("Error getting role leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// GetClassLeaderboard returns the leaderboard for a specific class
func (h *GlobalLeaderboardHandler) GetClassLeaderboard(c *gin.Context) {
	class := formatNameCase(c.Query("class"))
	if class == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class parameter is required"})
		return
	}

	// Validation de la classe
	validClasses := map[string]bool{
		"Warrior":     true,
		"Paladin":     true,
		"Hunter":      true,
		"Rogue":       true,
		"Priest":      true,
		"Shaman":      true,
		"Mage":        true,
		"Warlock":     true,
		"Monk":        true,
		"Druid":       true,
		"Demonhunter": true,
		"Deathknight": true,
		"Evoker":      true,
	}

	if !validClasses[class] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid class: %s", class),
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	orderBy, direction := getOrderParams(c)

	leaderboard, err := h.rankingsService.GetGlobalLeaderboardByClass(c.Request.Context(), class, limit, orderBy, direction)
	if err != nil {
		log.Printf("Error getting class leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get class leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// GetSpecLeaderboard returns the leaderboard for a specific spec
func (h *GlobalLeaderboardHandler) GetSpecLeaderboard(c *gin.Context) {

	class := formatNameCase(c.Query("class"))
	if class == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class parameter is required"})
		return
	}

	spec := formatNameCase(c.Query("spec"))
	if spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spec parameter is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	orderBy, direction := getOrderParams(c)

	leaderboard, err := h.rankingsService.GetGlobalLeaderboardBySpec(c.Request.Context(), class, spec, limit, orderBy, direction)
	if err != nil {
		log.Printf("Error getting spec leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}
