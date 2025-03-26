package warcraftlogsBuilds

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ClassRanking represents the ranking of a specific class in a specific spec
type ClassRanking struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	// Player info
	PlayerName string `gorm:"type:varchar(255);not null"`
	Class      string `gorm:"type:varchar(255);not null"`
	Spec       string `gorm:"type:varchar(255);not null"`

	// Dungeon info
	DungeonID   uint `gorm:"-"`
	EncounterID uint `gorm:"index"`

	// Run info
	Amount        float64
	HardModeLevel int
	Duration      int64
	StartTime     int64
	RankPosition  int
	Score         float64
	Medal         string `gorm:"type:varchar(50)"`
	Leaderboard   int

	// Server info
	ServerID     int
	ServerName   string `gorm:"type:varchar(255)"`
	ServerRegion string `gorm:"type:varchar(50)"`

	// Guild info
	GuildID      *int
	GuildName    *string `gorm:"type:varchar(255)"`
	GuildFaction *int

	// Report info
	ReportCode      string `gorm:"type:varchar(255)"`
	ReportFightID   int
	ReportStartTime int64

	// Other
	Faction int
	Affixes pq.Int64Array `gorm:"type:integer[]"`
}

// TableName specifies the table name for GORM
func (ClassRanking) TableName() string {
	return "class_rankings"
}
