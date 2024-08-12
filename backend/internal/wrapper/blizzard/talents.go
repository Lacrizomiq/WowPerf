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

	// Fetch talent tree data and enrich talent nodes
	treeData, err := gameDataClient.GetTalentTree(talentLoadout.TreeID, talentLoadout.LoadoutSpecID, region, namespace, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve talent tree data: %w", err)
	}

	err = enrichTalentNodes(talentLoadout, treeData, gameDataClient, region, namespace, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich talent nodes: %w", err)
	}

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

func enrichTalentNodes(talentLoadout *models.TalentLoadout, treeData map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string) error {
	enrichNodesForType := func(nodes []models.TalentNode, nodeType string) error {
		treeNodes, ok := treeData[nodeType].([]interface{})
		if !ok {
			log.Printf("Warning: '%s' field not found in talent tree data", nodeType)
			return nil
		}

		for i := range nodes {
			enrichNode(&nodes[i], treeNodes, gameDataClient, region, namespace, locale)
		}
		return nil
	}

	if err := enrichNodesForType(talentLoadout.ClassTalents, "class_talent_nodes"); err != nil {
		return err
	}
	if err := enrichNodesForType(talentLoadout.SpecTalents, "spec_talent_nodes"); err != nil {
		return err
	}

	return nil
}

func enrichNode(node *models.TalentNode, treeNodes []interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string) {
	for _, treeNode := range treeNodes {
		nodeData := treeNode.(map[string]interface{})
		if int(nodeData["id"].(float64)) == node.Node.ID {
			node.Node.PosX = int(nodeData["raw_position_x"].(float64))
			node.Node.PosY = int(nodeData["raw_position_y"].(float64))
			node.Node.Row = int(nodeData["display_row"].(float64))
			node.Node.Col = int(nodeData["display_col"].(float64))

			entries, ok := nodeData["ranks"].([]interface{})
			if !ok || len(entries) == 0 {
				continue
			}

			firstRank := entries[0].(map[string]interface{})
			tooltip, ok := firstRank["tooltip"].(map[string]interface{})
			if !ok {
				continue
			}

			spellTooltip, ok := tooltip["spell_tooltip"].(map[string]interface{})
			if !ok {
				continue
			}

			talentEntry, err := transformTalentEntry(spellTooltip, gameDataClient, region, namespace, locale)
			if err != nil {
				log.Printf("Error transforming talent entry: %v", err)
				continue
			}

			node.Node.Entries = append(node.Node.Entries, talentEntry)
			break
		}
	}
}

func transformTalentNode(talentMap map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string, treeID int) (models.TalentNode, error) {
	node := models.TalentNode{}

	nodeInfo, ok := talentMap["node"].(map[string]interface{})
	if !ok {
		nodeInfo = talentMap
	}

	node.Node.ID = safeGetInt(nodeInfo, "id")
	node.Node.TreeID = treeID
	node.Node.Type = safeGetInt(nodeInfo, "type")
	node.EntryIndex = safeGetInt(talentMap, "entryIndex")
	node.Rank = safeGetInt(talentMap, "rank")

	entryMap, ok := nodeInfo["entries"].([]interface{})
	if !ok || len(entryMap) == 0 {
		return node, nil
	}

	firstEntry := entryMap[0].(map[string]interface{})
	entry, err := transformTalentEntry(firstEntry, gameDataClient, region, namespace, locale)
	if err != nil {
		return node, err
	}

	node.Node.Entries = []models.TalentEntry{entry}

	return node, nil
}

// func to handle safe conversions
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
func transformTalentEntry(entryMap map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string) (models.TalentEntry, error) {
	entry := models.TalentEntry{}

	entry.ID = safeGetInt(entryMap, "id")
	entry.TraitDefinitionID = safeGetInt(entryMap, "trait_definition_id")
	entry.Type = safeGetInt(entryMap, "type")
	entry.MaxRanks = safeGetInt(entryMap, "max_ranks")

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
