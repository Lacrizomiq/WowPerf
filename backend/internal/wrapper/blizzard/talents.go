package wrapper

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard"
)

// TransformCharacterTalents transforms the talent data from the Blizzard API into an easier to use Talent struct.
// Using a channel to handle the concurrency of the requests.
func TransformCharacterTalents(data map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string) (*models.TalentLoadout, error) {
	talentLoadout := &models.TalentLoadout{}

	activeSpec, ok := data["active_specialization"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("active specialization not found")
	}

	talentLoadout.LoadoutSpecID = int(activeSpec["id"].(float64))

	specializations, ok := data["specializations"].([]interface{})
	if !ok || len(specializations) == 0 {
		return nil, fmt.Errorf("specializations not found or empty")
	}

	// Find the active loadout
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

	var wg sync.WaitGroup
	talentChan := make(chan models.TalentNode, 100)
	errorChan := make(chan error, 2)

	wg.Add(2)
	go processTalents(activeLoadout, "selected_class_talents", gameDataClient, region, namespace, locale, talentChan, errorChan, &wg)
	go processTalents(activeLoadout, "selected_spec_talents", gameDataClient, region, namespace, locale, talentChan, errorChan, &wg)

	go func() {
		wg.Wait()
		close(talentChan)
		close(errorChan)
	}()

	var classTalents, specTalents []models.TalentNode

	for talent := range talentChan {
		if talent.Node.Type == 0 {
			classTalents = append(classTalents, talent)
		} else {
			specTalents = append(specTalents, talent)
		}
	}

	if len(errorChan) > 0 {
		return nil, <-errorChan
	}

	talentLoadout.ClassTalents = classTalents
	talentLoadout.SpecTalents = specTalents

	return talentLoadout, nil
}

// processTalents is a helper function to process the talents for a specific key
func processTalents(data map[string]interface{}, key string, gameDataClient *blizzard.GameDataClient, region, namespace, locale string, talentChan chan<- models.TalentNode, errorChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	talents, ok := data[key].([]interface{})
	if !ok {
		log.Printf("%s not found or not a slice", key)
		return
	}

	for _, talent := range talents {
		talentMap, ok := talent.(map[string]interface{})
		if !ok {
			continue
		}

		node, err := transformTalentNode(talentMap, gameDataClient, region, namespace, locale)
		if err != nil {
			errorChan <- err
			return
		}

		talentChan <- node
	}
}

// transformTalentNode transforms a single node from the Blizzard API into a struct.
func transformTalentNode(talentMap map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string) (models.TalentNode, error) {
	node := models.TalentNode{}

	nodeInfo, ok := talentMap["node"].(map[string]interface{})
	if !ok {
		nodeInfo = talentMap
	}

	node.Node.ID = safeGetInt(nodeInfo, "id")
	node.Node.TreeID = safeGetInt(nodeInfo, "treeId")
	node.Node.Type = safeGetInt(nodeInfo, "type")
	node.Node.Important = safeGetBool(nodeInfo, "important")
	node.Node.PosX = safeGetInt(nodeInfo, "posX")
	node.Node.PosY = safeGetInt(nodeInfo, "posY")
	node.Node.Row = safeGetInt(nodeInfo, "row")
	node.Node.Col = safeGetInt(nodeInfo, "col")

	node.EntryIndex = safeGetInt(talentMap, "entryIndex")
	node.Rank = safeGetInt(talentMap, "rank")

	node.IncludeInSummary = safeGetBool(talentMap, "includeInSummary")

	entries, ok := nodeInfo["entries"].([]interface{})
	if !ok {
		entries, ok = talentMap["entries"].([]interface{})
		if !ok {
			return node, nil
		}
	}

	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		entryNode, err := transformTalentEntry(entryMap, gameDataClient, region, namespace, locale)
		if err != nil {
			return node, err
		}

		node.Node.Entries = append(node.Node.Entries, entryNode)
	}

	if len(node.Node.Entries) > 0 {
		talentID := node.Node.Entries[0].ID
		talentDetails, err := gameDataClient.GetTalent(talentID, region, namespace, locale)
		if err == nil {
			if description, ok := talentDetails["description"].(string); ok {
				node.Node.Entries[0].Description = description
			}
		}
	}

	return node, nil
}

// Fonctions d'aide pour extraire les valeurs en toute sécurité
func safeGetInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		if intVal, ok := val.(float64); ok {
			return int(intVal)
		}
	}
	return 0
}

func safeGetBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}

// transformTalentEntry transforms a single entry from the Blizzard API into a struct.
func transformTalentEntry(entryMap map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string) (models.TalentEntry, error) {
	entry := models.TalentEntry{}

	entry.ID = safeGetInt(entryMap, "id")
	entry.TraitDefinitionID = safeGetInt(entryMap, "traitDefinitionId")
	entry.Type = safeGetInt(entryMap, "type")
	entry.MaxRanks = safeGetInt(entryMap, "maxRanks")

	spell, ok := entryMap["spell"].(map[string]interface{})
	if !ok {
		spell = entryMap
	}

	entry.Spell.ID = safeGetInt(spell, "id")
	entry.Spell.Name = safeGetString(spell, "name")
	entry.Spell.Icon = safeGetString(spell, "icon")
	entry.Spell.School = safeGetInt(spell, "school")
	entry.Spell.HasCooldown = safeGetBool(spell, "hasCooldown")

	entry.Spell.Rank = safeGetString(entryMap, "rank")

	// retrieve spell media
	if entry.Spell.ID != 0 {
		mediaData, err := gameDataClient.GetSpellMedia(entry.Spell.ID, region, namespace, locale)
		if err == nil {
			if assets, ok := mediaData["assets"].([]interface{}); ok && len(assets) > 0 {
				if asset, ok := assets[0].(map[string]interface{}); ok {
					if value, ok := asset["value"].(string); ok {
						entry.Spell.Icon = value
						parts := strings.Split(value, "/")
						if len(parts) > 0 {
							entry.Spell.Icon = strings.TrimSuffix(parts[len(parts)-1], ".jpg")
						}
					}
				}
			}
		}
	}

	return entry, nil
}

func safeGetString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}
