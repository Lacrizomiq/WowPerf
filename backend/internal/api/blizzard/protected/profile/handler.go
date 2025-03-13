// services/blizzard/protected/profile/handler.go

package protectedProfile

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		protected.GET("/protected-character", h.GetProtectedCharacterProfile)

		protected.GET("/characters", h.ListAccountCharacters)
		protected.POST("/characters/sync", h.SyncSelectedCharacters)
		protected.GET("/user/characters", h.GetUserCharacters)
		protected.PUT("/characters/:id/favorite", h.SetFavoriteCharacter)
		protected.PUT("/characters/:id/display", h.ToggleCharacterDisplay)

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

	// Get the region (required)
	region := c.GetHeader("Region")
	if region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Region header is required"})
		return
	}

	// Verify that the battlenet token is here
	_, exists := c.Get("blizzard_token")
	if !exists {
		log.Printf("Battle.net token not found in context")
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

// GetProtectedCharacterProfile retrieves the WoW profile of a protected character
func (h *Handler) GetProtectedCharacterProfile(c *gin.Context) {
	// Get userID from context (set by auth middleware)
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get the params from query
	realmId := c.Query("realmId")
	characterId := c.Query("characterId")
	if realmId == "" || characterId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Realm ID and character ID are required"})
		return
	}

	// Get the region (required)
	region := c.GetHeader("Region")
	if region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Region header is required"})
		return
	}

	// create service params
	params := types.ProfileServiceParams{
		Region:    region,
		Namespace: fmt.Sprintf("profile-%s", region),
		Locale:    c.DefaultQuery("locale", "en_US"),
	}

	// get the protected character profile
	profile, err := h.service.GetProtectedCharacterProfile(c.Request.Context(), userID, realmId, characterId, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get protected character profile: %v", err)})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) ListAccountCharacters(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	region := c.GetHeader("Region")
	if region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Region header is required"})
		return
	}

	characters, err := h.service.ListAccountCharacters(c.Request.Context(), userID, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to list account characters: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, characters)
}

func (h *Handler) SyncSelectedCharacters(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var selections []protectedProfile.CharacterSelection
	if err := c.ShouldBindJSON(&selections); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request format: %v", err)})
		return
	}

	results, err := h.service.SyncSelectedCharacters(c.Request.Context(), userID, selections)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to sync characters: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) GetUserCharacters(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	characters, err := h.service.GetUserCharacters(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get user characters: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, characters)
}

func (h *Handler) SetFavoriteCharacter(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	characterIDStr := c.Param("id")
	characterID, err := strconv.ParseUint(characterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	if err := h.service.SetFavoriteCharacter(c.Request.Context(), userID, uint(characterID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to set favorite character: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Character set as favorite successfully"})
}

// ToggleCharacterDisplay active ou désactive l'affichage d'un personnage
func (h *Handler) ToggleCharacterDisplay(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	characterIDStr := c.Param("id")
	characterID, err := strconv.ParseUint(characterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	// Récupérer la valeur d'affichage depuis le corps de la requête
	var requestBody struct {
		Display bool `json:"display"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.service.ToggleCharacterDisplay(c.Request.Context(), userID, uint(characterID), requestBody.Display); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to toggle character display: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Character display updated successfully"})
}
