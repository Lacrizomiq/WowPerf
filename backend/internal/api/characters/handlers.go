package characters

import (
	"net/http"
	"strconv"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/character"
	"wowperf/internal/services/character/enrichers"

	"github.com/gin-gonic/gin"
)

type CharactersHandler struct {
	orchestrator *character.CharacterOrchestrator
}

// NewCharactersHandler crée un nouveau handler avec orchestrateur
func NewCharactersHandler(
	characterService character.CharacterServiceInterface,
	blizzardService *blizzard.Service,
) *CharactersHandler {
	// Créer l'orchestrateur
	orchestrator := character.NewCharacterOrchestrator(
		characterService,
		blizzardService.ProtectedProfile,
	)

	// Enregistrer les enrichisseurs
	summaryEnricher := enrichers.NewSummaryEnricher(blizzardService.Profile)
	orchestrator.RegisterEnricher(summaryEnricher)

	return &CharactersHandler{
		orchestrator: orchestrator,
	}
}

// RegisterRoutes enregistre les routes characters
func (h *CharactersHandler) RegisterRoutes(r *gin.RouterGroup) {
	characters := r.Group("/characters")
	{
		// Synchronisation et enrichissement
		characters.POST("/sync-and-enrich", h.SyncAndEnrichCharacters)
		characters.POST("/refresh-and-enrich", h.RefreshAndEnrichCharacters)

		// Récupération
		characters.GET("", h.GetUserCharacters)

		// Enrichissement individuel
		characters.POST("/:id/enrich", h.EnrichSingleCharacter)

		// Debug/info
		characters.GET("/status", h.GetStatus)
	}
}

// SyncAndEnrichCharacters - Sync + enrichissement automatique
func (h *CharactersHandler) SyncAndEnrichCharacters(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	region := c.GetHeader("Region")
	if region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Region header required"})
		return
	}

	// Vérifier que le token Battle.net est présent
	_, exists := c.Get("blizzard_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Battle.net token not found"})
		return
	}

	// Lancer la synchronisation et enrichissement complets
	result, err := h.orchestrator.SyncAndEnrichUserCharacters(c.Request.Context(), userID, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Characters synchronized and enriched successfully",
		"result":  result,
	})
}

// RefreshAndEnrichCharacters - Refresh + enrichissement automatique
func (h *CharactersHandler) RefreshAndEnrichCharacters(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	region := c.GetHeader("Region")
	if region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Region header required"})
		return
	}

	_, exists := c.Get("blizzard_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Battle.net token not found"})
		return
	}

	// Lancer le refresh et enrichissement
	result, err := h.orchestrator.RefreshAndEnrichUserCharacters(c.Request.Context(), userID, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Characters refreshed and enriched successfully",
		"result":  result,
	})
}

// GetUserCharacters - Récupère tous les personnages enrichis
func (h *CharactersHandler) GetUserCharacters(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	characters, err := h.orchestrator.GetUserCharacters(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": characters,
		"count":      len(characters),
	})
}

// EnrichSingleCharacter - Enrichit un seul personnage à la demande
func (h *CharactersHandler) EnrichSingleCharacter(c *gin.Context) {
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

	// 🔒 Vérifier que le personnage appartient à l'utilisateur
	characters, err := h.orchestrator.GetUserCharacters(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify character ownership",
		})
		return
	}

	// Chercher le personnage dans la liste des personnages de l'utilisateur
	characterFound := false
	for _, char := range characters {
		if char.ID == uint(characterID) {
			characterFound = true
			break
		}
	}

	if !characterFound {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Character does not belong to this user",
		})
		return
	}

	// Enrichir le personnage
	err = h.orchestrator.EnrichSingleCharacter(c.Request.Context(), userID, uint(characterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Character enriched successfully",
	})
}

// GetStatus - Info sur l'état du système d'enrichissement
func (h *CharactersHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "active",
		"description": "Character enrichment system",
		"version":     "1.0",
	})
}

/*

Synchronisation et enrichissement automatique :
POST /api/characters/sync-and-enrich           # characters

Rafraîchissement et enrichissement des personnages :
POST /api/characters/refresh-and-enrich        # characters

Récupération des personnages enrichis :
GET  /api/characters                           # characters

Enrichissement individuel d'un personnage :
POST /api/characters/:id/enrich                # characters

*/
