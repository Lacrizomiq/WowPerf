import api from "./api";
import {
  GlobalLeaderboardEntry,
  RoleLeaderboardEntry,
  ClassLeaderboardEntry,
  SpecLeaderboardEntry,
  Role,
  WowClass,
} from "../types/warcraftlogs/globalLeaderboard";
import {
  SpecAverageGlobalScore,
  BestTenPlayerPerSpec,
  MaxKeyLevelsPerSpecAndDungeon,
  DungeonMedia,
} from "@/types/warcraftlogs/globalLeaderboardAnalysis";
import { DungeonLeaderboardResponse } from "../types/warcraftlogs/dungeonRankings";
import { MythicPlusPlayerRankings } from "@/types/warcraftlogs/character/mythicplusPlayerRankings";
import { RaidRankingsResponse } from "@/types/warcraftlogs/character/raidPlayerRankings";
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
  page: number = 1,
  options?: {
    serverSlug?: string;
    serverRegion?: string;
    className?: WowClass;
    specName?: string;
  }
) => {
  try {
    const { data } = await api.get<DungeonLeaderboardResponse>(
      `/warcraftlogs/mythicplus/rankings/dungeon/player`,
      {
        params: {
          encounterID,
          page,
          ...options,
        },
      }
    );
    return data;
  } catch (error) {
    console.error("Error fetching dungeon leaderboard:", error);
    throw error;
  }
};

// Get player Mythic+ rankings
export const getPlayerMythicPlusRankings = async (
  characterName: string,
  serverSlug: string,
  serverRegion: string,
  zoneID: number
) => {
  try {
    const { data } = await api.get<MythicPlusPlayerRankings>(
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

// Get player raid rankings
export const getPlayerRaidRankings = async (
  characterName: string,
  serverSlug: string,
  serverRegion: string,
  zoneID: number
) => {
  try {
    const { data } = await api.get<RaidRankingsResponse>(
      `/warcraftlogs/character/ranking/player`,
      { params: { characterName, serverSlug, serverRegion, zoneID } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching player raid rankings:", error);
    throw error;
  }
};

// getSpecAverageGlobalScore get the average global score for a spec
export const getSpecAverageGlobalScore = async () => {
  try {
    const { data } = await api.get<SpecAverageGlobalScore[]>(
      "/warcraftlogs/mythicplus/analysis/specs/avg-scores"
    );
    return data;
  } catch (error) {
    console.error("Error fetching the average global score for a spec");
    throw error;
  }
};

// getBestTenPlayerPerSpec get the average global score for a spec
export const getBestTenPlayerPerSpec = async () => {
  try {
    const { data } = await api.get<BestTenPlayerPerSpec[]>(
      "/warcraftlogs/mythicplus/analysis/players/top-specs"
    );
    return data;
  } catch (error) {
    console.error("Error fetching the top 10 best players for a spec");
    throw error;
  }
};

// getMaxKeyLevelPerSpecAndDungeon get the max level key for a spec for each dungeon
export const getMaxKeyLevelPerSpecAndDungeon = async () => {
  try {
    const { data } = await api.get<MaxKeyLevelsPerSpecAndDungeon[]>(
      "/warcraftlogs/mythicplus/analysis/specs/dungeons/max-levels-key"
    );
    return data;
  } catch (error) {
    console.error(
      "Error fetching the max level key for a spec for each dungeon"
    );
    throw error;
  }
};

// getDungeonMedia get the dungeon media
export const getDungeonMedia = async () => {
  try {
    const { data } = await api.get<DungeonMedia[]>(
      "/warcraftlogs/mythicplus/analysis/specs/dungeons/media"
    );
    return data;
  } catch (error) {
    console.error("Error fetching the dungeons media");
    throw error;
  }
};
