import * as apiServices from "@/libs/apiServices";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { MythicPlusSeasonInfo } from "@/types/mythicPlusRuns";
import { StaticRaid } from "@/types/raids";

// useGetBlizzardCharacterProfile retrieves the profile for a character
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
    staleTime: 30000, // 30 seconds
  });
};

// useSearchBlizzardCharacter is a hook that searches for a character without auto-fetching
export const useSearchBlizzardCharacter = (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  return useQuery({
    queryKey: [
      "characterSearch",
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
    enabled: false, // Disable auto-fetch for search
    retry: false,
    staleTime: 30000,
  });
};

// useGetBlizzardCharacterMythicPlusBestRuns retrieves the best runs for a character in a specific season
export const useGetBlizzardCharacterMythicPlusBestRuns = (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string,
  seasonId: string
) => {
  return useQuery<MythicPlusSeasonInfo | null>({
    queryKey: ["mythic-plus-runs", region, realmSlug, characterName, seasonId],
    queryFn: async () => {
      try {
        return await apiServices.getBlizzardCharacterMythicPlusBestRuns(
          region,
          realmSlug,
          characterName,
          namespace,
          locale,
          seasonId
        );
      } catch (error: any) {
        if (error.response && error.response.status === 500) {
          console.warn("No Mythic+ data available for this season");
          return null;
        }
        throw error;
      }
    },
  });
};

// useGetBlizzardCharacterEquipment retrieves the equipment for a character
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

// useGetBlizzardCharacterSpecializations retrieves the specializations for a character
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

// useGetBlizzardTalentTree retrieves the talent tree for a character
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

// useGetBlizzardMythicDungeonPerSeason retrieves the mythic dungeon per season
export const useGetBlizzardMythicDungeonPerSeason = (seasonSlug: string) => {
  return useQuery({
    queryKey: ["mythicDungeonPerSeason", seasonSlug],
    queryFn: () => apiServices.getBlizzardMythicDungeonPerSeason(seasonSlug),
  });
};

// useGetBlizzardRaidsByExpansion retrieves the raids by expansion
export const useGetBlizzardRaidsByExpansion = (expansion: string) => {
  return useQuery<StaticRaid[]>({
    queryKey: ["raidsByExpansion", expansion],
    queryFn: () => apiServices.getBlizzardRaidsByExpansion(expansion),
  });
};

// useGetBlizzardCharacterEncounterRaid retrieves a character's raid encounters.
export const useGetBlizzardCharacterEncounterRaid = (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  return useQuery({
    queryKey: [
      "characterEncounterRaid",
      region,
      realmSlug,
      characterName,
      namespace,
      locale,
    ],
    queryFn: async () => {
      console.log(
        `Fetching raid encounters for ${characterName} on ${realmSlug}`
      );
      const data = await apiServices.getBlizzardCharacterEncounterRaid(
        region,
        realmSlug,
        characterName,
        namespace,
        locale
      );
      if (data === null) {
        console.log("No raid encounter data available for this character");
        return { expansions: [] };
      }
      return data;
    },
  });
};

// useGetBlizzardCharacterStats retrieves the stats for a character
export const useGetBlizzardCharacterStats = (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  return useQuery({
    queryKey: [
      "characterStats",
      region,
      realmSlug,
      characterName,
      namespace,
      locale,
    ],
    queryFn: () =>
      apiServices.getBlizzardCharacterStats(
        region,
        realmSlug,
        characterName,
        namespace,
        locale
      ),
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 30 * 60 * 1000, // 30 minutes
    retry: 2, // Retry 2 times
  });
};
