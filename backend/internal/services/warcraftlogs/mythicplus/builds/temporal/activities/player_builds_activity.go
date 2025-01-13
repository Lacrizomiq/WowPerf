package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/datatypes"

	"go.temporal.io/sdk/activity"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
)

type PlayerBuildsActivity struct {
	repository *playerBuildsRepository.PlayerBuildsRepository
}

func NewPlayerBuildsActivity(repository *playerBuildsRepository.PlayerBuildsRepository) *PlayerBuildsActivity {
	return &PlayerBuildsActivity{
		repository: repository,
	}
}

type PlayerDetails struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	Specs         []string        `json:"specs"`
	MaxItemLevel  float64         `json:"maxItemLevel"`
	CombatantInfo json.RawMessage `json:"combatantInfo"`
}

func (a *PlayerBuildsActivity) ProcessBuilds(ctx context.Context, reports []*warcraftlogsBuilds.Report) (*workflows.BuildsProcessingResult, error) {
	logger := activity.GetLogger(ctx)
	result := &workflows.BuildsProcessingResult{
		ProcessedAt: time.Now(),
	}

	if len(reports) == 0 {
		logger.Info("No reports to process")
		return result, nil
	}

	logger.Info("Starting player builds processing", "reportsCount", len(reports))

	// Batch processing for a better performance
	const batchSize = 10
	totalReports := len(reports)
	processedReports := 0

	for i := 0; i < totalReports; i += batchSize {
		end := i + batchSize
		if end > totalReports {
			end = totalReports
		}
		batch := reports[i:end]

		// Heartbeat for the temporal monitoring
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing builds for reports %d-%d of %d", i+1, end, totalReports))

		playerBuilds, err := a.extractPlayerBuilds(batch)
		if err != nil {
			logger.Error("Failed to extract player builds from batch", "error", err)
			result.FailureCount++
			continue
		}

		if err := a.repository.StoreManyPlayerBuilds(ctx, playerBuilds); err != nil {
			logger.Error("Failed to store player builds", "error", err)
			result.FailureCount++
			continue
		}

		result.SuccessCount++
		processedReports++
		result.ProcessedBuilds += len(playerBuilds)

		logger.Info("Batch processed",
			"processedReports", processedReports,
			"totalReports", totalReports,
			"buildsProcessed", len(playerBuilds))
	}

	return result, nil
}

func (a *PlayerBuildsActivity) extractPlayerBuilds(reports []*warcraftlogsBuilds.Report) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	var allPlayerBuilds []*warcraftlogsBuilds.PlayerBuild

	for _, report := range reports {
		// Extract DPS players
		var dpsPlayers []json.RawMessage
		if err := json.Unmarshal(report.PlayerDetailsDps, &dpsPlayers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dps players: %w", err)
		}
		dpsBuilds, err := a.extractBuildsFromPlayers(report, dpsPlayers)
		if err != nil {
			log.Printf("Error extracting DPS builds: %v", err)
			continue
		}
		allPlayerBuilds = append(allPlayerBuilds, dpsBuilds...)

		// Extract healers players
		var healersPlayers []json.RawMessage
		if err := json.Unmarshal(report.PlayerDetailsHealers, &healersPlayers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal healers players: %w", err)
		}
		healersBuilds, err := a.extractBuildsFromPlayers(report, healersPlayers)
		if err != nil {
			log.Printf("Error extracting healers builds: %v", err)
			continue
		}
		allPlayerBuilds = append(allPlayerBuilds, healersBuilds...)

		// Extract tanks players
		var tanksPlayers []json.RawMessage
		if err := json.Unmarshal(report.PlayerDetailsTanks, &tanksPlayers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tanks players: %w", err)
		}
		tanksBuilds, err := a.extractBuildsFromPlayers(report, tanksPlayers)
		if err != nil {
			log.Printf("Error extracting tanks builds: %v", err)
			continue
		}
		allPlayerBuilds = append(allPlayerBuilds, tanksBuilds...)
	}

	return allPlayerBuilds, nil
}

func (a *PlayerBuildsActivity) extractBuildsFromPlayers(report *warcraftlogsBuilds.Report, players []json.RawMessage) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	var playerBuilds []*warcraftlogsBuilds.PlayerBuild

	for _, playerData := range players {
		build, err := a.extractPlayerDetails(report, playerData)
		if err != nil {
			log.Printf("failed to extract player details: %v", err)
			continue
		}
		playerBuilds = append(playerBuilds, build)
	}
	return playerBuilds, nil
}

func (a *PlayerBuildsActivity) extractPlayerDetails(report *warcraftlogsBuilds.Report, playerData json.RawMessage) (*warcraftlogsBuilds.PlayerBuild, error) {
	var player PlayerDetails
	if err := json.Unmarshal(playerData, &player); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player data: %w", err)
	}

	specName := ""
	if len(player.Specs) > 0 {
		specName = player.Specs[0]
	}

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
		Affixes:       report.Affixes,
	}

	// Extract talent code
	playerBuild.TalentCode = a.getTalentCode(report, player)

	// Extract gear and stats
	if err := a.extractGearAndStats(playerBuild, player.CombatantInfo); err != nil {
		return nil, fmt.Errorf("failed to extract gear and stats: %w", err)
	}

	return playerBuild, nil
}

func (a *PlayerBuildsActivity) getTalentCode(report *warcraftlogsBuilds.Report, player PlayerDetails) string {
	var talentCodes map[string]string
	if err := json.Unmarshal(report.TalentCodes, &talentCodes); err != nil {
		log.Printf("failed to unmarshal talent codes: %v", err)
		return ""
	}
	// Construct the talent code key (e.g, "Priest_Discipline_talents")
	key := fmt.Sprintf("%s_%s_talents", player.Type, player.Specs[0])
	return talentCodes[key]
}

func (a *PlayerBuildsActivity) extractGearAndStats(playerBuild *warcraftlogsBuilds.PlayerBuild, combatantInfo json.RawMessage) error {
	var rawInfo map[string]json.RawMessage
	if err := json.Unmarshal(combatantInfo, &rawInfo); err != nil {
		return fmt.Errorf("failed to unmarshal raw combatant info: %w", err)
	}

	if stats, ok := rawInfo["stats"]; ok {
		playerBuild.Stats = datatypes.JSON(stats)
	}

	if gear, ok := rawInfo["gear"]; ok {
		playerBuild.Gear = datatypes.JSON(gear)
	}

	if talentTree, ok := rawInfo["talentTree"]; ok {
		playerBuild.TalentTree = datatypes.JSON(talentTree)
	}

	return nil
}
