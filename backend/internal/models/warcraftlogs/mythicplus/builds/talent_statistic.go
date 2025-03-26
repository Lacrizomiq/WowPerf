package warcraftlogsBuilds

import (
	"time"

	"gorm.io/gorm"
)

// TalentStatistic represents statistics for a specific talent in a build
type TalentStatistic struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`

	// Player classification
	Class string `gorm:"type:varchar(255);not null;index"`
	Spec  string `gorm:"type:varchar(255);not null;index"`

	// Encounter information
	EncounterID uint `gorm:"index"`

	// Talent import
	TalentImport string `gorm:"type:text;index"`

	// Usage statistics
	UsageCount      int     `gorm:"default:0"` // Number of players using this configuration
	UsagePercentage float64 `gorm:"default:0"` // Percentage of players using this configuration

	// Metrics for the keystone level where this build is used
	AvgKeystoneLevel float64 `gorm:"default:0"` // Average keystone level
	MinKeystoneLevel int     `gorm:"default:0"` // Minimum keystone level
	MaxKeystoneLevel int     `gorm:"default:0"` // Maximum keystone level

	// Item level statistics
	AvgItemLevel float64 `gorm:"default:0"` // Average item level
	MinItemLevel float64 `gorm:"default:0"` // Minimum item level
	MaxItemLevel float64 `gorm:"default:0"` // Maximum item level
}

func (TalentStatistic) TableName() string {
	return "talent_statistics"
}
