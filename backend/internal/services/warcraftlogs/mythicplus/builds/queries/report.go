package warcraftlogsBuildsQueries

import (
	"encoding/json"
	"fmt"
	"log"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

const ReportDetailsQuery = `
query getReportDetails($code: String!, $fightID: Int!, $encounterID: Int!) {
    reportData {
        report(code: $code) {
            table(
                fightIDs: [$fightID]
                encounterID: $encounterID
            )
            fights(fightIDs: [$fightID]) {
                id
                encounterID
                friendlyPlayers
                keystoneTime
                keystoneLevel
                keystoneAffixes
            }
        }
    }
}
`

func ParseReportDetailsResponse(response []byte, code string, fightID int, encounterID uint) (*warcraftlogsBuilds.Report, error) {
	log.Printf("Parsing report details data for report %s, fight %d, encounter %d", code, fightID, encounterID)

	var result struct {
		Data struct {
			ReportData struct {
				Report struct {
					Table struct {
						Data struct {
							TotalTime     int64             `json:"totalTime"`
							ItemLevel     float64           `json:"itemLevel"`
							Composition   []json.RawMessage `json:"composition"`
							DamageDone    []json.RawMessage `json:"damageDone"`
							HealingDone   []json.RawMessage `json:"healingDone"`
							DamageTaken   []json.RawMessage `json:"damageTaken"`
							DeathEvents   []json.RawMessage `json:"deathEvents"`
							PlayerDetails struct {
								Dps     []json.RawMessage `json:"dps"`
								Healers []json.RawMessage `json:"healers"`
								Tanks   []json.RawMessage `json:"tanks"`
							} `json:"playerDetails"`
							LogVersion  int `json:"logVersion"`
							GameVersion int `json:"gameVersion"`
						} `json:"data"`
					} `json:"table"`
					Fights []struct {
						ID              int   `json:"id"`
						KeystoneTime    int64 `json:"keystoneTime"`
						KeystoneLevel   int   `json:"keystoneLevel"`
						KeystoneAffixes []int `json:"keystoneAffixes"`
						FriendlyPlayers []int `json:"friendlyPlayers"`
					} `json:"fights"`
				} `json:"report"`
			} `json:"reportData"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal report details response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	reportData := result.Data.ReportData.Report
	tableData := reportData.Table.Data
	fight := reportData.Fights[0]

	// Building the report
	report := &warcraftlogsBuilds.Report{
		Code:          code,
		FightID:       fightID,
		EncounterID:   encounterID,
		TotalTime:     tableData.TotalTime,
		ItemLevel:     tableData.ItemLevel,
		KeystoneLevel: fight.KeystoneLevel,
		KeystoneTime:  fight.KeystoneTime,
		Affixes:       fight.KeystoneAffixes,
		LogVersion:    tableData.LogVersion,
		GameVersion:   tableData.GameVersion,
	}

	// Convert the JSON data
	var err error
	if report.Composition, err = json.Marshal(tableData.Composition); err != nil {
		return nil, fmt.Errorf("failed to marshal composition data: %w", err)
	}
	if report.DamageDone, err = json.Marshal(tableData.DamageDone); err != nil {
		return nil, fmt.Errorf("failed to marshal damage done data: %w", err)
	}
	if report.HealingDone, err = json.Marshal(tableData.HealingDone); err != nil {
		return nil, fmt.Errorf("failed to marshal healing done data: %w", err)
	}
	if report.DamageTaken, err = json.Marshal(tableData.DamageTaken); err != nil {
		return nil, fmt.Errorf("failed to marshal damage taken data: %w", err)
	}
	if report.DeathEvents, err = json.Marshal(tableData.DeathEvents); err != nil {
		return nil, fmt.Errorf("failed to marshal death events data: %w", err)
	}
	if report.PlayerDetailsDps, err = json.Marshal(tableData.PlayerDetails.Dps); err != nil {
		return nil, fmt.Errorf("failed to marshal player details dps data: %w", err)
	}
	if report.PlayerDetailsHealers, err = json.Marshal(tableData.PlayerDetails.Healers); err != nil {
		return nil, fmt.Errorf("failed to marshal player details healers data: %w", err)
	}
	if report.PlayerDetailsTanks, err = json.Marshal(tableData.PlayerDetails.Tanks); err != nil {
		return nil, fmt.Errorf("failed to marshal player details tanks data: %w", err)
	}

	report.FriendlyPlayers = fight.FriendlyPlayers

	return report, nil
}
