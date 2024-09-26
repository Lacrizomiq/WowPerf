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
