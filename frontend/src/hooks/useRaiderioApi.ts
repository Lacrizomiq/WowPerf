import * as apiServices from "@/libs/apiServices";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

export const useRaiderIoCharacterProfile = (
  region: string,
  realm: string,
  name: string,
  fields?: string[]
) => {
  return useQuery({
    queryKey: ["characterProfile", region, realm, name, fields],
    queryFn: async () => {
      try {
        const data = await apiServices.getCharacterProfile(
          region,
          realm,
          name,
          fields
        );
        if (!data) {
          throw new Error("No data returned from API");
        }
        return data;
      } catch (error) {
        console.error("Error in useRaiderIoCharacterProfile:", error);
        throw error;
      }
    },
  });
};

export const useGetRaiderIoCharacterMythicPlusScores = (
  region: string,
  realm: string,
  name: string
) => {
  return useQuery({
    queryKey: ["mythic-plus-scores"],
    queryFn: async () => {
      const { data } = await apiServices.getCharacterMythicPlusScores(
        region,
        realm,
        name
      );
      return data;
    },
  });
};

export const useGetRaiderIoCharacterRaidProgression = (
  region: string,
  realm: string,
  name: string
) => {
  return useQuery({
    queryKey: ["raid-progression"],
    queryFn: async () => {
      const { data } = await apiServices.getCharacterRaidProgression(
        region,
        realm,
        name
      );
      return data;
    },
  });
};
