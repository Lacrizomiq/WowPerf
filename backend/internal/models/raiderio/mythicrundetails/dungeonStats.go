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
	TeamComp    TeamCompMap               `gorm:"serializer:json"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

type TeamMember struct {
	Class string `json:"class"`
	Spec  string `json:"spec"`
}

type TeamComposition struct {
	Tank   TeamMember `json:"tank"`
	Healer TeamMember `json:"healer"`
	Dps1   TeamMember `json:"dps_1"`
	Dps2   TeamMember `json:"dps_2"`
	Dps3   TeamMember `json:"dps_3"`
}

type TeamCompStats struct {
	Count       int             `json:"count"`
	Composition TeamComposition `json:"composition"`
}

// TeamCompMap is a map of team compositions to their stats
type TeamCompMap map[string]TeamCompStats

func (DungeonStats) TableName() string {
	return "dungeon_stats"
}
