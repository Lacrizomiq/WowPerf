package wrapper

import (
	"fmt"
	"log"
	profile "wowperf/internal/models"
	talents "wowperf/internal/models/talents"

	"gorm.io/gorm"
)

// TransformCharacterTalents transforme les données de talents du personnage
func TransformCharacterTalents(blizzardData map[string]interface{}, db *gorm.DB, treeID, specID int) (*profile.TalentLoadout, error) {
	// Récupérer l'arbre de talents de la base de données
	talentTree, err := getTalentTreeFromDB(db, treeID, specID)
	if err != nil {
		log.Printf("Error getting talent tree from DB: %v", err)
		return nil, fmt.Errorf("failed to get talent tree from database: %w", err)
	}

	// Créer le TalentLoadout de base avec les données de Blizzard
	talentLoadout := &profile.TalentLoadout{
		LoadoutSpecID:      specID,
		TreeID:             treeID,
		LoadoutText:        getStringValue(blizzardData, "loadout_text"),
		EncodedLoadoutText: getStringValue(blizzardData, "encoded_loadout_text"),
	}

	// Transformer les talents en utilisant les données de la BDD et de Blizzard
	talentLoadout.ClassTalents = transformTalents(talentTree.ClassNodes, getSelectedTalents(blizzardData, "selected_class_talents"))
	talentLoadout.SpecTalents = transformTalents(talentTree.SpecNodes, getSelectedTalents(blizzardData, "selected_spec_talents"))
	talentLoadout.HeroTalents = transformTalents(talentTree.HeroNodes, getSelectedTalents(blizzardData, "selected_hero_talents"))

	return talentLoadout, nil
}

// getTalentTreeFromDB récupère l'arbre de talents depuis la base de données
func getTalentTreeFromDB(db *gorm.DB, treeID, specID int) (*talents.TalentTree, error) {
	var talentTrees []talents.TalentTree
	err := db.Find(&talentTrees).Error
	if err != nil {
		log.Printf("Error fetching all talent trees: %v", err)
	} else {
		log.Printf("Available talent trees: %v", talentTrees)
	}

	var talentTree talents.TalentTree
	err = db.Where("trait_tree_id = ? AND spec_id = ?", treeID, specID).
		Preload("ClassNodes").
		Preload("SpecNodes").
		Preload("HeroNodes").
		First(&talentTree).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get talent tree: %w", err)
	}
	return &talentTree, nil
}

func getSelectedTalents(data map[string]interface{}, key string) map[int]int {
	selected := make(map[int]int)
	if talents, ok := data[key].([]interface{}); ok {
		for _, talent := range talents {
			if t, ok := talent.(map[string]interface{}); ok {
				id := int(t["id"].(float64))
				rank := int(t["rank"].(float64))
				selected[id] = rank
			}
		}
	}
	return selected
}

func transformTalents(dbNodes []talents.TalentNode, selectedTalents map[int]int) []profile.TalentNode {
	var transformedNodes []profile.TalentNode
	for _, dbNode := range dbNodes {
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
			Next:      make([]int, len(dbNode.Next)),
			Prev:      make([]int, len(dbNode.Prev)),
			Entries:   make([]profile.TalentEntry, len(dbNode.Entries)),
			Rank:      selectedTalents[dbNode.NodeID], // Utilisez le rang de l'API Blizzard
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
