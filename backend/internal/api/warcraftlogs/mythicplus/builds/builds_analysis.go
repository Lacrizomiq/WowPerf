package WarcraftLogsMythicPlusBuildsAnalysis

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	service "wowperf/internal/services/warcraftlogs/mythicplus/analytics"
)

// BuildsAnalysisHandler handles API endpoints for the mythic + builds analysis
type MythicPlusBuildsAnalysisHandler struct {
	MythicPlusBuildsAnalysisService *service.BuildAnalysisService
}

// NewMythicPlusBuildsAnalysisHandler creates a new MythicPlusBuildsAnalysisHandler
func NewMythicPlusBuildsAnalysisHandler(analysisService *service.BuildAnalysisService) *MythicPlusBuildsAnalysisHandler {
	return &MythicPlusBuildsAnalysisHandler{MythicPlusBuildsAnalysisService: analysisService}
}

// GetPopularItemsBySlot returns the most popular items for each slot for a specific class and spec
// @Summary Get popular items by slot
// @Description Returns the most popular items for each slot for a specific class and spec
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Param spec query string true "Specialization name"
// @Param encounter_id query int false "Encounter ID to filter results"
// @Success 200 {array} service.ItemPopularity
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/items [get]
func (h *MythicPlusBuildsAnalysisHandler) GetPopularItemsBySlot(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))
	spec := normalizeWoWTerms(c.Query("spec"))

	if class == "" || spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class and spec parameters are required"})
		return
	}

	var encounterID *int
	if encIDStr := c.Query("encounter_id"); encIDStr != "" {
		encID, err := strconv.Atoi(encIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid encounter_id format"})
			return
		}
		encounterID = &encID
	}

	items, err := h.MythicPlusBuildsAnalysisService.GetPopularItemsBySlot(c.Request.Context(), class, spec, encounterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetGlobalPopularItemsBySlot returns the most popular items for each slot across all encounters
// @Summary Get global popular items by slot
// @Description Returns the most popular items for each slot across all encounters
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Param spec query string true "Specialization name"
// @Success 200 {array} service.GlobalItemPopularity
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/items/global [get]
func (h *MythicPlusBuildsAnalysisHandler) GetGlobalPopularItemsBySlot(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))
	spec := normalizeWoWTerms(c.Query("spec"))

	if class == "" || spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class and spec parameters are required"})
		return
	}

	items, err := h.MythicPlusBuildsAnalysisService.GetGlobalPopularItemsBySlot(c.Request.Context(), class, spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetEnchantUsage returns enchant usage statistics for a specific class and spec
// @Summary Get enchant usage
// @Description Returns enchant usage statistics for a specific class and spec
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Param spec query string true "Specialization name"
// @Success 200 {array} service.EnchantUsage
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/enchants [get]
func (h *MythicPlusBuildsAnalysisHandler) GetEnchantUsage(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))
	spec := normalizeWoWTerms(c.Query("spec"))

	if class == "" || spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class and spec parameters are required"})
		return
	}

	enchants, err := h.MythicPlusBuildsAnalysisService.GetEnchantUsage(c.Request.Context(), class, spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, enchants)
}

// GetGemUsage returns gem usage statistics for a specific class and spec
// @Summary Get gem usage
// @Description Returns gem usage statistics for a specific class and spec
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Param spec query string true "Specialization name"
// @Success 200 {array} service.GemUsage
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/gems [get]
func (h *MythicPlusBuildsAnalysisHandler) GetGemUsage(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))
	spec := normalizeWoWTerms(c.Query("spec"))

	if class == "" || spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class and spec parameters are required"})
		return
	}

	gems, err := h.MythicPlusBuildsAnalysisService.GetGemUsage(c.Request.Context(), class, spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gems)
}

// GetTopTalentBuilds returns the top talent builds for a specific class and spec
// @Summary Get top talent builds
// @Description Returns the top talent builds for a specific class and spec
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Param spec query string true "Specialization name"
// @Success 200 {array} service.TalentBuild
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/talents/top [get]
func (h *MythicPlusBuildsAnalysisHandler) GetTopTalentBuilds(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))
	spec := normalizeWoWTerms(c.Query("spec"))

	if class == "" || spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class and spec parameters are required"})
		return
	}

	builds, err := h.MythicPlusBuildsAnalysisService.GetTopTalentBuilds(c.Request.Context(), class, spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, builds)
}

