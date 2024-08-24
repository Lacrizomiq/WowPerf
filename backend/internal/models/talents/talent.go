package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ClassTalent represents a class talent in the talent tree for a specific Class
type ClassTalent struct {
	gorm.Model  `json:"-"`
	TraitTreeID int          `json:"traitTreeId"`
	ClassName   string       `json:"className"`
	ClassID     int          `json:"classId"`
	Nodes       []TalentNode `gorm:"foreignKey:ClassTalentID"`
}

// SpecTalent represents a spec talent in the talent tree for a specific Spec
type SpecTalent struct {
	gorm.Model  `json:"-"`
	TraitTreeID int          `json:"traitTreeId"`
	ClassName   string       `json:"className"`
	ClassID     int          `json:"classId"`
	SpecName    string       `json:"specName"`
	SpecID      int          `json:"specId"`
	Nodes       []TalentNode `gorm:"foreignKey:SpecTalentID"`
}

// HeroTalent represents a hero talent in the talent tree for a specific class and spec
type HeroTalent struct {
	gorm.Model  `json:"-"`
	TraitTreeID int          `json:"traitTreeId"`
	ClassName   string       `json:"className"`
	ClassID     int          `json:"classId"`
	SpecName    string       `json:"specName"`
	SpecID      int          `json:"specId"`
	SubTreeID   int          `json:"subTreeId"`
	Nodes       []TalentNode `gorm:"foreignKey:HeroTalentID"`
}

// SubTreeTalent represents a sub-tree talent in the talent tree for a specific hero spec
type SubTreeTalent struct {
	gorm.Model      `json:"-"`
	TraitTreeID     int    `json:"traitTreeId"`
	TraitSubTreeID  int    `json:"traitSubTreeId"`
	Name            string `json:"name"`
	AtlasMemberName string `json:"atlasMemberName"`
	Nodes           []int  `gorm:"type:integer[]"`
}

// Talent represents a talent in the talent tree
type TalentNode struct {
	gorm.Model    `json:"-"`
	NodeID        int           `json:"id"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	PosX          int           `json:"posX"`
	PosY          int           `json:"posY"`
	MaxRanks      int           `json:"maxRanks"`
	EntryNode     bool          `json:"entryNode"`
	ReqPoints     int           `json:"reqPoints,omitempty"`
	FreeNode      bool          `json:"freeNode,omitempty"`
	Next          pq.Int64Array `gorm:"type:integer[]"`
	Prev          pq.Int64Array `gorm:"type:integer[]"`
	ClassTalentID *uint
	SpecTalentID  *uint
	HeroTalentID  *uint
	Entries       []TalentEntry `gorm:"foreignKey:NodeID"`
}

// TalentEntry represents an entry in a talent node
type TalentEntry struct {
	gorm.Model   `json:"-"`
	NodeID       int    `json:"nodeId"`
	EntryID      int    `json:"id"`
	DefinitionID int    `json:"definitionId"`
	MaxRanks     int    `json:"maxRanks"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	SpellID      int    `json:"spellId"`
	Icon         string `json:"icon"`
	Index        int    `json:"index"`
}

// FullNodeOrder represents the full node order for a specific talent tree
type FullNodeOrder struct {
	gorm.Model  `json:"-"`
	TraitTreeID int           `json:"traitTreeId"`
	NodeOrder   pq.Int64Array `gorm:"type:integer[]"`
}
