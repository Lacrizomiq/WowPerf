package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DungeonStats struct {
	gorm.Model  `json:"-"`
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

type TeamCompMap map[string]TeamCompStats

// Scan implements the Scanner interface for TeamCompMap
func (t *TeamCompMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp map[string]TeamCompStats
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*t = TeamCompMap(temp)
	return nil
}

// Value implements the driver Valuer interface for TeamCompMap
func (t TeamCompMap) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (DungeonStats) TableName() string {
	return "dungeon_stats"
}
