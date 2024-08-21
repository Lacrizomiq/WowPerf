package models

import (
	"time"
)

// MythicPlusRun represents a run in the Mythic+ dungeon
// It includes the seasons, the dungeon, the affixes, the members, the score, the duration, etc.
type MythicPlusRun struct {
	CompletedTimestamp    time.Time
	DungeonID             uint
	Dungeon               Dungeon
	ShortName             string
	Duration              int64
	IsCompletedWithinTime bool
	KeyStoneUpgrades      int
	KeystoneLevel         int
	MythicRating          float64
	SeasonID              uint
	Season                Season
	Affixes               []Affix
	Members               []MythicPlusRunMember
}

// MythicPlusRunMember represents a member in the Mythic+ dungeon run
type MythicPlusRunMember struct {
	CharacterID       uint
	CharacterName     string
	RealmID           uint
	RealmName         string `json:"-"`
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
