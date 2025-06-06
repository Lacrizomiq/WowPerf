package enrichers

import (
	"context"
	"wowperf/internal/models"
)

// CharacterEnricher interface commune pour tous les enrichisseurs de personnages
type CharacterEnricher interface {
	// EnrichCharacter enrichit un personnage avec des données spécifiques
	EnrichCharacter(ctx context.Context, character *models.UserCharacter) error

	// GetName retourne le nom de l'enrichisseur (pour les logs et debug)
	GetName() string

	// GetPriority retourne la priorité d'exécution (1 = premier, 10 = dernier)
	GetPriority() int

	// CanEnrich vérifie si cet enrichisseur peut traiter ce personnage
	CanEnrich(character *models.UserCharacter) bool
}

// EnrichmentResult représente le résultat d'un enrichissement individuel
type EnrichmentResult struct {
	EnricherName string `json:"enricher_name"`
	CharacterID  uint   `json:"character_id"`
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
	Duration     int64  `json:"duration_ms"` // Durée en millisecondes
}
