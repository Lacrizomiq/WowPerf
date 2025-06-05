package character

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"
	"wowperf/internal/models"
	protectedProfile "wowperf/internal/services/blizzard/protected/profile"
	"wowperf/internal/services/character/enrichers"
)

// CharacterOrchestrator coordonne la synchronisation et l'enrichissement des personnages
type CharacterOrchestrator struct {
	characterService        CharacterServiceInterface
	protectedProfileService ProtectedProfileServiceInterface
	enrichersList           []enrichers.CharacterEnricher
	rateLimiter             *RateLimiter
}

// ProtectedProfileServiceInterface interface pour découpler le service protected profile
type ProtectedProfileServiceInterface interface {
	SyncAllAccountCharacters(ctx context.Context, userID uint, region string) (int, error)
	RefreshUserCharacters(ctx context.Context, userID uint, region string) (int, int, error)
}

// NewCharacterOrchestrator crée un nouvel orchestrateur
func NewCharacterOrchestrator(
	characterService CharacterServiceInterface,
	protectedProfileService *protectedProfile.ProtectedProfileService,
) *CharacterOrchestrator {
	orchestrator := &CharacterOrchestrator{
		characterService:        characterService,
		protectedProfileService: protectedProfileService,
		enrichersList:           []enrichers.CharacterEnricher{},
		rateLimiter:             NewRateLimiter(),
	}

	// Démarrer le nettoyage périodique du rate limiter
	go orchestrator.startCleanupRoutine()

	return orchestrator
}

// RegisterEnricher ajoute un enrichisseur à la liste
func (o *CharacterOrchestrator) RegisterEnricher(enricher enrichers.CharacterEnricher) {
	o.enrichersList = append(o.enrichersList, enricher)

	// Trier par priorité après chaque ajout
	sort.Slice(o.enrichersList, func(i, j int) bool {
		return o.enrichersList[i].GetPriority() < o.enrichersList[j].GetPriority()
	})

	log.Printf("Registered enricher: %s (priority: %d)", enricher.GetName(), enricher.GetPriority())
}

// GetRegisteredEnrichers retourne la liste des enrichisseurs enregistrés
func (o *CharacterOrchestrator) GetRegisteredEnrichers() []string {
	names := make([]string, len(o.enrichersList))
	for i, enricher := range o.enrichersList {
		names[i] = enricher.GetName()
	}
	return names
}

// SyncAndEnrichUserCharacters synchronise et enrichit tous les personnages d'un utilisateur
func (o *CharacterOrchestrator) SyncAndEnrichUserCharacters(ctx context.Context, userID uint, region string) (*SyncResult, error) {
	startTime := time.Now()
	result := &SyncResult{
		SyncedCount:   0,
		EnrichedCount: 0,
		UpdatedCount:  0,
		Errors:        []string{},
	}

	log.Printf("Starting character orchestration for user %d in region %s", userID, region)

	// Vérifier le rate limiting pour sync
	canSync, message := o.rateLimiter.CanSyncUser(userID)
	if !canSync {
		return nil, fmt.Errorf("rate limit exceeded: %s", message)
	}

	// Enregistrer la tentative de sync
	o.rateLimiter.RecordSync(userID)

	log.Printf("Registered enrichers: %v", o.GetRegisteredEnrichers())

	// 1. Synchroniser les personnages depuis l'API protected profile
	log.Printf("Step 1: Syncing characters from Blizzard API")
	syncCount, err := o.protectedProfileService.SyncAllAccountCharacters(ctx, userID, region)
	if err != nil {
		return nil, fmt.Errorf("failed to sync characters: %w", err)
	}
	result.SyncedCount = syncCount
	log.Printf("Synced %d characters", syncCount)

	// 2. Récupérer tous les personnages de l'utilisateur
	log.Printf("Step 2: Retrieving user characters from database")
	characters, err := o.characterService.GetCharactersByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user characters: %w", err)
	}
	log.Printf("Found %d characters to enrich", len(characters))

	// 3. Enrichir chaque personnage avec tous les enrichisseurs
	log.Printf("Step 3: Enriching characters")
	enrichmentResults := o.enrichAllCharacters(ctx, characters, result)

	// 4. Log des résultats détaillés
	o.logEnrichmentResults(enrichmentResults)

	duration := time.Since(startTime)
	log.Printf("Character orchestration completed in %v", duration)
	log.Printf("Final result: %d synced, %d enriched, %d errors",
		result.SyncedCount, result.EnrichedCount, len(result.Errors))

	return result, nil
}

// RefreshAndEnrichUserCharacters utilise la méthode refresh pour les personnages existants
func (o *CharacterOrchestrator) RefreshAndEnrichUserCharacters(ctx context.Context, userID uint, region string) (*SyncResult, error) {
	startTime := time.Now()
	result := &SyncResult{
		SyncedCount:   0,
		EnrichedCount: 0,
		UpdatedCount:  0,
		Errors:        []string{},
	}

	log.Printf("Starting character refresh and enrichment for user %d in region %s", userID, region)

	// Vérifier le rate limiting pour sync
	canSync, message := o.rateLimiter.CanSyncUser(userID)
	if !canSync {
		return nil, fmt.Errorf("rate limit exceeded: %s", message)
	}

	// Enregistrer la tentative de sync
	o.rateLimiter.RecordSync(userID)

	// 1. Refresh des personnages (sync nouveaux + update existants)
	log.Printf("Step 1: Refreshing characters from Blizzard API")
	newCount, updatedCount, err := o.protectedProfileService.RefreshUserCharacters(ctx, userID, region)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh characters: %w", err)
	}
	result.SyncedCount = newCount
	result.UpdatedCount = updatedCount
	log.Printf("Refreshed characters: %d new, %d updated", newCount, updatedCount)

	// 2. Récupérer et enrichir comme dans SyncAndEnrichUserCharacters
	characters, err := o.characterService.GetCharactersByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user characters: %w", err)
	}

	log.Printf("Step 2: Enriching %d characters", len(characters))
	enrichmentResults := o.enrichAllCharacters(ctx, characters, result)
	o.logEnrichmentResults(enrichmentResults)

	duration := time.Since(startTime)
	log.Printf("Character refresh and enrichment completed in %v", duration)

	return result, nil
}

