package raiderioMythicPlusRunsQueries

import (
	"encoding/json"
	"fmt"
	models "wowperf/internal/models/raiderio/mythicplus_runs"
	service "wowperf/internal/services/raiderio"
)

type MythicPlusRunsParams struct {
	Season    string
	Region    string
	Dungeon   string
	Affixes   string
	Page      int
	AccessKey string
}

// GetMythicPlusRunsComplete récupère et parse les runs M+ depuis l'API
func GetMythicPlusRunsComplete(s *service.RaiderIOService, params MythicPlusRunsParams) (*models.MythicPlusRunsResponse, error) {
	// 1. Préparation des paramètres d'appel
	apiParams := buildAPIParams(params)

	// 2. Appel API (réutilise votre infrastructure existante)
	rawData, err := s.Client.Get("/mythic-plus/runs", apiParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mythic plus runs: %w", err)
	}

	// 3. Parsing et validation
	response, err := parseMythicPlusRunsResponse(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 4. Validation business
	if err := validateMythicPlusResponse(response); err != nil {
		return nil, fmt.Errorf("invalid response data: %w", err)
	}

	return response, nil
}

// buildAPIParams convertit les paramètres en format API
func buildAPIParams(params MythicPlusRunsParams) map[string]string {
	apiParams := map[string]string{
		"season":  params.Season,
		"region":  params.Region,
		"dungeon": params.Dungeon,
		"page":    fmt.Sprintf("%d", params.Page),
	}

	// Paramètres optionnels
	if params.Affixes != "" {
		apiParams["affixes"] = params.Affixes
	}

	if params.AccessKey != "" {
		apiParams["access_key"] = params.AccessKey
	}

	return apiParams
}

// parseMythicPlusRunsResponse parse la réponse JSON brute
func parseMythicPlusRunsResponse(data map[string]interface{}) (*models.MythicPlusRunsResponse, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var response models.MythicPlusRunsResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// validateMythicPlusResponse valide la cohérence des données
func validateMythicPlusResponse(response *models.MythicPlusRunsResponse) error {
	if response.Rankings == nil {
		return fmt.Errorf("no rankings found in response")
	}

	// Validation des runs individuels
	for i, ranking := range response.Rankings {
		if ranking.Run.KeystoneRunID == 0 {
			return fmt.Errorf("ranking %d: missing keystone_run_id", i)
		}

		if len(ranking.Run.Roster) != 5 {
			return fmt.Errorf("ranking %d: invalid roster size %d, expected 5", i, len(ranking.Run.Roster))
		}

		// Validation de la composition (1 tank, 1 heal, 3 dps)
		if err := validateTeamComposition(ranking.Run.Roster); err != nil {
			return fmt.Errorf("ranking %d: %w", i, err)
		}
	}

	return nil
}

// validateTeamComposition vérifie qu'on a bien 1 tank, 1 heal, 3 dps
func validateTeamComposition(roster []models.RosterMember) error {
	roleCount := map[string]int{}

	for _, member := range roster {
		roleCount[member.Role]++
	}

	if roleCount["tank"] != 1 {
		return fmt.Errorf("expected 1 tank, got %d", roleCount["tank"])
	}

	if roleCount["healer"] != 1 {
		return fmt.Errorf("expected 1 healer, got %d", roleCount["healer"])
	}

	if roleCount["dps"] != 3 {
		return fmt.Errorf("expected 3 dps, got %d", roleCount["dps"])
	}

	return nil
}
