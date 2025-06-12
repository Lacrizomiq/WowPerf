// internal/api/warcraftlogs/mythicplus/specEvolutionMetricsAnalysis.go
package warcraftlogs

import (
	"log"
	"net/http"
	"strconv"
	"time"

	dungeons "wowperf/internal/services/warcraftlogs/dungeons"

	"github.com/gin-gonic/gin"
)

// SpecEvolutionMetricsAnalysisHandler handles evolution-specific endpoints for spec metrics
type SpecEvolutionMetricsAnalysisHandler struct {
	evolutionService *dungeons.SpecEvolutionMetricsAnalysisService
}

// NewSpecEvolutionMetricsAnalysisHandler creates a new instance of SpecEvolutionMetricsAnalysisHandler
func NewSpecEvolutionMetricsAnalysisHandler(evolutionService *dungeons.SpecEvolutionMetricsAnalysisService) *SpecEvolutionMetricsAnalysisHandler {
	return &SpecEvolutionMetricsAnalysisHandler{evolutionService: evolutionService}
}

// GetSpecEvolution retrieves the evolution of metrics for a specialization
// GET /warcraftlogs/mythicplus/evolution/spec?spec=Restoration&class=Druid&period=7&dungeon_id=1&date=2023-05-11
func (h *SpecEvolutionMetricsAnalysisHandler) GetSpecEvolution(c *gin.Context) {
	spec := c.Query("spec")
	if spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "spec parameter is required"})
		return
	}

	var class *string
	classParam := c.Query("class")
	if classParam != "" {
		class = &classParam
	}

	periodStr := c.DefaultQuery("period", "7")
	dungeonIDStr := c.Query("dungeon_id")
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	// Validate parameters
	period, err := strconv.Atoi(periodStr)
	if err != nil || (period != 7 && period != 30) {
		period = 7 // Default value
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = time.Now().Truncate(24 * time.Hour)
	}

	var dungeonID *int
	if dungeonIDStr != "" {
		dungeonIDInt, err := strconv.Atoi(dungeonIDStr)
		if err == nil {
			dungeonID = &dungeonIDInt
		}
	}

	isGlobal := dungeonID == nil

	// Get evolution data
	results, err := h.evolutionService.GetSpecEvolution(c.Request.Context(), spec, class, period, dungeonID, date, isGlobal)
	if err != nil {
		log.Printf("Error getting spec evolution: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec evolution data"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetCurrentRanking retrieves the current ranking of all specializations
// GET /warcraftlogs/mythicplus/evolution/ranking?role=dps&date=2023-05-11&is_global=true
func (h *SpecEvolutionMetricsAnalysisHandler) GetCurrentRanking(c *gin.Context) {
	roleStr := c.Query("role")
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	isGlobalStr := c.DefaultQuery("is_global", "true")

	// Validate parameters
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = time.Now().Truncate(24 * time.Hour)
	}

	isGlobal, err := strconv.ParseBool(isGlobalStr)
	if err != nil {
		isGlobal = true
	}

	var rolePtr *string
	if roleStr != "" {
		rolePtr = &roleStr
	}

	// Get ranking data
	results, err := h.evolutionService.GetCurrentRanking(c.Request.Context(), rolePtr, date, isGlobal)
	if err != nil {
		log.Printf("Error getting current ranking: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current ranking data"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetLatestMetricsDate retrieves the latest date for which metrics are available
// GET /warcraftlogs/mythicplus/evolution/latest-date
func (h *SpecEvolutionMetricsAnalysisHandler) GetLatestMetricsDate(c *gin.Context) {
	latestDate, err := h.evolutionService.GetLatestMetricsDate(c.Request.Context())
	if err != nil {
		log.Printf("Error getting latest metrics date: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest metrics date"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"latest_date": latestDate.Format("2006-01-02")})
}

// GetClassSpecs retrieves all available specs for a class
// GET /warcraftlogs/mythicplus/evolution/class/specs?class=Druid
func (h *SpecEvolutionMetricsAnalysisHandler) GetClassSpecs(c *gin.Context) {
	className := c.Query("class")
	if className == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class parameter is required"})
		return
	}

	specs, err := h.evolutionService.GetClassSpecs(c.Request.Context(), className)
	if err != nil {
		log.Printf("Error getting class specs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get class specs"})
		return
	}

	c.JSON(http.StatusOK, specs)
}

// GetAvailableClasses retrieves all available classes
// GET /warcraftlogs/mythicplus/evolution/classes
func (h *SpecEvolutionMetricsAnalysisHandler) GetAvailableClasses(c *gin.Context) {
	classes, err := h.evolutionService.GetAvailableClasses(c.Request.Context())
	if err != nil {
		log.Printf("Error getting available classes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get available classes"})
		return
	}

	c.JSON(http.StatusOK, classes)
}

// GetRoleSpecs retrieves all specs for a given role
// GET /warcraftlogs/mythicplus/evolution/role/specs?role=dps
func (h *SpecEvolutionMetricsAnalysisHandler) GetRoleSpecs(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role parameter is required"})
		return
	}

	specs, err := h.evolutionService.GetRoleSpecs(c.Request.Context(), role)
	if err != nil {
		log.Printf("Error getting role specs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role specs"})
		return
	}

	c.JSON(http.StatusOK, specs)
}

