// services/blizzard/protected/profile/service.go

package protectedProfile

import (
	"context"
	"encoding/json"
	"fmt"
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

// SyncSelectedCharacters synchronize the selectect character with the database
func (s *ProtectedProfileService) SyncSelectedCharacters(ctx context.Context, userID uint, selections []CharacterSelection) ([]CharacterSyncResult, error) {
	results := make([]CharacterSyncResult, 0, len(selections))

	for _, selection := range selections {
		// Create a model for the character
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

		// Save the character in the database
		err := s.repository.CreateCharacter(character)
		if err != nil {
			results = append(results, CharacterSyncResult{
				Success:     false,
				CharacterID: selection.CharacterID,
				Message:     fmt.Sprintf("Error creating character: %v", err),
			})
			continue
		}

		// Set character favorite if its the case
		if selection.IsFavorite {
			if err := s.repository.SetFavoriteCharacter(userID, character.ID); err != nil {
				// Continue even if there is an issue with favorite
				results = append(results, CharacterSyncResult{
					Success:     true,
					CharacterID: selection.CharacterID,
					Message:     "Character synchronized but failed to set as favorite",
				})
				continue
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

// ToggleCharacterDisplay active ou désactive l'affichage d'un personnage
func (s *ProtectedProfileService) ToggleCharacterDisplay(ctx context.Context, userID uint, characterID uint, display bool) error {
	// Vérifier que le personnage appartient à l'utilisateur
	character, err := s.repository.GetCharacterByID(characterID)
	if err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return fmt.Errorf("character does not belong to this user")
	}

	return s.repository.ToggleCharacterDisplay(characterID, display)
}
