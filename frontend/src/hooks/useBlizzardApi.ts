import * as apiServices from "@/libs/apiServices";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

export const useGetBlizzardCharacterProfile = (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  return useQuery({
    queryKey: [
      "characters",
      region,
      realmSlug,
      characterName,
      namespace,
      locale,
    ],
    queryFn: () =>
      apiServices.getBlizzardCharacterProfile(
        region,
        realmSlug,
        characterName,
        namespace,
        locale
      ),
  });
};

export const useGetBlizzardCharacterMythicPlusBestRuns = (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string,
  seasonId: string
) => {
  return useQuery({
    queryKey: [
      "mythic-plus-best-runs",
      seasonId,
      region,
      realmSlug,
      characterName,
      namespace,
      locale,
    ],
    queryFn: async () => {
      apiServices.getBlizzardCharacterMythicPlusBestRuns(
        region,
        realmSlug,
        characterName,
        namespace,
        locale,
        seasonId
      );
    },
  });
};

export const useGetBlizzardCharacterEquipment = (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  return useQuery({
    queryKey: [
      "equipment",
      region,
      realmSlug,
      characterName,
      namespace,
      locale,
    ],
    queryFn: () =>
      apiServices.getBlizzardCharacterEquipment(
        region,
        realmSlug,
        characterName,
        namespace,
        locale
      ),
  });
};

export const useGetBlizzardCharacterSpecializations = (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  return useQuery({
    queryKey: [
      "specializations",
      region,
      realmSlug,
      characterName,
      namespace,
      locale,
    ],
    queryFn: () =>
      apiServices.getBlizzardCharacterSpecializations(
        region,
        realmSlug,
        characterName,
        namespace,
        locale
      ),
  });
};

export const useGetBlizzardTalentTree = (
  talentTreeId: number,
  specId: number,
  region: string,
  namespace: string,
  locale: string
) => {
  return useQuery({
    queryKey: ["talentTree", talentTreeId, specId, region, namespace, locale],
    queryFn: () =>
      apiServices.getBlizzardTalentTree(
        talentTreeId,
        specId,
        region,
        namespace,
        locale
      ),
  });
};

export const useGetBlizzardMythicDungeonPerSeason = (seasonSlug: string) => {
  return useQuery({
    queryKey: ["mythicDungeonPerSeason", seasonSlug],
    queryFn: () => apiServices.getBlizzardMythicDungeonPerSeason(seasonSlug),
  });
};