// GetSpecHistoricalData retrieves historical data for a spec over multiple dates
// GET /warcraftlogs/mythicplus/evolution/spec/history?spec=Restoration&class=Druid&days=30&is_global=true&dungeon_id=1
func (h *SpecEvolutionMetricsAnalysisHandler) GetSpecHistoricalData(c *gin.Context) {
	spec := c.Query("spec")
	if spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "spec parameter is required"})
		return
	}

	var class *string
	classParam := c.Query("class")
	if classParam != "" {
		class = &classParam
	}

	daysStr := c.DefaultQuery("days", "30")
	isGlobalStr := c.DefaultQuery("is_global", "true")
	dungeonIDStr := c.Query("dungeon_id")

	// Validate parameters
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30 // Default value
	}

	isGlobal, err := strconv.ParseBool(isGlobalStr)
	if err != nil {
		isGlobal = true
	}

	var dungeonID *int
	if dungeonIDStr != "" {
		dungeonIDInt, err := strconv.Atoi(dungeonIDStr)
		if err == nil {
			dungeonID = &dungeonIDInt
		}
	}

	// Get historical data
	results, err := h.evolutionService.GetSpecHistoricalData(c.Request.Context(), spec, class, days, isGlobal, dungeonID)
	if err != nil {
		log.Printf("Error getting spec historical data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec historical data"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetAvailableDungeons retrieves all available dungeons with metrics
// GET /warcraftlogs/mythicplus/evolution/dungeons
func (h *SpecEvolutionMetricsAnalysisHandler) GetAvailableDungeons(c *gin.Context) {
	dungeons, err := h.evolutionService.GetAvailableDungeons(c.Request.Context())
	if err != nil {
		log.Printf("Error getting available dungeons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get available dungeons"})
		return
	}

	c.JSON(http.StatusOK, dungeons)
}

// GetTopSpecsForDungeon retrieves the top performing specs for a specific dungeon
// GET /warcraftlogs/mythicplus/evolution/dungeons/top-specs?dungeon_id=1&limit=10
func (h *SpecEvolutionMetricsAnalysisHandler) GetTopSpecsForDungeon(c *gin.Context) {
	dungeonIDStr := c.Query("dungeon_id")
	if dungeonIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dungeon_id parameter is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")

	// Validate parameters
	dungeonID, err := strconv.Atoi(dungeonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dungeon_id format"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 // Default value
	}

	// Get top specs
	results, err := h.evolutionService.GetTopSpecsForDungeon(c.Request.Context(), dungeonID, limit)
	if err != nil {
		log.Printf("Error getting top specs for dungeon: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top specs for dungeon"})
		return
	}

	c.JSON(http.StatusOK, results)
}

