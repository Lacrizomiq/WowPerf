package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	models "wowperf/internal/models/talents"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

const (
	talentsFilePath = "./static/talents.json"
)

// TalentTreeInput represents the input structure for a talent tree
type TalentTreeInput struct {
	TraitTreeID   int                  `json:"traitTreeId"`
	ClassName     string               `json:"className"`
	ClassID       int                  `json:"classId"`
	ClassIcon     string               `json:"classIcon"`
	SpecName      string               `json:"specName"`
	SpecID        int                  `json:"specId"`
	SpecIcon      string               `json:"specIcon"`
	ClassNodes    []models.TalentNode  `json:"classNodes"`
	SpecNodes     []models.TalentNode  `json:"specNodes"`
	HeroNodes     []models.HeroNode    `json:"heroNodes"` // Changed to HeroNode
	SubTreeNodes  []models.SubTreeNode `json:"subTreeNodes"`
	FullNodeOrder []int                `json:"fullNodeOrder"`
}

// SeedTalents seeds the database with talent tree data
func SeedTalents(db *gorm.DB) error {

	var count int64
	if err := db.Model(&models.TalentTree{}).Count(&count).Error; err != nil {
		return fmt.Errorf("error checking existing talent trees: %v", err)
	}
	if count > 0 {
		log.Println("Talent trees already seeded, skipping...")
		return nil
	}

	log.Println("Seeding talents...")
	db.Logger = db.Logger.LogMode(logger.Silent)

	fileContent, err := os.ReadFile(talentsFilePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", talentsFilePath, err)
	}

	var talentTrees []TalentTreeInput
	if err := json.Unmarshal(fileContent, &talentTrees); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", talentsFilePath, err)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, tree := range talentTrees {
			log.Printf("Processing talent tree for %s - %s", tree.ClassName, tree.SpecName)

			fullNodeOrder := make(pq.Int64Array, len(tree.FullNodeOrder))
			for i, v := range tree.FullNodeOrder {
				fullNodeOrder[i] = int64(v)
			}

			newTree := models.TalentTree{
				TraitTreeID:   tree.TraitTreeID,
				SpecID:        tree.SpecID,
				ClassName:     tree.ClassName,
				ClassID:       tree.ClassID,
				ClassIcon:     tree.ClassIcon,
				SpecName:      tree.SpecName,
				SpecIcon:      tree.SpecIcon,
				FullNodeOrder: fullNodeOrder,
			}

			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "trait_tree_id"}, {Name: "spec_id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"class_name", "class_id", "spec_name", "full_node_order",
				}),
			}).Create(&newTree).Error; err != nil {
				return fmt.Errorf("error upserting talent tree: %v", err)
			}

			if err := processNodes(tx, tree.ClassNodes, tree.TraitTreeID, tree.SpecID, "class"); err != nil {
				return err
			}
			if err := processNodes(tx, tree.SpecNodes, tree.TraitTreeID, tree.SpecID, "spec"); err != nil {
				return err
			}
			// Changed to process HeroNodes
			if err := processHeroNodes(tx, tree.HeroNodes, tree.TraitTreeID, tree.SpecID); err != nil {
				return err
			}
			if err := processSubTreeNodes(tx, tree.SubTreeNodes, tree.TraitTreeID, tree.SpecID); err != nil {
				return err
			}
		}

		log.Println("Talent seeding completed successfully")
		return nil
	})
}

