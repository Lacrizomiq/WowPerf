package gamedata

import (
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetJournalInstancesIndex retrieves an index of journal instances
func GetJournalInstancesIndex(s *blizzard.GameDataService, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/journal-instance/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetJournalInstanceByID retrieves a journal instance by ID
func GetJournalInstanceByID(s *blizzard.GameDataService, instanceID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/journal-instance/%d", region, instanceID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetJournalInstanceMedia retrieves the media assets for a journal instance
func GetJournalInstanceMedia(s *blizzard.GameDataService, instanceID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/media/journal-instance/%d", region, instanceID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}
