package rankings

import (
	"time"

	"gorm.io/gorm"
)

type PlayerRanking struct {
	gorm.Model
	DungeonID int       `json:"dungeon_id"`
	PlayerID  int       `json:"player_id"`
	Name      string    `json:"name"`
	Class     string    `json:"class"`
	Spec      string    `json:"spec"`
	Role      string    `json:"role"`
	Score     float64   `json:"score"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RankingsUpdateState struct {
	gorm.Model
	LastUpdateTime time.Time
}

func (RankingsUpdateState) TableName() string {
	return "rankings_update_states"
}
