package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TalentTree struct {
	gorm.Model    `json:"-"`
	TraitTreeID   int           `gorm:"uniqueIndex:idx_trait_tree_spec;not null" json:"traitTreeId"`
	SpecID        int           `gorm:"uniqueIndex:idx_trait_tree_spec;not null" json:"specId"`
	ClassName     string        `json:"className"`
	ClassID       int           `json:"classId"`
	ClassIcon     string        `json:"classIcon"`
	SpecName      string        `json:"specName"`
	SpecIcon      string        `json:"specIcon"`
	ClassNodes    []TalentNode  `gorm:"foreignKey:TalentTreeID,SpecID;references:TraitTreeID,SpecID" json:"classNodes"`
	SpecNodes     []TalentNode  `gorm:"foreignKey:TalentTreeID,SpecID;references:TraitTreeID,SpecID" json:"specNodes"`
	HeroNodes     []HeroNode    `gorm:"foreignKey:TalentTreeID,SpecID;references:TraitTreeID,SpecID" json:"heroNodes"`
	SubTreeNodes  []SubTreeNode `gorm:"foreignKey:TalentTreeID,SpecID;references:TraitTreeID,SpecID" json:"subTreeNodes"`
	FullNodeOrder pq.Int64Array `gorm:"type:integer[]" json:"fullNodeOrder"`
}

type TalentNode struct {
	gorm.Model   `json:"-"`
	TalentTreeID int           `gorm:"uniqueIndex:idx_node_tree_spec" json:"-"`
	SpecID       int           `gorm:"uniqueIndex:idx_node_tree_spec" json:"specId"`
	NodeID       int           `gorm:"uniqueIndex:idx_node_tree_spec" json:"id"`
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
	Entries      []TalentEntry `gorm:"foreignKey:NodeID,TalentTreeID,SpecID;references:NodeID,TalentTreeID,SpecID" json:"entries"`
}

type TalentEntry struct {
	gorm.Model   `json:"-"`
	NodeID       int    `gorm:"uniqueIndex:idx_entry_node_tree_spec" json:"-"`
	TalentTreeID int    `gorm:"uniqueIndex:idx_entry_node_tree_spec" json:"-"`
	SpecID       int    `gorm:"uniqueIndex:idx_entry_node_tree_spec" json:"-"`
	EntryID      int    `gorm:"uniqueIndex:idx_entry_node_tree_spec" json:"id"`
	DefinitionID int    `json:"definitionId"`
	MaxRanks     int    `json:"maxRanks"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	SpellID      int    `json:"spellId"`
	Icon         string `json:"icon"`
	Index        int    `json:"index"`
}

type HeroNode struct {
	gorm.Model   `json:"-"`
	TalentTreeID int           `gorm:"uniqueIndex:idx_hero_node_tree_spec" json:"-"`
	SpecID       int           `gorm:"uniqueIndex:idx_hero_node_tree_spec" json:"specId"`
	NodeID       int           `gorm:"uniqueIndex:idx_hero_node_tree_spec" json:"id"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	PosX         int           `json:"posX"`
	PosY         int           `json:"posY"`
	MaxRanks     int           `json:"maxRanks"`
	EntryNode    bool          `json:"entryNode"`
	SubTreeID    int           `json:"subTreeId"`
	RequiresNode int           `json:"requiresNode,omitempty"`
	Next         pq.Int64Array `gorm:"type:integer[]" json:"next"`
	Prev         pq.Int64Array `gorm:"type:integer[]" json:"prev"`
	Entries      []HeroEntry   `gorm:"foreignKey:NodeID,TalentTreeID,SpecID;references:NodeID,TalentTreeID,SpecID" json:"entries"`
	FreeNode     bool          `json:"freeNode,omitempty"`
}

type HeroEntry struct {
	gorm.Model   `json:"-"`
	NodeID       int    `gorm:"uniqueIndex:idx_hero_entry_node_tree_spec" json:"-"`
	TalentTreeID int    `gorm:"uniqueIndex:idx_hero_entry_node_tree_spec" json:"-"`
	SpecID       int    `gorm:"uniqueIndex:idx_hero_entry_node_tree_spec" json:"-"`
	EntryID      int    `gorm:"uniqueIndex:idx_hero_entry_node_tree_spec" json:"id"`
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
	TalentTreeID  int            `json:"talentTreeId"`
	SpecID        int            `json:"specId"`
	SubTreeNodeID int            `json:"id" gorm:"uniqueIndex:idx_sub_tree_node_id_talent_tree_id_spec_id"`
	Name          string         `json:"name"`
	Type          string         `json:"type"`
	PosX          int            `json:"posX"`
	PosY          int            `json:"posY"`
	EntryNode     bool           `json:"entryNode"`
	Entries       []SubTreeEntry `gorm:"foreignKey:SubTreeNodeID;references:SubTreeNodeID" json:"entries"`
	Nodes         []TalentNode   `gorm:"many2many:sub_tree_node_talents;"`
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
