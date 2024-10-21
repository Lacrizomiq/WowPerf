package models

import (
	"time"

	"gorm.io/gorm"
)

type DungeonStats struct {
	gorm.Model
	Season      string                    `json:"season"`
	Region      string                    `json:"region"`
	DungeonSlug string                    `json:"dungeon_slug"`
	RoleStats   map[string]map[string]int `gorm:"serializer:json"`
	SpecStats   map[string]map[string]int `gorm:"serializer:json"`
	LevelStats  map[int]int               `gorm:"serializer:json"`
	TeamComp    map[string]int            `gorm:"serializer:json"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

func (DungeonStats) TableName() string {
	return "dungeon_stats"
}
