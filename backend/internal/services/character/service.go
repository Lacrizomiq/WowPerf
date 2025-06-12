package character

import (
	"context"
	"fmt"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/character/enrichers"

	"gorm.io/gorm"
)

// CharacterService handles operations related to characters
// Implements the CharacterServiceInterface to be used by the CharacterService
type CharacterService struct {
	db             *gorm.DB
	profileService *blizzard.ProfileService
	repository     CharacterRepositoryInterface
}

// NewCharacterService creates a new character service
// Returns a CharacterServiceInterface to be used by the CharacterService
// This allows for mocking the service in tests and dependency injection
func NewCharacterService(db *gorm.DB, profileService *blizzard.ProfileService) CharacterServiceInterface {
	return &CharacterService{
		db:             db,
		profileService: profileService,
		repository:     NewCharacterRepository(db),
	}
}

// GetCharacterDetails retrieves the details of a character and saves them
func (s *CharacterService) GetCharacterDetails(character *models.UserCharacter) error {
	// Utilise le nouvel enrichisseur summary
	summaryEnricher := enrichers.NewSummaryEnricher(s.profileService)
	if err := summaryEnricher.EnrichCharacter(context.Background(), character); err != nil {
		return fmt.Errorf("failed to get character details: %w", err)
	}

	// Save updated character to database
	return s.repository.UpdateCharacterSummary(character)
}

// UpdateCharacterDetails retrieves and saves the details of a character
func (s *CharacterService) UpdateCharacterDetails(character *models.UserCharacter) error {
	return s.GetCharacterDetails(character)
}

// UpdateAllCharactersForUser updates all characters for a specific user
func (s *CharacterService) UpdateAllCharactersForUser(userID uint) (int, error) {
	characters, err := s.repository.GetCharactersByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve user characters: %w", err)
	}

	successCount := 0
	for i := range characters {
		if err := s.UpdateCharacterDetails(&characters[i]); err != nil {
			// Log error but continue with other characters
			fmt.Printf("Error updating character %s: %v\n", characters[i].Name, err)
			continue
		}
		successCount++
	}

	return successCount, nil
}

// GetCharacterByID retrieves a character by its ID
func (s *CharacterService) GetCharacterByID(characterID uint) (*models.UserCharacter, error) {
	return s.repository.GetCharacterByID(characterID)
}

// GetCharactersByUserID retrieves all characters belonging to a user
func (s *CharacterService) GetCharactersByUserID(userID uint) ([]models.UserCharacter, error) {
	return s.repository.GetCharactersByUserID(userID)
}

// GetFavoriteCharacter gets the user's favorite character
func (s *CharacterService) GetFavoriteCharacter(userID uint) (*models.UserCharacter, error) {
	return s.repository.GetFavoriteCharacter(userID)
}

// SetFavoriteCharacter sets a character as the user's favorite
func (s *CharacterService) SetFavoriteCharacter(userID uint, characterID uint) error {
	// Verify the character belongs to the user
	character, err := s.GetCharacterByID(characterID)
	if err != nil {
		return err
	}

	if character.UserID != userID {
		return fmt.Errorf("character does not belong to this user")
	}

	return s.repository.SetFavoriteCharacter(userID, characterID)
}

// ToggleCharacterDisplay toggles the visibility of a character
func (s *CharacterService) ToggleCharacterDisplay(userID uint, characterID uint, display bool) error {
	// Verify the character belongs to the user
	character, err := s.GetCharacterByID(characterID)
	if err != nil {
		return err
	}

	if character.UserID != userID {
		return fmt.Errorf("character does not belong to this user")
	}

	return s.repository.ToggleCharacterDisplay(characterID, display)
}

// CreateCharacter creates a new character in the database
func (s *CharacterService) CreateCharacter(character *models.UserCharacter) error {
	return s.repository.CreateCharacter(character)
}

// CreateOrUpdateCharacter creates or updates a character with all fields
func (s *CharacterService) CreateOrUpdateCharacter(character *models.UserCharacter) error {
	return s.repository.CreateOrUpdateCharacter(character)
}

// GetCharacterByGameID retrieves a character by its game ID, realm and region
func (s *CharacterService) GetCharacterByGameID(userID uint, characterID int64, realm, region string) (*models.UserCharacter, error) {
	return s.repository.GetCharacterByGameID(userID, characterID, realm, region)
}

// DeleteCharacter removes a character from the database
func (s *CharacterService) DeleteCharacter(characterID uint) error {
	return s.repository.DeleteCharacter(characterID)
}
