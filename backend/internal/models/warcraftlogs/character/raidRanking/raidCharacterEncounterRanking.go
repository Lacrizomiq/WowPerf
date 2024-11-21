// package models/warcraftlogs/character/raidRanking/raidCharacterEncounterRanking.go
package characterRaidRankingByEncounter

import "time"

type WarcraftLogsResponse struct {
	Data struct {
		CharacterData CharacterData `json:"characterData"`
	} `json:"data"`
}

type CharacterData struct {
	Character Character `json:"character"`
}

type Character struct {
	Name              string            `json:"name"`
	ClassID           int               `json:"classID"`
	ID                int               `json:"id"`
	EncounterRankings EncounterRankings `json:"encounterRankings"`
}

type EncounterRankings struct {
	BestAmount         float64 `json:"bestAmount"`
	MedianPerformance  float64 `json:"medianPerformance"`
	AveragePerformance float64 `json:"averagePerformance"`
	TotalKills         int     `json:"totalKills"`
	FastestKill        int     `json:"fastestKill"`
	Difficulty         int     `json:"difficulty"`
	Metric             string  `json:"metric"`
	Partition          int     `json:"partition"`
	Zone               int     `json:"zone"`
	Ranks              []Rank  `json:"ranks"`
}

type Rank struct {
	LockedIn              bool      `json:"lockedIn"`
	RankPercent           float64   `json:"rankPercent"`
	HistoricalPercent     float64   `json:"historicalPercent"`
	TodayPercent          float64   `json:"todayPercent"`
	RankTotalParses       int       `json:"rankTotalParses"`
	HistoricalTotalParses int       `json:"historicalTotalParses"`
	TodayTotalParses      int       `json:"todayTotalParses"`
	Guild                 Guild     `json:"guild"`
	Report                Report    `json:"report"`
	Duration              int       `json:"duration"`
	StartTime             time.Time `json:"startTime"`
	Amount                float64   `json:"amount"`
	BracketData           int       `json:"bracketData"`
	Spec                  string    `json:"spec"`
	BestSpec              string    `json:"bestSpec"`
	Class                 int       `json:"class"`
	Faction               int       `json:"faction"`
	Talents               Talents   `json:"talents"`
	Gear                  []Item    `json:"gear"`
	AzeritePowers         []string  `json:"azeritePowers"` // Peut être modifié selon le type réel
}

type Guild struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Faction int    `json:"faction"`
}

type Report struct {
	Code      string    `json:"code"`
	StartTime time.Time `json:"startTime"`
	FightID   int       `json:"fightID"`
}

type Talents struct {
	Class map[string][]TalentNode `json:"class"`
	Spec  map[string][]TalentNode `json:"spec"`
}

type TalentNode struct {
	Node struct {
		NodeID            int       `json:"nodeId"`
		Name              string    `json:"name"`
		Type              string    `json:"type"`
		PosX              int       `json:"posX"`
		PosY              int       `json:"posY"`
		RequiredPoints    int       `json:"requiredPoints"`
		MaxRanks          int       `json:"maxRanks"`
		ChildNodes        []int     `json:"childNodes"`
		ParentNodes       []int     `json:"parentNodes"`
		IsCapstone        bool      `json:"isCapstone"`
		Row               int       `json:"row"`
		Column            int       `json:"column"`
		IsPreselectedNode bool      `json:"isPreselectedNode"`
		TreeType          string    `json:"treeType"`
		Abilities         []Ability `json:"abilities"`
		SubTreeId         *int      `json:"subTreeId"`
	} `json:"node"`
	SelectedEntryId int `json:"selectedEntryId"`
	PointsInvested  int `json:"pointsInvested"`
}

type Ability struct {
	ID           int     `json:"id"`
	DefinitionId int     `json:"definitionId"`
	MaxRanks     int     `json:"maxRanks"`
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	SpellId      int     `json:"spellId"`
	Icon         string  `json:"icon"`
	Index        int     `json:"index"`
	SpellSchool  *string `json:"spellSchool"`
	SubTreeId    *int    `json:"subTreeId"`
}

type Item struct {
	Name             string   `json:"name"`
	Quality          string   `json:"quality"`
	ID               int      `json:"id"`
	Icon             string   `json:"icon"`
	ItemLevel        string   `json:"itemLevel"`
	PermanentEnchant string   `json:"permanentEnchant,omitempty"`
	TemporaryEnchant string   `json:"temporaryEnchant,omitempty"`
	BonusIDs         []string `json:"bonusIDs"`
	Gems             []Gem    `json:"gems,omitempty"`
}

type Gem struct {
	ID        string `json:"id"`
	ItemLevel string `json:"itemLevel"`
}
