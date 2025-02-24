package warcraftlogs

import (
	"log"
	"net/http"

	dungeons "wowperf/internal/services/warcraftlogs/dungeons"

	"github.com/gin-gonic/gin"
)

// GlobalLeaderboardAnalysisHandler handles analysis-specific endpoints for high-key M+ data
type GlobalLeaderboardAnalysisHandler struct {
	analysisService *dungeons.GlobalLeaderboardAnalysisService
}

// NewGlobalLeaderboardAnalysisHandler creates a new instance of GlobalLeaderboardAnalysisHandler
func NewGlobalLeaderboardAnalysisHandler(analysisService *dungeons.GlobalLeaderboardAnalysisService) *GlobalLeaderboardAnalysisHandler {
	return &GlobalLeaderboardAnalysisHandler{analysisService: analysisService}
}

// GetSpecGlobalScores returns the average global scores per spec
func (h *GlobalLeaderboardAnalysisHandler) GetSpecGlobalScores(c *gin.Context) {
	scores, err := h.analysisService.GetSpecGlobalScores(c.Request.Context())
	if err != nil {
		log.Printf("Error getting spec global scores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec global scores"})
		return
	}

	c.JSON(http.StatusOK, scores)
}

// GetClassGlobalScores returns the average global scores per class
func (h *GlobalLeaderboardAnalysisHandler) GetClassGlobalScores(c *gin.Context) {
	scores, err := h.analysisService.GetClassGlobalScores(c.Request.Context())
	if err != nil {
		log.Printf("Error getting class global scores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get class global scores"})
		return
	}

	c.JSON(http.StatusOK, scores)
}

// GetSpecDungeonMaxKeyLevels returns the max key levels per spec and dungeon
func (h *GlobalLeaderboardAnalysisHandler) GetSpecDungeonMaxKeyLevels(c *gin.Context) {
	levels, err := h.analysisService.GetSpecDungeonMaxKeyLevels(c.Request.Context())
	if err != nil {
		log.Printf("Error getting spec dungeon max key levels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec dungeon max key levels"})
		return
	}

	c.JSON(http.StatusOK, levels)
}

// GetDungeonAvgKeyLevels returns the average key levels per dungeon
func (h *GlobalLeaderboardAnalysisHandler) GetDungeonAvgKeyLevels(c *gin.Context) {
	levels, err := h.analysisService.GetDungeonAvgKeyLevels(c.Request.Context())
	if err != nil {
		log.Printf("Error getting dungeon avg key levels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dungeon avg key levels"})
		return
	}

	c.JSON(http.StatusOK, levels)
}

// GetTop10PlayersPerSpec returns the top 10 players per spec
func (h *GlobalLeaderboardAnalysisHandler) GetTop10PlayersPerSpec(c *gin.Context) {
	players, err := h.analysisService.GetTop10PlayersPerSpec(c.Request.Context())
	if err != nil {
		log.Printf("Error getting top 10 players per spec: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top 10 players per spec"})
		return
	}

	c.JSON(http.StatusOK, players)
}

// GetTop5PlayersPerRole returns the top 5 players per role
func (h *GlobalLeaderboardAnalysisHandler) GetTop5PlayersPerRole(c *gin.Context) {
	players, err := h.analysisService.GetTop5PlayersPerRole(c.Request.Context())
	if err != nil {
		log.Printf("Error getting top 5 players per role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top 5 players per role"})
		return
	}

	c.JSON(http.StatusOK, players)
}
