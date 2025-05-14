package warcraftlogs

import (
	"log"
	"net/http"
	"strconv"

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
// @Summary Get average global scores per spec
// @Description Get average global scores per spec (excludes CN players)
// @Tags Mythic+ Performance Analysis
// @Accept json
// @Produce json
// @Success 200 {array} dungeons.SpecGlobalScore
// @Failure 500 {object} gin.H
// @Router /warcraftlogs/mythicplus/analysis/specs/avg-scores [get]
func (h *GlobalLeaderboardAnalysisHandler) GetSpecGlobalScores(c *gin.Context) {
	scores, err := h.analysisService.GetSpecGlobalScores(c.Request.Context())
	if err != nil {
		log.Printf("Error getting spec global scores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec global scores"})
		return
	}

	c.JSON(http.StatusOK, scores)
}

// GetSpecDungeonScoreAverages returns the average scores per spec and dungeon
// @Summary Get average scores per spec and dungeon
// @Description Get average scores per spec and dungeon with optional filters (excludes CN players)
// @Tags Mythic+ Performance Analysis
// @Accept json
// @Produce json
// @Param class query string false "Filter by class name"
// @Param spec query string false "Filter by spec name"
// @Param encounter_id query int false "Filter by dungeon encounter ID"
// @Success 200 {array} dungeons.SpecDungeonScoreAverage
// @Failure 500 {object} gin.H
// @Router /warcraftlogs/mythicplus/analysis/specs/dungeons/avg-scores [get]
func (h *GlobalLeaderboardAnalysisHandler) GetSpecDungeonScoreAverages(c *gin.Context) {
	// Parse query parameters
	class := c.Query("class")
	spec := c.Query("spec")

	encounterID := int64(0)
	if encID := c.Query("encounter_id"); encID != "" {
		if parsedID, err := strconv.ParseInt(encID, 10, 64); err == nil {
			encounterID = parsedID
		}
	}

	scores, err := h.analysisService.GetSpecDungeonScoreAverages(c.Request.Context(), class, spec, encounterID)
	if err != nil {
		log.Printf("Error getting spec dungeon score averages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec dungeon score averages"})
		return
	}

	c.JSON(http.StatusOK, scores)
}

// GetClassGlobalScores returns the average global scores per class
// @Summary Get average global scores per class
// @Description Get average global scores per class
// @Tags Mythic+ Performance Analysis
// @Accept json
// @Produce json
// @Success 200 {array} dungeons.ClassGlobalScore
// @Failure 500 {object} gin.H
// @Router /warcraftlogs/mythicplus/analysis/classes/avg-scores [get]
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
// @Summary Get max key levels per spec and dungeon
// @Description Get maximum key levels achieved per spec and dungeon (excludes CN players)
// @Tags Mythic+ Performance Analysis
// @Accept json
// @Produce json
// @Success 200 {array} dungeons.SpecDungeonMaxKeyLevel
// @Failure 500 {object} gin.H
// @Router /warcraftlogs/mythicplus/analysis/specs/dungeons/max-levels-key [get]
func (h *GlobalLeaderboardAnalysisHandler) GetSpecDungeonMaxKeyLevels(c *gin.Context) {
	levels, err := h.analysisService.GetSpecDungeonMaxKeyLevels(c.Request.Context())
	if err != nil {
		log.Printf("Error getting spec dungeon max key levels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec dungeon max key levels"})
		return
	}

	c.JSON(http.StatusOK, levels)
}

// GetDungeonMedia retrieves the dungeons media
// @Summary Get dungeon media information
// @Description Get media information for dungeons including icons and images
// @Tags Mythic+ Performance Analysis
// @Accept json
// @Produce json
// @Success 200 {array} dungeons.DungeonMedia
// @Failure 500 {object} gin.H
// @Router /warcraftlogs/mythicplus/analysis/specs/dungeons/media [get]
func (h *GlobalLeaderboardAnalysisHandler) GetDungeonMedia(c *gin.Context) {
	media, err := h.analysisService.GetDungeonMedia(c.Request.Context())
	if err != nil {
		log.Printf("Error getting the dungeon media: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dungeon media"})
		return
	}

	c.JSON(http.StatusOK, media)
}

