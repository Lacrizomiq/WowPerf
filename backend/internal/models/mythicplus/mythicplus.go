package models

import (
	"time"

	"gorm.io/gorm"
)

// Season represents a season in the Mythic+ dungeon pool
// Exeample: TWW Season 1
type Season struct {
	gorm.Model
	Slug            string `gorm:"uniqueIndex"`
	Name            string
	ShortName       string
	SeasonalAffixID *uint
	SeasonalAffix   *SeasonalAffix
	StartsUS        time.Time
	StartsEU        time.Time
	StartsTW        time.Time
	StartsKR        time.Time
	StartsCN        time.Time
	EndsUS          time.Time
	EndsEU          time.Time
	EndsTW          time.Time
	EndsKR          time.Time
	EndsCN          time.Time
	Dungeons        []Dungeon `gorm:"many2many:season_dungeons;"`
}

// SeasonalAffix represents a seasonal affix in the Mythic+ dungeon
// No more SeasonalAffix since DF season 1 but here for reference to older seasons
// Example: Seasonal affix for Dragonflight S1
type SeasonalAffix struct {
	gorm.Model
	ID   uint `gorm:"primaryKey"`
	Name string
	Icon string
}

// Dungeon represents a dungeon in the Mythic+ dungeon pool
// Example: Mists of Tirna Scithe
type Dungeon struct {
	gorm.Model
	ID               uint   `gorm:"primaryKey"`
	ChallengeModeID  *uint  `gorm:"index"`
	Slug             string `gorm:"uniqueIndex"`
	Name             string
	ShortName        string
	MediaURL         string
	Icon             *string
	KeyStoneUpgrades []KeyStoneUpgrade
	Seasons          []Season `gorm:"many2many:season_dungeons;"`
}

// SeasonDungeon represents a season-dungeon relationship in the Mythic+ dungeon pool
type SeasonDungeon struct {
	SeasonID  uint
	DungeonID uint
}

// Affix represents an affix in the Mythic+ dungeon hunt
type Affix struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	Name        string
	Icon        string
	Description string
	WowheadURL  string
}

// KeyStoneUpgrades represents the number of keystone upgrades for a dungeon
type KeyStoneUpgrade struct {
	gorm.Model
	DungeonID          uint
	QualifyingDuration int
	UpgradeLevel       int
}
