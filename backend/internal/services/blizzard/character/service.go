// services/blizzard/character/service.go

package character

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard/types"

	"gorm.io/gorm"
)

// CharacterService handles operations related to WoW characters
type CharacterService struct {
	db              *gorm.DB
	protectedClient types.ProtectedClientInterface
	repo            *CharacterRepository
}

// NewCharacterService creates a new instance of the character service
func NewCharacterService(db *gorm.DB, protectedClient types.ProtectedClientInterface) *CharacterService {
	repo := NewCharacterRepository(db)
	return &CharacterService{
		db:              db,
		protectedClient: protectedClient,
		repo:            repo,
	}
}

// ListAccountCharacters retrieves all characters linked to the user's Battle.net account
func (s *CharacterService) ListAccountCharacters(ctx context.Context, userID uint) ([]CharacterBasicInfo, error) {
	// Battle.net API requires profile namespace format: profile-{region}
	// The region is determined by the OAuth token associated with the user
	namespace := "profile-eu" // This will be used properly by the protectedClient
	locale := "en_GB"         // Can be configurable if needed

	// Call the API to retrieve account characters using the existing protected client
	data, err := s.protectedClient.MakeProtectedRequest(ctx, userID, "/profile/user/wow", namespace, locale)
	if err != nil {
		return nil, fmt.Errorf("error retrieving characters: %w", err)
	}

	// Parse the response
	var response AccountCharactersResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	// Convert to a simpler format and filter level 80+ characters
	characters := make([]CharacterBasicInfo, 0)
	for _, account := range response.WowAccounts {
		for _, char := range account.Characters {
			// Only include max level characters (customize this filter as needed)
			if char.Level >= 80 {
				characters = append(characters, CharacterBasicInfo{
					CharacterID: char.ID,
					Name:        char.Name,
					Realm:       char.Realm.Slug,
					Region:      "eu", // This comes from the token, but we hardcode for simplicity
					Class:       char.PlayableClass.Name,
					Race:        char.PlayableRace.Name,
					Level:       char.Level,
					Faction:     char.Faction.Type,
				})
			}
		}
	}

	return characters, nil
}

// SyncSelectedCharacters synchronizes the selected characters with the database
func (s *CharacterService) SyncSelectedCharacters(ctx context.Context, userID uint, selections []CharacterSelection) ([]CharacterSyncResult, error) {
	results := make([]CharacterSyncResult, 0, len(selections))

	for _, selection := range selections {
		// Create a base model for the character
		character := &models.UserCharacter{
			UserID:      userID,
			CharacterID: selection.CharacterID,
			Name:        selection.Name,
			Realm:       selection.Realm,
			Region:      selection.Region,
			Class:       selection.Class,
			Race:        selection.Race,
			Level:       selection.Level,
			Faction:     selection.Faction,
			IsDisplayed: true,
		}

		// Save to database
		err := s.repo.CreateCharacter(character)
		if err != nil {
			results = append(results, CharacterSyncResult{
				Success:     false,
				CharacterID: selection.CharacterID,
				Message:     fmt.Sprintf("Error creating character: %v", err),
			})
			continue
		}

		// If it's the favorite character, set it as such
		if selection.IsFavorite {
			if err := s.repo.SetFavoriteCharacter(userID, character.ID); err != nil {
				log.Printf("Error setting favorite character %d: %v", character.ID, err)
			}
		}

		results = append(results, CharacterSyncResult{
			Success:     true,
			CharacterID: selection.CharacterID,
			Message:     "Character synchronized successfully",
		})
	}

	return results, nil
}

// GetUserCharacters retrieves all characters belonging to a user from the database
func (s *CharacterService) GetUserCharacters(ctx context.Context, userID uint) ([]models.UserCharacter, error) {
	return s.repo.GetCharactersByUserID(userID)
}

// GetCharacterDetails retrieves detailed information for a specific character
func (s *CharacterService) GetCharacterDetails(ctx context.Context, userID uint, characterID uint) (*models.UserCharacter, error) {
	// First verify ownership
	isOwner, err := s.IsCharacterOwner(ctx, userID, characterID)
	if err != nil {
		return nil, fmt.Errorf("error checking character ownership: %w", err)
	}

	if !isOwner {
		return nil, fmt.Errorf("character does not belong to this user")
	}

	return s.repo.GetCharacterByID(characterID)
}

// RefreshCharacter updates a character's information from the Battle.net API
func (s *CharacterService) RefreshCharacter(ctx context.Context, userID uint, characterID uint) error {
	// Get character from database first
	character, err := s.GetCharacterDetails(ctx, userID, characterID)
	if err != nil {
		return fmt.Errorf("error retrieving character: %w", err)
	}

	// Fetch updated character data
	// This would involve multiple API calls to different endpoints:
	// - Character profile summary
	// - Equipment
	// - Stats
	// - Talents
	// - Mythic+ profile
	// - Raid progress

	// For now, just update the last API update timestamp
	character.LastAPIUpdate = time.Now()
	return s.repo.UpdateCharacter(character)
}

// SetFavoriteCharacter sets a character as the user's favorite
func (s *CharacterService) SetFavoriteCharacter(ctx context.Context, userID uint, characterID uint) error {
	// First check if the character belongs to the user
	isOwner, err := s.IsCharacterOwner(ctx, userID, characterID)
	if err != nil {
		return fmt.Errorf("error checking character ownership: %w", err)
	}

	if !isOwner {
		return fmt.Errorf("character does not belong to this user")
	}

	return s.repo.SetFavoriteCharacter(userID, characterID)
}

// IsCharacterOwner checks if a character belongs to a user
func (s *CharacterService) IsCharacterOwner(ctx context.Context, userID uint, characterID uint) (bool, error) {
	character, err := s.repo.GetCharacterByID(characterID)
	if err != nil {
		return false, err
	}

	return character.UserID == userID, nil
}