// processNodes handles the creation or update of TalentNodes
func processNodes(tx *gorm.DB, nodes []models.TalentNode, talentTreeID int, specID int, nodeType string) error {
	for _, node := range nodes {
		var existingNode models.TalentNode
		err := tx.Where("node_id = ? AND talent_tree_id = ? AND spec_id = ?", node.NodeID, talentTreeID, specID).First(&existingNode).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		currentNode := models.TalentNode{
			TalentTreeID: talentTreeID,
			SpecID:       specID,
			NodeID:       node.NodeID,
			NodeType:     nodeType,
			Name:         node.Name,
			Type:         node.Type,
			PosX:         node.PosX,
			PosY:         node.PosY,
			MaxRanks:     node.MaxRanks,
			EntryNode:    node.EntryNode,
			ReqPoints:    node.ReqPoints,
			FreeNode:     node.FreeNode,
			Next:         node.Next,
			Prev:         node.Prev,
		}

		if err == gorm.ErrRecordNotFound {
			if err := tx.Create(&currentNode).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&existingNode).Updates(currentNode).Error; err != nil {
				return err
			}
		}

		// Process entries
		for _, entry := range node.Entries {
			var existingEntry models.TalentEntry
			err := tx.Where("node_id = ? AND talent_tree_id = ? AND spec_id = ? AND entry_id = ?", currentNode.NodeID, talentTreeID, specID, entry.EntryID).First(&existingEntry).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			newEntry := models.TalentEntry{
				NodeID:       currentNode.NodeID,
				TalentTreeID: talentTreeID,
				SpecID:       specID,
				EntryID:      entry.EntryID,
				DefinitionID: entry.DefinitionID,
				MaxRanks:     entry.MaxRanks,
				Type:         entry.Type,
				Name:         entry.Name,
				SpellID:      entry.SpellID,
				Icon:         entry.Icon,
				Index:        entry.Index,
			}

			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&newEntry).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&existingEntry).Updates(newEntry).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// processHeroNodes handles the creation or update of HeroNodes
func processHeroNodes(tx *gorm.DB, nodes []models.HeroNode, talentTreeID int, specID int) error {
	for _, node := range nodes {
		var existingNode models.HeroNode
		err := tx.Where("node_id = ? AND talent_tree_id = ? AND spec_id = ?", node.NodeID, talentTreeID, specID).First(&existingNode).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		currentNode := models.HeroNode{
			TalentTreeID: talentTreeID,
			SpecID:       specID,
			NodeID:       node.NodeID,
			Name:         node.Name,
			Type:         node.Type,
			PosX:         node.PosX,
			PosY:         node.PosY,
			MaxRanks:     node.MaxRanks,
			EntryNode:    node.EntryNode,
			SubTreeID:    node.SubTreeID,
			RequiresNode: node.RequiresNode,
			Next:         node.Next,
			Prev:         node.Prev,
			FreeNode:     node.FreeNode,
		}

		if err == gorm.ErrRecordNotFound {
			if err := tx.Create(&currentNode).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&existingNode).Updates(currentNode).Error; err != nil {
				return err
			}
		}

		// Process entries
		for _, entry := range node.Entries {
			var existingEntry models.HeroEntry
			err := tx.Where("node_id = ? AND talent_tree_id = ? AND spec_id = ? AND entry_id = ?", currentNode.NodeID, talentTreeID, specID, entry.EntryID).First(&existingEntry).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			newEntry := models.HeroEntry{
				NodeID:       currentNode.NodeID,
				TalentTreeID: talentTreeID,
				SpecID:       specID,
				EntryID:      entry.EntryID,
				DefinitionID: entry.DefinitionID,
				MaxRanks:     entry.MaxRanks,
				Type:         entry.Type,
				Name:         entry.Name,
				SpellID:      entry.SpellID,
				Icon:         entry.Icon,
				Index:        entry.Index,
			}

			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&newEntry).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&existingEntry).Updates(newEntry).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// processSubTreeNodes handles the creation or update of SubTreeNodes
