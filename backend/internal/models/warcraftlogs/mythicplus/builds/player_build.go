package warcraftlogsBuilds

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PlayerBuild struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`

	// Player information
	PlayerName string `gorm:"type:varchar(255);not null"`
	Class      string `gorm:"type:varchar(255);not null"`
	Spec       string `gorm:"type:varchar(255);not null"`

	// Report information
	ReportCode string `gorm:"type:varchar(255)"`
	FightID    int    `gorm:"not null"`

	// Talent information
	TalentImport string         `gorm:"column:talent_import;type:text"`
	TalentTree   datatypes.JSON `gorm:"type:jsonb"`
	ActorID      int            `gorm:"index"`

	// Equipment and stats
	ItemLevel float64        `gorm:"type:numeric"`
	Gear      datatypes.JSON `gorm:"type:jsonb"`
	Stats     datatypes.JSON `gorm:"type:jsonb"`

	// Dungeon information
	EncounterID uint `gorm:"index"`

	// Mythic+ information
	KeystoneLevel int           `gorm:"index"`
	Affixes       pq.Int64Array `gorm:"type:integer[]"`
}

func (PlayerBuild) TableName() string {
	return "player_builds"
}
