package raiderioMythicPlusRunsModels

import (
	"time"

	"gorm.io/gorm"
)

type MythicPlusTeamComposition struct {
	ID              uint   `gorm:"primaryKey"`
	CompositionHash string `gorm:"uniqueIndex;not null"` // Hash pour grouper compositions identiques

	// Tank
	TankClass string `gorm:"not null"`
	TankSpec  string `gorm:"not null"`

	// Healer
	HealerClass string `gorm:"not null"`
	HealerSpec  string `gorm:"not null"`

	// DPS (ordonnés alphabétiquement)
	Dps1Class string `gorm:"not null"`
	Dps1Spec  string `gorm:"not null"`
	Dps2Class string `gorm:"not null"`
	Dps2Spec  string `gorm:"not null"`
	Dps3Class string `gorm:"not null"`
	Dps3Spec  string `gorm:"not null"`

	// Relations
	MythicRuns []MythicPlusRuns      `gorm:"foreignKey:TeamCompositionID"`
	RunRoster  []MythicPlusRunRoster `gorm:"foreignKey:TeamCompositionID"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (MythicPlusTeamComposition) TableName() string {
	return "mythicplus_team_compositions"
}
