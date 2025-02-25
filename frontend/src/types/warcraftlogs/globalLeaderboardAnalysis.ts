// types/globalLeaderboardAnalysis.ts

// SpecAverageGlobalScore is the average global score for a spec
export interface SpecAverageGlobalScore {
  class: string;
  spec: string;
  avg_global_score: number;
  player_count: number;
}

// ClassAverageGlobalScore is the average global score for a class
export interface ClassAverageGlobalScore {
  class: string;
  avg_global_score: number;
  player_count: number;
}

// MaxKeyLevelsPerSpecAndDungeon is the max key levels per spec and dungeon
export interface MaxKeyLevelsPerSpecAndDungeon {
  class: string;
  spec: string;
  dungeon_name: string;
  dungeon_slug: string;
  max_key_level: number;
}

// AverageKeyLevelsPerDungeon is the average key levels per dungeon
export interface AverageKeyLevelsPerDungeon {
  dungeon_name: string;
  dungeon_slug: string;
  avg_key_level: number;
  run_count: number;
}

// BestTenPlayerPerSpec is the best ten player per class with global score
export interface BestTenPlayerPerSpec {
  class: string;
  spec: string;
  name: string;
  server_name: string;
  server_region: string;
  total_score: number;
  rank: number;
}
