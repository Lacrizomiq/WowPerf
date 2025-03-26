package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserCharacter represents a WoW character linked to a user account
type UserCharacter struct {
	gorm.Model
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"index;not null" json:"user_id"`

	// Character identifiers
	CharacterID int64  `gorm:"index;not null" json:"character_id"`
	Name        string `gorm:"not null" json:"name"`
	Realm       string `gorm:"not null" json:"realm"`
	Region      string `gorm:"not null" json:"region"`

	// Basic information
	Class          string `gorm:"not null" json:"class"`
	Race           string `gorm:"not null" json:"race"`
	Gender         string `gorm:"not null" json:"gender"`
	Faction        string `gorm:"not null" json:"faction"`
	ActiveSpecName string `gorm:"not null" json:"active_spec_name"`
	ActiveSpecID   int    `gorm:"not null" json:"active_spec_id"`
	ActiveSpecRole string `gorm:"not null" json:"active_spec_role"`

	// Level and main stats
	Level                 int     `json:"level"`
	ItemLevel             float64 `json:"item_level"`
	MythicPlusRating      float64 `json:"mythic_plus_rating"`
	MythicPlusRatingColor string  `json:"mythic_plus_rating_color"`
	AchievementPoints     int     `json:"achievement_points"`
	HonorableKills        int     `json:"honorable_kills"`

	// Image URLs
	AvatarURL      string `json:"avatar_url"`
	InsetAvatarURL string `json:"inset_avatar_url"`
	MainRawURL     string `json:"main_raw_url"`
	ProfileURL     string `json:"profile_url"`

	// JSON structured data
	EquipmentJSON  datatypes.JSON `gorm:"type:jsonb" json:"equipment_json,omitempty"`
	StatsJSON      datatypes.JSON `gorm:"type:jsonb" json:"stats_json,omitempty"`
	TalentsJSON    datatypes.JSON `gorm:"type:jsonb" json:"talents_json,omitempty"`
	MythicPlusJSON datatypes.JSON `gorm:"type:jsonb" json:"mythic_plus_json,omitempty"`
	RaidsJSON      datatypes.JSON `gorm:"type:jsonb" json:"raids_json,omitempty"`

	// Metadata
	IsDisplayed   bool      `gorm:"default:true" json:"is_displayed"`
	LastAPIUpdate time.Time `json:"last_api_update"`
}

// TableName overrides the table name
func (UserCharacter) TableName() string {
	return "user_characters"
}

// BeforeCreate ensures uniqueness of character across realm and region
func (uc *UserCharacter) BeforeCreate(tx *gorm.DB) error {
	// Set unique composite index on character_id, realm, region
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   []clause.Column{{Name: "character_id"}, {Name: "realm"}, {Name: "region"}},
		DoNothing: true,
	})
	return nil
}
