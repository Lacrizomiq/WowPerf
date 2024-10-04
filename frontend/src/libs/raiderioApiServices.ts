import axios from "axios";
import api from "./api";

// getRaiderioMythicPlusBestRuns retrieves the best runs of all characters in a specific season and region and dungeon, by page
export const getRaiderioMythicPlusBestRuns = async (
  season: string,
  region: string,
  dungeon: string,
  page: number
) => {
  try {
    const response = await api.get(`/mythic-plus/best-runs`, {
      params: { season, region, dungeon, page },
    });
    return response.data;
  } catch (error) {
    console.error("Error in getRaiderioMythicPlusBestRuns:", error);
    throw error;
  }
};

// getRaiderioMythicPlusRunDetails retrieves the details of a specific mythic plus run by its ID
export const getRaiderioMythicPlusRunDetails = async (
  season: string,
  id: number
) => {
  try {
    const response = await api.get(`/mythic-plus/run-details`, {
      params: { season, id },
    });
    return response.data;
  } catch (error) {
    console.error("Error in getRaiderioMythicPlusRunDetails:", error);
    throw error;
  }
};

// getRaiderioRaidLeaderboard retrieves the leaderboard of a specific raid by season and region
export const getRaiderioRaidLeaderboard = async (
  raid: string,
  difficulty: string,
  region: string,
  limit: number,
  page: number
) => {
  try {
    const response = await api.get(`/raids/leaderboard`, {
      params: { raid, difficulty, region, limit, page },
    });
    return response.data;
  } catch (error) {
    console.error("Error in getRaiderioRaidLeaderboard:", error);
    throw error;
  }
};

// getDungeonStats retrieves the dungeon stats of a specific season and region
export const getDungeonStats = async (season: string, region: string) => {
  try {
    const response = await api.get(`/mythic-plus/dungeon-stats`, {
      params: { season, region },
    });
    return response.data;
  } catch (error) {
    console.error("Error in getDungeonStats:", error);
    throw error;
  }
};
