package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	models "wowperf/internal/models/talents"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	talentsFilePath = "./static/talents.json"
)

func SeedTalents(db *gorm.DB) error {
	log.Println("Seeding talents...")
	db.Logger = db.Logger.LogMode(logger.Silent)

	fileContent, err := os.ReadFile(talentsFilePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", talentsFilePath, err)
	}

	var talentTrees []models.TalentTree
	if err := json.Unmarshal(fileContent, &talentTrees); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", talentsFilePath, err)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := clearExistingData(tx); err != nil {
			return err
		}

		for _, tree := range talentTrees {
			var existingTree models.TalentTree
			err := tx.Where("trait_tree_id = ?", tree.TraitTreeID).First(&existingTree).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			fullNodeOrder := make(pq.Int64Array, len(tree.FullNodeOrder))
			for i, v := range tree.FullNodeOrder {
				fullNodeOrder[i] = int64(v)
			}

			if err == gorm.ErrRecordNotFound {
				// Create new tree
				newTree := models.TalentTree{
					TraitTreeID:   tree.TraitTreeID,
					ClassName:     tree.ClassName,
					ClassID:       tree.ClassID,
					SpecName:      tree.SpecName,
					SpecID:        tree.SpecID,
					FullNodeOrder: fullNodeOrder,
				}
				if err := tx.Create(&newTree).Error; err != nil {
					return err
				}
			} else {
				// Update existing tree
				existingTree.ClassName = tree.ClassName
				existingTree.ClassID = tree.ClassID
				existingTree.SpecName = tree.SpecName
				existingTree.SpecID = tree.SpecID
				existingTree.FullNodeOrder = fullNodeOrder
				if err := tx.Save(&existingTree).Error; err != nil {
					return err
				}
			}

			if err := processNodes(tx, tree.ClassNodes, tree.TraitTreeID, "class"); err != nil {
				return err
			}
			if err := processNodes(tx, tree.SpecNodes, tree.TraitTreeID, "spec"); err != nil {
				return err
			}
			if err := processNodes(tx, tree.HeroNodes, tree.TraitTreeID, "hero"); err != nil {
				return err
			}
			if err := processSubTreeNodes(tx, tree.SubTreeNodes, tree.TraitTreeID); err != nil {
				return err
			}
		}

		return nil
	})
}

func clearExistingData(tx *gorm.DB) error {
	tables := []interface{}{
		&models.TalentTree{},
		&models.TalentNode{},
		&models.TalentEntry{},
		&models.SubTreeNode{},
		&models.SubTreeEntry{},
	}

	for _, table := range tables {
		if err := tx.Unscoped().Where("1 = 1").Delete(table).Error; err != nil {
			return fmt.Errorf("error clearing table %T: %v", table, err)
		}
	}
	return nil
}

func processNodes(tx *gorm.DB, nodes []models.TalentNode, talentTreeID int, nodeType string) error {
	for _, node := range nodes {
		var existingNode models.TalentNode
		err := tx.Where("node_id = ? AND talent_tree_id = ?", node.NodeID, talentTreeID).First(&existingNode).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		var currentNode models.TalentNode
		if err == gorm.ErrRecordNotFound {
			// Create new node
			currentNode = models.TalentNode{
				TalentTreeID: talentTreeID,
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
			if err := tx.Create(&currentNode).Error; err != nil {
				return err
			}
		} else {
			// Update existing node
			currentNode = existingNode
			currentNode.NodeType = nodeType
			currentNode.Name = node.Name
			currentNode.Type = node.Type
			currentNode.PosX = node.PosX
			currentNode.PosY = node.PosY
			currentNode.MaxRanks = node.MaxRanks
			currentNode.EntryNode = node.EntryNode
			currentNode.ReqPoints = node.ReqPoints
			currentNode.FreeNode = node.FreeNode
			currentNode.Next = node.Next
			currentNode.Prev = node.Prev
			if err := tx.Save(&currentNode).Error; err != nil {
				return err
			}
		}

		// Process entries
		for _, entry := range node.Entries {
			var existingEntry models.TalentEntry
			err := tx.Where("node_id = ? AND entry_id = ?", currentNode.NodeID, entry.EntryID).First(&existingEntry).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			if err == gorm.ErrRecordNotFound {
				// Create new entry
				newEntry := models.TalentEntry{
					NodeID:       currentNode.NodeID,
					EntryID:      entry.EntryID,
					DefinitionID: entry.DefinitionID,
					MaxRanks:     entry.MaxRanks,
					Type:         entry.Type,
					Name:         entry.Name,
					SpellID:      entry.SpellID,
					Icon:         entry.Icon,
					Index:        entry.Index,
				}
				if err := tx.Create(&newEntry).Error; err != nil {
					return err
				}
			} else {
				// Update existing entry
				existingEntry.DefinitionID = entry.DefinitionID
				existingEntry.MaxRanks = entry.MaxRanks
				existingEntry.Type = entry.Type
				existingEntry.Name = entry.Name
				existingEntry.SpellID = entry.SpellID
				existingEntry.Icon = entry.Icon
				existingEntry.Index = entry.Index
				if err := tx.Save(&existingEntry).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func processSubTreeNodes(tx *gorm.DB, subTreeNodes []models.SubTreeNode, talentTreeID int) error {
	log.Printf("Processing %d SubTreeNodes for TalentTree %d", len(subTreeNodes), talentTreeID)

	// First, create or update all SubTreeNodes
	for _, subTreeNode := range subTreeNodes {
		var existingSubTreeNode models.SubTreeNode
		err := tx.Where("sub_tree_node_id = ? AND talent_tree_id = ?", subTreeNode.SubTreeNodeID, talentTreeID).First(&existingSubTreeNode).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		var currentSubTreeNode models.SubTreeNode
		if err == gorm.ErrRecordNotFound {
			// Create new sub tree node
			currentSubTreeNode = models.SubTreeNode{
				TalentTreeID:  talentTreeID,
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

		// Process SubTreeEntries immediately after creating/updating the SubTreeNode
		log.Printf("Processing entries for SubTreeNode %d", currentSubTreeNode.SubTreeNodeID)
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
		}

		// Process associated nodes
		for _, nodeID := range subTreeNode.Nodes {
			var talentNode models.TalentNode
			if err := tx.Where("node_id = ? AND talent_tree_id = ?", nodeID, talentTreeID).First(&talentNode).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Printf("TalentNode %v not found for SubTreeNode %d", nodeID, currentSubTreeNode.SubTreeNodeID)
					continue
				}
				return err
			}

			if err := tx.Exec("INSERT INTO sub_tree_node_talents (sub_tree_node_id, talent_node_id) VALUES (?, ?) ON CONFLICT DO NOTHING", currentSubTreeNode.ID, talentNode.ID).Error; err != nil {
				log.Printf("Error associating TalentNode %d with SubTreeNode %d: %v", talentNode.ID, currentSubTreeNode.ID, err)
				return err
			}
			log.Printf("Associated TalentNode %d with SubTreeNode %d", talentNode.ID, currentSubTreeNode.ID)
		}
	}
	return nil
}
