package raiderioMythicPlusRunsModels

import (
	"time"

	"gorm.io/gorm"
)

type MythicPlusRuns struct {
	ID                uint   `gorm:"primaryKey"`
	KeystoneRunID     int64  `gorm:"uniqueIndex;not null"`
	Season            string `gorm:"index;not null"`
	Region            string `gorm:"index;not null"`
	DungeonSlug       string `gorm:"index;not null"`
	DungeonName       string
	MythicLevel       int     `gorm:"index;not null"`
	Score             float64 `gorm:"index"`
	Status            string  `gorm:"index"`
	ClearTimeMs       int64
	KeystoneTimeMs    int64
	CompletedAt       time.Time `gorm:"index"`
	NumChests         int
	TimeRemainingMs   int64
	TeamCompositionID *uint `gorm:"index"`

	// Relations
	TeamComposition *MythicPlusTeamComposition `gorm:"foreignKey:TeamCompositionID"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (MythicPlusRuns) TableName() string {
	return "mythicplus_runs"
}
