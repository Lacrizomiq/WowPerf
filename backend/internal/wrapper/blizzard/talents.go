package wrapper

import (
	"fmt"
	"log"
	"strings"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard"
)

func TransformCharacterTalents(data map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string, treeID, specID int) (*models.TalentLoadout, error) {
	talentLoadout := &models.TalentLoadout{
		LoadoutSpecID: specID,
		TreeID:        treeID,
	}

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

	classTalents := processTalents(activeLoadout, "selected_class_talents", gameDataClient, region, namespace, locale, treeID)
	specTalents := processTalents(activeLoadout, "selected_spec_talents", gameDataClient, region, namespace, locale, specID)

	talentLoadout.ClassTalents = classTalents
	talentLoadout.SpecTalents = specTalents

	return talentLoadout, nil
}

func processTalents(data map[string]interface{}, key string, gameDataClient *blizzard.GameDataClient, region, namespace, locale string, treeID int) []models.TalentNode {
	var talents []models.TalentNode

	talentsData, ok := data[key].([]interface{})
	if !ok {
		log.Printf("No talents found for key: %s", key)
		return talents
	}

	for _, talent := range talentsData {
		talentMap, ok := talent.(map[string]interface{})
		if !ok {
			log.Printf("Invalid talent data structure")
			continue
		}

		node, err := transformTalentNode(talentMap, gameDataClient, region, namespace, locale, treeID)
		if err != nil {
			log.Printf("Error transforming talent node: %v", err)
			continue
		}

		talents = append(talents, node)
	}

	return talents
}

func transformTalentNode(talentMap map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string, treeID int) (models.TalentNode, error) {
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

	entry, err := transformTalentEntry(spellTooltip, gameDataClient, region, namespace, locale)
	if err != nil {
		return node, err
	}

	entry.TalentID = safeGetInt(talentMap, "talent_id")
	entry.Type = safeGetInt(talentMap, "type")
	entry.Rank = safeGetInt(talentMap, "rank")

	node.Node.Entries = append(node.Node.Entries, entry)

	return node, nil
}

func transformTalentEntry(entryMap map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string) (models.TalentEntry, error) {
	entry := models.TalentEntry{}

	spell, ok := entryMap["spell"].(map[string]interface{})
	if !ok {
		return entry, fmt.Errorf("spell information not found")
	}

	entry.Spell.ID = safeGetInt(spell, "id")
	entry.Spell.Name = safeGetString(spell, "name")
	entry.Spell.Icon = safeGetString(spell, "icon")

	// Retrieve spell media only if the icon is empty
	if entry.Spell.Icon == "" {
		mediaData, err := gameDataClient.GetSpellMedia(entry.Spell.ID, region, namespace, locale)
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
		}
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

func safeGetString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}
