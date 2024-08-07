import * as apiServices from "@/libs/apiServices";
import { getCharacterGear } from "@/libs/apiServices";
import { getCharacterTalents } from "@/libs/apiServices";
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

export const useGetRaiderIoCharacterGear = (
  region: string,
  realm: string,
  name: string
) => {
  return useQuery({
    queryKey: ["gear", region, realm, name],
    queryFn: async () => {
      try {
        const data = await getCharacterGear(region, realm, name);
        if (!data || !data.gear) {
          throw new Error("No gear data found");
        }
        return data;
      } catch (error) {
        console.error("Error fetching gear data:", error);
        throw error;
      }
    },
  });
};

export const useGetRaiderIoCharacterTalents = (
  region: string,
  realm: string,
  name: string
) => {
  return useQuery({
    queryKey: ["talents", region, realm, name],
    queryFn: async () => {
      try {
        const data = await getCharacterTalents(region, realm, name);
        return data;
      } catch (error) {
        console.error("Error fetching talents data:", error);
        throw error;
      }
    },
    retry: 3,
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
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
