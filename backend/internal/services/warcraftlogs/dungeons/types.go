package warcraftlogs

import (
	"time"
	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"

	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

// Constants for rate limiting
const (
	requestsPerSecond = 2
	burstLimit        = 5
	requestTimeout    = 10 * time.Second
)

// Constants for dungeons
const (
	DungeonAraKara       = 12660
	DungeonCityOfThreads = 12669
	DungeonGrimBatol     = 60670
	DungeonMists         = 62290
	DungeonSiege         = 61822
	DungeonDawnbreaker   = 12662
	DungeonNecroticWake  = 62286
	DungeonStonevault    = 12652
)

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

// RankingsService structure
type RankingsService struct {
	dungeonService *DungeonService
	db             *gorm.DB
	rateLimiter    *rate.Limiter
}

// NewRankingsService creates a new rankings service with a rate limiter initialized
func NewRankingsService(dungeonService *DungeonService, db *gorm.DB) *RankingsService {
	return &RankingsService{
		dungeonService: dungeonService,
		db:             db,
		rateLimiter:    rate.NewLimiter(rate.Limit(requestsPerSecond), burstLimit),
	}
}

// Temporary data structure
type playerData struct {
	ranking   playerRankingModels.PlayerRanking
	dungeonID int
}

// RankingsUpdater manages periodic updates of rankings
type RankingsUpdater struct {
	db             *gorm.DB
	rankingService *RankingsService
}

// NewRankingsUpdater creates a new rankings updater
func NewRankingsUpdater(db *gorm.DB, rankingService *RankingsService) *RankingsUpdater {
	return &RankingsUpdater{
		db:             db,
		rankingService: rankingService,
	}
}