// GetDungeonAvgKeyLevels returns the average key levels per dungeon
// @Summary Get average key levels per dungeon
// @Description Get average key levels across all dungeons
// @Tags Mythic+ Performance Analysis
// @Accept json
// @Produce json
// @Success 200 {array} dungeons.DungeonAvgKeyLevel
// @Failure 500 {object} gin.H
// @Router /warcraftlogs/mythicplus/analysis/dungeons/avg-levels-key [get]
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
// @Summary Get top 10 players per spec
// @Description Get the top 10 highest-scoring players for each spec
// @Tags Mythic+ Performance Analysis
// @Accept json
// @Produce json
// @Success 200 {array} dungeons.TopPlayerPerSpec
// @Failure 500 {object} gin.H
// @Router /warcraftlogs/mythicplus/analysis/players/top-specs [get]
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
// @Summary Get top 5 players per role
// @Description Get the top 5 highest-scoring players for each role (Tank, Healer, DPS)
// @Tags Mythic+ Performance Analysis
// @Accept json
// @Produce json
// @Success 200 {array} dungeons.TopPlayerPerRole
// @Failure 500 {object} gin.H
// @Router /warcraftlogs/mythicplus/analysis/players/top-roles [get]
func (h *GlobalLeaderboardAnalysisHandler) GetTop5PlayersPerRole(c *gin.Context) {
	players, err := h.analysisService.GetTop5PlayersPerRole(c.Request.Context())
	if err != nil {
		log.Printf("Error getting top 5 players per role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top 5 players per role"})
		return
	}

	c.JSON(http.StatusOK, players)
}

/*
Mythic+ Analysis API Endpoints Documentation

1. GET /warcraftlogs/mythicplus/analysis/specs/avg-scores
   Desc: Get average global scores per spec (excludes CN players)
   Response: Includes avg_global_score, max_global_score, min_global_score, player_count, overall_rank, role_rank
   Example: GET /warcraftlogs/mythicplus/analysis/specs/avg-scores

2. GET /warcraftlogs/mythicplus/analysis/specs/dungeons/avg-scores
	Desc: Get average scores per spec and dungeon (excludes CN players)
	Query params:
		- class (optional): Filter by class name
		- spec (optional): Filter by spec name
		- encounter_id (optional): Filter by dungeon encounter ID
	Response: Includes avg_dungeon_score, max_score, min_score, player_count, overall_rank, role_rank
	Examples:
		- GET /warcraftlogs/mythicplus/analysis/specs/dungeons/avg-scores
		- GET /warcraftlogs/mythicplus/analysis/specs/dungeons/avg-scores?class=Warrior&spec=Arms
		- GET /warcraftlogs/mythicplus/analysis/specs/dungeons/avg-scores?encounter_id=12651
		- GET /warcraftlogs/mythicplus/analysis/specs/dungeons/avg-scores?class=Warrior&spec=Arms&encounter_id=12651

3. GET /warcraftlogs/mythicplus/analysis/classes/avg-scores
   Desc: Get average global scores per class
   Example: GET /warcraftlogs/mythicplus/analysis/classes/avg-scores

4. GET /warcraftlogs/mythicplus/analysis/specs/dungeons/max-levels-key
   Desc: Get max key levels per spec and dungeon (excludes CN players)
   Response: Includes encounter_id
   Example: GET /warcraftlogs/mythicplus/analysis/specs/dungeons/max-levels-key

5. GET /warcraftlogs/mythicplus/analysis/specs/dungeons/media
   Desc: Get dungeon media (icons, images)
   Example: GET /warcraftlogs/mythicplus/analysis/specs/dungeons/media

6. GET /warcraftlogs/mythicplus/analysis/dungeons/avg-levels-key
   Desc: Get average key levels per dungeon
   Example: GET /warcraftlogs/mythicplus/analysis/dungeons/avg-levels-key

7. GET /warcraftlogs/mythicplus/analysis/players/top-specs
   Desc: Get top 10 players per spec
   Example: GET /warcraftlogs/mythicplus/analysis/players/top-specs

8. GET /warcraftlogs/mythicplus/analysis/players/top-roles
   Desc: Get top 5 players per role
   Example: GET /warcraftlogs/mythicplus/analysis/players/top-roles
*/
