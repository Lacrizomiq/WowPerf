// services/blizzard/protected/profile/handler.go

package protectedProfile

import (
	"fmt"
	"log"
	"net/http"
	protectedProfile "wowperf/internal/services/blizzard/protected/profile"
	"wowperf/internal/services/blizzard/types"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *protectedProfile.ProtectedProfileService
}

func NewHandler(service *protectedProfile.ProtectedProfileService) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes enregistre les routes pour le profil protégé
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	protected := r.Group("/wow/profile")
	{
		protected.GET("", h.GetAccountProfile)
	}
}

// GetAccountProfile récupère le profil WoW de l'utilisateur
func (h *Handler) GetAccountProfile(c *gin.Context) {
	// Get userID from context (set by auth middleware)
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	log.Printf("User ID: %d", userID)

	// Get the region (required)
	region := c.GetHeader("Region")
	if region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Region header is required"})
		return
	}

	// Verify that the battlenet token is here
	_, exists := c.Get("battlenet_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Battle.net token not found"})
		return
	}

	// Create service params
	params := types.ProfileServiceParams{
		Region:    region,
		Namespace: fmt.Sprintf("profile-%s", region),
		Locale:    c.DefaultQuery("locale", "en_US"),
	}

	// Get the profile data
	profile, err := h.service.GetAccountProfile(c.Request.Context(), userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get WoW profile: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}
