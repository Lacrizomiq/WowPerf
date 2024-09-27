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
