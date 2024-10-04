export type ClassStats = {
  [className: string]: number;
};

export type RoleStats = {
  dps: ClassStats;
  healer: ClassStats;
  tank: ClassStats;
};

export type DungeonStat = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  season: string;
  region: string;
  dungeon_slug: string;
  RoleStats: RoleStats;
  updated_at: string;
};

export type DungeonStatsResponse = DungeonStat[];
