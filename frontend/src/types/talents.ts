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
