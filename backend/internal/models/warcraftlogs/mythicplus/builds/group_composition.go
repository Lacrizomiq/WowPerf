package warcraftlogsBuilds

import (
	"time"

	"gorm.io/gorm"
)

type GroupComposition struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`

	ReportCode  string `gorm:"type:varchar(255)"`
	FightID     int
	TankSpecs   []string `gorm:"type:text[]"`
	HealerSpecs []string `gorm:"type:text[]"`
	DpsSpecs    []string `gorm:"type:text[]"`
	Success     bool

	DungeonID     uint `gorm:"index"`
	EncounterID   uint `gorm:"index"`
	KeystoneLevel int
}

func (GroupComposition) TableName() string {
	return "group_compositions"
}
