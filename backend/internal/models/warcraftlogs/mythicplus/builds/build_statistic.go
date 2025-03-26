package warcraftlogsBuilds

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// BuildStatistic represents statistics for a specific build
type BuildStatistic struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`

	// Class and spec of the build
	Class string `gorm:"type:varchar(255);not null;index"`
	Spec  string `gorm:"type:varchar(255);not null;index"`

	// Encounter ID of the build
	EncounterID uint `gorm:"index"`

	// base item informations
	ItemSlot int `gorm:"index"`
	ItemID   int `gorm:"index"`

	// item detailled informations
	ItemName    string  `gorm:"type:varchar(255)"`
	ItemIcon    string  `gorm:"type:varchar(255)"`
	ItemQuality int     `gorm:"default:0"` // 0-4 quality
	ItemLevel   float64 `gorm:"default:0"`

	// Class Set informations
	HasSetBonus bool `gorm:"default:false"`
	SetID       int  `gorm:"default:0"`

	// Bonus IDs
	BonusIDs pq.Int64Array `gorm:"type:integer[]"`

	// Gems informations
	HasGems   bool            `gorm:"default:false"`
	GemsCount int             `gorm:"default:0"`
	GemIDs    pq.Int64Array   `gorm:"type:integer[]"`
	GemIcons  pq.StringArray  `gorm:"type:text[]"`
	GemLevels pq.Float64Array `gorm:"type:numeric[]"`

	// Enchant informations
	HasPermanentEnchant  bool   `gorm:"default:false"`
	PermanentEnchantID   int    `gorm:"default:0"`
	PermanentEnchantName string `gorm:"type:varchar(255)"`

	HasTemporaryEnchant  bool   `gorm:"default:false"`
	TemporaryEnchantID   int    `gorm:"default:0"`
	TemporaryEnchantName string `gorm:"type:varchar(255)"`

	// Usage statistics
	UsageCount      int     `gorm:"default:0"`
	UsagePercentage float64 `gorm:"default:0"`
	AvgItemLevel    float64 `gorm:"default:0"`
	MinItemLevel    float64 `gorm:"default:0"`
	MaxItemLevel    float64 `gorm:"default:0"`

	// Keystone level metrics
	AvgKeystoneLevel float64 `gorm:"default:0"`
	MinKeystoneLevel int     `gorm:"default:0"`
	MaxKeystoneLevel int     `gorm:"default:0"`
}

func (BuildStatistic) TableName() string {
	return "build_statistics"
}
