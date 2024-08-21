package wrapper

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/gamedata"
)

func TransformCharacterTalents(data map[string]interface{}, gameDataService *blizzard.GameDataService, region, namespace, locale string, treeID, specID int) (*models.TalentLoadout, error) {
	talentLoadout := &models.TalentLoadout{
		LoadoutSpecID: specID,
		TreeID:        treeID,
	}

	specializations, ok := data["specializations"].([]interface{})
	if !ok || len(specializations) == 0 {
		return nil, fmt.Errorf("specializations not found or empty")
	}

	var activeLoadout map[string]interface{}
	for _, spec := range specializations {
		specMap := spec.(map[string]interface{})
		loadouts := specMap["loadouts"].([]interface{})
		for _, loadout := range loadouts {
			loadoutMap := loadout.(map[string]interface{})
			if loadoutMap["is_active"].(bool) {
				activeLoadout = loadoutMap
				break
			}
		}
		if activeLoadout != nil {
			break
		}
	}

	if activeLoadout == nil {
		return nil, fmt.Errorf("active loadout not found")
	}

	talentLoadout.LoadoutText = activeLoadout["talent_loadout_code"].(string)
	talentLoadout.EncodedLoadoutText = EncodeLoadoutString(talentLoadout.LoadoutText)

	var wg sync.WaitGroup
	classTalentsChan := make(chan []models.TalentNode, 1)
	specTalentsChan := make(chan []models.TalentNode, 1)
	errorChan := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		classTalents, err := processTalents(activeLoadout, "selected_class_talents", gameDataService, region, namespace, locale, treeID)
		if err != nil {
			errorChan <- err
		} else {
			classTalentsChan <- classTalents
		}
	}()

	go func() {
		defer wg.Done()
		specTalents, err := processTalents(activeLoadout, "selected_spec_talents", gameDataService, region, namespace, locale, specID)
		if err != nil {
			errorChan <- err
		} else {
			specTalentsChan <- specTalents
		}
	}()

	go func() {
		wg.Wait()
		close(classTalentsChan)
		close(specTalentsChan)
		close(errorChan)
	}()

	select {
	case err := <-errorChan:
		return nil, err
	case classTalents := <-classTalentsChan:
		talentLoadout.ClassTalents = classTalents
	}

	select {
	case err := <-errorChan:
		return nil, err
	case specTalents := <-specTalentsChan:
		talentLoadout.SpecTalents = specTalents
	}

	return talentLoadout, nil
}

func processTalents(data map[string]interface{}, key string, gameDataService *blizzard.GameDataService, region, namespace, locale string, treeID int) ([]models.TalentNode, error) {
	var talents []models.TalentNode

	talentsData, ok := data[key].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no talents found for key: %s", key)
	}

	var wg sync.WaitGroup
	talentChan := make(chan models.TalentNode, len(talentsData))
	errorChan := make(chan error, len(talentsData))

	for _, talent := range talentsData {
		wg.Add(1)
		go func(talent interface{}) {
			defer wg.Done()
			talentMap, ok := talent.(map[string]interface{})
			if !ok {
				errorChan <- fmt.Errorf("invalid talent data structure")
				return
			}

			node, err := transformTalentNode(talentMap, gameDataService, region, namespace, locale, treeID)
			if err != nil {
				errorChan <- err
				return
			}
			talentChan <- node
		}(talent)
	}

	go func() {
		wg.Wait()
		close(talentChan)
		close(errorChan)
	}()

	for talent := range talentChan {
		talents = append(talents, talent)
	}

	if len(errorChan) > 0 {
		return nil, <-errorChan
	}

	return talents, nil
}

func transformTalentNode(talentMap map[string]interface{}, gameDataService *blizzard.GameDataService, region, namespace, locale string, treeID int) (models.TalentNode, error) {
	node := models.TalentNode{}

	node.Node.ID = safeGetInt(talentMap, "id")
	node.Node.TreeID = treeID
	node.EntryIndex = safeGetInt(talentMap, "entryIndex")

	tooltip, ok := talentMap["tooltip"].(map[string]interface{})
	if !ok {
		return node, nil
	}

	spellTooltip, ok := tooltip["spell_tooltip"].(map[string]interface{})
	if !ok {
		return node, nil
	}

	entry, err := transformTalentEntry(spellTooltip, gameDataService, region, namespace, locale)

	if err != nil {
		return node, err
	}

	if talentInfo, ok := tooltip["talent"].(map[string]interface{}); ok {
		entry.TalentID = safeGetInt(talentInfo, "id")
	}

	entry.Type = safeGetInt(talentMap, "type")
	entry.Rank = safeGetInt(talentMap, "rank")

	node.Node.Entries = append(node.Node.Entries, entry)

	return node, nil
}

func transformTalentEntry(entryMap map[string]interface{}, gameDataService *blizzard.GameDataService, region, namespace, locale string) (models.TalentEntry, error) {
	entry := models.TalentEntry{}

	spell, ok := entryMap["spell"].(map[string]interface{})
	if !ok {
		return entry, fmt.Errorf("spell information not found")
	}

	entry.Spell.ID = safeGetInt(spell, "id")
	entry.Spell.Name = safeGetString(spell, "name")
	entry.Spell.Icon = safeGetString(spell, "icon")

	// Always try to retrieve spell media
	mediaData, err := gamedata.GetSpellMedia(gameDataService, entry.Spell.ID, region, namespace, locale)
	if err == nil {
		if assets, ok := mediaData["assets"].([]interface{}); ok && len(assets) > 0 {
			if asset, ok := assets[0].(map[string]interface{}); ok {
				if value, ok := asset["value"].(string); ok {
					entry.Spell.IconURL = value
					// Extract icon name from URL
					parts := strings.Split(value, "/")
					if len(parts) > 0 {
						entry.Spell.Icon = strings.TrimSuffix(parts[len(parts)-1], ".jpg")
					}
				}
			}
		}
	} else {
		log.Printf("Failed to get spell media for spell ID %d: %v", entry.Spell.ID, err)
	}

	return entry, nil
}

func safeGetInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		if intVal, ok := val.(float64); ok {
			return int(intVal)
		}
	}
	return 0
}

// safeGetString returns a string value from a map, or an empty string if the value is not found or is not a string
func safeGetString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

// EncodeLoadoutString encodes a loadout string for use in a URL
func EncodeLoadoutString(loadout string) string {
	if loadout == "" {
		return ""
	}

	// Encode the loadout string
	encoded := url.QueryEscape(loadout)

	// Remove any double percent signs
	re := regexp.MustCompile(`%2`)
	return re.ReplaceAllString(encoded, "")
}
