package warcraftlogsBuilds

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PlayerBuild struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`

	PlayerName string `gorm:"type:varchar(255);not null"`
	Class      string `gorm:"type:varchar(255);not null"`
	Spec       string `gorm:"type:varchar(255);not null"`

	ReportCode string `gorm:"type:varchar(255)"`
	FightID    int    `gorm:"not null"`

	TalentImport string         `gorm:"type:text"`
	TalentTree   datatypes.JSON `gorm:"type:jsonb"`
	TalentTreeID int
	ActorID      int `gorm:"index"`

	Gear          datatypes.JSON `gorm:"type:jsonb"`
	Stats         datatypes.JSON `gorm:"type:jsonb"`
	CombatantInfo datatypes.JSON `gorm:"type:jsonb"`

	DungeonID   uint `gorm:"index"`
	EncounterID uint `gorm:"index"`
}

func (PlayerBuild) TableName() string {
	return "player_builds"
}
