import { useQuery } from "@tanstack/react-query";
import {
  Role,
  SpecByRole,
  TopTeamCompositionsGlobal,
  SpecByDungeonAndRole,
  TopTeamCompositionsByDungeon,
  MetaByKeyLevels,
  MetaByRegion,
  OverallStatsData,
  KeyLevelDistribution,
} from "../types/raiderio/mythicplus_runs/mythicPlusRuns";
import {
  getSpecByRole,
  getTopTeamCompositionsGlobal,
  getSpecByDungeonAndRole,
  getTopTeamCompositionsByDungeon,
  getMetaByKeyLevels,
  getMetaByRegion,
  getOverallStats,
  getKeyLevelDistribution,
} from "../libs/mythicPlusRunsAnalysisApiService";

/* === Analyses globales === */

const defaultQueryConfig = {
  staleTime: 1000 * 60 * 5, // 5 minutes
  gcTime: 1000 * 60 * 30, // 30 minutes
  retry: 2,
  refetchOnWindowFocus: false,
};

// useSpecsByRole - Récupère les stats par spécialisation et rôle
export const useSpecsByRole = (role: Role, topN?: number, options?: any) => {
  return useQuery({
    queryKey: ["raiderio", "specs", role, topN],
    queryFn: () => getSpecByRole(role, topN ? { top_n: topN } : undefined),
    ...defaultQueryConfig,
    enabled: !!role,
    ...options,
  });
};

// useTopTeamCompositionsGlobal - Récupère les compositions de team les plus utilisées
export const useTopTeamCompositionsGlobal = (
  params?: {
    limit?: number;
    min_usage?: number;
  },
  options?: any
) => {
  return useQuery({
    queryKey: ["raiderio", "top-team-compositions-global", params],
    queryFn: () => getTopTeamCompositionsGlobal(params),
    ...defaultQueryConfig,
    ...options,
  });
};

/* === Analyses par donjon === */

// useSpecsByDungeonAndRole - Récupère les stats par spécialisation et rôle pour un donjon
export const useSpecsByDungeonAndRole = (
  dungeonSlug: string,
  role: Role,
  topN?: number,
  options?: any
) => {
  return useQuery({
    queryKey: ["raiderio", "specs", dungeonSlug, role, topN],
    queryFn: () =>
      getSpecByDungeonAndRole(
        dungeonSlug,
        role,
        topN ? { top_n: topN } : undefined
      ),
    ...defaultQueryConfig,
    enabled: !!dungeonSlug && !!role,
    ...options,
  });
};

// useTopTeamCompositionsByDungeon - Récupère les compositions de team les plus utilisées pour un donjon
export const useTopTeamCompositionsByDungeon = (
  options?: any,
  params?: {
    top_n?: number;
    min_usage?: number;
  }
) => {
  return useQuery({
    queryKey: ["raiderio", "top-team-compositions-by-dungeon", params],
    queryFn: () => getTopTeamCompositionsByDungeon(params),
    ...defaultQueryConfig,
    ...options,
  });
};

/* === Analyses par niveau de clé === */

// useMetaByKeyLevels - Récupère les stats par spécialisation et niveau de clé
export const useMetaByKeyLevels = (min_usage?: number, options?: any) => {
  return useQuery({
    queryKey: ["raiderio", "meta-by-key-levels", min_usage],
    queryFn: () => getMetaByKeyLevels(min_usage),
    ...defaultQueryConfig,
    ...options,
  });
};

// useMetaByRegion - Récupère les stats par spécialisation et région
export const useMetaByRegion = (options?: any) => {
  return useQuery({
    queryKey: ["raiderio", "meta-by-region"],
    queryFn: () => getMetaByRegion(),
    ...defaultQueryConfig,
    ...options,
  });
};

/* === Analyses utilitaires === */

// useOverallStats - Récupère les stats globales
export const useOverallStats = (options?: any) => {
  return useQuery({
    queryKey: ["raiderio", "overall-stats"],
    queryFn: () => getOverallStats(),
    ...defaultQueryConfig,
    ...options,
  });
};

// useKeyLevelDistribution - Récupère la distribution des niveaux de clé
export const useKeyLevelDistribution = (options?: any) => {
  return useQuery({
    queryKey: ["raiderio", "key-level-distribution"],
    queryFn: () => getKeyLevelDistribution(),
    ...defaultQueryConfig,
    ...options,
  });
};
