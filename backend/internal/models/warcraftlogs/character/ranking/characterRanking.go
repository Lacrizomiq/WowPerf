// package models/warcraftlogs/character/ranking/characterRanking.go
package characterRanking

// WarcraftLogsResponse is the response from the WarcraftLogs API
type WarcraftLogsResponse struct {
	Data struct {
		CharacterData struct {
			Character CharacterData `json:"character"`
		} `json:"characterData"`
	} `json:"data"`
}

// CharacterData is the data for a character
type CharacterData struct {
	Name         string       `json:"name"`
	ClassID      int          `json:"classID"`
	ZoneRankings ZoneRankings `json:"zoneRankings"`
}

// ZoneRankings is the rankings for a zone
type ZoneRankings struct {
	BestPerformanceAverage   float64   `json:"bestPerformanceAverage"`
	MedianPerformanceAverage float64   `json:"medianPerformanceAverage"`
	Difficulty               int       `json:"difficulty"`
	Metric                   string    `json:"metric"`
	Partition                int       `json:"partition"`
	Zone                     int       `json:"zone"`
	AllStars                 []AllStar `json:"allStars"`
	Rankings                 []Ranking `json:"rankings"`
}

// AllStar is the all stars of a character
type AllStar struct {
	Partition      int     `json:"partition"`
	Spec           string  `json:"spec"`
	Points         float64 `json:"points"`
	PossiblePoints int     `json:"possiblePoints"`
	Rank           int     `json:"rank"`
	RegionRank     int     `json:"regionRank"`
	ServerRank     int     `json:"serverRank"`
	RankPercent    float64 `json:"rankPercent"`
	Total          int     `json:"total"`
}

// Ranking is the ranking of a character
type Ranking struct {
	Encounter struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"encounter"`
	RankPercent   float64 `json:"rankPercent"`
	MedianPercent float64 `json:"medianPercent"`
	LockedIn      bool    `json:"lockedIn"`
	TotalKills    int     `json:"totalKills"`
	FastestKill   int     `json:"fastestKill"`
	AllStars      AllStar `json:"allStars"`
	Spec          string  `json:"spec"`
	BestSpec      string  `json:"bestSpec"`
	BestAmount    float64 `json:"bestAmount"`
}
