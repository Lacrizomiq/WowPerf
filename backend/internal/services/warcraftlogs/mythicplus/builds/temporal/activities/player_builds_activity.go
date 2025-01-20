// Package warcraftlogsBuildsTemporalActivities handles Temporal activities for WarcraftLogs builds processing
package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/activity"
	"gorm.io/datatypes"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
)

const (
	// ProcessingBatchSize defines the number of builds to process in a single batch
	ProcessingBatchSize = 5
	// DelayBetweenBatches defines the delay between processing batches
	DelayBetweenBatches = 100 * time.Millisecond
)

// PlayerBuildsActivity handles all player build related operations
type PlayerBuildsActivity struct {
	repository *playerBuildsRepository.PlayerBuildsRepository
}

// NewPlayerBuildsActivity creates a new instance of PlayerBuildsActivity
func NewPlayerBuildsActivity(repository *playerBuildsRepository.PlayerBuildsRepository) *PlayerBuildsActivity {
	return &PlayerBuildsActivity{
		repository: repository,
	}
}

// PlayerDetails represents the detailed information about a player from WarcraftLogs
type PlayerDetails struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	Specs         []string        `json:"specs"`
	MaxItemLevel  float64         `json:"maxItemLevel"`
	CombatantInfo json.RawMessage `json:"combatantInfo"`
}

// ProcessBuilds processes builds from a batch of reports and stores them in the database
func (a *PlayerBuildsActivity) ProcessBuilds(ctx context.Context, reports []*warcraftlogsBuilds.Report) (*workflows.BuildsProcessingResult, error) {
	logger := activity.GetLogger(ctx)
	result := &workflows.BuildsProcessingResult{
		ProcessedAt: time.Now(),
	}

	if len(reports) == 0 {
		logger.Info("No reports to process")
		return result, nil
	}

	logger.Info("Starting build processing", "reportsCount", len(reports))

	// Process each report individually
	for i := 0; i < len(reports); i++ {
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing report %d of %d", i+1, len(reports)))

		report := reports[i]
		builds, err := a.extractPlayerBuilds(report)
		if err != nil {
			logger.Error("Failed to extract builds from report",
				"reportCode", report.Code,
				"error", err)
			result.FailureCount++
			continue
		}

		// Process builds in smaller batches
		for j := 0; j < len(builds); j += ProcessingBatchSize {
			end := j + ProcessingBatchSize
			if end > len(builds) {
				end = len(builds)
			}
			buildsBatch := builds[j:end]

			err := a.repository.StoreManyPlayerBuilds(ctx, buildsBatch)
			if err != nil {
				logger.Error("Failed to store builds batch",
					"reportCode", report.Code,
					"batchSize", len(buildsBatch),
					"error", err)
				result.FailureCount++
				continue
			}

			result.SuccessCount++
			result.ProcessedBuilds += len(buildsBatch)

			// Add delay between batches to prevent overload
			time.Sleep(DelayBetweenBatches)
		}

		logger.Info("Processed report builds",
			"reportCode", report.Code,
			"buildsProcessed", len(builds))
	}

	return result, nil
}

// extractPlayerBuilds extracts all player builds from a report
func (a *PlayerBuildsActivity) extractPlayerBuilds(report *warcraftlogsBuilds.Report) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	var builds []*warcraftlogsBuilds.PlayerBuild

	// Extract builds from each role (DPS, Healers, Tanks)
	dpsBuilds, err := a.extractBuildsFromPlayers(report, report.PlayerDetailsDps)
	if err != nil {
		log.Printf("Error extracting DPS builds: %v", err)
	}
	builds = append(builds, dpsBuilds...)

	healerBuilds, err := a.extractBuildsFromPlayers(report, report.PlayerDetailsHealers)
	if err != nil {
		log.Printf("Error extracting healer builds: %v", err)
	}
	builds = append(builds, healerBuilds...)

	tankBuilds, err := a.extractBuildsFromPlayers(report, report.PlayerDetailsTanks)
	if err != nil {
		log.Printf("Error extracting tank builds: %v", err)
	}
	builds = append(builds, tankBuilds...)

	return builds, nil
}

// extractBuildsFromPlayers processes player data and creates builds
func (a *PlayerBuildsActivity) extractBuildsFromPlayers(report *warcraftlogsBuilds.Report, playerData []byte) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	if len(playerData) == 0 {
		return nil, nil
	}

	var players []PlayerDetails
	if err := json.Unmarshal(playerData, &players); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player details: %w", err)
	}

	var builds []*warcraftlogsBuilds.PlayerBuild
	for _, player := range players {
		build, err := a.createPlayerBuild(report, player)
		if err != nil {
			log.Printf("Error creating build for player %s: %v", player.Name, err)
			continue
		}
		if build != nil {
			builds = append(builds, build)
		}
	}

	return builds, nil
}

// createPlayerBuild creates a PlayerBuild from player details
func (a *PlayerBuildsActivity) createPlayerBuild(report *warcraftlogsBuilds.Report, player PlayerDetails) (*warcraftlogsBuilds.PlayerBuild, error) {
	var combatInfo struct {
		Stats      json.RawMessage `json:"stats"`
		Gear       json.RawMessage `json:"gear"`
		TalentTree json.RawMessage `json:"talentTree"`
	}

	if err := json.Unmarshal(player.CombatantInfo, &combatInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal combatant info: %w", err)
	}

	specName := ""
	if len(player.Specs) > 0 {
		specName = player.Specs[0]
	}

	// Extract talent code from the report's talent codes
	var talentCodes map[string]string
	if err := json.Unmarshal(report.TalentCodes, &talentCodes); err != nil {
		log.Printf("Error unmarshaling talent codes: %v", err)
	}
	talentCode := talentCodes[fmt.Sprintf("%s_%s_talents", player.Type, specName)]

	build := &warcraftlogsBuilds.PlayerBuild{
		PlayerName:    player.Name,
		Class:         player.Type,
		Spec:          specName,
		ReportCode:    report.Code,
		FightID:       report.FightID,
		ActorID:       player.ID,
		ItemLevel:     player.MaxItemLevel,
		TalentImport:  talentCode,
		TalentTree:    datatypes.JSON(combatInfo.TalentTree),
		Gear:          datatypes.JSON(combatInfo.Gear),
		Stats:         datatypes.JSON(combatInfo.Stats),
		EncounterID:   report.EncounterID,
		KeystoneLevel: report.KeystoneLevel,
		Affixes:       report.Affixes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return build, nil
}

// CountPlayerBuilds returns the total count of player builds in the database
func (a *PlayerBuildsActivity) CountPlayerBuilds(ctx context.Context) (int64, error) {
	logger := activity.GetLogger(ctx)

	count, err := a.repository.CountPlayerBuilds(ctx)
	if err != nil {
		logger.Error("Failed to count player builds", "error", err)
		return 0, fmt.Errorf("failed to count player builds: %w", err)
	}

	logger.Info("Counted player builds", "count", count)
	return count, nil
}
