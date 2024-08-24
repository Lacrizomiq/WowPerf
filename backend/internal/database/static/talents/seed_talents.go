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

type TalentData struct {
	TraitTreeID   int           `json:"traitTreeId"`
	ClassName     string        `json:"className"`
	ClassID       int           `json:"classId"`
	SpecName      string        `json:"specName"`
	SpecID        int           `json:"specId"`
	ClassNodes    []Node        `json:"classNodes"`
	SpecNodes     []Node        `json:"specNodes"`
	HeroNodes     []Node        `json:"heroNodes"`
	SubTreeNodes  []SubTreeNode `json:"subTreeNodes"`
	FullNodeOrder []int         `json:"fullNodeOrder"`
}

type Node struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	PosX      int     `json:"posX"`
	PosY      int     `json:"posY"`
	MaxRanks  int     `json:"maxRanks"`
	EntryNode bool    `json:"entryNode"`
	ReqPoints int     `json:"reqPoints,omitempty"`
	FreeNode  bool    `json:"freeNode,omitempty"`
	Next      []int   `json:"next"`
	Prev      []int   `json:"prev"`
	Entries   []Entry `json:"entries"`
}

type SubTreeNode struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	PosX            int    `json:"posX"`
	PosY            int    `json:"posY"`
	EntryNode       bool   `json:"entryNode"`
	TraitSubTreeID  int    `json:"traitSubTreeId"`
	AtlasMemberName string `json:"atlasMemberName"`
	Nodes           []int  `json:"nodes"`
}

type Entry struct {
	ID           int    `json:"id"`
	DefinitionID int    `json:"definitionId"`
	MaxRanks     int    `json:"maxRanks"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	SpellID      int    `json:"spellId"`
	Icon         string `json:"icon"`
	Index        int    `json:"index"`
}

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

	var talentData []TalentData
	if err := json.Unmarshal(fileContent, &talentData); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", talentsFilePath, err)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Clear existing data
		if err := clearExistingData(tx); err != nil {
			return err
		}

		for _, data := range talentData {
			log.Printf("Processing talent data for %s %s", data.ClassName, data.SpecName)

			// Create or update ClassTalent
			classTalent := models.ClassTalent{
				TraitTreeID: data.TraitTreeID,
				ClassName:   data.ClassName,
				ClassID:     data.ClassID,
			}
			if err := tx.Where(models.ClassTalent{ClassID: data.ClassID}).FirstOrCreate(&classTalent).Error; err != nil {
				return fmt.Errorf("error creating/updating class talent: %v", err)
			}

			// Create or update SpecTalent
			specTalent := models.SpecTalent{
				TraitTreeID: data.TraitTreeID,
				ClassName:   data.ClassName,
				ClassID:     data.ClassID,
				SpecName:    data.SpecName,
				SpecID:      data.SpecID,
			}
			if err := tx.Where(models.SpecTalent{SpecID: data.SpecID}).FirstOrCreate(&specTalent).Error; err != nil {
				return fmt.Errorf("error creating/updating spec talent: %v", err)
			}

			// Process nodes
			if err := processNodes(tx, data.ClassNodes, &classTalent.ID, nil, nil); err != nil {
				return err
			}
			if err := processNodes(tx, data.SpecNodes, nil, &specTalent.ID, nil); err != nil {
				return err
			}
			if err := processNodes(tx, data.HeroNodes, nil, nil, &data.SpecID); err != nil {
				return err
			}

			// Process SubTreeNodes
			if err := processSubTreeNodes(tx, data); err != nil {
				return err
			}

			// Process FullNodeOrder
			if err := processFullNodeOrder(tx, data); err != nil {
				return err
			}
		}
		return nil
	})
}

func clearExistingData(tx *gorm.DB) error {
	tables := []interface{}{
		&models.ClassTalent{},
		&models.SpecTalent{},
		&models.HeroTalent{},
		&models.SubTreeTalent{},
		&models.TalentNode{},
		&models.TalentEntry{},
		&models.FullNodeOrder{},
	}

	for _, table := range tables {
		if err := tx.Unscoped().Where("1 = 1").Delete(table).Error; err != nil {
			return fmt.Errorf("error clearing table %T: %v", table, err)
		}
	}
	return nil
}

func processNodes(tx *gorm.DB, nodes []Node, classTalentID, specTalentID *uint, heroTalentID *int) error {
	for _, node := range nodes {
		talentNode := models.TalentNode{
			NodeID:        node.ID,
			Name:          node.Name,
			Type:          node.Type,
			PosX:          node.PosX,
			PosY:          node.PosY,
			MaxRanks:      node.MaxRanks,
			EntryNode:     node.EntryNode,
			ReqPoints:     node.ReqPoints,
			FreeNode:      node.FreeNode,
			Next:          pq.Int64Array(convertToInt64Slice(node.Next)),
			Prev:          pq.Int64Array(convertToInt64Slice(node.Prev)),
			ClassTalentID: classTalentID,
			SpecTalentID:  specTalentID,
		}

		if heroTalentID != nil {
			heroTalent := models.HeroTalent{SpecID: *heroTalentID}
			if err := tx.Where(&heroTalent).FirstOrCreate(&heroTalent).Error; err != nil {
				return fmt.Errorf("error creating/updating hero talent: %v", err)
			}
			talentNode.HeroTalentID = &heroTalent.ID
		}

		if err := tx.Create(&talentNode).Error; err != nil {
			return fmt.Errorf("error creating talent node: %v", err)
		}

		for _, entry := range node.Entries {
			talentEntry := models.TalentEntry{
				NodeID:       int(talentNode.ID),
				EntryID:      entry.ID,
				DefinitionID: entry.DefinitionID,
				MaxRanks:     entry.MaxRanks,
				Type:         entry.Type,
				Name:         entry.Name,
				SpellID:      entry.SpellID,
				Icon:         entry.Icon,
				Index:        entry.Index,
			}

			if err := tx.Create(&talentEntry).Error; err != nil {
				return fmt.Errorf("error creating talent entry: %v", err)
			}
		}
	}
	return nil
}

func processSubTreeNodes(tx *gorm.DB, data TalentData) error {
	for _, subTreeNode := range data.SubTreeNodes {
		subTreeTalent := models.SubTreeTalent{
			TraitTreeID:     data.TraitTreeID,
			TraitSubTreeID:  subTreeNode.TraitSubTreeID,
			Name:            subTreeNode.Name,
			AtlasMemberName: subTreeNode.AtlasMemberName,
			Nodes:           subTreeNode.Nodes,
		}

		if err := tx.Create(&subTreeTalent).Error; err != nil {
			return fmt.Errorf("error creating sub tree talent: %v", err)
		}
	}
	return nil
}

func processFullNodeOrder(tx *gorm.DB, data TalentData) error {
	fullNodeOrder := models.FullNodeOrder{
		TraitTreeID: data.TraitTreeID,
		NodeOrder:   pq.Int64Array(convertToInt64Slice(data.FullNodeOrder)),
	}

	if err := tx.Create(&fullNodeOrder).Error; err != nil {
		return fmt.Errorf("error creating full node order: %v", err)
	}
	return nil
}

// convertToInt64Slice converts a slice of ints to a slice of int64s
func convertToInt64Slice(intSlice []int) []int64 {
	int64Slice := make([]int64, len(intSlice))
	for i, v := range intSlice {
		int64Slice[i] = int64(v)
	}
	return int64Slice
}
