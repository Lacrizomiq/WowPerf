// characterMythicPlusBestRuns.go

package raiderioMythicPlus

import (
	"fmt"
	"strings"
	"wowperf/internal/services/raiderio"
)

// RunInfo contains the essential information of a run
type RunInfo struct {
	DungeonName     string  `json:"dungeon"`
	ChallengeModeID int     `json:"challenge_mode_id"`
	RunID           string  `json:"run_id"`
	Level           int     `json:"mythic_level"`
	Score           float64 `json:"score"`
	URL             string  `json:"url"`
}

// GetCharacterMythicPlusBestRuns returns the best runs for a character with extracted run IDs
func GetCharacterMythicPlusBestRuns(s *raiderio.RaiderIOService, region, realm, name, fields string) ([]RunInfo, error) {
	params := map[string]string{
		"region": region,
		"realm":  realm,
		"name":   name,
		"fields": fields,
	}

	endpoint := "/characters/profile"
	data, err := s.Client.Get(endpoint, params)
	if err != nil {
		return nil, err
	}

	// Extract mythic_plus_best_runs from the result
	bestRuns, ok := data["mythic_plus_best_runs"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse mythic_plus_best_runs")
	}

	var runInfos []RunInfo
	for _, run := range bestRuns {
		runMap, ok := run.(map[string]interface{})
		if !ok {
			continue
		}

		url, ok := runMap["url"].(string)
		if !ok {
			continue
		}

		// Extract the run ID from the URL
		// Format: https://raider.io/mythic-plus-runs/season-tww-1/15736236-18-arakara-city-of-echoes
		// The ID is 15736236
		parts := strings.Split(url, "/")
		if len(parts) < 5 {
			continue
		}

		lastPart := parts[len(parts)-1]
		runIDParts := strings.Split(lastPart, "-")
		if len(runIDParts) < 2 {
			continue
		}

		runInfo := RunInfo{
			DungeonName: runMap["dungeon"].(string),
			RunID:       runIDParts[0],
			Level:       int(runMap["mythic_level"].(float64)),
			Score:       runMap["score"].(float64),
			URL:         url,
		}

		runInfos = append(runInfos, runInfo)
	}

	return runInfos, nil
}

// ExtractRunID is a utility function that can be used separately
func ExtractRunID(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}

	lastPart := parts[len(parts)-1]
	runIDParts := strings.Split(lastPart, "-")
	if len(runIDParts) < 2 {
		return ""
	}

	return runIDParts[0]
}
