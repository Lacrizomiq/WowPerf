package mythicPlusRunsAnalysis

import (
	"net/http"
	"strconv"

	analyticsService "wowperf/internal/services/raiderio/mythicplus/analytics"

	"github.com/gin-gonic/gin"
)

// MythicPlusRunsAnalysisHandler gère les endpoints d'analyse des runs Mythic+
type MythicPlusRunsAnalysisHandler struct {
	mythicPlusRunsAnalysisService *analyticsService.MythicPlusRunsAnalysisService
}

// NewMythicPlusRunsAnalysisHandler crée un nouveau handler d'analyse
func NewMythicPlusRunsAnalysisHandler(service *analyticsService.MythicPlusRunsAnalysisService) *MythicPlusRunsAnalysisHandler {
	return &MythicPlusRunsAnalysisHandler{
		mythicPlusRunsAnalysisService: service,
	}
}

// ========================================
// ENDPOINTS - SPÉCIALISATIONS GLOBALES
// ========================================

// GetTankSpecializations returns the most popular Tank specializations globally
// @Summary Get popular Tank specializations
// @Description Returns the most popular Tank specializations across all dungeons with usage statistics
// @Tags Mythic+ Analytics - Global
// @Accept json
// @Produce json
// @Param top_n query int false "Number of top specs to return (0 = all specs, default: 0)"
// @Success 200 {array} analyticsService.SpecializationStats
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/specs/tank [get]
func (h *MythicPlusRunsAnalysisHandler) GetTankSpecializations(c *gin.Context) {
	topN := 0 // default: return all
	if topNStr := c.Query("top_n"); topNStr != "" {
		if parsedTopN, err := strconv.Atoi(topNStr); err == nil && parsedTopN >= 0 {
			topN = parsedTopN
		}
	}

	results, err := h.mythicPlusRunsAnalysisService.GetTankSpecializations(topN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetHealerSpecializations returns the most popular Healer specializations globally
// @Summary Get popular Healer specializations
// @Description Returns the most popular Healer specializations across all dungeons with usage statistics
// @Tags Mythic+ Analytics - Global
// @Accept json
// @Produce json
// @Param top_n query int false "Number of top specs to return (0 = all specs, default: 0)"
// @Success 200 {array} analyticsService.SpecializationStats
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/specs/healer [get]
func (h *MythicPlusRunsAnalysisHandler) GetHealerSpecializations(c *gin.Context) {
	topN := 0 // default: return all
	if topNStr := c.Query("top_n"); topNStr != "" {
		if parsedTopN, err := strconv.Atoi(topNStr); err == nil && parsedTopN >= 0 {
			topN = parsedTopN
		}
	}

	results, err := h.mythicPlusRunsAnalysisService.GetHealerSpecializations(topN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetDPSSpecializations returns the most popular DPS specializations globally
// @Summary Get popular DPS specializations
// @Description Returns the most popular DPS specializations across all dungeons with usage statistics
// @Tags Mythic+ Analytics - Global
// @Accept json
// @Produce json
// @Param top_n query int false "Number of top specs to return (0 = all specs, default: 0)"
// @Success 200 {array} analyticsService.SpecializationStats
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/specs/dps [get]
func (h *MythicPlusRunsAnalysisHandler) GetDPSSpecializations(c *gin.Context) {
	topN := 0 // default: return all
	if topNStr := c.Query("top_n"); topNStr != "" {
		if parsedTopN, err := strconv.Atoi(topNStr); err == nil && parsedTopN >= 0 {
			topN = parsedTopN
		}
	}
	results, err := h.mythicPlusRunsAnalysisService.GetDPSSpecializations(topN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetSpecializationsByRole returns the most popular specializations for a specific role
// @Summary Get popular specializations by role
// @Description Returns the most popular specializations for a specific role (tank, healer, dps)
// @Tags Mythic+ Analytics - Global
// @Accept json
// @Produce json
// @Param role path string true "Role (tank, healer, dps)"
// @Param top_n query int false "Number of top specs to return (0 = all specs, default: 0)"
// @Success 200 {array} analyticsService.SpecializationStats
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/specs/{role} [get]
func (h *MythicPlusRunsAnalysisHandler) GetSpecializationsByRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role parameter is required"})
		return
	}

	topN := 0 // default: return all
	if topNStr := c.Query("top_n"); topNStr != "" {
		if parsedTopN, err := strconv.Atoi(topNStr); err == nil && parsedTopN >= 0 {
			topN = parsedTopN
		}
	}
	results, err := h.mythicPlusRunsAnalysisService.GetSpecializationsByRole(role, topN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ========================================
// ENDPOINTS - COMPOSITIONS GLOBALES
// ========================================

// GetTopCompositions returns the most popular team compositions globally
// @Summary Get popular team compositions
// @Description Returns the most popular team compositions with usage statistics and average scores
// @Tags Mythic+ Analytics - Global
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of results to return (default: 20)"
// @Param min_usage query int false "Minimum usage count to filter results (default: 5)"
// @Success 200 {array} analyticsService.CompositionStats
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/compositions [get]
func (h *MythicPlusRunsAnalysisHandler) GetTopCompositions(c *gin.Context) {
	// Parse query parameters
	limit := 5 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	minUsage := 5 // default
	if minUsageStr := c.Query("min_usage"); minUsageStr != "" {
		if parsedMinUsage, err := strconv.Atoi(minUsageStr); err == nil && parsedMinUsage > 0 {
			minUsage = parsedMinUsage
		}
	}

	results, err := h.mythicPlusRunsAnalysisService.GetTopCompositions(limit, minUsage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ========================================
// ENDPOINTS - ANALYSES PAR DONJON
// ========================================

// GetTankSpecsByDungeon returns Tank specializations ranked by dungeon
// @Summary Get Tank specs by dungeon
// @Description Returns Tank specializations usage statistics for each dungeon
// @Tags Mythic+ Analytics - By Dungeon
// @Accept json
// @Produce json
// @Param top_n query int false "Number of top specs per dungeon (0 = all specs, default: 0)"
// @Success 200 {array} analyticsService.DungeonSpecStats
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/dungeons/specs/tank [get]
func (h *MythicPlusRunsAnalysisHandler) GetTankSpecsByDungeon(c *gin.Context) {
	topN := 0 // default: return all
	if topNStr := c.Query("top_n"); topNStr != "" {
		if parsedTopN, err := strconv.Atoi(topNStr); err == nil && parsedTopN >= 0 {
			topN = parsedTopN
		}
	}

	results, err := h.mythicPlusRunsAnalysisService.GetTankSpecsByDungeon(topN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetHealerSpecsByDungeon returns Healer specializations ranked by dungeon
// @Summary Get Healer specs by dungeon
// @Description Returns Healer specializations usage statistics for each dungeon
// @Tags Mythic+ Analytics - By Dungeon
// @Accept json
// @Produce json
// @Param top_n query int false "Number of top specs per dungeon (0 = all specs, default: 0)"
// @Success 200 {array} analyticsService.DungeonSpecStats
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/dungeons/specs/healer [get]
func (h *MythicPlusRunsAnalysisHandler) GetHealerSpecsByDungeon(c *gin.Context) {
	topN := 0 // default: return all
	if topNStr := c.Query("top_n"); topNStr != "" {
		if parsedTopN, err := strconv.Atoi(topNStr); err == nil && parsedTopN >= 0 {
			topN = parsedTopN
		}
	}

	results, err := h.mythicPlusRunsAnalysisService.GetHealerSpecsByDungeon(topN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetDPSSpecsByDungeon returns DPS specializations ranked by dungeon
// @Summary Get DPS specs by dungeon
// @Description Returns DPS specializations usage statistics for each dungeon
// @Tags Mythic+ Analytics - By Dungeon
// @Accept json
// @Produce json
// @Param top_n query int false "Number of top specs per dungeon (0 = all specs, default: 0)"
// @Success 200 {array} analyticsService.DungeonSpecStats
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/dungeons/specs/dps [get]
func (h *MythicPlusRunsAnalysisHandler) GetDPSSpecsByDungeon(c *gin.Context) {
	topN := 0 // default: return all
	if topNStr := c.Query("top_n"); topNStr != "" {
		if parsedTopN, err := strconv.Atoi(topNStr); err == nil && parsedTopN >= 0 {
			topN = parsedTopN
		}
	}

	results, err := h.mythicPlusRunsAnalysisService.GetDPSSpecsByDungeon(topN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetSpecsByDungeonAndRole returns specializations for a specific role and dungeon
// @Summary Get specs by dungeon and role
// @Description Returns specializations usage statistics for a specific role in a specific dungeon
// @Tags Mythic+ Analytics - By Dungeon
// @Accept json
// @Produce json
// @Param dungeon_slug path string true "Dungeon slug"
// @Param role path string true "Role (tank, healer, dps)"
// @Param top_n query int false "Number of top specs (0 = all specs, default: 0)"
// @Success 200 {array} analyticsService.DungeonSpecStats
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/dungeons/{dungeon_slug}/specs/{role} [get]
func (h *MythicPlusRunsAnalysisHandler) GetSpecsByDungeonAndRole(c *gin.Context) {
	dungeonSlug := c.Param("dungeon_slug")
	role := c.Param("role")

	if dungeonSlug == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dungeon_slug and role parameters are required"})
		return
	}

	topN := 0 // default: return all
	if topNStr := c.Query("top_n"); topNStr != "" {
		if parsedTopN, err := strconv.Atoi(topNStr); err == nil && parsedTopN >= 0 {
			topN = parsedTopN
		}
	}

	results, err := h.mythicPlusRunsAnalysisService.GetSpecsByDungeonAndRole(dungeonSlug, role, topN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetCompositionsByDungeon returns team compositions ranked by dungeon
// @Summary Get compositions by dungeon
// @Description Returns team compositions usage statistics for each dungeon
// @Tags Mythic+ Analytics - By Dungeon
// @Accept json
// @Produce json
// @Param top_n query int false "Number of top compositions per dungeon (0 = all compositions, default: 0)"
// @Param min_usage query int false "Minimum usage count to filter results (default: 3)"
// @Success 200 {array} analyticsService.DungeonCompositionStats
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/dungeons/compositions [get]
func (h *MythicPlusRunsAnalysisHandler) GetCompositionsByDungeon(c *gin.Context) {
	topN := 0 // default: return all
	if topNStr := c.Query("top_n"); topNStr != "" {
		if parsedTopN, err := strconv.Atoi(topNStr); err == nil && parsedTopN >= 0 {
			topN = parsedTopN
		}
	}

	minUsage := 3 // default
	if minUsageStr := c.Query("min_usage"); minUsageStr != "" {
		if parsedMinUsage, err := strconv.Atoi(minUsageStr); err == nil && parsedMinUsage > 0 {
			minUsage = parsedMinUsage
		}
	}

	results, err := h.mythicPlusRunsAnalysisService.GetTopCompositionsByDungeon(topN, minUsage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ========================================
// ENDPOINTS - ANALYSES AVANCÉES
// ========================================

// GetSpecsByKeyLevel returns specializations usage by key level brackets
// @Summary Get specs by key level
// @Description Returns specializations usage statistics grouped by key level brackets (Very High 20+, High 18-19, Mid 16-17)
// @Tags Mythic+ Analytics - Advanced
// @Accept json
// @Produce json
// @Param min_usage query int false "Minimum usage count to filter results (default: 5)"
// @Success 200 {array} analyticsService.KeyLevelStats
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/key-levels [get]
func (h *MythicPlusRunsAnalysisHandler) GetSpecsByKeyLevel(c *gin.Context) {
	minUsage := 5 // default
	if minUsageStr := c.Query("min_usage"); minUsageStr != "" {
		if parsedMinUsage, err := strconv.Atoi(minUsageStr); err == nil && parsedMinUsage > 0 {
			minUsage = parsedMinUsage
		}
	}

	results, err := h.mythicPlusRunsAnalysisService.GetSpecsByKeyLevel(minUsage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetSpecsByRegion returns specializations usage by region
// @Summary Get specs by region
// @Description Returns specializations usage statistics for each region (US, EU, KR, TW)
// @Tags Mythic+ Analytics - Advanced
// @Accept json
// @Produce json
// @Success 200 {array} analyticsService.RegionStats
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/regions [get]
func (h *MythicPlusRunsAnalysisHandler) GetSpecsByRegion(c *gin.Context) {
	results, err := h.mythicPlusRunsAnalysisService.GetSpecsByRegion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ========================================
// ENDPOINTS - STATISTIQUES UTILITAIRES
// ========================================

// GetOverallStats returns general statistics about the dataset
// @Summary Get overall statistics
// @Description Returns general statistics about the Mythic+ runs dataset (total runs, score averages, etc.)
// @Tags Mythic+ Analytics - Utilities
// @Accept json
// @Produce json
// @Success 200 {object} analyticsService.OverallStats
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/stats/overall [get]
func (h *MythicPlusRunsAnalysisHandler) GetOverallStats(c *gin.Context) {
	results, err := h.mythicPlusRunsAnalysisService.GetOverallStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetKeyLevelDistribution returns the distribution of runs by key level
// @Summary Get key level distribution
// @Description Returns the distribution of runs by mythic level with statistics
// @Tags Mythic+ Analytics - Utilities
// @Accept json
// @Produce json
// @Success 200 {array} analyticsService.KeyLevelDistribution
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/stats/key-levels [get]
func (h *MythicPlusRunsAnalysisHandler) GetKeyLevelDistribution(c *gin.Context) {
	results, err := h.mythicPlusRunsAnalysisService.GetKeyLevelDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetDungeonDistribution returns the distribution of runs by dungeon
// @Summary Get dungeon distribution
// @Description Returns the distribution of runs by dungeon with statistics
// @Tags Mythic+ Analytics - Utilities
// @Accept json
// @Produce json
// @Success 200 {array} analyticsService.DungeonDistribution
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/stats/dungeons [get]
func (h *MythicPlusRunsAnalysisHandler) GetDungeonDistribution(c *gin.Context) {
	results, err := h.mythicPlusRunsAnalysisService.GetDungeonDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetRegionDistribution returns the distribution of runs by region
// @Summary Get region distribution
// @Description Returns the distribution of runs by region with statistics
// @Tags Mythic+ Analytics - Utilities
// @Accept json
// @Produce json
// @Success 200 {array} analyticsService.RegionDistribution
// @Failure 500 {object} string "Internal server error"
// @Router /raiderio/mythicplus/analytics/stats/regions [get]
func (h *MythicPlusRunsAnalysisHandler) GetRegionDistribution(c *gin.Context) {
	results, err := h.mythicPlusRunsAnalysisService.GetRegionDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

/*
=======================================================
MYTHIC+ ANALYTICS API ROUTES
=======================================================

🌍 ANALYSES GLOBALES
GET /raiderio/mythicplus/analytics/specs/tank                    // Top spécs Tank globales
GET /raiderio/mythicplus/analytics/specs/healer                  // Top spécs Healer globales
GET /raiderio/mythicplus/analytics/specs/dps                     // Top spécs DPS globales
GET /raiderio/mythicplus/analytics/specs/{role}                  // Spécs par rôle (tank/healer/dps)
GET /raiderio/mythicplus/analytics/compositions                  // Top compositions globales

🏰 ANALYSES PAR DONJON
GET /raiderio/mythicplus/analytics/dungeons/specs/tank           // Spécs Tank par donjon
GET /raiderio/mythicplus/analytics/dungeons/specs/healer         // Spécs Healer par donjon
GET /raiderio/mythicplus/analytics/dungeons/specs/dps            // Spécs DPS par donjon
GET /raiderio/mythicplus/analytics/dungeons/compositions         // Compositions par donjon
GET /raiderio/mythicplus/analytics/dungeons/{slug}/specs/{role}  // Spécs pour donjon+rôle spécifique

🔥 ANALYSES AVANCÉES
GET /raiderio/mythicplus/analytics/key-levels                    // Méta par niveau de clé (20+, 18-19, 16-17)
GET /raiderio/mythicplus/analytics/regions                       // Méta par région (US/EU/KR/TW)

📈 STATISTIQUES UTILITAIRES
GET /raiderio/mythicplus/analytics/stats/overall                 // Stats générales du dataset (total runs, score averages, etc.)
GET /raiderio/mythicplus/analytics/stats/key-levels              // Distribution par niveau de clé
GET /raiderio/mythicplus/analytics/stats/dungeons                // Distribution par donjon
GET /raiderio/mythicplus/analytics/stats/regions                 // Distribution par région

PARAMÈTRES PRINCIPAUX:
- top_n=0        : Retourne TOUS les résultats (défaut pour analyses complètes)
- top_n>0        : Limite aux top N résultats (Par exemple les 3 meilleurs specs)
- min_usage=X    : Filtre les données avec <X utilisations (évite le bruit)
- limit=X        : Nombre max de résultats (compositions uniquement)

CACHE: 24h sur tous les endpoints
FORMAT: JSON avec structures SpecializationStats, CompositionStats, etc.
=======================================================
*/
