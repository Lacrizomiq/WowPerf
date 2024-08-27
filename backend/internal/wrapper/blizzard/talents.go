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

	talentLoadout := &profile.TalentLoadout{
		LoadoutSpecID:      specID,
		TreeID:             treeID,
		LoadoutText:        getStringValue(blizzardData, "loadout_text"),
		EncodedLoadoutText: getStringValue(blizzardData, "encoded_loadout_text"),
		ClassIcon:          talentTree.ClassIcon,
		SpecIcon:           talentTree.SpecIcon,
	}

	selectedClassTalents := getSelectedTalents(blizzardData, "selected_class_talents")
	selectedSpecTalents := getSelectedTalents(blizzardData, "selected_spec_talents")
	selectedSubTreeNodes := getSelectedTalents(blizzardData, "selected_sub_tree_nodes")
	selectedHeroTalents := getSelectedTalents(blizzardData, "selected_hero_talents")
	log.Printf("Selected hero talents: %+v", selectedHeroTalents)

	if len(selectedHeroTalents) == 0 {
		log.Printf("No hero talents selected, using all available talents with rank 0")
		for _, node := range talentTree.HeroNodes {
			selectedHeroTalents[node.NodeID] = 0
		}
	}

	talentLoadout.ClassTalents = transformTalents(talentTree.ClassNodes, selectedClassTalents)
	talentLoadout.SpecTalents = transformTalents(talentTree.SpecNodes, selectedSpecTalents)
	talentLoadout.SubTreeNodes = transformSubTreeNodes(talentTree.SubTreeNodes, selectedSubTreeNodes)
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
		Preload("SubTreeNodes.Entries").
		First(&talentTree).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get talent tree: %w", err)
	}

	log.Printf("Fetched talent tree: %+v", talentTree)
	log.Printf("Number of HeroNodes: %d", len(talentTree.HeroNodes))

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
	log.Printf("Transforming hero talents. Number of dbNodes: %d", len(dbNodes))
	for _, dbNode := range dbNodes {
		rank := selectedTalents[dbNode.NodeID]
		log.Printf("Processing hero node: ID=%d, Name=%s, Rank=%d", dbNode.NodeID, dbNode.Name, rank)
		heroTalent := profile.HeroTalent{
			ID:             dbNode.NodeID,
			Type:           dbNode.Type,
			Name:           dbNode.Name,
			TraitSubTreeID: dbNode.SubTreeID,
			Nodes:          convertInt64ArrayToIntSlice(dbNode.Next),
			Rank:           rank,
			Entries:        transformHeroEntries(dbNode.Entries, rank),
		}
		log.Printf("Transformed hero talent: %+v", heroTalent)
		transformedHeroTalents = append(transformedHeroTalents, heroTalent)
	}
	log.Printf("Number of transformed hero talents: %d", len(transformedHeroTalents))
	return transformedHeroTalents
}

// transformHeroEntries transforms database hero entries to profile hero nodes
func transformHeroEntries(dbEntries []talents.HeroEntry, rank int) []profile.HeroEntry {
	var transformedEntries []profile.HeroEntry
	for _, entry := range dbEntries {
		transformedEntry := profile.HeroEntry{
			ID:        entry.EntryID,
			Name:      entry.Name,
			Type:      entry.Type,
			MaxRanks:  rank,
			EntryNode: true,
			SubTreeID: 0,
			FreeNode:  true,
		}
		transformedEntries = append(transformedEntries, transformedEntry)
	}
	return transformedEntries
}

// transformSubTreeNodes transforms database subtree nodes to profile subtree nodes
func transformSubTreeNodes(dbNodes []talents.SubTreeNode, selectedTalents map[int]int) []profile.SubTreeNode {
	var transformedNodes []profile.SubTreeNode
	for _, dbNode := range dbNodes {
		rank := selectedTalents[dbNode.SubTreeNodeID]
		profileNode := profile.SubTreeNode{
			SubTreeNodeID: dbNode.SubTreeNodeID,
			Name:          dbNode.Name,
			Type:          dbNode.Type,
			PosX:          dbNode.PosX,
			PosY:          dbNode.PosY,
			EntryNode:     dbNode.EntryNode,
			Entries:       transformSubTreeEntries(dbNode.Entries),
			Rank:          rank,
		}
		transformedNodes = append(transformedNodes, profileNode)
	}
	return transformedNodes
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
