package models

import (
	"time"

	"gorm.io/gorm"
)

// Season represents a season in the Mythic+ dungeon pool
// Exeample: TWW Season 1
type Season struct {
	*gorm.Model     `json:"-"`
	Slug            string `gorm:"uniqueIndex"`
	Name            string
	ShortName       string
	SeasonalAffixID *uint          `json:"-"`
	SeasonalAffix   *SeasonalAffix `json:"-"`
	StartsUS        time.Time      `json:"-"`
	StartsEU        time.Time      `json:"-"`
	StartsTW        time.Time      `json:"-"`
	StartsKR        time.Time      `json:"-"`
	StartsCN        time.Time      `json:"-"`
	EndsUS          time.Time      `json:"-"`
	EndsEU          time.Time      `json:"-"`
	EndsTW          time.Time      `json:"-"`
	EndsKR          time.Time      `json:"-"`
	EndsCN          time.Time      `json:"-"`
	Dungeons        []Dungeon      `gorm:"many2many:season_dungeons;"`
}

// SeasonalAffix represents a seasonal affix in the Mythic+ dungeon
// No more SeasonalAffix since DF season 1 but here for reference to older seasons
// Example: Seasonal affix for Dragonflight S1
type SeasonalAffix struct {
	*gorm.Model `json:"-"`
	ID          uint `gorm:"primaryKey"`
	Name        string
	Icon        string
}

// Dungeon represents a dungeon in the Mythic+ dungeon pool
// Example: Mists of Tirna Scithe
type Dungeon struct {
	*gorm.Model      `json:"-"`
	ID               uint   `gorm:"primaryKey"`
	ChallengeModeID  *uint  `gorm:"uniqueIndex"`
	Slug             string `gorm:"uniqueIndex"`
	Name             string
	ShortName        string
	MediaURL         string
	Icon             *string
	KeyStoneUpgrades []KeyStoneUpgrade `gorm:"foreignKey:ChallengeModeID;references:ChallengeModeID"`
	Seasons          []Season          `gorm:"many2many:season_dungeons;"`
}

// SeasonDungeon represents a season-dungeon relationship in the Mythic+ dungeon pool
type SeasonDungeon struct {
	SeasonID  uint
	DungeonID uint
}

// Affix represents an affix in the Mythic+ dungeon hunt
type Affix struct {
	*gorm.Model `json:"-"`
	ID          uint `gorm:"primaryKey"`
	Name        string
	Icon        string
	Description string `json:"-"`
	WowheadURL  string
}

// KeyStoneUpgrades represents the number of keystone upgrades for a dungeon
type KeyStoneUpgrade struct {
	*gorm.Model        `json:"-"`
	ChallengeModeID    *uint `gorm:"index"`
	QualifyingDuration int
	UpgradeLevel       int
}
