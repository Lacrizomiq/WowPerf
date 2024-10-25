package warcraftlogs

import (
	"log"
	"net/http"
	"strconv"

	ranking "wowperf/internal/services/warcraftlogs/dungeons"

	"github.com/gin-gonic/gin"
)

type GlobalLeaderboardHandler struct {
	globalLeaderboardService *ranking.RankingsService
}

func NewGlobalLeaderboardHandler(globalLeaderboardService *ranking.RankingsService) *GlobalLeaderboardHandler {
	return &GlobalLeaderboardHandler{
		globalLeaderboardService: globalLeaderboardService,
	}
}

// Get the global leaderboard in every role
func (h *GlobalLeaderboardHandler) GetRoleLeaderboard(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role parameter is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	entries, err := h.globalLeaderboardService.GetGlobalLeaderboardByRole(c.Request.Context(), role, limit)
	if err != nil {
		log.Printf("Error getting role leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role leaderboard"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// Get the global leaderboard in every class
func (h *GlobalLeaderboardHandler) GetClassLeaderboard(c *gin.Context) {
	class := c.Query("class")
	if class == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class parameter is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	entries, err := h.globalLeaderboardService.GetGlobalLeaderboardByClass(c.Request.Context(), class, limit)
	if err != nil {
		log.Printf("Error getting class leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get class leaderboard"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// Get the global leaderboard in every spec
func (h *GlobalLeaderboardHandler) GetSpecLeaderboard(c *gin.Context) {
	class := c.Query("class")
	if class == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class parameter is required"})
		return
	}

	spec := c.Query("spec")
	if spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spec parameter is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	entries, err := h.globalLeaderboardService.GetGlobalLeaderboardBySpec(c.Request.Context(), class, spec, limit)
	if err != nil {
		log.Printf("Error getting spec leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec leaderboard"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// Get the global leaderboard in every role, class and spec
func (h *GlobalLeaderboardHandler) GetGlobalLeaderboard(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	entries, err := h.globalLeaderboardService.GetGlobalLeaderboard(c.Request.Context(), limit)
	if err != nil {
		log.Printf("Error getting global leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get global leaderboard"})
		return
	}

	c.JSON(http.StatusOK, entries)
}
