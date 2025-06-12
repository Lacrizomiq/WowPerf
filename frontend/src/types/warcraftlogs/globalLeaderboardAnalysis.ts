// types/globalLeaderboardAnalysis.ts

// SpecAverageGlobalScore is the average global score for a spec
export interface SpecAverageGlobalScore {
  class: string;
  spec: string;
  avg_global_score: number;
  max_global_score: number;
  min_global_score: number;
  player_count: number;
  role: string;
  overall_rank: number;
  role_rank: number;
  slug: string;
}

// SpecDungeonScoreAverage is the average score for a spec per dungeon
export interface SpecDungeonScoreAverage {
  class: string;
  spec: string;
  slug: string;
  encounter_id: number;
  avg_dungeon_score: number;
  max_score: number;
  min_score: number;
  player_count: number;
  role: string;
  overall_rank: number;
  role_rank: number;
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
  encounter_id: number;
  max_key_level: number;
}

// DungeonMedia is the media of dungeons
export interface DungeonMedia {
  dungeon_slug: string;
  encounter_id: number;
  icon: string;
  media_url: string;
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
