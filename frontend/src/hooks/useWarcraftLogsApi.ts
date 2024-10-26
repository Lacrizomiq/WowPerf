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
  page: number = 1
) => {
  return useQuery<DungeonLeaderboardResponse, Error>({
    queryKey: ["warcraftlogs-dungeon-leaderboard", encounterID, page],
    queryFn: () =>
      warcraftLogsApiService.getDungeonLeaderboard(encounterID, page),
    enabled: !!encounterID, // Only fetch if encounterID is provided
  });
};
