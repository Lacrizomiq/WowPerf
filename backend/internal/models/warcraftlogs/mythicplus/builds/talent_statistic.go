package warcraftlogsBuilds

import (
	"time"

	"gorm.io/gorm"
)

type TalentStatistic struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`

	Class       string `gorm:"type:varchar(255);not null"`
	Spec        string `gorm:"type:varchar(255);not null"`
	DungeonID   uint   `gorm:"index"`
	EncounterID uint   `gorm:"index"`

	TalentImport    string  `gorm:"type:text"`
	UsageCount      int     `gorm:"default:0"`
	UsagePercentage float64 `gorm:"default:0"`
	AvgScore        float64 `gorm:"default:0"`

	PeriodStart time.Time
	PeriodEnd   time.Time
}

func (TalentStatistic) TableName() string {
	return "talent_statistics"
}
