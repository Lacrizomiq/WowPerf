import { useQuery } from "@tanstack/react-query";
import * as raiderioApiServices from "@/libs/raiderioApiServices";

// useGetRaiderioMythicPlusBestRuns retrieves the best runs for a character in a specific season
export const useGetRaiderioMythicPlusBestRuns = (
  season: string,
  region: string,
  dungeon: string,
  page: number
) => {
  return useQuery({
    queryKey: ["raiderio-mythic-plus-best-runs", season, region, dungeon, page],
    queryFn: () =>
      raiderioApiServices.getRaiderioMythicPlusBestRuns(
        season,
        region,
        dungeon,
        page
      ),
  });
};

// useGetRaiderioMythicPlusRunDetails retrieves the details of a specific mythic plus run by its ID
export const useGetRaiderioMythicPlusRunDetails = (
  season: string,
  id: number
) => {
  return useQuery({
    queryKey: ["raiderio-mythic-plus-run-details", season, id],
    queryFn: () =>
      raiderioApiServices.getRaiderioMythicPlusRunDetails(season, id),
  });
};

// useGetRaiderioRaidLeaderboard retrieves the leaderboard of a specific raid by season and region
export const useGetRaiderioRaidLeaderboard = (
  raid: string,
  difficulty: string,
  region: string,
  limit: number,
  page: number
) => {
  return useQuery({
    queryKey: [
      "raiderio-raid-leaderboard",
      raid,
      difficulty,
      region,
      limit,
      page,
    ],
    queryFn: () =>
      raiderioApiServices.getRaiderioRaidLeaderboard(
        raid,
        difficulty,
        region,
        limit,
        page
      ),
  });
};
