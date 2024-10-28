package models

import (
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
	gorm.Model
	LastUpdateTime time.Time
}

func (RankingsUpdateState) TableName() string {
	return "rankings_update_states"
}
