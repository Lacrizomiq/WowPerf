package models

import "gorm.io/gorm"

// MythicDj represents a mythic dungeon
// For example the "Atal'Dazar" mythic dungeon and every static data related to it
type MythicDj struct {
	ID               uint              `gorm:"primaryKey"`
	Name             string            `gorm:"unique;not null"`
	ShortName        string            `gorm:"unique;not null"`
	IconURL          string            `gorm:"not null"`
	KeystoneUpgrades []KeystoneUpgrade `gorm:"foreignKey:MythicDjID"`
	Affixes          []Affix           `gorm:"foreignKey:ID"`
}

// KeystoneUpgrade represents the different keystone upgrades timer and upgrade level for a mythic dungeon
// For example the time needed for a +2 keystone upgrade
type KeystoneUpgrade struct {
	gorm.Model
	MythicDjID         uint `gorm:"not null"`
	QualifyingDuration int  `gorm:"not null"`
	UpgradeLevel       int  `gorm:"not null"`
}

// Affix represents an affix for a mythic dungeon
// For example the "Blessing of the Ancients" affix
type Affix struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"unique;not null"`
	IconURL string `gorm:"not null"`
}

// Period represents a period for a mythic dungeon
type Period struct {
	gorm.Model
	MythicDjID uint   `gorm:"not null"`
	PeriodID   uint   `gorm:"not null"`
	StartTime  string `gorm:"not null"`
	EndTime    string `gorm:"not null"`
}
