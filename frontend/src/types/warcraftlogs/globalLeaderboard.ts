// types/globalLeaderboard.ts

// Interface for the global leaderboard entry
export interface GlobalLeaderboardEntry {
  player_id: number;
  name: string;
  class: WowClass;
  spec: string;
  role: Role;
  total_score: number;
  rank: number;
  dungeon_count: number;
  server_name: string;
  server_region: string;
}

// Enum for the roles
export type Role = "tank" | "healer" | "dps";

// Enum for the classes
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

// Types for the different leaderboards that use the same structure
export type RoleLeaderboardEntry = GlobalLeaderboardEntry;
export type ClassLeaderboardEntry = GlobalLeaderboardEntry;
export type SpecLeaderboardEntry = GlobalLeaderboardEntry;
