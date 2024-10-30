package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PlayerRanking représente un classement de joueur dans un donjon
type PlayerRanking struct {
	gorm.Model
	DungeonID     int     `json:"dungeon_id" gorm:"index"`
	Name          string  `json:"name" gorm:"index"`
	Class         string  `json:"class"`
	Spec          string  `json:"spec"`
	Role          string  `json:"role"`
	Amount        float64 `json:"amount"`
	HardModeLevel int     `json:"hard_mode_level"`
	Duration      int64   `json:"duration"`
	StartTime     int64   `json:"start_time"`

	// Informations du rapport
	ReportCode      string `json:"report_code"`
	ReportFightID   int    `json:"report_fight_id"`
	ReportStartTime int64  `json:"report_start_time"`

	// Informations de la guilde
	GuildID      int    `json:"guild_id"`
	GuildName    string `json:"guild_name"`
	GuildFaction int    `json:"guild_faction"`

	// Informations du serveur
	ServerID     int    `json:"server_id"`
	ServerName   string `json:"server_name"`
	ServerRegion string `json:"server_region"`

	// Autres informations
	BracketData int       `json:"bracket_data"`
	Faction     int       `json:"faction"`
	Affixes     []int     `json:"affixes" gorm:"type:integer[]"`
	Medal       string    `json:"medal"`
	Score       float64   `json:"score"`
	Leaderboard int       `json:"leaderboard"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RankingsUpdateState struct {
	ID             uint `gorm:"primaryKey;default:1"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	LastUpdateTime time.Time
}

func (RankingsUpdateState) TableName() string {
	return "rankings_update_states"
}

// InitializeRankingsUpdateState initializes the rankings update state
func InitializeRankingsUpdateState(db *gorm.DB) error {
	return db.Exec(`
			INSERT INTO rankings_update_states (id, created_at, updated_at, last_update_time)
			VALUES (1, NOW(), NOW(), NOW() - INTERVAL '25 hours')
			ON CONFLICT (id) DO NOTHING
	`).Error
}

// GetOrCreateRankingsUpdateState gets or creates the rankings update state
func GetOrCreateRankingsUpdateState(db *gorm.DB) (*RankingsUpdateState, error) {
	var state RankingsUpdateState

	err := db.Transaction(func(tx *gorm.DB) error {
		// Trying to get the existing state with ID = 1
		result := tx.First(&state, 1)
		if result.Error == nil {
			return nil // State found
		}

		if result.Error != gorm.ErrRecordNotFound {
			return result.Error // Unexpected error
		}

		// Create a new state with ID = 1
		state = RankingsUpdateState{
			ID:             1,
			LastUpdateTime: time.Now().Add(-25 * time.Hour),
		}
		return tx.Create(&state).Error
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get or create rankings update state: %w", err)
	}

	return &state, nil
}
