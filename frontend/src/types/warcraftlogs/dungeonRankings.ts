// types/warcraftlogs/dungeonRankings.ts

// Re-export WowClass if needed
export type WowClass =
  | "Warrior"
  | "Paladin"
  | "Hunter"
  | "Rogue"
  | "Priest"
  | "DeathKnight"
  | "Shaman"
  | "Mage"
  | "Warlock"
  | "Monk"
  | "Druid"
  | "DemonHunter"
  | "Evoker";

export interface DungeonGuild {
  faction: number;
  id: number;
  name: string;
}

export interface DungeonServer {
  id: number;
  name: string;
  region: string;
}

export interface DungeonReport {
  code: string;
  fightID: number;
  startTime: number;
}

export interface DungeonRanking {
  affixes: number[];
  amount: number;
  bracketData: number;
  class: WowClass;
  duration: number;
  faction: number;
  guild: DungeonGuild;
  hardModeLevel: number;
  leaderboard: number;
  medal: "gold" | "silver" | "bronze";
  name: string;
  report: DungeonReport;
  score: number;
  server: DungeonServer;
  spec: string;
  startTime: number;
}

export interface DungeonLeaderboardResponse {
  count: number;
  hasMorePages: boolean;
  page: number;
  rankings: DungeonRanking[];
}
