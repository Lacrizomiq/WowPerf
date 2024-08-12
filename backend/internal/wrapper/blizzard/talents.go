package wrapper

import (
	"fmt"
	"log"
	"strings"
	"sync"
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

	var wg sync.WaitGroup
	talentChan := make(chan models.TalentNode, 100)
	errorChan := make(chan error, 2)

	wg.Add(2)
	go processTalents(activeLoadout, "selected_class_talents", gameDataClient, region, namespace, locale, talentChan, errorChan, &wg, treeID)
	go processTalents(activeLoadout, "selected_spec_talents", gameDataClient, region, namespace, locale, talentChan, errorChan, &wg, specID)

	go func() {
		wg.Wait()
		close(talentChan)
		close(errorChan)
	}()

	var classTalents, specTalents []models.TalentNode

	for talent := range talentChan {
		if talent.Node.TreeID == treeID {
			classTalents = append(classTalents, talent)
		} else if talent.Node.TreeID == specID {
			specTalents = append(specTalents, talent)
		}
	}

	if len(errorChan) > 0 {
		return nil, <-errorChan
	}

	talentLoadout.ClassTalents = classTalents
	talentLoadout.SpecTalents = specTalents

	treeData, err := gameDataClient.GetTalentTree(talentLoadout.TreeID, talentLoadout.LoadoutSpecID, region, namespace, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve talent tree data: %w", err)
	}
	log.Printf("Talent tree data: %+v", treeData)

	err = enrichTalentNodes(talentLoadout, treeData, gameDataClient, region, namespace, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich talent nodes: %w", err)
	}

	return talentLoadout, nil
}

func enrichTalentNodes(talentLoadout *models.TalentLoadout, treeData map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string) error {
	nodes, ok := treeData["class_talent_nodes"].([]interface{})
	if !ok {
		nodes, ok = treeData["spec_talent_nodes"].([]interface{})
		if !ok {
			log.Printf("Warning: neither 'class_talent_nodes' nor 'spec_talent_nodes' field found in talent tree data")
			return nil
		}
	}

	for i := range talentLoadout.ClassTalents {
		enrichNode(&talentLoadout.ClassTalents[i], nodes, gameDataClient, region, namespace, locale)
	}

	for i := range talentLoadout.SpecTalents {
		enrichNode(&talentLoadout.SpecTalents[i], nodes, gameDataClient, region, namespace, locale)
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

func processTalents(data map[string]interface{}, key string, gameDataClient *blizzard.GameDataClient, region, namespace, locale string, talentChan chan<- models.TalentNode, errorChan chan<- error, wg *sync.WaitGroup, treeID int) {
	defer wg.Done()

	talents, ok := data[key].([]interface{})
	if !ok {
		errorChan <- fmt.Errorf("%s not found or not a slice", key)
		return
	}

	for _, talent := range talents {
		talentMap, ok := talent.(map[string]interface{})
		if !ok {
			continue
		}

		node, err := transformTalentNode(talentMap, gameDataClient, region, namespace, locale, treeID)
		if err != nil {
			errorChan <- err
			return
		}

		talentChan <- node
	}
}

func transformTalentNode(talentMap map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string, treeID int) (models.TalentNode, error) {
	node := models.TalentNode{}

	nodeInfo, ok := talentMap["node"].(map[string]interface{})
	if !ok {
		nodeInfo = talentMap
	}

	// Utiliser une fonction d'aide pour gérer les conversions de manière sûre
	node.Node.ID = safeGetInt(nodeInfo, "id")
	node.Node.TreeID = treeID
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

// func to handle safe conversions
func safeGetBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}

func transformTalentEntry(spellTooltip map[string]interface{}, gameDataClient *blizzard.GameDataClient, region, namespace, locale string) (models.TalentEntry, error) {
	entry := models.TalentEntry{
		ID:                0,
		TraitDefinitionID: 0,
		Type:              0,
		MaxRanks:          1, // Assuming 1 as default
	}

	spell, ok := spellTooltip["spell"].(map[string]interface{})
	if !ok {
		return entry, fmt.Errorf("spell information not found")
	}

	entry.Spell.ID = int(spell["id"].(float64))
	entry.Spell.Name = spell["name"].(string)

	// Retrieve spell media
	mediaData, err := gameDataClient.GetSpellMedia(entry.Spell.ID, region, namespace, locale)
	if err == nil {
		if assets, ok := mediaData["assets"].([]interface{}); ok && len(assets) > 0 {
			if asset, ok := assets[0].(map[string]interface{}); ok {
				if value, ok := asset["value"].(string); ok {
					entry.Spell.IconURL = value
					parts := strings.Split(value, "/")
					if len(parts) > 0 {
						entry.Spell.Icon = strings.TrimSuffix(parts[len(parts)-1], ".jpg")
					}
				}
			}
		}
	}

	return entry, nil
}
