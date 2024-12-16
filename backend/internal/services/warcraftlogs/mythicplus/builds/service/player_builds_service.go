package warcraftlogsBuildsService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/datatypes"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	warcraftlogsBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
)

// PlayerBuildsService handles the business logic for player builds
type PlayerBuildsService struct {
	repository *warcraftlogsBuildsRepository.PlayerBuildsRepository
}

// NewPlayerBuildsService creates a new instance of PlayerBuildsService
func NewPlayerBuildsService(repository *warcraftlogsBuildsRepository.PlayerBuildsRepository) *PlayerBuildsService {
	return &PlayerBuildsService{repository: repository}
}

// PlayerDetails represents the structure of player details in the report
type PlayerDetails struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Specs []struct {
		Spec string `json:"spec"`
	} `json:"specs"`
	Server        string          `json:"server"`
	MaxItemLevel  float64         `json:"maxItemLevel"`
	CombatantInfo json.RawMessage `json:"combatantInfo"`
}

// ProcessReportBuilds process all players from a report and stores them in the database
// It handles DPS, Healers and Tanks separately but stores them all in the same table
func (s *PlayerBuildsService) ProcessReportBuilds(ctx context.Context, report *warcraftlogsBuilds.Report) error {
	var playerBuilds []*warcraftlogsBuilds.PlayerBuild

	// Extract players builds from each role category
	var dpsPlayers []json.RawMessage
	var healersPlayers []json.RawMessage
	var tanksPlayers []json.RawMessage

	// Unmarshal the players builds
	if err := json.Unmarshal(report.PlayerDetailsDps, &dpsPlayers); err != nil {
		return fmt.Errorf("failed to unmarshal dps players: %w", err)
	}
	if err := json.Unmarshal(report.PlayerDetailsHealers, &healersPlayers); err != nil {
		return fmt.Errorf("failed to unmarshal healers players: %w", err)
	}
	if err := json.Unmarshal(report.PlayerDetailsTanks, &tanksPlayers); err != nil {
		return fmt.Errorf("failed to unmarshal tanks players: %w", err)
	}

	// Process each role category
	dpsBuilds, err := s.extractPlayerBuilds(report, dpsPlayers)
	if err != nil {
		return fmt.Errorf("failed to extract dps builds: %w", err)
	}
	playerBuilds = append(playerBuilds, dpsBuilds...)

	healersBuilds, err := s.extractPlayerBuilds(report, healersPlayers)
	if err != nil {
		return fmt.Errorf("failed to extract healers builds: %w", err)
	}
	playerBuilds = append(playerBuilds, healersBuilds...)

	tanksBuilds, err := s.extractPlayerBuilds(report, tanksPlayers)
	if err != nil {
		return fmt.Errorf("failed to extract tanks builds: %w", err)
	}
	playerBuilds = append(playerBuilds, tanksBuilds...)

	// Store all builds in the database
	return s.repository.StoreManyPlayerBuilds(ctx, playerBuilds)
}

// extractPlayerBuilds processes a list of players and creates PlayerBuild objects
func (s *PlayerBuildsService) extractPlayerBuilds(report *warcraftlogsBuilds.Report, playersData []json.RawMessage) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	var playerBuilds []*warcraftlogsBuilds.PlayerBuild

	for _, playerData := range playersData {
		build, err := s.extractPlayerDetails(report, playerData)
		if err != nil {
			log.Printf("failed to extract player details: %v", err)
			continue // Skip this player if there's an error but continue processing other players
		}
		playerBuilds = append(playerBuilds, build)
	}
	return playerBuilds, nil
}

// extractPlayerDetails extracts the details of a player from the report
func (s *PlayerBuildsService) extractPlayerDetails(report *warcraftlogsBuilds.Report, playerData json.RawMessage) (*warcraftlogsBuilds.PlayerBuild, error) {
	var player struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Specs []struct {
			Spec string `json:"spec"`
		} `json:"specs"`
		Server        string          `json:"server"`
		MaxItemLevel  float64         `json:"maxItemLevel"`
		CombatantInfo json.RawMessage `json:"combatantInfo"`
	}

	if err := json.Unmarshal(playerData, &player); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player data: %w", err)
	}

	// Create basic player build struct
	playerBuild := &warcraftlogsBuilds.PlayerBuild{
		PlayerName:    player.Name,
		Class:         player.Type,
		Spec:          player.Specs[0].Spec, // Using first spec
		ReportCode:    report.Code,
		FightID:       report.FightID,
		ActorID:       player.ID,
		ItemLevel:     player.MaxItemLevel,
		KeystoneLevel: report.KeystoneLevel,
		EncounterID:   report.EncounterID,
	}

	// Extract talent code
	playerBuild.TalentCode = s.getTalentCode(report, player)

	// Convert CombatantInfo to JSONB
	playerBuild.CombatantInfo = datatypes.JSON(player.CombatantInfo)

	// Extract gear and stats from CombatantInfo
	if err := s.extractGearAndStats(playerBuild, player.CombatantInfo); err != nil {
		return nil, fmt.Errorf("failed to extract gear and stats: %w", err)
	}

	// Set affixes from report
	playerBuild.Affixes = report.Affixes

	return playerBuild, nil
}

// getTalentCode extracts the talent code for a specific player
func (s *PlayerBuildsService) getTalentCode(report *warcraftlogsBuilds.Report, player PlayerDetails) string {
	var talentCodes map[string]string

	if err := json.Unmarshal(report.TalentCodes, &talentCodes); err != nil {
		log.Printf("failed to unmarshal talent codes: %v", err)
		return ""
	}

	// Construct the talent code key (e.g, "Priest_Discipline_talents")
	key := fmt.Sprintf("%s_%s_talents", player.Type, player.Specs[0].Spec)

	return talentCodes[key]
}

// extractGearAndStats extracts the gear and stats from the CombatantInfo
func (s *PlayerBuildsService) extractGearAndStats(playerBuild *warcraftlogsBuilds.PlayerBuild, combatantInfo json.RawMessage) error {
	var info struct {
		Gear  json.RawMessage `json:"gear"`
		Stats json.RawMessage `json:"stats"`
	}

	if err := json.Unmarshal(combatantInfo, &info); err != nil {
		return fmt.Errorf("failed to unmarshal combatant info: %w", err)
	}

	playerBuild.Gear = datatypes.JSON(info.Gear)
	playerBuild.Stats = datatypes.JSON(info.Stats)

	return nil
}