// EnrichSingleCharacter enrichit un seul personnage (pour les enrichissements à la demande)
func (o *CharacterOrchestrator) EnrichSingleCharacter(ctx context.Context, userID uint, characterID uint) error {
	// Vérifier le rate limiting pour enrichissement
	canEnrich, message := o.rateLimiter.CanEnrichUser(userID)
	if !canEnrich {
		return fmt.Errorf("rate limit exceeded: %s", message)
	}

	// Enregistrer la tentative d'enrichissement
	o.rateLimiter.RecordEnrich(userID)

	character, err := o.characterService.GetCharacterByID(characterID)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}

	log.Printf("Enriching single character: %s", character.Name)
	results := o.enrichSingleCharacterInternal(ctx, character)

	// Vérifier qu'au moins un enrichissement a réussi
	hasSuccess := false
	for _, result := range results {
		if result.Success {
			hasSuccess = true
			break
		}
	}

	if hasSuccess {
		if err := o.characterService.CreateOrUpdateCharacter(character); err != nil {
			return fmt.Errorf("failed to save enriched character: %w", err)
		}
	}

	return nil
}

// enrichAllCharacters applique tous les enrichisseurs sur tous les personnages
func (o *CharacterOrchestrator) enrichAllCharacters(ctx context.Context, characters []models.UserCharacter, result *SyncResult) []enrichers.EnrichmentResult {
	var allResults []enrichers.EnrichmentResult

	for i := range characters {
		character := &characters[i]
		log.Printf("Enriching character: %s (%s-%s)", character.Name, character.Realm, character.Region)

		characterResults := o.enrichSingleCharacterInternal(ctx, character)
		allResults = append(allResults, characterResults...)

		// Compter les enrichissements réussis
		successCount := 0
		for _, res := range characterResults {
			if res.Success {
				successCount++
			}
		}

		if successCount > 0 {
			result.EnrichedCount++
			// Sauvegarder les changements après enrichissement
			if err := o.characterService.CreateOrUpdateCharacter(character); err != nil {
				result.Errors = append(result.Errors,
					fmt.Sprintf("Failed to save character %s: %v", character.Name, err))
			}
		}
	}

	return allResults
}

// enrichSingleCharacterInternal applique tous les enrichisseurs sur un seul personnage
func (o *CharacterOrchestrator) enrichSingleCharacterInternal(ctx context.Context, character *models.UserCharacter) []enrichers.EnrichmentResult {
	var results []enrichers.EnrichmentResult

	for _, enricher := range o.enrichersList {
		startTime := time.Now()
		result := enrichers.EnrichmentResult{
			EnricherName: enricher.GetName(),
			CharacterID:  character.ID,
			Success:      false,
		}

		// Vérifier si l'enrichisseur peut traiter ce personnage
		if !enricher.CanEnrich(character) {
			result.Error = "enricher cannot process this character"
			result.Duration = time.Since(startTime).Milliseconds()
			results = append(results, result)
			continue
		}

		// Appliquer l'enrichissement
		if err := enricher.EnrichCharacter(ctx, character); err != nil {
			result.Error = err.Error()
			log.Printf("Enricher %s failed for character %s: %v",
				enricher.GetName(), character.Name, err)
		} else {
			result.Success = true
			log.Printf("Enricher %s succeeded for character %s",
				enricher.GetName(), character.Name)
		}

		result.Duration = time.Since(startTime).Milliseconds()
		results = append(results, result)
	}

	return results
}

// logEnrichmentResults log les résultats détaillés des enrichissements
func (o *CharacterOrchestrator) logEnrichmentResults(results []enrichers.EnrichmentResult) {
	if len(results) == 0 {
		return
	}

	log.Printf("=== Enrichment Results Summary ===")

	// Compter par enrichisseur
	enricherStats := make(map[string]struct {
		total         int
		success       int
		failed        int
		totalDuration int64
	})

	for _, result := range results {
		stats := enricherStats[result.EnricherName]
		stats.total++
		stats.totalDuration += result.Duration
		if result.Success {
			stats.success++
		} else {
			stats.failed++
		}
		enricherStats[result.EnricherName] = stats
	}

	for enricherName, stats := range enricherStats {
		avgDuration := float64(stats.totalDuration) / float64(stats.total)
		log.Printf("Enricher %s: %d/%d successful (%.1fms avg)",
			enricherName, stats.success, stats.total, avgDuration)
	}
}

// GetUserCharacters récupère les personnages d'un utilisateur
func (o *CharacterOrchestrator) GetUserCharacters(ctx context.Context, userID uint) ([]models.UserCharacter, error) {
	return o.characterService.GetCharactersByUserID(userID)
}

// startCleanupRoutine démarre une routine de nettoyage périodique
func (o *CharacterOrchestrator) startCleanupRoutine() {
	ticker := time.NewTicker(1 * time.Hour) // Nettoyer toutes les heures
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			o.rateLimiter.CleanupOldEntries()
			log.Printf("Rate limiter cleanup completed")
		}
	}
}
