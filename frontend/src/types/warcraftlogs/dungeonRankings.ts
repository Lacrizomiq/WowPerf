// types/dungeonRankings.ts

import { WowClass } from "./globalLeaderboard";

// Interface for the dungeon server
export interface DungeonServer {
  id: number;
  name: string;
  region: string;
}

// Interface for the team member
export interface TeamMember {
  id: number;
  name: string;
  class: WowClass;
  spec: string;
  role: "Tank" | "Healer" | "DPS";
}

// Interface for the dungeon ranking
export interface DungeonRanking {
  server: DungeonServer;
  duration: number;
  startTime: number;
  deaths: number;
  tanks: number;
  healers: number;
  melee: number;
  ranged: number;
  bracketData: number;
  affixes: number[];
  team: TeamMember[];
  medal: "gold" | "silver" | "bronze";
  score: number;
  leaderboard: number;
}

// Interface for the dungeon leaderboard response
export interface DungeonLeaderboardResponse {
  page: number;
  hasMorePages: boolean;
  count: number;
  rankings: DungeonRanking[];
}

// Re-export WowClass from globalLeaderboard if needed
export type { WowClass } from "./globalLeaderboard";
