package warcraftlogsBuilds

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Report struct {
	Code          string         `gorm:"primaryKey;type:varchar(255)"`
	FightID       int            `gorm:"not null"`
	EncounterID   uint           `gorm:"index"`
	RawData       datatypes.JSON `gorm:"type:jsonb;not null"`
	PlayerDetails datatypes.JSON `gorm:"type:jsonb"`
	KeystoneLevel int
	KeystoneTime  int64
	Affixes       []int `gorm:"type:integer[]"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *gorm.DeletedAt `gorm:"index"`
}

func (Report) TableName() string {
	return "warcraft_logs_reports"
}
