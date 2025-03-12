package warcraftlogs

import (
	"time"
	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"
)

// Constants for rate limiting
const (
	requestsPerSecond = 2
	burstLimit        = 5
	requestTimeout    = 10 * time.Second
)

// Constants for dungeons
const (
	// Season 2 dungeons
	DungeonCinderbrew  = 12661  // Cinderbrew Meadery
	DungeonDarkflame   = 12651  // Darkflame Cleft
	DungeonFloodgate   = 12773  // Operation: Floodgate
	DungeonMechagon    = 112098 // Operation: Mechagon - Workshop
	DungeonPriory      = 12649  // Priory of the Sacred Flame
	DungeonMotherlode  = 61594  // The MOTHERLODE!!
	DungeonRookery     = 12648  // The Rookery
	DungeonTheaterPain = 62293  // Theater of Pain

	// Season 1 dungeons
	DungeonAraKara       = 12660
	DungeonCityOfThreads = 12669
	DungeonGrimBatol     = 60670
	DungeonMists         = 62290
	DungeonSiege         = 61822
	DungeonDawnbreaker   = 12662
	DungeonNecroticWake  = 62286
	DungeonStonevault    = 12652
)

// Specialization is a struct that represents a class specialization
type Specialization struct {
	ClassName string
	SpecName  string
}

// List of all specializations mapped to their classes
var Specializations = []Specialization{
	// Priest
	{"Priest", "Discipline"},
	{"Priest", "Holy"},
	{"Priest", "Shadow"},

	// Death Knight
	{"DeathKnight", "Blood"},
	{"DeathKnight", "Frost"},
	{"DeathKnight", "Unholy"},

	// Druid
	{"Druid", "Balance"},
	{"Druid", "Feral"},
	{"Druid", "Guardian"},
	{"Druid", "Restoration"},

	// Hunter
	{"Hunter", "BeastMastery"},
	{"Hunter", "Marksmanship"},
	{"Hunter", "Survival"},

	// Mage
	{"Mage", "Arcane"},
	{"Mage", "Fire"},
	{"Mage", "Frost"},

	// Monk
	{"Monk", "Brewmaster"},
	{"Monk", "Mistweaver"},
	{"Monk", "Windwalker"},

	// Paladin
	{"Paladin", "Holy"},
	{"Paladin", "Protection"},
	{"Paladin", "Retribution"},

	// Rogue
	{"Rogue", "Assassination"},
	{"Rogue", "Subtlety"},
	{"Rogue", "Outlaw"},

	// Shaman
	{"Shaman", "Elemental"},
	{"Shaman", "Enhancement"},
	{"Shaman", "Restoration"},

	// Warlock
	{"Warlock", "Affliction"},
	{"Warlock", "Demonology"},
	{"Warlock", "Destruction"},

	// Warrior
	{"Warrior", "Arms"},
	{"Warrior", "Fury"},
	{"Warrior", "Protection"},

	// Demon Hunter
	{"DemonHunter", "Havoc"},
	{"DemonHunter", "Vengeance"},

	// Evoker
	{"Evoker", "Devastation"},
	{"Evoker", "Preservation"},
	{"Evoker", "Augmentation"},
}

// Structures for API data
type Report struct {
	Code      string `json:"code"`
	FightID   int    `json:"fightID"`
	StartTime int64  `json:"startTime"`
}

type Guild struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Faction int    `json:"faction"`
}

type Server struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type Run struct {
	DungeonID     int     `json:"dungeonId"`
	Score         float64 `json:"score"`
	Duration      int64   `json:"duration"`
	StartTime     int64   `json:"startTime"`
	HardModeLevel int     `json:"hardModeLevel"`
	BracketData   int     `json:"bracketData"`
	Medal         string  `json:"medal"`
	Affixes       []int   `json:"affixes"`
	Report        Report  `json:"report"`
}

type PlayerScore struct {
	Name       string  `json:"name"`
	Class      string  `json:"class"`
	Spec       string  `json:"spec"`
	Role       string  `json:"role"`
	TotalScore float64 `json:"totalScore"`
	Amount     float64 `json:"amount"`
	Guild      Guild   `json:"guild"`
	Server     Server  `json:"server"`
	Faction    int     `json:"faction"`
	Runs       []Run   `json:"runs"`
}

type RoleRankings struct {
	Players []PlayerScore `json:"players"`
	Count   int           `json:"count"`
}

type GlobalRankings struct {
	Tanks   RoleRankings `json:"tanks"`
	Healers RoleRankings `json:"healers"`
	DPS     RoleRankings `json:"dps"`
}

// Temporary data structure
type playerData struct {
	ranking   playerRankingModels.PlayerRanking
	dungeonID int
}
