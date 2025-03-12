package character

import (
	"fmt"
	"wowperf/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CharacterRepository struct {
	db *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) *CharacterRepository {
	return &CharacterRepository{db}
}

// CreateCharacter create a new character in the database
func (r *CharacterRepository) CreateCharacter(character *models.UserCharacter) error {
	// Use clause OnConflit to avoid duplicates
	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "character_id"}, {Name: "realm"}, {Name: "region"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "class", "race", "level", "faction", "is_displayed"}),
	}).Create(character)

	return result.Error
}

// SetFavoriteCharacter define a character as favorite
func (r *CharacterRepository) SetFavoriteCharacter(userID uint, characterID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("favorite_character_id", characterID).Error
}

// GetCharacterByID retrieves a character by its ID
func (r *CharacterRepository) GetCharacterByID(characterID uint) (*models.UserCharacter, error) {
	var character models.UserCharacter
	if err := r.db.First(&character, characterID).Error; err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	return &character, nil
}

// GetCharactersByUserID retrieves all characters belonging to a user
func (r *CharacterRepository) GetCharactersByUserID(userID uint) ([]models.UserCharacter, error) {
	var characters []models.UserCharacter
	if err := r.db.Where("user_id = ?", userID).Find(&characters).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve user characters: %w", err)
	}
	return characters, nil
}

// GetCharacterByGameID retrieves a character by its game ID, realm and region
func (r *CharacterRepository) GetCharacterByGameID(userID uint, characterID int64, realm, region string) (*models.UserCharacter, error) {
	var character models.UserCharacter
	if err := r.db.Where("user_id = ? AND character_id = ? AND realm = ? AND region = ?",
		userID, characterID, realm, region).First(&character).Error; err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	return &character, nil
}

// UpdateCharacter updates an existing character in the database
func (r *CharacterRepository) UpdateCharacter(character *models.UserCharacter) error {
	return r.db.Save(character).Error
}

// DeleteCharacter removes a character from the database
func (r *CharacterRepository) DeleteCharacter(characterID uint) error {
	return r.db.Delete(&models.UserCharacter{}, characterID).Error
}

// GetFavoriteCharacter gets the user's favorite character
func (r *CharacterRepository) GetFavoriteCharacter(userID uint) (*models.UserCharacter, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.FavoriteCharacterID == nil {
		return nil, nil
	}

	return r.GetCharacterByID(*user.FavoriteCharacterID)
}

// ToggleCharacterDisplay toggles the visibility of a character
func (r *CharacterRepository) ToggleCharacterDisplay(characterID uint, display bool) error {
	return r.db.Model(&models.UserCharacter{}).Where("id = ?", characterID).
		Update("is_displayed", display).Error
}
