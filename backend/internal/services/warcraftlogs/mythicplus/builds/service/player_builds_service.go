package warcraftlogsBuildsService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/datatypes"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	warcraftlogsBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

	warcraftlogsBuildsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
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
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	Specs         []string        `json:"specs"`
	Server        string          `json:"server"`
	MaxItemLevel  float64         `json:"maxItemLevel"`
	CombatantInfo json.RawMessage `json:"combatantInfo"`
}

type CombatantInfo struct {
	Stats struct {
		Speed       warcraftlogsBuildsQueries.StatRange `json:"Speed"`
		Intellect   warcraftlogsBuildsQueries.StatRange `json:"Intellect,omitempty"`
		Strength    warcraftlogsBuildsQueries.StatRange `json:"Strength,omitempty"`
		Agility     warcraftlogsBuildsQueries.StatRange `json:"Agility,omitempty"`
		Mastery     warcraftlogsBuildsQueries.StatRange `json:"Mastery"`
		Stamina     warcraftlogsBuildsQueries.StatRange `json:"Stamina"`
		Haste       warcraftlogsBuildsQueries.StatRange `json:"Haste"`
		Leech       warcraftlogsBuildsQueries.StatRange `json:"Leech"`
		Crit        warcraftlogsBuildsQueries.StatRange `json:"Crit"`
		Versatility warcraftlogsBuildsQueries.StatRange `json:"Versatility"`
	} `json:"stats,omitempty"`
	Talents    []json.RawMessage `json:"talents"`
	TalentTree json.RawMessage   `json:"talentTree"`
	Gear       []json.RawMessage `json:"gear"`
	SpecIDs    []int             `json:"specIDs"`
}

// ProcessReportBuilds process all players from a report and stores them in the database
// It handles DPS, Healers and Tanks separately but stores them all in the same table
func (s *PlayerBuildsService) ProcessReportBuilds(ctx context.Context, reports []*warcraftlogsBuilds.Report) error {

	log.Printf("[DEBUG] Processing %d reports", len(reports))

	for _, report := range reports {
		log.Printf("[DEBUG] Processing player builds for report %s", report.Code)
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
		log.Printf("[DEBUG] Storing %d player builds for report %s", len(playerBuilds), report.Code)
		if err := s.repository.StoreManyPlayerBuilds(ctx, playerBuilds); err != nil {
			log.Printf("Error storing builds for report %s: %v", report.Code, err)
			continue
		}
	}

	return nil
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
	var player PlayerDetails

	if err := json.Unmarshal(playerData, &player); err != nil {
		log.Printf("[ERROR] Unmarshal failed, player data structure: %s", string(playerData))
		return nil, fmt.Errorf("failed to unmarshal player data: %w", err)
	}

	// Validate player details
	if err := s.validatePlayerDetails(&player); err != nil {
		return nil, fmt.Errorf("invalid player details: %w", err)
	}

	// Log player details
	log.Printf("[DEBUG] Processing player %s (%s-%s)", player.Name, player.Type, player.Specs[0])

	specName := ""
	if len(player.Specs) > 0 {
		specName = player.Specs[0]
	}

	// Create basic player build struct
	playerBuild := &warcraftlogsBuilds.PlayerBuild{
		PlayerName:    player.Name,
		Class:         player.Type,
		Spec:          specName,
		ReportCode:    report.Code,
		FightID:       report.FightID,
		ActorID:       player.ID,
		ItemLevel:     player.MaxItemLevel,
		KeystoneLevel: report.KeystoneLevel,
		EncounterID:   report.EncounterID,
	}

	// Extract talent code
	playerBuild.TalentImport = s.getTalentCode(report, player)

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
	key := fmt.Sprintf("%s_%s_talents", player.Type, player.Specs[0])

	return talentCodes[key]
}

// extractGearAndStats extracts the gear and stats from the CombatantInfo
func (s *PlayerBuildsService) extractGearAndStats(playerBuild *warcraftlogsBuilds.PlayerBuild, combatantInfo json.RawMessage) error {
	// Unmarshal the combatant info
	var rawInfo map[string]json.RawMessage
	if err := json.Unmarshal(combatantInfo, &rawInfo); err != nil {
		return fmt.Errorf("failed to unmarshal raw combatant info: %w", err)
	}

	// Extract stats and log
	if stats, ok := rawInfo["stats"]; ok {
		playerBuild.Stats = datatypes.JSON(stats)
		log.Printf("[DEBUG] Stats extracted successfully for player %s", playerBuild.PlayerName)
	}

	// Extract gear and log
	if gear, ok := rawInfo["gear"]; ok {
		playerBuild.Gear = datatypes.JSON(gear)
		log.Printf("[DEBUG] Gear extracted successfully for player %s", playerBuild.PlayerName)
	}

	// Extract talent tree and log
	if talentTree, ok := rawInfo["talentTree"]; ok {
		playerBuild.TalentTree = datatypes.JSON(talentTree)
		log.Printf("[DEBUG] TalentTree extracted successfully for player %s", playerBuild.PlayerName)
	} else {
		log.Printf("[DEBUG] No talentTree data found in combatantInfo for player %s", playerBuild.PlayerName)
	}

	return nil
}

// validatePlayerDetails validates the player details
func (s *PlayerBuildsService) validatePlayerDetails(player *PlayerDetails) error {
	if len(player.Specs) == 0 {
		return fmt.Errorf("player %s has no specs", player.Name)
	}
	if player.Name == "" {
		return fmt.Errorf("player name is empty")
	}
	if player.Type == "" {
		return fmt.Errorf("player type is empty for %s", player.Name)
	}
	return nil
}
