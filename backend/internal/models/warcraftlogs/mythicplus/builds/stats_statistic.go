package warcraftlogsBuilds

import (
	"time"

	"gorm.io/gorm"
)

// StatStatistic represents statistics for a specific stat in a build
type StatStatistic struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`

	// Player classification
	Class string `gorm:"type:varchar(255);not null;index"`
	Spec  string `gorm:"type:varchar(255);not null;index"`

	// Encounter information
	EncounterID uint `gorm:"index"`

	// Stat identification
	StatName     string `gorm:"type:varchar(50);not null;index"`
	StatCategory string `gorm:"type:varchar(50);not null"` // "secondary" or "minor"

	// Stat value
	AvgValue float64 `gorm:"default:0"` // Average value
	MinValue float64 `gorm:"default:0"` // Minimum value observed
	MaxValue float64 `gorm:"default:0"` // Maximum value observed

	// Sample information
	SampleSize int `gorm:"default:0"` // Number of builds analyzed

	// Metrics related to keystone levels
	AvgKeystoneLevel float64 `gorm:"default:0"` // Average keystone level
	MinKeystoneLevel int     `gorm:"default:0"` // Minimum keystone level
	MaxKeystoneLevel int     `gorm:"default:0"` // Maximum keystone level

	// Item level correlation
	AvgItemLevel float64 `gorm:"default:0"` // Average item level
	MinItemLevel float64 `gorm:"default:0"` // Minimum item level
	MaxItemLevel float64 `gorm:"default:0"` // Maximum item level
}

func (StatStatistic) TableName() string {
	return "stat_statistics"
}
