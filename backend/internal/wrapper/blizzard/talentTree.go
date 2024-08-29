package wrapper

import (
	"sort"
	talents "wowperf/internal/models/talents"

	"gorm.io/gorm"
)

func GetFullTalentTree(db *gorm.DB, traitTreeID, specID int) (*talents.TalentTree, error) {
	var dbTalentTree talents.TalentTree
	err := db.Where("trait_tree_id = ? AND spec_id = ?", traitTreeID, specID).
		Preload("ClassNodes.Entries").
		Preload("SpecNodes.Entries").
		Preload("HeroNodes.Entries").
		Preload("SubTreeNodes.Entries").
		First(&dbTalentTree).Error
	if err != nil {
		return nil, err
	}

	allNodes := append(dbTalentTree.ClassNodes, dbTalentTree.SpecNodes...)

	talentTree := &talents.TalentTree{
		TraitTreeID:   dbTalentTree.TraitTreeID,
		SpecID:        dbTalentTree.SpecID,
		ClassName:     dbTalentTree.ClassName,
		ClassID:       dbTalentTree.ClassID,
		ClassIcon:     dbTalentTree.ClassIcon,
		SpecName:      dbTalentTree.SpecName,
		SpecIcon:      dbTalentTree.SpecIcon,
		ClassNodes:    filterAndTransformTalentNodes(allNodes, "class"),
		SpecNodes:     filterAndTransformTalentNodes(allNodes, "spec"),
		HeroNodes:     transformAllHeroNodes(dbTalentTree.HeroNodes),
		SubTreeNodes:  transformAllSubTreeNodes(dbTalentTree.SubTreeNodes),
		FullNodeOrder: dbTalentTree.FullNodeOrder,
	}

	sortAllTalentNodes(talentTree.ClassNodes)
	sortAllTalentNodes(talentTree.SpecNodes)

	return talentTree, nil
}

func filterAndTransformTalentNodes(dbNodes []talents.TalentNode, nodeType string) []talents.TalentNode {
	uniqueNodes := make(map[int]talents.TalentNode)
	for _, dbNode := range dbNodes {
		if dbNode.NodeType == nodeType {
			if _, exists := uniqueNodes[dbNode.NodeID]; !exists {
				uniqueNodes[dbNode.NodeID] = talents.TalentNode{
					NodeID:    dbNode.NodeID,
					SpecID:    dbNode.SpecID,
					NodeType:  dbNode.NodeType,
					Name:      dbNode.Name,
					Type:      dbNode.Type,
					PosX:      dbNode.PosX,
					PosY:      dbNode.PosY,
					MaxRanks:  dbNode.MaxRanks,
					EntryNode: dbNode.EntryNode,
					ReqPoints: dbNode.ReqPoints,
					FreeNode:  dbNode.FreeNode,
					Next:      dbNode.Next,
					Prev:      dbNode.Prev,
					Entries:   transformAllTalentEntries(dbNode.Entries),
				}
			}
		}
	}

	filteredNodes := make([]talents.TalentNode, 0, len(uniqueNodes))
	for _, node := range uniqueNodes {
		filteredNodes = append(filteredNodes, node)
	}
	return filteredNodes
}

func sortAllTalentNodes(nodes []talents.TalentNode) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].PosY == nodes[j].PosY {
			return nodes[i].PosX < nodes[j].PosX
		}
		return nodes[i].PosY < nodes[j].PosY
	})
}

func transformAllTalentEntries(dbEntries []talents.TalentEntry) []talents.TalentEntry {
	var entries []talents.TalentEntry
	for _, dbEntry := range dbEntries {
		entries = append(entries, talents.TalentEntry{
			EntryID:      dbEntry.EntryID,
			DefinitionID: dbEntry.DefinitionID,
			MaxRanks:     dbEntry.MaxRanks,
			Type:         dbEntry.Type,
			Name:         dbEntry.Name,
			SpellID:      dbEntry.SpellID,
			Icon:         dbEntry.Icon,
			Index:        dbEntry.Index,
		})
	}
	return entries
}

func transformAllHeroNodes(dbNodes []talents.HeroNode) []talents.HeroNode {
	var nodes []talents.HeroNode
	for _, dbNode := range dbNodes {
		nodes = append(nodes, talents.HeroNode{
			NodeID:       dbNode.NodeID,
			SpecID:       dbNode.SpecID,
			Name:         dbNode.Name,
			Type:         dbNode.Type,
			PosX:         dbNode.PosX,
			PosY:         dbNode.PosY,
			MaxRanks:     dbNode.MaxRanks,
			EntryNode:    dbNode.EntryNode,
			SubTreeID:    dbNode.SubTreeID,
			RequiresNode: dbNode.RequiresNode,
			Next:         dbNode.Next,
			Prev:         dbNode.Prev,
			Entries:      transformAllHeroEntries(dbNode.Entries),
			FreeNode:     dbNode.FreeNode,
		})
	}
	return nodes
}

func transformAllHeroEntries(dbEntries []talents.HeroEntry) []talents.HeroEntry {
	var entries []talents.HeroEntry
	for _, dbEntry := range dbEntries {
		entries = append(entries, talents.HeroEntry{
			EntryID:      dbEntry.EntryID,
			DefinitionID: dbEntry.DefinitionID,
			MaxRanks:     dbEntry.MaxRanks,
			Type:         dbEntry.Type,
			Name:         dbEntry.Name,
			SpellID:      dbEntry.SpellID,
			Icon:         dbEntry.Icon,
			Index:        dbEntry.Index,
		})
	}
	return entries
}

func transformAllSubTreeNodes(dbNodes []talents.SubTreeNode) []talents.SubTreeNode {
	var nodes []talents.SubTreeNode
	for _, dbNode := range dbNodes {
		nodes = append(nodes, talents.SubTreeNode{
			TalentTreeID:  dbNode.TalentTreeID,
			SpecID:        dbNode.SpecID,
			SubTreeNodeID: dbNode.SubTreeNodeID,
			Name:          dbNode.Name,
			Type:          dbNode.Type,
			PosX:          dbNode.PosX,
			PosY:          dbNode.PosY,
			EntryNode:     dbNode.EntryNode,
			Entries:       transformAllSubTreeEntries(dbNode.Entries),
		})
	}
	return nodes
}

func transformAllSubTreeEntries(dbEntries []talents.SubTreeEntry) []talents.SubTreeEntry {
	var entries []talents.SubTreeEntry
	for _, dbEntry := range dbEntries {
		entries = append(entries, talents.SubTreeEntry{
			SubTreeNodeID:   dbEntry.SubTreeNodeID,
			EntryID:         dbEntry.EntryID,
			Type:            dbEntry.Type,
			Name:            dbEntry.Name,
			TraitSubTreeID:  dbEntry.TraitSubTreeID,
			TraitTreeID:     dbEntry.TraitTreeID,
			AtlasMemberName: dbEntry.AtlasMemberName,
			Nodes:           dbEntry.Nodes,
		})
	}
	return entries
}
