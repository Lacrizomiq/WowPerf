package raiderioMythicPlusRunsQueries

import (
	"encoding/json"
	"fmt"
	models "wowperf/internal/models/raiderio/mythicplus_runs"
	service "wowperf/internal/services/raiderio"
)

// MythicPlusRunsParams contient les paramètres pour l'appel API
type MythicPlusRunsParams struct {
	Season    string
	Region    string
	Dungeon   string
	Page      int
	AccessKey string
}

// APIResponse représente la structure exacte de la réponse API
type APIResponse struct {
	Rankings []struct {
		Rank  int        `json:"rank"`
		Score float64    `json:"score"`
		Run   models.Run `json:"run"`
	} `json:"rankings"`
	LeaderboardURL string `json:"leaderboard_url"`
	Params         struct {
		Season  string `json:"season"`
		Region  string `json:"region"`
		Dungeon string `json:"dungeon"`
		Page    int    `json:"page"`
	} `json:"params"`
}

// GetMythicPlusRuns récupère les runs depuis l'API et les retourne parsés
// VERSION OPTIMISÉE - utilise GetRaw pour éviter le double parsing
func GetMythicPlusRuns(s *service.RaiderIOService, params MythicPlusRunsParams) ([]*models.Run, error) {
	// 1. Appel API direct avec GetRaw (évite le double parsing)
	jsonData, err := s.GetRaw("/mythic-plus/runs", buildAPIParams(params))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mythic plus runs: %w", err)
	}

	// 2. Parse directement vers notre structure
	var apiResponse APIResponse
	if err := json.Unmarshal(jsonData, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// 3. Extrait et retourne les runs
	runs := make([]*models.Run, len(apiResponse.Rankings))
	for i, ranking := range apiResponse.Rankings {
		run := ranking.Run
		run.Score = ranking.Score
		runs[i] = &run
	}

	return runs, nil
}

// GetMythicPlusRunsWithScore récupère les runs avec leur score de ranking
// VERSION OPTIMISÉE - utilise GetRaw pour éviter le double parsing
func GetMythicPlusRunsWithScore(s *service.RaiderIOService, params MythicPlusRunsParams) ([]*RunWithScore, error) {
	// 1. Appel API direct avec GetRaw
	jsonData, err := s.GetRaw("/mythic-plus/runs", buildAPIParams(params))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mythic plus runs: %w", err)
	}

	// 2. Parse directement vers notre structure
	var apiResponse APIResponse
	if err := json.Unmarshal(jsonData, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// 3. Extrait et retourne les runs avec score
	runsWithScore := make([]*RunWithScore, len(apiResponse.Rankings))
	for i, ranking := range apiResponse.Rankings {
		runsWithScore[i] = &RunWithScore{
			Run:   &ranking.Run,
			Score: ranking.Score,
			Rank:  ranking.Rank,
		}
	}

	return runsWithScore, nil
}

// buildAPIParams convertit les paramètres en format API
func buildAPIParams(params MythicPlusRunsParams) map[string]string {
	apiParams := map[string]string{
		"season":  params.Season,
		"region":  params.Region,
		"dungeon": params.Dungeon,
		"page":    fmt.Sprintf("%d", params.Page),
	}

	if params.AccessKey != "" {
		apiParams["access_key"] = params.AccessKey
	}

	return apiParams
}

// RunWithScore combine un run avec son score de ranking
type RunWithScore struct {
	Run   *models.Run `json:"run"`
	Score float64     `json:"score"`
	Rank  int         `json:"rank"`
}
