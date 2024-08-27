package wrapper

import (
	"fmt"
	"log"
	profile "wowperf/internal/models"
	talents "wowperf/internal/models/talents"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TransformCharacterTalents transforms character talent data using Blizzard API data and the local database.
func TransformCharacterTalents(blizzardData map[string]interface{}, db *gorm.DB, treeID, specID int) (*profile.TalentLoadout, error) {
	talentTree, err := getTalentTreeFromDB(db, treeID, specID)
	if err != nil {
		log.Printf("Error getting talent tree from DB: %v", err)
		return nil, fmt.Errorf("failed to get talent tree from database: %w", err)
	}

	selectedHeroTalentTree := getSelectedHeroTalentTree(blizzardData, db)

	talentLoadout := &profile.TalentLoadout{
		LoadoutSpecID:      specID,
		TreeID:             treeID,
		LoadoutText:        getStringValue(blizzardData, "loadout_text"),
		EncodedLoadoutText: getStringValue(blizzardData, "encoded_loadout_text"),
		ClassIcon:          talentTree.ClassIcon,
		SpecIcon:           talentTree.SpecIcon,
		SubTreeNodes:       []profile.SubTreeNode{selectedHeroTalentTree},
	}

	selectedClassTalents := getSelectedTalents(blizzardData, "selected_class_talents")
	selectedSpecTalents := getSelectedTalents(blizzardData, "selected_spec_talents")
	selectedHeroTalents := getSelectedTalents(blizzardData, "selected_hero_talents")

	talentLoadout.ClassTalents = transformTalents(talentTree.ClassNodes, selectedClassTalents)
	talentLoadout.SpecTalents = transformTalents(talentTree.SpecNodes, selectedSpecTalents)
	talentLoadout.HeroTalents = transformHeroTalents(talentTree.HeroNodes, selectedHeroTalents)

	return talentLoadout, nil
}

// getTalentTreeFromDB retrieves the talent tree from the database
func getTalentTreeFromDB(db *gorm.DB, treeID, specID int) (*talents.TalentTree, error) {
	var talentTree talents.TalentTree
	err := db.Where("trait_tree_id = ? AND spec_id = ?", treeID, specID).
		Preload("ClassNodes.Entries").
		Preload("SpecNodes.Entries").
		Preload("HeroNodes.Entries").
		Preload("SubTreeNodes").
		First(&talentTree).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get talent tree: %w", err)
	}

	log.Printf("Number of SubTreeNodes: %d", len(talentTree.SubTreeNodes))
	for _, node := range talentTree.SubTreeNodes {
		log.Printf("SubTreeNode %d: %s", node.SubTreeNodeID, node.Name)
	}

	return &talentTree, nil
}

// getSelectedTalents extracts the selected talents from Blizzard data
func getSelectedTalents(data map[string]interface{}, key string) map[int]int {
	selected := make(map[int]int)
	log.Printf("Getting selected talents for key: %s", key)

	specializations, ok := data["specializations"].([]interface{})
	if !ok {
		log.Printf("Specializations not found or incorrect type")
		return selected
	}

	for _, spec := range specializations {
		specMap, ok := spec.(map[string]interface{})
		if !ok {
			continue
		}

		loadouts, ok := specMap["loadouts"].([]interface{})
		if !ok {
			continue
		}

		for _, loadout := range loadouts {
			loadoutMap, ok := loadout.(map[string]interface{})
			if !ok {
				continue
			}

			talents, ok := loadoutMap[key].([]interface{})
			if !ok {
				continue
			}

			for _, talent := range talents {
				t, ok := talent.(map[string]interface{})
				if !ok {
					continue
				}

				id, ok := t["id"].(float64)
				if !ok {
					continue
				}

				rank, ok := t["rank"].(float64)
				if !ok {
					continue
				}

				selected[int(id)] = int(rank)
				log.Printf("Added talent: id=%d, rank=%d", int(id), int(rank))
			}
		}
	}

	return selected
}

func getSelectedHeroTalentTree(data map[string]interface{}, db *gorm.DB) profile.SubTreeNode {
	var selectedTree map[string]interface{}
	if specializations, ok := data["specializations"].([]interface{}); ok && len(specializations) > 0 {
		if spec, ok := specializations[0].(map[string]interface{}); ok {
			if loadouts, ok := spec["loadouts"].([]interface{}); ok && len(loadouts) > 0 {
				if loadout, ok := loadouts[0].(map[string]interface{}); ok {
					selectedTree = loadout["selected_hero_talent_tree"].(map[string]interface{})
				}
			}
		}
	}

	if selectedTree == nil {
		log.Println("No selected hero talent tree found")
		return profile.SubTreeNode{}
	}

	subTreeNodeID := int(selectedTree["id"].(float64))

	var subTreeEntries []talents.SubTreeEntry
	err := db.Where("trait_sub_tree_id = ?", subTreeNodeID).Find(&subTreeEntries).Error
	if err != nil {
		log.Printf("Error fetching sub tree entries for TraitSubTreeID %d: %v", subTreeNodeID, err)
		return profile.SubTreeNode{}
	}

	if len(subTreeEntries) == 0 {
		log.Printf("No sub tree entries found for TraitSubTreeID: %d", subTreeNodeID)
		return profile.SubTreeNode{}
	}

	log.Printf("Found %d sub tree entries for TraitSubTreeID: %d", len(subTreeEntries), subTreeNodeID)

	return profile.SubTreeNode{
		SubTreeNodeID: subTreeNodeID,
		Name:          selectedTree["name"].(string),
		Type:          "subtree",
		Entries:       transformSubTreeEntries(subTreeEntries),
	}
}