func processSubTreeNodes(tx *gorm.DB, subTreeNodes []models.SubTreeNode, talentTreeID int, specID int) error {
	log.Printf("Processing %d SubTreeNodes for TalentTree %d, SpecID %d", len(subTreeNodes), talentTreeID, specID)

	for _, subTreeNode := range subTreeNodes {
		var existingSubTreeNode models.SubTreeNode
		err := tx.Where("sub_tree_node_id = ? AND talent_tree_id = ? AND spec_id = ?", subTreeNode.SubTreeNodeID, talentTreeID, specID).First(&existingSubTreeNode).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		var currentSubTreeNode models.SubTreeNode
		if err == gorm.ErrRecordNotFound {
			// Create new sub tree node
			currentSubTreeNode = models.SubTreeNode{
				TalentTreeID:  talentTreeID,
				SpecID:        specID,
				SubTreeNodeID: subTreeNode.SubTreeNodeID,
				Name:          subTreeNode.Name,
				Type:          subTreeNode.Type,
				PosX:          subTreeNode.PosX,
				PosY:          subTreeNode.PosY,
				EntryNode:     subTreeNode.EntryNode,
			}
			if err := tx.Create(&currentSubTreeNode).Error; err != nil {
				log.Printf("Error creating SubTreeNode %d: %v", subTreeNode.SubTreeNodeID, err)
				return err
			}
			log.Printf("Created SubTreeNode %d with ID %d", currentSubTreeNode.SubTreeNodeID, currentSubTreeNode.ID)
		} else {
			// Update existing sub tree node
			currentSubTreeNode = existingSubTreeNode
			currentSubTreeNode.Name = subTreeNode.Name
			currentSubTreeNode.Type = subTreeNode.Type
			currentSubTreeNode.PosX = subTreeNode.PosX
			currentSubTreeNode.PosY = subTreeNode.PosY
			currentSubTreeNode.EntryNode = subTreeNode.EntryNode
			if err := tx.Save(&currentSubTreeNode).Error; err != nil {
				log.Printf("Error updating SubTreeNode %d: %v", currentSubTreeNode.SubTreeNodeID, err)
				return err
			}
			log.Printf("Updated SubTreeNode %d with ID %d", currentSubTreeNode.SubTreeNodeID, currentSubTreeNode.ID)
		}

		// Process SubTreeEntries
		for _, entry := range subTreeNode.Entries {
			var existingEntry models.SubTreeEntry
			err := tx.Where("sub_tree_node_id = ? AND entry_id = ?", currentSubTreeNode.SubTreeNodeID, entry.EntryID).First(&existingEntry).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			if err == gorm.ErrRecordNotFound {
				// Create new entry
				newEntry := models.SubTreeEntry{
					SubTreeNodeID:   currentSubTreeNode.SubTreeNodeID,
					EntryID:         entry.EntryID,
					Type:            entry.Type,
					Name:            entry.Name,
					TraitSubTreeID:  entry.TraitSubTreeID,
					TraitTreeID:     talentTreeID,
					AtlasMemberName: entry.AtlasMemberName,
					Nodes:           entry.Nodes,
				}
				if err := tx.Create(&newEntry).Error; err != nil {
					log.Printf("Error creating SubTreeEntry for SubTreeNode %d: %v", currentSubTreeNode.SubTreeNodeID, err)
					return err
				}
				log.Printf("Created SubTreeEntry %d for SubTreeNode %d", entry.EntryID, currentSubTreeNode.SubTreeNodeID)
			} else {
				// Update existing entry
				existingEntry.Type = entry.Type
				existingEntry.Name = entry.Name
				existingEntry.TraitSubTreeID = entry.TraitSubTreeID
				existingEntry.TraitTreeID = talentTreeID
				existingEntry.AtlasMemberName = entry.AtlasMemberName
				existingEntry.Nodes = entry.Nodes
				if err := tx.Save(&existingEntry).Error; err != nil {
					log.Printf("Error updating SubTreeEntry for SubTreeNode %d: %v", currentSubTreeNode.SubTreeNodeID, err)
					return err
				}
				log.Printf("Updated SubTreeEntry %d for SubTreeNode %d", entry.EntryID, currentSubTreeNode.SubTreeNodeID)
			}

			// Process associated nodes for each SubTreeEntry
			for _, nodeID := range entry.Nodes {
				var talentNode models.TalentNode
				if err := tx.Where("node_id = ? AND talent_tree_id = ? AND spec_id = ?", nodeID, talentTreeID, specID).First(&talentNode).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						log.Printf("TalentNode %v not found for SubTreeEntry %d", nodeID, entry.EntryID)
						continue
					}
					return err
				}

				log.Printf("Attempting to associate TalentNode %d with SubTreeNode %d", talentNode.ID, currentSubTreeNode.ID)
				if err := tx.Exec("INSERT INTO sub_tree_node_talents (sub_tree_node_id, talent_node_id) VALUES (?, ?) ON CONFLICT DO NOTHING", currentSubTreeNode.ID, talentNode.ID).Error; err != nil {
					log.Printf("Error inserting association between SubTreeNode %d and TalentNode %d: %v", currentSubTreeNode.ID, talentNode.ID, err)
					return err
				}
				log.Printf("Successfully associated TalentNode %d with SubTreeNode %d", talentNode.ID, currentSubTreeNode.ID)
			}
		}
	}
	return nil
}
