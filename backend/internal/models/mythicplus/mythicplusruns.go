package models

import (
	"time"

	"gorm.io/gorm"
)

// MythicPlusRun represents a run in the Mythic+ dungeon
// It includes the seasons, the dungeon, the affixes, the members, the score, the duration, etc.
type MythicPlusRun struct {
	gorm.Model
	CompletedTimestamp    time.Time
	DungeonID             uint
	Dungeon               Dungeon
	Duration              int64
	IsCompletedWithinTime bool
	KeystoneLevel         int
	MythicRating          float64
	SeasonID              uint
	Season                Season
	Affixes               []Affix `gorm:"many2many:mythic_run_affixes;"`
	Members               []MythicPlusRunMember
}

// MythicPlusRunMember represents a member in the Mythic+ dungeon run
type MythicPlusRunMember struct {
	gorm.Model
	MythicPlusRunID   uint
	CharacterID       uint
	CharacterName     string
	RealmID           uint
	RealmName         string
	RealmSlug         string
	EquippedItemLevel int
	RaceID            uint
	RaceName          string
	SpecializationID  uint
	Specialization    string
}

// MythicPlusRunAffix represents an affix in the Mythic+ dungeon
type MythicPlusRunAffix struct {
	MythicPlusRunID uint
	AffixID         uint
}
