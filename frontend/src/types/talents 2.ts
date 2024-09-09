// types/talents.ts

export interface TalentEntry {
  id: number;
  definitionId: number;
  maxRanks: number;
  type: string;
  name: string;
  spellId: number;
  icon: string;
  index: number;
}

export interface TalentNode {
  id: number;
  nodeType: string;
  name: string;
  type: string;
  posX: number;
  posY: number;
  maxRanks: number;
  entryNode: boolean;
  reqPoints?: number;
  freeNode?: boolean;
  next: number[];
  prev: number[];
  entries: TalentEntry[];
  rank: number;
}

export interface CharacterTalentProps {
  region: string;
  realm: string;
  name: string;
  namespace: string;
  locale: string;
}

export interface TalentLoadout {
  loadout_spec_id: number;
  tree_id: number;
  loadout_text: string;
  encoded_loadout_text: string;
  class_icon: string;
  spec_icon: string;
  class_talents: TalentNode[];
  spec_talents: TalentNode[];
  sub_tree_nodes: SubTreeNode[];
  hero_talents: HeroTalent[];
}

export interface SubTreeNode {
  id: number;
  name: string;
  type: string;
  entries: SubTreeEntry[];
}

export interface SubTreeEntry {
  id: number;
  type: string;
  name: string;
  traitSubTreeId: number;
  atlasMemberName: string;
  nodes: number[];
}

export interface HeroTalent {
  id: number;
  type: string;
  name: string;
  traitSubTreeId: number;
  posX: number;
  posY: number;
  nodes: number[];
  rank: number;
  entries: HeroEntry[];
}

export interface HeroEntry {
  id: number;
  name: string;
  type: string;
  maxRanks: number;
  entryNode: boolean;
  subTreeId: number;
  freeNode: boolean;
  spellId: number;
  icon: string;
}
