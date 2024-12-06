package warcraftlogsBuildsQueries

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

// Struct for the player
type PlayerSpec struct {
	Name  string   `json:"name"`
	ID    int      `json:"id"`
	Type  string   `json:"type"` // class
	Icon  string   `json:"icon"`
	Specs []string `json:"spec"`
}

const GetReportTableQuery = `
query getReportTableQuery($code: String!, $fightID: Int!, $encounterID: Int!) {
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

// buildTalentsQuery to build the talents query
func buildTalentsQuery(code string, fightID int, players []PlayerSpec) string {
	var talentFields []string
	for _, player := range players {
		if len(player.Specs) > 0 {
			alias := fmt.Sprintf("%s_%s_talents: talentImportCode(actorID: %d)",
				player.Type, player.Specs[0], player.ID)
			talentFields = append(talentFields, alias)
		}
	}

	return fmt.Sprintf(`
	query {
			reportData {
					report(code: "%s") {
							fights(fightIDs: [%d]) {
									%s
							}
					}
			}
	}`, code, fightID, strings.Join(talentFields, "\n"))
}

// ParseReportDetailsResponse to parse the report details response
func ParseReportDetailsResponse(response []byte, code string, fightID int, encounterID uint) (*warcraftlogsBuilds.Report, string, error) {
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
								Dps     []PlayerSpec `json:"dps"`
								Healers []PlayerSpec `json:"healers"`
								Tanks   []PlayerSpec `json:"tanks"`
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
		return nil, "", fmt.Errorf("failed to unmarshal report details response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, "", fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	// Collect all players
	var allPlayers []PlayerSpec
	reportData := result.Data.ReportData.Report
	tableData := reportData.Table.Data
	allPlayers = append(allPlayers, tableData.PlayerDetails.Dps...)
	allPlayers = append(allPlayers, tableData.PlayerDetails.Healers...)
	allPlayers = append(allPlayers, tableData.PlayerDetails.Tanks...)

	// build the request for the talents
	talentsQuery := buildTalentsQuery(code, fightID, allPlayers)
	log.Printf("Build talents query: %s", talentsQuery)

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
		return nil, "", fmt.Errorf("failed to marshal composition data: %w", err)
	}
	if report.DamageDone, err = json.Marshal(tableData.DamageDone); err != nil {
		return nil, "", fmt.Errorf("failed to marshal damage done data: %w", err)
	}
	if report.HealingDone, err = json.Marshal(tableData.HealingDone); err != nil {
		return nil, "", fmt.Errorf("failed to marshal healing done data: %w", err)
	}
	if report.DamageTaken, err = json.Marshal(tableData.DamageTaken); err != nil {
		return nil, "", fmt.Errorf("failed to marshal damage taken data: %w", err)
	}
	if report.DeathEvents, err = json.Marshal(tableData.DeathEvents); err != nil {
		return nil, "", fmt.Errorf("failed to marshal death events data: %w", err)
	}
	if report.PlayerDetailsDps, err = json.Marshal(tableData.PlayerDetails.Dps); err != nil {
		return nil, "", fmt.Errorf("failed to marshal player details dps data: %w", err)
	}
	if report.PlayerDetailsHealers, err = json.Marshal(tableData.PlayerDetails.Healers); err != nil {
		return nil, "", fmt.Errorf("failed to marshal player details healers data: %w", err)
	}
	if report.PlayerDetailsTanks, err = json.Marshal(tableData.PlayerDetails.Tanks); err != nil {
		return nil, "", fmt.Errorf("failed to marshal player details tanks data: %w", err)
	}

	report.FriendlyPlayers = fight.FriendlyPlayers

	return report, talentsQuery, nil
}

// ParseReportTalentsResponse to parse the report talents response
func ParseReportTalentsResponse(response []byte) (map[string]string, error) {
	log.Printf("Parsing report talents data")

	var result struct {
		Data struct {
			ReportData struct {
				Report struct {
					Fights []map[string]interface{} `json:"fights"`
				} `json:"report"`
			} `json:"reportData"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal report talents response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if len(result.Data.ReportData.Report.Fights) == 0 {
		return nil, fmt.Errorf("no fights found in the report")
	}

	// Extract talents for the first fight
	fightData := result.Data.ReportData.Report.Fights[0]
	talentCodes := make(map[string]string)

	for key, value := range fightData {
		if strings.HasSuffix(key, "_talents") {
			if talentCode, ok := value.(string); ok {
				talentCodes[key] = talentCode
			}
		}
	}

	return talentCodes, nil
}
