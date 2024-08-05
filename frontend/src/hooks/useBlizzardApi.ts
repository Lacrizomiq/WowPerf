import * as apiServices from "@/libs/apiServices";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

export const useGetBlizzardCharacterProfile = (
  region: string,
  realmSlug: string,
  characterName: string
) => {
  return useQuery({
    queryKey: ["characters"],
    queryFn: async () => {
      const { data } = await apiServices.getBlizzardCharacterProfile(
        region,
        realmSlug,
        characterName
      );
      return data;
    },
  });
};

export const useGetBlizzardCharacterMythicKeystoneProfile = (
  region: string,
  realmSlug: string,
  characterName: string
) => {
  return useQuery({
    queryKey: ["mythic-keystone-profile"],
    queryFn: async () => {
      const { data } =
        await apiServices.getBlizzardCharacterMythicKeystoneProfile(
          region,
          realmSlug,
          characterName
        );
      return data;
    },
  });
};

export const useGetBlizzardCharacterEquipment = (
  region: string,
  realmSlug: string,
  characterName: string
) => {
  return useQuery({
    queryKey: ["equipment"],
    queryFn: async () => {
      const { data } = await apiServices.getBlizzardCharacterEquipment(
        region,
        realmSlug,
        characterName
      );
      return data;
    },
  });
};

export const useGetBlizzardCharacterSpecializations = (
  region: string,
  realmSlug: string,
  characterName: string
) => {
  return useQuery({
    queryKey: ["specializations"],
    queryFn: async () => {
      const { data } = await apiServices.getBlizzardCharacterSpecializations(
        region,
        realmSlug,
        characterName
      );
      return data;
    },
  });
};
