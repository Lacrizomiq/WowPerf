package wrapper

import (
	"fmt"
	"log"
	"sort"
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

	selectedClassTalents := getSelectedTalents(blizzardData, "selected_class_talents")
	selectedSpecTalents := getSelectedTalents(blizzardData, "selected_spec_talents")
	selectedHeroTalents := getSelectedTalents(blizzardData, "selected_hero_talents")
	encodedLoadoutText := getEncodedLoadoutText(blizzardData)

	classTalents := make([]profile.TalentNode, 0)
	specTalents := make([]profile.TalentNode, 0)

	for _, dbNode := range talentTree.ClassNodes {
		if node := transformSingleTalent(dbNode, selectedClassTalents[dbNode.NodeID]); node != nil {
			classTalents = append(classTalents, *node)
		}
	}

	for _, dbNode := range talentTree.SpecNodes {
		if node := transformSingleTalent(dbNode, selectedSpecTalents[dbNode.NodeID]); node != nil {
			specTalents = append(specTalents, *node)
		}
	}

	classTalents = filterTalentsByType(classTalents, "class")
	specTalents = filterTalentsByType(specTalents, "spec")

	talentLoadout := &profile.TalentLoadout{
		LoadoutSpecID:      specID,
		TreeID:             treeID,
		LoadoutText:        getStringValue(blizzardData, "loadout_text"),
		EncodedLoadoutText: encodedLoadoutText,
		ClassIcon:          talentTree.ClassIcon,
		SpecIcon:           talentTree.SpecIcon,
		ClassTalents:       classTalents,
		SpecTalents:        specTalents,
		HeroTalents:        transformHeroTalents(talentTree.HeroNodes, selectedHeroTalents),
		SubTreeNodes:       []profile.SubTreeNode{getSelectedHeroTalentTree(blizzardData, db)},
	}

	sortTalentNodes(talentLoadout.ClassTalents)
	sortTalentNodes(talentLoadout.SpecTalents)

	return talentLoadout, nil
}

// getEncodedLoadoutText extracts the encoded loadout text from Blizzard data
func getEncodedLoadoutText(data map[string]interface{}) string {
	specializations, ok := data["specializations"].([]interface{})
	if !ok {
		log.Printf("Specializations not found or incorrect type")
		return ""
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

			isActive, ok := loadoutMap["is_active"].(bool)
			if !ok || !isActive {
				continue
			}

			talentLoadoutCode, ok := loadoutMap["talent_loadout_code"].(string)
			if !ok {
				log.Printf("talent_loadout_code not found or incorrect type in active loadout")
				return ""
			}

			return talentLoadoutCode
		}
	}

	log.Printf("No active loadout found")
	return ""
}

// transformSingleTalent transforms a single talent node from the Blizzard API into a struct.
func transformSingleTalent(dbNode talents.TalentNode, rank int) *profile.TalentNode {
	return &profile.TalentNode{
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
}

// filterTalentsByType filters the talents by node type
func filterTalentsByType(talents []profile.TalentNode, nodeType string) []profile.TalentNode {
	filtered := make([]profile.TalentNode, 0)
	for _, talent := range talents {
		if talent.NodeType == nodeType {
			filtered = append(filtered, talent)
		}
	}
	return filtered
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

// getSelectedHeroTalentTree extracts the selected hero talent tree from Blizzard data
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
				PosX:           dbNode.PosX,
				PosY:           dbNode.PosY,
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

// sortTalentNodes sorts the talent nodes by their position in the tree
func sortTalentNodes(nodes []profile.TalentNode) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].PosY == nodes[j].PosY {
			return nodes[i].PosX < nodes[j].PosX
		}
		return nodes[i].PosY < nodes[j].PosY
	})
}
