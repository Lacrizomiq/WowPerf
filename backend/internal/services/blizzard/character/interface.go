package character

import (
	"wowperf/internal/models"
)

/*

	CharacterServiceInterface defines the interface for character-related operations
	CharacterRepositoryInterface defines the interface for character database operations

	These interfaces are used to decouple the business logic from the database operations
	They are implemented by the CharacterService and CharacterRepository structs
	They are used to allow for mocking the service and repository in tests
	They are used to allow for dependency injection of the service and repository

*/

// CharacterServiceInterface defines the interface for character-related operations
type CharacterServiceInterface interface {
	// Core character operations
	GetCharacterDetails(character *models.UserCharacter) error
	UpdateCharacterDetails(character *models.UserCharacter) error
	UpdateAllCharactersForUser(userID uint) (int, error)

	// Character retrieval
	GetCharacterByID(characterID uint) (*models.UserCharacter, error)
	GetCharactersByUserID(userID uint) ([]models.UserCharacter, error)
	GetCharacterByGameID(userID uint, characterID int64, realm, region string) (*models.UserCharacter, error)
	GetFavoriteCharacter(userID uint) (*models.UserCharacter, error)

	// Character management
	CreateCharacter(character *models.UserCharacter) error
	CreateOrUpdateCharacter(character *models.UserCharacter) error
	DeleteCharacter(characterID uint) error
	SetFavoriteCharacter(userID uint, characterID uint) error
	ToggleCharacterDisplay(userID uint, characterID uint, display bool) error
}

// CharacterRepositoryInterface defines the interface for character database operations
type CharacterRepositoryInterface interface {
	// CRUD operations
	CreateCharacter(character *models.UserCharacter) error
	CreateOrUpdateCharacter(character *models.UserCharacter) error
	GetCharacterByID(characterID uint) (*models.UserCharacter, error)
	GetCharactersByUserID(userID uint) ([]models.UserCharacter, error)
	GetCharacterByGameID(userID uint, characterID int64, realm, region string) (*models.UserCharacter, error)
	UpdateCharacter(character *models.UserCharacter) error
	UpdateCharacterSummary(character *models.UserCharacter) error
	DeleteCharacter(characterID uint) error

	// Special operations
	SetFavoriteCharacter(userID uint, characterID uint) error
	GetFavoriteCharacter(userID uint) (*models.UserCharacter, error)
	ToggleCharacterDisplay(characterID uint, display bool) error
}

// Ensure that the concrete types implement these interfaces
var _ CharacterServiceInterface = (*CharacterService)(nil)
var _ CharacterRepositoryInterface = (*CharacterRepository)(nil)
