package gamedata

import (
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetTalentTreeIndex retrieves an index of talent trees
func GetTalentTreeIndex(s *blizzard.GameDataService, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent-tree/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetTalentTree retrieves a talent tree by spec ID
func GetTalentTree(s *blizzard.GameDataService, talentTreeID, specID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent-tree/%d/playable-specialization/%d", region, talentTreeID, specID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetTalentTreeNodes retrieves the nodes of a talent tree as well as links to associated playable specializations given a talent tree id
func GetTalentTreeNodes(s *blizzard.GameDataService, talentTreeID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent-tree/%d", region, talentTreeID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetTalentIndex retrieves an index of talents
func GetTalentIndex(s *blizzard.GameDataService, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetTalentByID retrieves a talent by ID
func GetTalentByID(s *blizzard.GameDataService, talentID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent/%d", region, talentID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}
