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

	talentTreeID := treeID

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
		SubTreeNodes:       getSelectedHeroTalentTree(blizzardData, db, specID, talentTreeID),
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

func transformSingleTalent(dbNode talents.TalentNode, rank int) profile.TalentNode {
	transformed := profile.TalentNode{
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
		Rank:      rank,
	}

	if dbNode.Type == "choice" || rank > 0 {
		transformed.Entries = transformTalentEntries(dbNode.Entries)
	}

	return transformed
}

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

func filterTalentsByType(talents []profile.TalentNode, nodeType string) []profile.TalentNode {
	filtered := make([]profile.TalentNode, 0)
	for _, talent := range talents {
		if talent.NodeType == nodeType {
			filtered = append(filtered, talent)
			log.Printf("Filtered talent: NodeID=%d, Name=%s, Type=%s, Rank=%d", talent.NodeID, talent.Name, talent.NodeType, talent.Rank)
		}
	}
	return filtered
}

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
		// sort by posY (top to bottom)
		if nodes[i].PosY != nodes[j].PosY {
			return nodes[i].PosY < nodes[j].PosY
		}
		// if posY is equal, sort by posX (left to right)
		if nodes[i].PosX != nodes[j].PosX {
			return nodes[i].PosX < nodes[j].PosX
		}
		// if posX and posY are equal, use NodeID as final criteria
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

func transformTalentEntry(entry talents.TalentEntry) profile.TalentEntry {
	return profile.TalentEntry{
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

func transformTalentEntries(entries []talents.TalentEntry) []profile.TalentEntry {
	transformed := make([]profile.TalentEntry, len(entries))
	for i, entry := range entries {
		transformed[i] = transformTalentEntry(entry)
	}
	return transformed
}

func getSelectedHeroTalentTree(data map[string]interface{}, db *gorm.DB, specID int, talentTreeID int) []profile.SubTreeNode {
	selectedHeroTalentTree := extractSelectedHeroTalentTree(data)
	if selectedHeroTalentTree == nil {
		log.Println("Selected hero talent tree is nil")
		return []profile.SubTreeNode{}
	}

	subTreeID, ok := selectedHeroTalentTree["id"].(float64)
	if !ok {
		log.Println("Failed to extract subTreeID")
		return []profile.SubTreeNode{}
	}

	log.Printf("Fetching SubTreeNode for specID: %d, talentTreeID: %d, subTreeID: %f", specID, talentTreeID, subTreeID)

	var subTreeNode talents.SubTreeNode
	err := db.Where("spec_id = ? AND talent_tree_id = ?", specID, talentTreeID).First(&subTreeNode).Error
	if err != nil {
		log.Printf("Error fetching SubTreeNode: %v", err)
		return []profile.SubTreeNode{}
	}

	var subTreeEntries []talents.SubTreeEntry
	err = db.Where("sub_tree_node_id = ? AND trait_sub_tree_id = ?", subTreeNode.SubTreeNodeID, int(subTreeID)).
		Find(&subTreeEntries).Error
	if err != nil {
		log.Printf("Error fetching SubTreeEntries: %v", err)
		return []profile.SubTreeNode{}
	}

	if len(subTreeEntries) == 0 {
		log.Println("No SubTreeEntries found")
		return []profile.SubTreeNode{}
	}

	log.Printf("Found SubTreeNode: %+v", subTreeNode)

	transformed := transformSubTreeNode(subTreeNode, subTreeEntries)
	log.Printf("Transformed SubTreeNode: %+v", transformed)

	return []profile.SubTreeNode{transformed}
}

func transformSubTreeNode(dbNode talents.SubTreeNode, dbEntries []talents.SubTreeEntry) profile.SubTreeNode {
	return profile.SubTreeNode{
		ID:      dbNode.SubTreeNodeID,
		Name:    dbNode.Name,
		Type:    dbNode.Type,
		PosX:    dbNode.PosX,
		PosY:    dbNode.PosY,
		Entries: transformSubTreeEntries(dbEntries),
	}
}

func extractSelectedHeroTalentTree(data map[string]interface{}) map[string]interface{} {
	specializations, ok := data["specializations"].([]interface{})
	if !ok || len(specializations) == 0 {
		log.Println("No specializations found")
		return nil
	}

	for _, spec := range specializations {
		specMap, ok := spec.(map[string]interface{})
		if !ok {
			continue
		}

		loadouts, ok := specMap["loadouts"].([]interface{})
		if !ok || len(loadouts) == 0 {
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

			selectedTree, ok := loadoutMap["selected_hero_talent_tree"].(map[string]interface{})
			if !ok {
				log.Println("No selected_hero_talent_tree found in active loadout")
				continue
			}

			log.Printf("Found selected hero talent tree: %+v", selectedTree)
			return selectedTree
		}
	}

	log.Println("No active loadout or selected hero talent tree found")
	return nil
}

func transformSubTreeEntries(dbEntries []talents.SubTreeEntry) []profile.SubTreeEntry {
	transformed := make([]profile.SubTreeEntry, len(dbEntries))
	for i, entry := range dbEntries {
		transformed[i] = profile.SubTreeEntry{
			ID:              entry.EntryID,
			Type:            entry.Type,
			Name:            entry.Name,
			TraitSubTreeID:  entry.TraitSubTreeID,
			AtlasMemberName: entry.AtlasMemberName,
			Nodes:           convertInt64ArrayToIntSlice(entry.Nodes),
		}
	}
	return transformed
}