/*

List of endpoints:


1. GetSpecEvolution

GET /warcraftlogs/mythicplus/evolution/spec

Parameters:
- spec: string, mandatory, name of the spec
- class: string, optional, name of the class
- period: int, optional, 7 or 30
- dungeon_id: int, optional, id of the dungeon
- date: string, optional, date in the format YYYY-MM-DD

Example :
GET /warcraftlogs/mythicplus/evolution/spec?spec=Restoration&class=Druid&period=7&dungeon_id=12648&date=2025-05-11

2. GetCurrentRanking

GET /warcraftlogs/mythicplus/evolution/ranking

Parameters:
- role: string, optional, name of the role
- date: string, optional, date in the format YYYY-MM-DD
- is_global: boolean, optional, true or false

Example :
GET /warcraftlogs/mythicplus/evolution/ranking?role=dps&date=2025-05-11&is_global=true

3. GetLatestMetricsDate

GET /warcraftlogs/mythicplus/evolution/latest-date

Parameters:
- None

Example :
GET /warcraftlogs/mythicplus/evolution/latest-date

4. GetClassSpecs

GET /warcraftlogs/mythicplus/evolution/class/specs

Parameters:
- class: string, mandatory, name of the class

Example :
GET /warcraftlogs/mythicplus/evolution/class/specs?class=Druid

5. GetAvailableClasses

GET /warcraftlogs/mythicplus/evolution/classes

Parameters:
- None

Example :
GET /warcraftlogs/mythicplus/evolution/classes

6. GetRoleSpecs

GET /warcraftlogs/mythicplus/evolution/role/specs

Parameters:
- role: string, mandatory, name of the role (dps, tank, healer)

Example :
GET /warcraftlogs/mythicplus/evolution/role/specs?role=dps

7. GetSpecHistoricalData

GET /warcraftlogs/mythicplus/evolution/spec/history

Parameters:
- spec: string, mandatory, name of the spec
- class: string, optional, name of the class
- days: int, optional, number of days to get data for
- is_global: boolean, optional, true or false
- dungeon_id: int, optional, id of the dungeon

Example :
GET /warcraftlogs/mythicplus/evolution/spec/history?spec=Restoration&class=Druid&days=30&is_global=true&dungeon_id=12648

8. GetAvailableDungeons

GET /warcraftlogs/mythicplus/evolution/dungeons

Parameters:
- None

Example :
GET /warcraftlogs/mythicplus/evolution/dungeons

9. GetTopSpecsForDungeon

GET /warcraftlogs/mythicplus/evolution/dungeons/top-specs

Parameters:
- dungeon_id: int, mandatory, id of the dungeon
- limit: int, optional, number of specs to get

Example :
GET /warcraftlogs/mythicplus/evolution/dungeons/top-specs?dungeon_id=12648&limit=10


*/

/*

Use cases:

1. Follow the spec evolution for Restoration Druid

# Get the global evolution over 7 days
GET /warcraftlogs/mythicplus/evolution/spec?spec=Restoration&class=Druid&period=7

# Get the evolution over a specific dungeon
GET /warcraftlogs/mythicplus/evolution/spec?spec=Restoration&class=Druid&period=7&dungeon_id=12648

# Get the complete history over 30 days
GET /warcraftlogs/mythicplus/evolution/spec/history?spec=Restoration&class=Druid&days=30

2. Analyse performance over time per role

# Get the DPS ranking
GET /warcraftlogs/mythicplus/evolution/ranking?role=dps

# List all DPS specs
GET /warcraftlogs/mythicplus/evolution/role/specs?role=dps

# Get the top specs for a specific dungeon
GET /warcraftlogs/mythicplus/evolution/dungeons/top-specs?dungeon_id=61594&limit=10

*/
