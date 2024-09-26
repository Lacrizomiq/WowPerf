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
