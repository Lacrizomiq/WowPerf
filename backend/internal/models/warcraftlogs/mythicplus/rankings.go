package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PlayerRanking represents a player ranking in a dungeon
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

	// Report information
	ReportCode      string `json:"report_code"`
	ReportFightID   int    `json:"report_fight_id"`
	ReportStartTime int64  `json:"report_start_time"`

	// Guild information
	GuildID      int    `json:"guild_id"`
	GuildName    string `json:"guild_name"`
	GuildFaction int    `json:"guild_faction"`

	// Server information
	ServerID     int    `json:"server_id"`
	ServerName   string `json:"server_name"`
	ServerRegion string `json:"server_region"`

	// Other information
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

// Structures for API data
type Report struct {
	Code      string `json:"code"`
	FightID   int    `json:"fightID"`
	StartTime int64  `json:"startTime"`
}

type Guild struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Faction int    `json:"faction"`
}

type Server struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type Run struct {
	DungeonID     int     `json:"dungeonId"`
	Score         float64 `json:"score"`
	Duration      int64   `json:"duration"`
	StartTime     int64   `json:"startTime"`
	HardModeLevel int     `json:"hardModeLevel"`
	BracketData   int     `json:"bracketData"`
	Medal         string  `json:"medal"`
	Affixes       []int   `json:"affixes"`
	Report        Report  `json:"report"`
}

type PlayerScore struct {
	Name       string  `json:"name"`
	Class      string  `json:"class"`
	Spec       string  `json:"spec"`
	Role       string  `json:"role"`
	TotalScore float64 `json:"totalScore"`
	Amount     float64 `json:"amount"`
	Guild      Guild   `json:"guild"`
	Server     Server  `json:"server"`
	Faction    int     `json:"faction"`
	Runs       []Run   `json:"runs"`
}

type RoleRankings struct {
	Players []PlayerScore `json:"players"`
	Count   int           `json:"count"`
}

type GlobalRankings struct {
	Tanks   RoleRankings `json:"tanks"`
	Healers RoleRankings `json:"healers"`
	DPS     RoleRankings `json:"dps"`
}

// Temporary data structure
type playerData struct {
	ranking   PlayerRanking
	dungeonID int
}
