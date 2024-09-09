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

	selectedTalents, selectedHeroTalents := getActiveLoadoutTalents(blizzardData, specID)

	classTalents := transformTalents(talentTree.ClassNodes, selectedTalents)
	specTalents := transformTalents(talentTree.SpecNodes, selectedTalents)
	heroTalents := transformHeroTalents(talentTree.HeroNodes, selectedHeroTalents)

	encodedLoadoutText := getEncodedLoadoutText(blizzardData, specID)

	talentLoadout := &profile.TalentLoadout{
		LoadoutSpecID:      specID,
		TreeID:             treeID,
		LoadoutText:        getStringValue(blizzardData, "loadout_text"),
		EncodedLoadoutText: encodedLoadoutText,
		ClassIcon:          talentTree.ClassIcon,
		SpecIcon:           talentTree.SpecIcon,
		ClassTalents:       filterTalentsByType(classTalents, "class"),
		SpecTalents:        filterTalentsByType(specTalents, "spec"),
		HeroTalents:        heroTalents,
		SubTreeNodes:       []profile.SubTreeNode{getSelectedHeroTalentTree(blizzardData, db, specID)},
	}

	sortTalentNodes(talentLoadout.ClassTalents)
	sortTalentNodes(talentLoadout.SpecTalents)

	return talentLoadout, nil
}

func transformTalents(dbNodes []talents.TalentNode, selectedTalents map[int]int) []profile.TalentNode {
	transformed := make([]profile.TalentNode, 0)
	for _, dbNode := range dbNodes {
		rank, isSelected := selectedTalents[dbNode.NodeID]
		if isSelected {
			transformed = append(transformed, transformSingleTalent(dbNode, rank))
		}
	}
	return transformed
}

// getActiveLoadoutTalents extracts the selected talents from the active loadout for a specific specialization
func getActiveLoadoutTalents(data map[string]interface{}, targetSpecID int) (map[int]int, map[int]int) {
	selectedTalents := make(map[int]int)
	selectedHeroTalents := make(map[int]int)
	specializations, ok := data["specializations"].([]interface{})
	if !ok {
		log.Printf("Specializations not found or incorrect type")
		return selectedTalents, selectedHeroTalents
	}

	for _, spec := range specializations {
		specMap, ok := spec.(map[string]interface{})
		if !ok {
			continue
		}

		specInfo, ok := specMap["specialization"].(map[string]interface{})
		if !ok {
			continue
		}

		specID, ok := specInfo["id"].(float64)
		if !ok || int(specID) != targetSpecID {
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

			processTalents(loadoutMap, "selected_class_talents", selectedTalents)
			processTalents(loadoutMap, "selected_spec_talents", selectedTalents)
			processTalents(loadoutMap, "selected_hero_talents", selectedHeroTalents)

			return selectedTalents, selectedHeroTalents
		}
	}

	return selectedTalents, selectedHeroTalents
}

// processTalents processes talents of a specific type and adds them to the selected talents map
func processTalents(loadoutMap map[string]interface{}, key string, selected map[int]int) {
	talents, ok := loadoutMap[key].([]interface{})
	if !ok {
		return
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

		if rank > 0 {
			selected[int(id)] = int(rank)
			log.Printf("Added selected talent: id=%d, rank=%d, type=%s", int(id), int(rank), key)
		}
	}
}