// transformTalents transforms database talent nodes to profile talent nodes
func transformTalents(dbNodes []talents.TalentNode, selectedTalents map[int]int) []profile.TalentNode {
	var transformedNodes []profile.TalentNode
	for _, dbNode := range dbNodes {
		rank := selectedTalents[dbNode.NodeID]
		profileNode := profile.TalentNode{
			NodeID:    dbNode.NodeID,
			NodeType:  dbNode.NodeType,
			Name:      dbNode.Name,
			Type:      dbNode.Type,
			PosX:      dbNode.PosX,
			PosY:      dbNode.PosY,
			MaxRanks:  dbNode.MaxRanks,
			EntryNode: dbNode.EntryNode,
			ReqPoints: dbNode.ReqPoints,
			FreeNode:  dbNode.FreeNode,
			Next:      convertInt64ArrayToIntSlice(dbNode.Next),
			Prev:      convertInt64ArrayToIntSlice(dbNode.Prev),
			Entries:   transformTalentEntries(dbNode.Entries),
			Rank:      rank,
		}
		transformedNodes = append(transformedNodes, profileNode)
	}
	return transformedNodes
}

// transformHeroTalents transforms hero nodes into hero talents
func transformHeroTalents(dbNodes []talents.HeroNode, selectedTalents map[int]int) []profile.HeroTalent {
	var transformedHeroTalents []profile.HeroTalent
	for _, dbNode := range dbNodes {
		rank := selectedTalents[dbNode.NodeID]
		if rank > 0 {
			heroTalent := profile.HeroTalent{
				ID:             dbNode.NodeID,
				Type:           dbNode.Type,
				Name:           dbNode.Name,
				TraitSubTreeID: dbNode.SubTreeID,
				Nodes:          convertInt64ArrayToIntSlice(dbNode.Next),
				Rank:           rank,
				Entries:        transformHeroEntries(dbNode.Entries),
			}
			transformedHeroTalents = append(transformedHeroTalents, heroTalent)
		}
	}
	return transformedHeroTalents
}

// transformHeroEntries transforms database hero entries to profile hero nodes
func transformHeroEntries(dbEntries []talents.HeroEntry) []profile.HeroEntry {
	var transformedEntries []profile.HeroEntry
	for _, entry := range dbEntries {
		transformedEntry := profile.HeroEntry{
			ID:        entry.EntryID,
			Name:      entry.Name,
			Type:      entry.Type,
			MaxRanks:  entry.MaxRanks,
			EntryNode: true,
			FreeNode:  true,
			SpellID:   entry.SpellID,
			Icon:      entry.Icon,
		}
		transformedEntries = append(transformedEntries, transformedEntry)
	}
	return transformedEntries
}

// transformSubTreeNodes transforms database subtree nodes to profile subtree nodes
func transformSubTreeNodes(dbNodes []talents.SubTreeNode, selectedTree profile.SubTreeNode) []profile.SubTreeNode {
	for _, dbNode := range dbNodes {
		if dbNode.SubTreeNodeID == selectedTree.SubTreeNodeID {
			return []profile.SubTreeNode{selectedTree}
		}
	}
	return []profile.SubTreeNode{}
}

// transformTalentEntries transforms database talent entries to profile talent entries
func transformTalentEntries(dbEntries []talents.TalentEntry) []profile.TalentEntry {
	var transformedEntries []profile.TalentEntry
	for _, entry := range dbEntries {
		transformedEntry := profile.TalentEntry{
			EntryID:      entry.EntryID,
			DefinitionID: entry.DefinitionID,
			MaxRanks:     entry.MaxRanks,
			Type:         entry.Type,
			Name:         entry.Name,
			SpellID:      entry.SpellID,
			Icon:         entry.Icon,
			Index:        entry.Index,
		}
		transformedEntries = append(transformedEntries, transformedEntry)
	}
	return transformedEntries
}

// transformSubTreeEntries transforms database subtree entries to profile subtree entries
func transformSubTreeEntries(dbEntries []talents.SubTreeEntry) []profile.SubTreeEntry {
	var transformedEntries []profile.SubTreeEntry
	for _, entry := range dbEntries {
		transformedEntry := profile.SubTreeEntry{
			EntryID:         entry.EntryID,
			Type:            entry.Type,
			Name:            entry.Name,
			TraitSubTreeID:  entry.TraitSubTreeID,
			AtlasMemberName: entry.AtlasMemberName,
			Nodes:           convertInt64ArrayToIntSlice(entry.Nodes),
		}
		transformedEntries = append(transformedEntries, transformedEntry)
	}

	return transformedEntries
}

// convertInt64ArrayToIntSlice converts pq.Int64Array to []int
func convertInt64ArrayToIntSlice(arr pq.Int64Array) []int {
	result := make([]int, len(arr))
	for i, v := range arr {
		result[i] = int(v)
	}
	return result
}