// GetTalentBuildsByDungeon returns talent build statistics per dungeon for a specific class and spec
// @Summary Get talent builds by dungeon
// @Description Returns talent build statistics per dungeon for a specific class and spec
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Param spec query string true "Specialization name"
// @Success 200 {array} service.DungeonTalentBuild
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/talents/dungeons [get]
func (h *MythicPlusBuildsAnalysisHandler) GetTalentBuildsByDungeon(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))
	spec := normalizeWoWTerms(c.Query("spec"))

	if class == "" || spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class and spec parameters are required"})
		return
	}

	talents, err := h.MythicPlusBuildsAnalysisService.GetTalentBuildsByDungeon(c.Request.Context(), class, spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, talents)
}

// GetStatPriorities returns stat priority statistics for a specific class and spec
// @Summary Get stat priorities
// @Description Returns stat priority statistics for a specific class and spec
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Param spec query string true "Specialization name"
// @Success 200 {array} service.StatPriority
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/stats [get]
func (h *MythicPlusBuildsAnalysisHandler) GetStatPriorities(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))
	spec := normalizeWoWTerms(c.Query("spec"))

	if class == "" || spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class and spec parameters are required"})
		return
	}

	stats, err := h.MythicPlusBuildsAnalysisService.GetStatPriorities(c.Request.Context(), class, spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetOptimalBuild returns the optimal build for a specific class and spec
// @Summary Get optimal build
// @Description Returns the optimal build for a specific class and spec
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Param spec query string true "Specialization name"
// @Success 200 {object} service.OptimalBuild
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/optimal [get]
func (h *MythicPlusBuildsAnalysisHandler) GetOptimalBuild(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))
	spec := normalizeWoWTerms(c.Query("spec"))

	if class == "" || spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class and spec parameters are required"})
		return
	}

	optimalBuild, err := h.MythicPlusBuildsAnalysisService.GetOptimalBuild(c.Request.Context(), class, spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, optimalBuild)
}

// GetSpecComparison returns comparison statistics for all specs of a class
// @Summary Get spec comparison
// @Description Returns comparison statistics for all specs of a class
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Success 200 {array} service.SpecComparison
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/specs/comparison [get]
func (h *MythicPlusBuildsAnalysisHandler) GetSpecComparison(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))

	if class == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class parameter is required"})
		return
	}

	specs, err := h.MythicPlusBuildsAnalysisService.GetSpecComparison(c.Request.Context(), class)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, specs)
}

// GetClassSpecSummary returns summary statistics for a specific class and spec
// @Summary Get class spec summary
// @Description Returns summary statistics for a specific class and spec
// @Tags builds, analysis
// @Accept json
// @Produce json
// @Param class query string true "Class name"
// @Param spec query string true "Specialization name"
// @Success 200 {object} service.ClassSpecSummary
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /warcraftlogs/mythicplus/builds/analysis/summary [get]
func (h *MythicPlusBuildsAnalysisHandler) GetClassSpecSummary(c *gin.Context) {
	class := normalizeWoWTerms(c.Query("class"))
	spec := normalizeWoWTerms(c.Query("spec"))

	if class == "" || spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class and spec parameters are required"})
		return
	}

	summary, err := h.MythicPlusBuildsAnalysisService.GetClassSpecSummary(c.Request.Context(), class, spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func normalizeWoWTerms(term string) string {
	caser := cases.Title(language.English)
	return caser.String(strings.ToLower(term))
}

/*

Popular items by slot
/warcraftlogs/mythicplus/builds/analysis/items?class=priest&spec=discipline&encounter_id=12648

Popular items by slot across all encounters
/warcraftlogs/mythicplus/builds/analysis/items/global?class=priest&spec=discipline

Enchant usage
/warcraftlogs/mythicplus/builds/analysis/enchants?class=priest&spec=discipline

Gem usage
/warcraftlogs/mythicplus/builds/analysis/gems?class=priest&spec=discipline

Top talent builds
/warcraftlogs/mythicplus/builds/analysis/talents/top?class=priest&spec=discipline

Talent builds by dungeon
/warcraftlogs/mythicplus/builds/analysis/talents/dungeons?class=priest&spec=discipline

Stat priorities
/warcraftlogs/mythicplus/builds/analysis/stats?class=priest&spec=discipline

Optimal build
/warcraftlogs/mythicplus/builds/analysis/optimal?class=priest&spec=discipline


Class spec summary
/warcraftlogs/mythicplus/builds/analysis/summary?class=priest&spec=discipline

Spec comparison
/warcraftlogs/mythicplus/builds/analysis/specs/comparison?class=priest

*/
