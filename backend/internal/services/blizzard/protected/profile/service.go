// services/blizzard/protected/profile/service.go

package protectedProfile

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

const (
	apiURL          = "https://%s.api.blizzard.com"
	accountEndpoint = "/profile/user/wow"
)

type ProtectedProfileService struct {
	protectedClient types.ProtectedClientInterface
	repository      *CharacterRepository
	db              *gorm.DB
}

func NewProtectedProfileService(protectedClient types.ProtectedClientInterface, db *gorm.DB) *ProtectedProfileService {
	return &ProtectedProfileService{
		protectedClient: protectedClient,
		repository:      NewCharacterRepository(db),
		db:              db,
	}
}

// GetAccountProfile get the account profile summary
func (s *ProtectedProfileService) GetAccountProfile(ctx context.Context, userID uint, params types.ProfileServiceParams) (map[string]interface{}, error) {
	// Validation
	if params.Region == "" {
		return nil, fmt.Errorf("region is required")
	}
	if params.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if params.Locale == "" {
		params.Locale = "en_US"
	}

	// Endpoint URL with the base URL
	endpoint := "/profile/user/wow"

	// Calling the protected client
	data, err := s.protectedClient.MakeProtectedRequest(
		ctx,
		userID,
		endpoint, // Just the endpoint, not the full URL
		params.Namespace,
		params.Locale,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get WoW profile: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse profile data: %w", err)
	}

	return result, nil
}

// GetProtectedCharacterProfile retrieve the profile of a specific character
func (s *ProtectedProfileService) GetProtectedCharacterProfile(
	ctx context.Context,
	userID uint,
	realmId string,
	characterId string,
	params types.ProfileServiceParams,
) (map[string]interface{}, error) {
	// Validation
	if params.Region == "" {
		return nil, fmt.Errorf("region is required")
	}
	if params.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	if params.Locale == "" {
		params.Locale = "en_US"
	}

	// Build the endpoint
	endpoint := fmt.Sprintf("/profile/user/wow/protected-character/%s-%s", realmId, characterId)

	// Call the protected client
	data, err := s.protectedClient.MakeProtectedRequest(
		ctx,
		userID,
		endpoint,
		params.Namespace,
		params.Locale,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get protected character profile: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse profile data: %w", err)
	}

	return result, nil
}

// ListAccountCharacters extract all the characters that have a 80+ level
func (s *ProtectedProfileService) ListAccountCharacters(ctx context.Context, userID uint, region string) ([]CharacterBasicInfo, error) {
	// Retrieve the full profil
	params := types.ProfileServiceParams{
		Region:    region,
		Namespace: fmt.Sprintf("profile-%s", region),
		Locale:    "en_US",
	}

	profileData, err := s.GetAccountProfile(ctx, userID, params)
	if err != nil {
		return nil, fmt.Errorf("error retrieving account profile: %w", err)
	}

	// Extract and filter 80+ level chars
	var characters []CharacterBasicInfo

	wowAccounts, ok := profileData["wow_accounts"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid profile data format")
	}

	for _, account := range wowAccounts {
		accountMap, ok := account.(map[string]interface{})
		if !ok {
			continue
		}

		chars, ok := accountMap["characters"].([]interface{})
		if !ok {
			continue
		}

		for _, char := range chars {
			charMap, ok := char.(map[string]interface{})
			if !ok {
				continue
			}

			// Extract char informations
			level, _ := charMap["level"].(float64)

			// Filter 80+ lvl char
			if int(level) >= 80 {
				// Extract realm data
				realmData, _ := charMap["realm"].(map[string]interface{})
				if realmData == nil {
					continue
				}

				// Extract class / race data
				classData, _ := charMap["playable_class"].(map[string]interface{})
				raceData, _ := charMap["playable_race"].(map[string]interface{})
				factionData, _ := charMap["faction"].(map[string]interface{})

				character := CharacterBasicInfo{
					CharacterID: int64(charMap["id"].(float64)),
					Name:        charMap["name"].(string),
					Realm:       realmData["slug"].(string),
					Region:      region,
					Class:       classData["name"].(string),
					Race:        raceData["name"].(string),
					Level:       int(level),
					Faction:     factionData["type"].(string),
				}

				characters = append(characters, character)
			}
		}
	}

	return characters, nil
}

// SyncAllAccountCharacters synchronizes all level 80+ characters from the account
func (s *ProtectedProfileService) SyncAllAccountCharacters(ctx context.Context, userID uint, region string) (int, error) {
	// Get all characters level 80+
	characters, err := s.ListAccountCharacters(ctx, userID, region)
	if err != nil {
		return 0, fmt.Errorf("error retrieving account characters: %w", err)
	}

	successCount := 0

	for _, char := range characters {
		// Create a model for the character
		character := &models.UserCharacter{
			UserID:      userID,
			CharacterID: char.CharacterID,
			Name:        char.Name,
			Realm:       char.Realm,
			Region:      char.Region,
			Class:       char.Class,
			Race:        char.Race,
			Level:       char.Level,
			Faction:     char.Faction,
			IsDisplayed: true,
		}

		// Save the character in the database
		err := s.repository.CreateCharacter(character)
		if err != nil {
			// Log but continue with other characters
			log.Printf("Error creating character %s: %v", char.Name, err)
			continue
		}

		successCount++
	}

	return successCount, nil
}

// RefreshUserCharacters refreshes all user characters and adds new ones if necessary
func (s *ProtectedProfileService) RefreshUserCharacters(ctx context.Context, userID uint, region string) (int, int, error) {
	// First step: synchronize all account characters (new and existing)
	newOrUpdatedCount, err := s.SyncAllAccountCharacters(ctx, userID, region)
	if err != nil {
		return 0, 0, err
	}

	// Second step: update all existing characters
	existingChars, err := s.repository.GetCharactersByUserID(userID)
	if err != nil {
		return newOrUpdatedCount, 0, err
	}

	updatedCount := 0
	for i := range existingChars {
		// Update timestamp and potentially other info
		existingChars[i].LastAPIUpdate = time.Now()
		if err := s.repository.UpdateCharacter(&existingChars[i]); err != nil {
			// Log but continue
			log.Printf("Error updating character %s: %v", existingChars[i].Name, err)
			continue
		}
		updatedCount++
	}

	return newOrUpdatedCount, updatedCount, nil
}

// GetUserCharacters get all the characters of a user
func (s *ProtectedProfileService) GetUserCharacters(ctx context.Context, userID uint) ([]models.UserCharacter, error) {
	return s.repository.GetCharactersByUserID(userID)
}

// SetFavoriteCharacter set a character as favorite
func (s *ProtectedProfileService) SetFavoriteCharacter(ctx context.Context, userID uint, characterID uint) error {
	// Check that the character belong to the user
	character, err := s.repository.GetCharacterByID(characterID)
	if err != nil {
		return err
	}

	if character.UserID != userID {
		return fmt.Errorf("character does not belong to this user")
	}

	return s.repository.SetFavoriteCharacter(userID, characterID)
}

// ToggleCharacterDisplay activates or deactivates the display of a character
func (s *ProtectedProfileService) ToggleCharacterDisplay(ctx context.Context, userID uint, characterID uint, display bool) error {
	// Check that the character belongs to the user
	character, err := s.repository.GetCharacterByID(characterID)
	if err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return fmt.Errorf("character does not belong to this user")
	}

	return s.repository.ToggleCharacterDisplay(characterID, display)
}
