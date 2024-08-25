package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TalentTree struct {
	gorm.Model    `json:"-"`
	TraitTreeID   int           `gorm:"uniqueIndex" json:"traitTreeId"`
	ClassName     string        `json:"className"`
	ClassID       int           `json:"classId"`
	SpecName      string        `json:"specName"`
	SpecID        int           `json:"specId"`
	ClassNodes    []TalentNode  `gorm:"foreignKey:TalentTreeID;references:TraitTreeID" json:"classNodes"`
	SpecNodes     []TalentNode  `gorm:"foreignKey:TalentTreeID;references:TraitTreeID" json:"specNodes"`
	HeroNodes     []TalentNode  `gorm:"foreignKey:TalentTreeID;references:TraitTreeID" json:"heroNodes"`
	SubTreeNodes  []SubTreeNode `gorm:"foreignKey:TalentTreeID;references:TraitTreeID" json:"subTreeNodes"`
	FullNodeOrder pq.Int64Array `gorm:"type:integer[]" json:"fullNodeOrder"`
}

type TalentNode struct {
	gorm.Model   `json:"-"`
	TalentTreeID int           `json:"-"`
	NodeID       int           `gorm:"uniqueIndex" json:"id"`
	NodeType     string        `json:"nodeType"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	PosX         int           `json:"posX"`
	PosY         int           `json:"posY"`
	MaxRanks     int           `json:"maxRanks"`
	EntryNode    bool          `json:"entryNode"`
	ReqPoints    int           `json:"reqPoints,omitempty"`
	FreeNode     bool          `json:"freeNode,omitempty"`
	Next         pq.Int64Array `gorm:"type:integer[]" json:"next"`
	Prev         pq.Int64Array `gorm:"type:integer[]" json:"prev"`
	Entries      []TalentEntry `gorm:"foreignKey:NodeID;references:NodeID" json:"entries"`
}

type TalentEntry struct {
	gorm.Model   `json:"-"`
	NodeID       int    `json:"-"`
	EntryID      int    `json:"id"`
	DefinitionID int    `json:"definitionId"`
	MaxRanks     int    `json:"maxRanks"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	SpellID      int    `json:"spellId"`
	Icon         string `json:"icon"`
	Index        int    `json:"index"`
}

type SubTreeNode struct {
	gorm.Model    `json:"-"`
	TalentTreeID  int            `json:"-"`
	SubTreeNodeID int            `json:"id" gorm:"uniqueIndex:idx_sub_tree_node_id_talent_tree_id"`
	Name          string         `json:"name"`
	Type          string         `json:"type"`
	PosX          int            `json:"posX"`
	PosY          int            `json:"posY"`
	EntryNode     bool           `json:"entryNode"`
	Entries       []SubTreeEntry `gorm:"foreignKey:SubTreeNodeID;references:SubTreeNodeID" json:"entries"`
	Nodes         []TalentNode   `gorm:"many2many:sub_tree_node_talents;" json:"-"`
}

type SubTreeEntry struct {
	gorm.Model      `json:"-"`
	SubTreeNodeID   int           `json:"subTreeNodeId"`
	EntryID         int           `json:"id"`
	Type            string        `json:"type"`
	Name            string        `json:"name"`
	TraitSubTreeID  int           `json:"traitSubTreeId"`
	TraitTreeID     int           `json:"traitTreeId"`
	AtlasMemberName string        `json:"atlasMemberName"`
	Nodes           pq.Int64Array `gorm:"type:integer[]" json:"nodes"`
}