// getEncodedLoadoutText extracts the encoded loadout text from the active loadout
func getEncodedLoadoutText(data map[string]interface{}, targetSpecID int) string {
	specializations, ok := data["specializations"].([]interface{})
	if !ok {
		return ""
	}

	for _, spec := range specializations {
		specMap, ok := spec.(map[string]interface{})
		if !ok {
			continue
		}

		specInfo, ok := specMap["specialization"].(map[string]interface{})
		if !ok {
			continue
		}

		specID, ok := specInfo["id"].(float64)
		if !ok || int(specID) != targetSpecID {
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

	return ""
}

// transformSingleTalent transforms a single talent node from the Blizzard API into a struct.
func transformSingleTalent(dbNode talents.TalentNode, rank int) profile.TalentNode {
	return profile.TalentNode{
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

// transformHeroTalent transforms a hero talent node from the Blizzard API into a struct.
func transformHeroTalents(dbNodes []talents.HeroNode, selectedHeroTalents map[int]int) []profile.HeroTalent {
	transformed := make([]profile.HeroTalent, 0)
	for _, dbNode := range dbNodes {
		rank, isSelected := selectedHeroTalents[dbNode.NodeID]
		if isSelected {
			transformed = append(transformed, transformHeroTalent(dbNode, rank))
		}
	}
	return transformed
}

func transformHeroTalent(dbNode talents.HeroNode, rank int) profile.HeroTalent {
	return profile.HeroTalent{
		ID:      dbNode.NodeID,
		Name:    dbNode.Name,
		Type:    dbNode.Type,
		PosX:    dbNode.PosX,
		PosY:    dbNode.PosY,
		Rank:    rank,
		Entries: transformHeroEntriesToProfileEntries(dbNode.Entries),
	}
}

func transformHeroEntriesToProfileEntries(entries []talents.HeroEntry) []profile.HeroEntry {
	transformed := make([]profile.HeroEntry, len(entries))
	for i, entry := range entries {
		transformed[i] = profile.HeroEntry{
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
	return transformed
}

// filterTalentsByType filters the talents by node type
func filterTalentsByType(talents []profile.TalentNode, nodeType string) []profile.TalentNode {
	filtered := make([]profile.TalentNode, 0)
	for _, talent := range talents {
		if talent.NodeType == nodeType && talent.Rank > 0 {
			filtered = append(filtered, talent)
			log.Printf("Filtered selected talent: NodeID=%d, Name=%s, Type=%s, Rank=%d", talent.NodeID, talent.Name, talent.NodeType, talent.Rank)
		} else {
			log.Printf("Skipped unselected talent: NodeID=%d, Name=%s, Type=%s, Rank=%d", talent.NodeID, talent.Name, talent.NodeType, talent.Rank)
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

func sortTalentNodes(nodes []profile.TalentNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].NodeID < nodes[j].NodeID
	})
}

func convertInt64ArrayToIntSlice(arr pq.Int64Array) []int {
	result := make([]int, len(arr))
	for i, v := range arr {
		result[i] = int(v)
	}
	return result
}

func transformTalentEntries(entries []talents.TalentEntry) []profile.TalentEntry {
	transformed := make([]profile.TalentEntry, len(entries))
	for i, entry := range entries {
		transformed[i] = profile.TalentEntry{
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
	return transformed
}

func getSelectedHeroTalentTree(data map[string]interface{}, db *gorm.DB, specID int) profile.SubTreeNode {
	selectedTalents, _ := getActiveLoadoutTalents(data, specID)
	selectedHeroTalentTree := extractSelectedHeroTalentTree(data)

	if selectedHeroTalentTree == nil {
		log.Println("No selected hero talent tree found")
		return profile.SubTreeNode{}
	}

	subTreeNodeID, ok := selectedHeroTalentTree["id"].(float64)
	if !ok {
		log.Println("Invalid id in selected hero talent tree")
		return profile.SubTreeNode{}
	}

	var subTreeNode talents.SubTreeNode
	err := db.Where("sub_tree_node_id = ? AND spec_id = ?", int(subTreeNodeID), specID).
		Preload("Entries").
		First(&subTreeNode).Error
	if err != nil {
		log.Printf("Error fetching SubTreeNode: %v", err)
		return profile.SubTreeNode{}
	}

	return transformSubTreeNode(subTreeNode, selectedTalents)
}

func extractSelectedHeroTalentTree(data map[string]interface{}) map[string]interface{} {
	specializations, ok := data["specializations"].([]interface{})
	if !ok || len(specializations) == 0 {
		return nil
	}

	spec, ok := specializations[0].(map[string]interface{})
	if !ok {
		return nil
	}

	loadouts, ok := spec["loadouts"].([]interface{})
	if !ok || len(loadouts) == 0 {
		return nil
	}

	loadout, ok := loadouts[0].(map[string]interface{})
	if !ok {
		return nil
	}

	selectedTree, ok := loadout["selected_hero_talent_tree"].(map[string]interface{})
	if !ok {
		return nil
	}

	return selectedTree
}

func transformSubTreeNode(dbNode talents.SubTreeNode, selectedTalents map[int]int) profile.SubTreeNode {
	return profile.SubTreeNode{
		SubTreeNodeID: dbNode.SubTreeNodeID,
		Name:          dbNode.Name,
		Type:          dbNode.Type,
		Entries:       transformSubTreeEntries(dbNode.Entries, selectedTalents),
	}
}

func transformSubTreeEntries(dbEntries []talents.SubTreeEntry, selectedTalents map[int]int) []profile.SubTreeEntry {
	var transformedEntries []profile.SubTreeEntry
	for _, entry := range dbEntries {
		rank := selectedTalents[entry.EntryID]
		if rank > 0 {
			transformedEntry := profile.SubTreeEntry{
				EntryID:         entry.EntryID,
				Type:            entry.Type,
				Name:            entry.Name,
				TraitSubTreeID:  entry.TraitSubTreeID,
				AtlasMemberName: entry.AtlasMemberName,
				Nodes:           convertInt64ArrayToIntSlice(entry.Nodes),
				Rank:            rank,
			}
			transformedEntries = append(transformedEntries, transformedEntry)
		}
	}
	return transformedEntries
}
