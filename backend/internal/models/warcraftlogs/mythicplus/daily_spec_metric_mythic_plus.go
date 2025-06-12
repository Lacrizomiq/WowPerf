package models

import (
	"time"

	"gorm.io/gorm"
)

// DailySpecMetricMythicPlus represents the daily metrics of a spec for a given dungeon
type DailySpecMetricMythicPlus struct {
	gorm.Model
	CaptureDate time.Time `json:"capture_date" gorm:"not null;index"`
	Spec        string    `json:"spec" gorm:"not null;index"`
	Class       string    `json:"class" gorm:"not null;index"`
	Role        string    `json:"role" gorm:"not null;index"`
	EncounterID int       `json:"encounter_id" gorm:"not null;index"`
	IsGlobal    bool      `json:"is_global" gorm:"not null;default:false;index"`
	AvgScore    float64   `json:"avg_score" gorm:"not null"`
	MaxScore    float64   `json:"max_score" gorm:"not null"`
	MinScore    float64   `json:"min_score" gorm:"not null"`
	AvgKeyLevel float64   `json:"avg_key_level" gorm:"not null"`
	MaxKeyLevel int       `json:"max_key_level" gorm:"not null"`
	MinKeyLevel int       `json:"min_key_level" gorm:"not null"`
	RoleRank    int       `json:"role_rank" gorm:"not null"`
	OverallRank int       `json:"overall_rank" gorm:"not null"`
}

// TableName specifies the table name for GORM
func (DailySpecMetricMythicPlus) TableName() string {
	return "daily_spec_metrics_mythic_plus"
}

// SpecEvolution represents the evolution of metrics between two dates
// Uses pointers to handle NULL values during initial executions
type SpecEvolutionMythicPlus struct {
	EndDate            time.Time  `json:"end_date"`
	StartDate          *time.Time `json:"start_date"`
	Spec               string     `json:"spec"`
	Class              string     `json:"class"`
	Role               string     `json:"role"`
	EncounterID        int        `json:"encounter_id"`
	IsGlobal           bool       `json:"is_global"`
	CurrentAvgScore    float64    `json:"current_avg_score"`
	PrevAvgScore       *float64   `json:"prev_avg_score"`
	ScoreChange        *float64   `json:"score_change"`
	ScoreChangePercent *float64   `json:"score_change_percent"`
	CurrentMaxScore    float64    `json:"current_max_score"`
	PrevMaxScore       *float64   `json:"prev_max_score"`
	MaxScoreChange     *float64   `json:"max_score_change"`
	CurrentAvgKeyLevel float64    `json:"current_avg_key"`
	PrevAvgKeyLevel    *float64   `json:"prev_avg_key"`
	KeyLevelChange     *float64   `json:"key_level_change"`
	CurrentOverallRank int        `json:"current_overall_rank"`
	PrevOverallRank    *int       `json:"prev_overall_rank"`
	OverallRankChange  *int       `json:"overall_rank_change"`
	CurrentRoleRank    int        `json:"current_role_rank"`
	PrevRoleRank       *int       `json:"prev_role_rank"`
	RoleRankChange     *int       `json:"role_rank_change"`
}
