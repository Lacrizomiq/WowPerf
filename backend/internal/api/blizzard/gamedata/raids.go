package gamedata

import (
	"net/http"
	"strings"
	raids "wowperf/internal/models/raids"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RaidsByExpansionHandler struct {
	DB *gorm.DB
}

func NewRaidsByExpansionHandler(db *gorm.DB) *RaidsByExpansionHandler {
	return &RaidsByExpansionHandler{DB: db}
}

// GetRaidsByExpansion retrieves all raids for a given expansion
// For example all raids for Dragonflight are retrieved with /raids/df
func (h *RaidsByExpansionHandler) GetRaidsByExpansion(c *gin.Context) {
	expansion := strings.ToUpper(c.Param("expansion"))

	if expansion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	var raids []raids.Raid
	if err := h.DB.Where("UPPER(expansion) LIKE ?", "%"+expansion+"%").Find(&raids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve raids"})
		return
	}

	c.JSON(http.StatusOK, raids)
}
