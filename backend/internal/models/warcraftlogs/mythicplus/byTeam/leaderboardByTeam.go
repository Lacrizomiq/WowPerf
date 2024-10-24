// package models/warcraftlogs/mythicplus/ByTeam/models.go
package leaderboard

type Server struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type TeamMember struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Class string `json:"class"`
	Spec  string `json:"spec"`
	Role  string `json:"role"`
}

type Ranking struct {
	Server      Server       `json:"server"`
	Duration    int64        `json:"duration"`
	StartTime   int64        `json:"startTime"`
	Deaths      int          `json:"deaths"`
	Tanks       int          `json:"tanks"`
	Healers     int          `json:"healers"`
	Melee       int          `json:"melee"`
	Ranged      int          `json:"ranged"`
	BracketData int          `json:"bracketData"`
	Affixes     []int        `json:"affixes"`
	Team        []TeamMember `json:"team"`
	Medal       string       `json:"medal"`
	Score       float64      `json:"score"`
	Leaderboard int          `json:"leaderboard"`
}

type DungeonLeaderboard struct {
	Page         int       `json:"page"`
	HasMorePages bool      `json:"hasMorePages"`
	Count        int       `json:"count"`
	Rankings     []Ranking `json:"rankings"`
}
