package wrapper

import (
	"fmt"
	"log"
	profile "wowperf/internal/models"
	talents "wowperf/internal/models/talents"

	"gorm.io/gorm"
)

// TransformCharacterTalents transforms the talents from the Blizzard API into an easier to use TalentLoadout struct
func TransformCharacterTalents(blizzardData map[string]interface{}, db *gorm.DB, treeID, specID int) (*profile.TalentLoadout, error) {
	log.Printf("Blizzard API raw data: %+v", blizzardData)

	talentTree, err := getTalentTreeFromDB(db, treeID, specID)
	if err != nil {
		log.Printf("Error getting talent tree from DB: %v", err)
		return nil, fmt.Errorf("failed to get talent tree from database: %w", err)
	}

	specializations, ok := blizzardData["specializations"].([]interface{})
	if !ok || len(specializations) == 0 {
		return nil, fmt.Errorf("no specializations found in Blizzard data")
	}

	activeSpec := specializations[0].(map[string]interface{})
	loadouts, ok := activeSpec["loadouts"].([]interface{})
	if !ok || len(loadouts) == 0 {
		return nil, fmt.Errorf("no loadouts found in active specialization")
	}

	activeLoadout := loadouts[0].(map[string]interface{})

	talentLoadout := &profile.TalentLoadout{
		LoadoutSpecID:      specID,
		TreeID:             treeID,
		LoadoutText:        getStringValue(activeLoadout, "talent_loadout_code"),
		EncodedLoadoutText: getStringValue(activeLoadout, "talent_loadout_code"),
	}

	classTalents := extractTalents(activeLoadout, "selected_class_talents")
	specTalents := extractTalents(activeLoadout, "selected_spec_talents")
	heroTalents := extractTalents(activeLoadout, "selected_hero_talents")

	log.Printf("Extracted class talents: %d", len(classTalents))
	log.Printf("Extracted spec talents: %d", len(specTalents))
	log.Printf("Extracted hero talents: %d", len(heroTalents))

	talentLoadout.ClassTalents = transformTalents(classTalents, talentTree.ClassNodes)
	talentLoadout.SpecTalents = transformTalents(specTalents, talentTree.SpecNodes)
	talentLoadout.HeroTalents = transformTalents(heroTalents, talentTree.HeroNodes)

	log.Printf("Transformed class talents: %d", len(talentLoadout.ClassTalents))
	log.Printf("Transformed spec talents: %d", len(talentLoadout.SpecTalents))
	log.Printf("Transformed hero talents: %d", len(talentLoadout.HeroTalents))

	return talentLoadout, nil
}

// getTalentTreeFromDB retrieves the talent tree from the database
func getTalentTreeFromDB(db *gorm.DB, treeID, specID int) (*talents.TalentTree, error) {
	var talentTrees []talents.TalentTree
	err := db.Find(&talentTrees).Error
	if err != nil {
		log.Printf("Error fetching all talent trees: %v", err)
	}

	var talentTree talents.TalentTree
	err = db.Where("trait_tree_id = ? AND spec_id = ?", treeID, specID).
		Preload("ClassNodes.Entries").
		Preload("SpecNodes.Entries").
		Preload("HeroNodes.Entries").
		First(&talentTree).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get talent tree: %w", err)
	}
	return &talentTree, nil
}

// extractTalents extracts the talents from the Blizzard API
func extractTalents(data map[string]interface{}, key string) []map[string]interface{} {
	var talents []map[string]interface{}
	if selectedTalents, ok := data[key].([]interface{}); ok {
		for _, talent := range selectedTalents {
			if t, ok := talent.(map[string]interface{}); ok {
				talents = append(talents, t)
			}
		}
	}
	return talents
}

// transformTalents transforms the selected talents from the Blizzard API into a struct
func transformTalents(selectedTalents []map[string]interface{}, dbNodes []talents.TalentNode) []profile.TalentNode {
	var transformedNodes []profile.TalentNode
	dbNodeMap := make(map[int]talents.TalentNode)
	for _, node := range dbNodes {
		dbNodeMap[node.NodeID] = node
	}

	for _, talent := range selectedTalents {
		id, ok := talent["id"].(float64)
		if !ok {
			log.Printf("Warning: talent id not found or not a number")
			continue
		}
		rank, ok := talent["rank"].(float64)
		if !ok {
			log.Printf("Warning: talent rank not found or not a number")
			continue
		}

		dbNode, exists := dbNodeMap[int(id)]
		if !exists {
			log.Printf("Warning: Node %d not found in database", int(id))
			continue
		}

		profileNode := profile.TalentNode{
			NodeID:    int(id),
			NodeType:  dbNode.NodeType,
			Name:      dbNode.Name,
			Type:      dbNode.Type,
			PosX:      dbNode.PosX,
			PosY:      dbNode.PosY,
			MaxRanks:  dbNode.MaxRanks,
			EntryNode: dbNode.EntryNode,
			ReqPoints: dbNode.ReqPoints,
			FreeNode:  dbNode.FreeNode,
			Next:      make([]int, len(dbNode.Next)),
			Prev:      make([]int, len(dbNode.Prev)),
			Entries:   make([]profile.TalentEntry, len(dbNode.Entries)),
			Rank:      int(rank),
		}

		for i, next := range dbNode.Next {
			profileNode.Next[i] = int(next)
		}
		for i, prev := range dbNode.Prev {
			profileNode.Prev[i] = int(prev)
		}

		for i, entry := range dbNode.Entries {
			profileNode.Entries[i] = profile.TalentEntry{
				EntryID:      entry.EntryID,
				DefinitionID: entry.DefinitionID,
				MaxRanks:     entry.MaxRanks,
				Type:         entry.Type,
				Name:         entry.Name,
				SpellID:      entry.SpellID,
				Icon:         entry.Icon,
				Index:        entry.Index,
			}
		}

		transformedNodes = append(transformedNodes, profileNode)
	}

	return transformedNodes
}
