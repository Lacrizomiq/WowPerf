package warcraftlogsBuilds

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Report struct {
	// primary key and metadata
	Code        string `gorm:"primaryKey;type:varchar(255)"`
	FightID     int    `gorm:"not null"`
	EncounterID uint   `gorm:"index"`

	// base data
	TotalTime int64
	ItemLevel float64

	// composition and performances
	Composition datatypes.JSON `gorm:"type:jsonb"`
	DamageDone  datatypes.JSON `gorm:"type:jsonb"`
	HealingDone datatypes.JSON `gorm:"type:jsonb"`
	DamageTaken datatypes.JSON `gorm:"type:jsonb"`
	DeathEvents datatypes.JSON `gorm:"type:jsonb"`

	// player details
	PlayerDetailsDps     datatypes.JSON `gorm:"type:jsonb;column:player_details_dps"`
	PlayerDetailsHealers datatypes.JSON `gorm:"type:jsonb;column:player_details_healers"`
	PlayerDetailsTanks   datatypes.JSON `gorm:"type:jsonb;column:player_details_tanks"`

	// combat data
	LogVersion  int
	GameVersion int

	// mythic+ data
	KeystoneLevel   int            `gorm:"column:keystonelevel"`
	KeystoneTime    int64          `gorm:"column:keystonetime"`
	Affixes         pq.Int64Array  `gorm:"type:integer[]"`
	FriendlyPlayers pq.Int64Array  `gorm:"type:integer[]"`
	TalentCodes     datatypes.JSON `gorm:"type:jsonb"`

	// raw data
	RawData datatypes.JSON `gorm:"type:jsonb"`

	// tracking data
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

func (Report) TableName() string {
	return "warcraft_logs_reports"
}
