import api from "./api";
import {
  GlobalLeaderboardEntry,
  RoleLeaderboardEntry,
  ClassLeaderboardEntry,
  SpecLeaderboardEntry,
  Role,
  WowClass,
} from "../types/warcraftlogs/globalLeaderboard";
import { DungeonLeaderboardResponse } from "../types/warcraftlogs/dungeonRankings";
import { PlayerRankings } from "@/types/warcraftlogs/playerRankings";

// Get global leaderboard
export const getGlobalLeaderboard = async (limit: number = 100) => {
  try {
    const { data } = await api.get<GlobalLeaderboardEntry[]>(
      `/warcraftlogs/mythicplus/global/leaderboard`,
      { params: { limit } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching global leaderboard:", error);
    throw error;
  }
};

// Get leaderboard by role
export const getRoleLeaderboard = async (role: Role, limit: number = 100) => {
  try {
    const { data } = await api.get<RoleLeaderboardEntry[]>(
      `/warcraftlogs/mythicplus/global/leaderboard/role`,
      { params: { role, limit } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching role leaderboard:", error);
    throw error;
  }
};

// Get leaderboard by class
export const getClassLeaderboard = async (
  className: WowClass,
  limit: number = 100
) => {
  try {
    const { data } = await api.get<ClassLeaderboardEntry[]>(
      `/warcraftlogs/mythicplus/global/leaderboard/class`,
      { params: { class: className, limit } }
    );
    return data;
  } catch (error) {
    console.error("Error in getClassLeaderboard:", error);
    throw error;
  }
};

// Get leaderboard by spec
export const getSpecLeaderboard = async (
  className: WowClass,
  spec: string,
  limit: number = 100
) => {
  try {
    const { data } = await api.get<SpecLeaderboardEntry[]>(
      `/warcraftlogs/mythicplus/global/leaderboard/spec`,
      { params: { class: className, spec, limit } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching spec leaderboard:", error);
    throw error;
  }
};

// Get dungeon leaderboard
export const getDungeonLeaderboard = async (
  encounterID: number,
  page: number = 1
) => {
  try {
    const { data } = await api.get<DungeonLeaderboardResponse>(
      `/warcraftlogs/mythicplus/rankings/dungeon`,
      { params: { encounterID, page } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching dungeon leaderboard:", error);
    throw error;
  }
};

// Get player rankings
export const getPlayerRankings = async (
  characterName: string,
  serverSlug: string,
  serverRegion: string,
  zoneID: number
) => {
  try {
    const { data } = await api.get<PlayerRankings>(
      `/warcraftlogs/character/ranking/player`,
      {
        params: { characterName, serverSlug, serverRegion, zoneID },
      }
    );
    return data;
  } catch (error) {
    console.error("Error fetching player rankings:", error);
    throw error;
  }
};
