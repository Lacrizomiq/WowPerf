package logs

type Server struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type Report struct {
	Code      string `json:"code"`
	FightID   int    `json:"fightID"`
	StartTime int64  `json:"startTime"`
}

type Ranking struct {
	Name          string  `json:"name"`
	Class         string  `json:"class"`
	Spec          string  `json:"spec"`
	Amount        float64 `json:"amount"`
	HardModeLevel int     `json:"hardModeLevel"`
	Duration      int64   `json:"duration"`
	StartTime     int64   `json:"startTime"`
	Report        Report  `json:"report"`
	Server        Server  `json:"server"`
	BracketData   int     `json:"bracketData"`
	Faction       int     `json:"faction"`
	Affixes       []int   `json:"affixes"`
	Medal         string  `json:"medal"`
	Score         float64 `json:"score"`
	Leaderboard   int     `json:"leaderboard"`
}

type DungeonLogs struct {
	Page         int       `json:"page"`
	HasMorePages bool      `json:"hasMorePages"`
	Count        int       `json:"count"`
	Rankings     []Ranking `json:"rankings"`
}
