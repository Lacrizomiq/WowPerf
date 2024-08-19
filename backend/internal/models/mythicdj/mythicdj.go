package models

// MythicDj represents a mythic dungeon
// For example the "Atal'Dazar" mythic dungeon and every static data related to it
type MythicDj struct {
	ID               uint              `gorm:"primaryKey"`
	Name             string            `gorm:"unique;not null"`
	IconURL          string            `gorm:"not null"`
	KeystoneUpgrades []KeystoneUpgrade `gorm:"foreignKey:MythicDjID"`
}

// KeystoneUpgrade represents the different keystone upgrades timer and upgrade level for a mythic dungeon
// For example the time needed for a +2 keystone upgrade
type KeystoneUpgrade struct {
	MythicDjID         uint `gorm:"uniqueIndex:idx_mythic_dj_upgrade"`
	QualifyingDuration int
	UpgradeLevel       int `gorm:"uniqueIndex:idx_mythic_dj_upgrade"`
}

// Affix represents an affix for a mythic dungeon
// For example the "Blessing of the Ancients" affix
type Affix struct {
	AffixID uint   `gorm:"primaryKey"`
	Name    string `gorm:"unique;not null"`
	IconURL string `gorm:"not null"`
}

// Period represents a period for a mythic dungeon
type Period struct {
	PeriodID  uint   `gorm:"primaryKey"`
	SeasonID  uint   `gorm:"not null"`
	StartTime string `gorm:"not null"`
	EndTime   string `gorm:"not null"`
}

// Season represents a season
// For example "Dragonflight S4"
type Season struct {
	SeasonID uint     `gorm:"primaryKey"`
	Name     string   `gorm:"not null"`
	Periods  []Period `gorm:"foreignKey:SeasonID"`
}

// To handle the many to many relationship between mythic dungeons, affixes and seasons
// Join table to store the relationship between the three tables
type MythicDjSeason struct {
	MythicDjID uint
	SeasonID   uint
}
