import { useQuery } from "@tanstack/react-query";
import * as warcraftLogsApiService from "../libs/warcraftlogsApiService";
import {
  GlobalLeaderboardEntry,
  RoleLeaderboardEntry,
  ClassLeaderboardEntry,
  SpecLeaderboardEntry,
  Role,
  WowClass,
} from "../types/warcraftlogs/globalLeaderboard";
import { DungeonLeaderboardResponse } from "../types/warcraftlogs/dungeonRankings";
import { MythicPlusPlayerRankings } from "@/types/warcraftlogs/character/mythicplusPlayerRankings";
import { RaidRankingsResponse } from "@/types/warcraftlogs/character/raidPlayerRankings";

// Hook for global leaderboard with required limit
export const useGetGlobalLeaderboard = (limit: number) => {
  return useQuery<GlobalLeaderboardEntry[], Error>({
    queryKey: ["warcraftlogs-global-leaderboard", limit],
    queryFn: () => warcraftLogsApiService.getGlobalLeaderboard(limit),
  });
};

// Hook for role leaderboard with required limit
export const useGetRoleLeaderboard = (role: Role, limit: number) => {
  return useQuery<RoleLeaderboardEntry[], Error>({
    queryKey: ["warcraftlogs-role-leaderboard", role, limit],
    queryFn: () => warcraftLogsApiService.getRoleLeaderboard(role, limit),
    enabled: !!role,
  });
};

// Hook for class leaderboard with required limit
export const useGetClassLeaderboard = (className: WowClass, limit: number) => {
  return useQuery<ClassLeaderboardEntry[], Error>({
    queryKey: ["warcraftlogs-class-leaderboard", className, limit],
    queryFn: () => warcraftLogsApiService.getClassLeaderboard(className, limit),
    enabled: !!className,
  });
};

// Hook for spec leaderboard with required limit
export const useGetSpecLeaderboard = (
  className: WowClass,
  spec: string,
  limit: number
) => {
  return useQuery<SpecLeaderboardEntry[], Error>({
    queryKey: ["warcraftlogs-spec-leaderboard", className, spec, limit],
    queryFn: () =>
      warcraftLogsApiService.getSpecLeaderboard(className, spec, limit),
    enabled: !!(className && spec),
  });
};

// Hook for the dungeon leaderboard
export const useGetDungeonLeaderboard = (
  encounterID: number,
  page: number = 1,
  options?: {
    serverSlug?: string;
    serverRegion?: string;
    className?: WowClass;
    specName?: string;
  }
) => {
  return useQuery<DungeonLeaderboardResponse, Error>({
    queryKey: [
      "warcraftlogs-dungeon-leaderboard",
      encounterID,
      page,
      options?.serverRegion,
      options?.className,
      options?.specName,
      options?.serverSlug,
    ],
    queryFn: () =>
      warcraftLogsApiService.getDungeonLeaderboard(encounterID, page, options),
    enabled: !!encounterID, // Only fetch if encounterID is provided
  });
};

// Hook for player Mythic+ rankings
export const useGetPlayerMythicPlusRankings = (
  characterName: string,
  serverSlug: string,
  serverRegion: string,
  zoneID: number,
  queryOptions = {}
) => {
  return useQuery<MythicPlusPlayerRankings, Error>({
    queryKey: [
      "warcraftlogs-player-rankings",
      characterName,
      serverSlug,
      serverRegion,
      zoneID,
    ],
    queryFn: () =>
      warcraftLogsApiService.getPlayerMythicPlusRankings(
        characterName,
        serverSlug,
        serverRegion,
        zoneID
      ),
    enabled: !!(characterName && serverSlug && serverRegion && zoneID), // Only fetch if all required parameters are provided
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 30 * 60 * 1000, // 30 minutes
    retry: 2, // Retry 2 times
    ...queryOptions, // Allow additional query options
  });
};

// Hook for player raid rankings
export const useGetPlayerRaidRankings = (
  characterName: string,
  serverSlug: string,
  serverRegion: string,
  zoneID: number,
  queryOptions = {}
) => {
  return useQuery<RaidRankingsResponse, Error>({
    queryKey: [
      "warcraftlogs-player-raid-rankings",
      characterName,
      serverSlug,
      serverRegion,
      zoneID,
    ],
    queryFn: () =>
      warcraftLogsApiService.getPlayerRaidRankings(
        characterName,
        serverSlug,
        serverRegion,
        zoneID
      ),
    enabled: !!(characterName && serverSlug && serverRegion && zoneID), // Only fetch if all required parameters are provided
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 30 * 60 * 1000, // 30 minutes
    retry: 2, // Retry 2 times
    ...queryOptions, // Allow additional query options
  });
};
