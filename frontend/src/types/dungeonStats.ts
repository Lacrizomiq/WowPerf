export type ClassStats = {
  [className: string]: number;
};

export type RoleStats = {
  dps: ClassStats;
  healer: ClassStats;
  tank: ClassStats;
};

export type SpecStats = {
  [className: string]: {
    [specName: string]: number;
  };
};

export type LevelStats = {
  [level: string]: number;
};

export type TeamCompPosition = {
  class: string;
  spec: string;
};

export type TeamCompData = {
  count: number;
  composition: {
    tank: TeamCompPosition;
    healer: TeamCompPosition;
    dps_1: TeamCompPosition;
    dps_2: TeamCompPosition;
    dps_3: TeamCompPosition;
  };
};

export type TeamComp = {
  [compName: string]: TeamCompData;
};

export type DungeonStat = {
  season: string;
  region: string;
  dungeon_slug: string;
  RoleStats: RoleStats;
  SpecStats: SpecStats;
  LevelStats: LevelStats;
  TeamComp: TeamComp;
  updated_at: string;
};

export type DungeonStatsResponse = DungeonStat[];
