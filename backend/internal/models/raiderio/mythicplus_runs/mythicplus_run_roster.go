package raiderioMythicPlusRunsModels

import (
	"time"

	"gorm.io/gorm"
)

type MythicPlusRunRoster struct {
	ID                uint   `gorm:"primaryKey"`
	TeamCompositionID uint   `gorm:"index;not null"`
	Role              string `gorm:"index;not null"` // tank, healer, dps
	ClassName         string `gorm:"index;not null"`
	SpecName          string `gorm:"index;not null"`

	// Relations
	TeamComposition MythicPlusTeamComposition `gorm:"foreignKey:TeamCompositionID"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (MythicPlusRunRoster) TableName() string {
	return "mythicplus_run_roster"
}
