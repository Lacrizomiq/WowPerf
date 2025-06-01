// frontend/src/libs/mythicPlusRunsAnalysisApiService.ts

import api from "./api";
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

/* === Analyses globales === */

interface GetSpecByRoleParams {
  top_n?: number;
}

// GetSpecByRole - Récupère les stats par spécialisation et rôle
export const getSpecByRole = async (
  role: Role,
  params?: GetSpecByRoleParams
): Promise<SpecByRole[]> => {
  try {
    const { data } = await api.get<SpecByRole[]>(
      `raiderio/mythicplus/analytics/specs/${role}`,
      { params }
    );
    return data;
  } catch (error) {
    console.error("Error fetching spec by role:", error);
    throw error;
  }
};

interface GetTopTeamCompositionsParams {
  limit?: number;
  min_usage?: number;
}

// GetTopTeamCompositionsGlobal - Récupère les compositions de team les plus utilisées
export const getTopTeamCompositionsGlobal = async (
  params?: GetTopTeamCompositionsParams
): Promise<TopTeamCompositionsGlobal[]> => {
  try {
    const { data } = await api.get<TopTeamCompositionsGlobal[]>(
      "raiderio/mythicplus/analytics/compositions",
      { params }
    );
    return data;
  } catch (error) {
    console.error("Error fetching top team compositions:", error);
    throw error;
  }
};

/* === Analyses par donjon === */

interface GetSpecByDungeonAndRoleParams {
  top_n?: number;
}

// GetSpecByDungeonAndRole - Récupère les stats par spécialisation et rôle pour un donjon
export const getSpecByDungeonAndRole = async (
  dungeonSlug: string,
  role: Role,
  params?: GetSpecByDungeonAndRoleParams
): Promise<SpecByDungeonAndRole[]> => {
  try {
    const { data } = await api.get<SpecByDungeonAndRole[]>(
      `raiderio/mythicplus/analytics/dungeons/${dungeonSlug}/specs/${role}`,
      { params }
    );
    return data;
  } catch (error) {
    console.error("Error fetching spec by dungeon and role:", error);
    throw error;
  }
};

interface GetTopTeamCompositionsByDungeonParams {
  top_n?: number;
  min_usage?: number;
}

// GetTopTeamCompositionsByDungeon - Récupère les compositions de team les plus utilisées par donjon
export const getTopTeamCompositionsByDungeon = async (
  params?: GetTopTeamCompositionsByDungeonParams
): Promise<TopTeamCompositionsByDungeon[]> => {
  try {
    const { data } = await api.get<TopTeamCompositionsByDungeon[]>(
      "raiderio/mythicplus/analytics/dungeons/compositions",
      { params }
    );
    return data;
  } catch (error) {
    console.error("Error fetching top team compositions by dungeon:", error);
    throw error;
  }
};

/* === Analyses par niveau de clé === */

// GetMetaByKeyLevels - Récupère les stats par spécialisation et niveau de clé
export const getMetaByKeyLevels = async (
  min_usage?: number
): Promise<MetaByKeyLevels[]> => {
  try {
    const { data } = await api.get<MetaByKeyLevels[]>(
      "raiderio/mythicplus/analytics/key-levels",
      { params: { min_usage } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching meta by key levels:", error);
    throw error;
  }
};

/* === Analyses par région === */

// GetMetaByRegion - Récupère les stats par spécialisation et région
export const getMetaByRegion = async (): Promise<MetaByRegion[]> => {
  try {
    const { data } = await api.get<MetaByRegion[]>(
      "raiderio/mythicplus/analytics/regions"
    );
    return data;
  } catch (error) {
    console.error("Error fetching meta by region:", error);
    throw error;
  }
};

/* === Analyses utilitaires === */

// GetOverallStats - Récupère les statistiques globales
export const getOverallStats = async (): Promise<OverallStatsData> => {
  try {
    const { data } = await api.get<OverallStatsData>(
      "raiderio/mythicplus/analytics/stats/overall"
    );
    return data;
  } catch (error) {
    console.error("Error fetching overall stats:", error);
    throw error;
  }
};

// GetKeyLevelDistribution - Récupère la distribution des niveaux de clé
export const getKeyLevelDistribution = async (): Promise<
  KeyLevelDistribution[]
> => {
  try {
    const { data } = await api.get<KeyLevelDistribution[]>(
      "raiderio/mythicplus/analytics/stats/key-levels"
    );
    return data;
  } catch (error) {
    console.error("Error fetching key level distribution:", error);
    throw error;
  }
};
