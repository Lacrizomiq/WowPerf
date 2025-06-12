import { useQuery, UseQueryOptions } from "@tanstack/react-query";
import * as buildsAnalysisApiService from "../libs/buildsAnalysisApiService";
import {
  WowClassParam,
  WowSpecParam,
} from "../types/warcraftlogs/builds/classSpec";
import {
  PopularItemsResponse,
  GlobalPopularItemsResponse,
  EnchantUsageResponse,
  GemUsageResponse,
  TopTalentBuildsResponse,
  TalentBuildsByDungeonResponse,
  StatPrioritiesResponse,
  OptimalBuildResponse,
  ClassSpecSummaryResponse,
  SpecComparisonResponse,
} from "../types/warcraftlogs/builds/buildsAnalysis";

// Default query config for all builds analysis hooks
const defaultQueryConfig = {
  staleTime: 1000 * 60 * 5, // 5 minutes
  gcTime: 1000 * 60 * 30, // 30 minutes
  retry: 2,
};

// Hook to get popular items by class and spec for a specific encounter (Or all if null)
export const useGetPopularItems = (
  className: WowClassParam,
  spec: WowSpecParam,
  encounterId?: number | string,
  options?: UseQueryOptions<PopularItemsResponse, Error>
) => {
  return useQuery<PopularItemsResponse, Error>({
    queryKey: [
      "warcraftlogs-builds-popular-items",
      className,
      spec,
      encounterId,
    ],
    queryFn: () =>
      buildsAnalysisApiService.getPopularItems(className, spec, encounterId),
    enabled: !!(className && spec),
    ...defaultQueryConfig,
    ...options,
  });
};

// Hook to get global popular items by class and spec
export const useGetGlobalPopularItems = (
  className: WowClassParam,
  spec: WowSpecParam,
  options?: UseQueryOptions<GlobalPopularItemsResponse, Error>
) => {
  return useQuery<GlobalPopularItemsResponse, Error>({
    queryKey: ["warcraftlogs-builds-global-popular-items", className, spec],
    queryFn: () =>
      buildsAnalysisApiService.getGlobalPopularItems(className, spec),
    enabled: !!(className && spec),
    ...defaultQueryConfig,
    ...options,
  });
};

// Hook to get enchant usage by class and spec
export const useGetEnchantUsage = (
  className: WowClassParam,
  spec: WowSpecParam,
  options?: UseQueryOptions<EnchantUsageResponse, Error>
) => {
  return useQuery<EnchantUsageResponse, Error>({
    queryKey: ["warcraftlogs-builds-enchant-usage", className, spec],
    queryFn: () => buildsAnalysisApiService.getEnchantUsage(className, spec),
    enabled: !!(className && spec),
    ...defaultQueryConfig,
    ...options,
  });
};

// Hook to get gem usage by class and spec
export const useGetGemUsage = (
  className: WowClassParam,
  spec: WowSpecParam,
  options?: UseQueryOptions<GemUsageResponse, Error>
) => {
  return useQuery<GemUsageResponse, Error>({
    queryKey: ["warcraftlogs-builds-gem-usage", className, spec],
    queryFn: () => buildsAnalysisApiService.getGemUsage(className, spec),
    enabled: !!(className && spec),
    ...defaultQueryConfig,
    ...options,
  });
};

// Hook for top talent builds by class and spec
export const useGetTopTalentBuilds = (
  className: WowClassParam,
  spec: WowSpecParam,
  options?: UseQueryOptions<TopTalentBuildsResponse, Error>
) => {
  return useQuery<TopTalentBuildsResponse, Error>({
    queryKey: ["warcraftlogs-builds-top-talents", className, spec],
    queryFn: () => buildsAnalysisApiService.getTopTalentBuilds(className, spec),
    enabled: !!(className && spec),
    ...defaultQueryConfig,
    ...options,
  });
};

// Hook for talent builds by dungeon for a specific class and spec
export const useGetTalentBuildsByDungeon = (
  className: WowClassParam,
  spec: WowSpecParam,
  options?: UseQueryOptions<TalentBuildsByDungeonResponse, Error>
) => {
  return useQuery<TalentBuildsByDungeonResponse, Error>({
    queryKey: ["warcraftlogs-builds-talents-by-dungeon", className, spec],
    queryFn: () =>
      buildsAnalysisApiService.getTalentBuildsByDungeon(className, spec),
    enabled: !!(className && spec),
    ...defaultQueryConfig,
    ...options,
  });
};

// Hook for stat priorities by class and spec
export const useGetStatPriorities = (
  className: WowClassParam,
  spec: WowSpecParam,
  options?: UseQueryOptions<StatPrioritiesResponse, Error>
) => {
  return useQuery<StatPrioritiesResponse, Error>({
    queryKey: ["warcraftlogs-builds-stat-priorities", className, spec],
    queryFn: () => buildsAnalysisApiService.getStatPriorities(className, spec),
    enabled: !!(className && spec),
    ...defaultQueryConfig,
    ...options,
  });
};

// Hook for optimal build by class and spec
export const useGetOptimalBuild = (
  className: WowClassParam,
  spec: WowSpecParam,
  options?: UseQueryOptions<OptimalBuildResponse, Error>
) => {
  return useQuery<OptimalBuildResponse, Error>({
    queryKey: ["warcraftlogs-builds-optimal", className, spec],
    queryFn: () => buildsAnalysisApiService.getOptimalBuild(className, spec),
    enabled: !!(className && spec),
    ...defaultQueryConfig,
    ...options,
  });
};

// Hook for class and spec summary by class and spec
export const useGetClassSpecSummary = (
  className: WowClassParam,
  spec: WowSpecParam,
  options?: UseQueryOptions<ClassSpecSummaryResponse, Error>
) => {
  return useQuery<ClassSpecSummaryResponse, Error>({
    queryKey: ["warcraftlogs-builds-class-spec-summary", className, spec],
    queryFn: () =>
      buildsAnalysisApiService.getClassSpecSummary(className, spec),
    enabled: !!(className && spec),
    ...defaultQueryConfig,
    ...options,
  });
};

// Hook for spec comparison by class
export const useGetSpecComparison = (
  className: WowClassParam,
  options?: UseQueryOptions<SpecComparisonResponse, Error>
) => {
  return useQuery<SpecComparisonResponse, Error>({
    queryKey: ["warcraftlogs-builds-spec-comparison", className],
    queryFn: () => buildsAnalysisApiService.getSpecComparison(className),
    enabled: !!className,
    ...defaultQueryConfig,
    ...options,
  });
};
