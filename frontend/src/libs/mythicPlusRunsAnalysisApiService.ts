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

// ========================================
// INTERFACES POUR LES PARAMÈTRES
// ========================================

/**
 * Paramètres pour récupérer les spécialisations par rôle
 */
interface GetSpecByRoleParams {
  /** Nombre de spécialisations à retourner (0 = toutes) */
  top_n?: number;
}

/**
 * Paramètres pour récupérer les compositions d'équipes globales
 */
interface GetTopTeamCompositionsParams {
  /** Nombre maximum de compositions à retourner */
  limit?: number;
  /** Nombre minimum d'utilisations pour filtrer les résultats */
  min_usage?: number;
}

/**
 * Paramètres pour récupérer les spécialisations par donjon et rôle
 */
interface GetSpecByDungeonAndRoleParams {
  /** Nombre de spécialisations à retourner par donjon (0 = toutes) */
  top_n?: number;
}

/**
 * Paramètres pour récupérer les compositions par donjon
 */
interface GetTopTeamCompositionsByDungeonParams {
  /** Nombre de compositions à retourner par donjon (0 = toutes) */
  top_n?: number;
  /** Nombre minimum d'utilisations pour filtrer les résultats */
  min_usage?: number;
}

/**
 * Paramètres pour récupérer les métadonnées par niveau de clé
 */
interface GetMetaByKeyLevelsParams {
  /** Nombre minimum d'utilisations pour filtrer les résultats */
  min_usage?: number;
  /** Nombre maximum de résultats à retourner (0 = toutes) */
  top_n?: number;
}

/**
 * Paramètres pour récupérer les métadonnées par région
 */
interface GetMetaByRegionParams {
  /** Nombre de spécialisations à retourner par région et par rôle (0 = toutes) */
  top_n?: number;
}

// ========================================
// ANALYSES GLOBALES
// ========================================

/**
 * Récupère les statistiques des spécialisations par rôle à travers tous les donjons
 *
 * @param role - Le rôle à analyser (tank, healer, dps)
 * @param params - Paramètres optionnels de la requête
 * @param params.top_n - Limite le nombre de résultats retournés (0 = tous)
 *
 * @returns Promise<SpecByRole[]> - Liste des spécialisations avec leurs statistiques
 */
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
    console.error(`Error fetching ${role} specializations:`, error);
    throw error;
  }
};

/**
 * Récupère les compositions d'équipes les plus populaires globalement
 *
 * @param params - Paramètres optionnels de la requête
 * @param params.limit - Nombre maximum de compositions à retourner (défaut: 20)
 * @param params.min_usage - Filtre les compositions avec moins d'utilisations (défaut: 5)
 *
 * @returns Promise<TopTeamCompositionsGlobal[]> - Liste des compositions populaires
 *
 */
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

// ========================================
// ANALYSES PAR DONJON
// ========================================

/**
 * Récupère les statistiques des spécialisations pour un donjon et rôle spécifiques
 *
 * @param dungeonSlug - Identifiant unique du donjon (ex: "mists-of-tirna-scithe")
 * @param role - Le rôle à analyser (tank, healer, dps)
 * @param params - Paramètres optionnels de la requête
 * @param params.top_n - Limite le nombre de résultats par donjon (0 = tous)
 *
 * @returns Promise<SpecByDungeonAndRole[]> - Spécialisations avec stats par donjon
 *
 */
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
    console.error(
      `Error fetching ${role} specs for dungeon ${dungeonSlug}:`,
      error
    );
    throw error;
  }
};

/**
 * Récupère les compositions d'équipes les plus utilisées groupées par donjon
 *
 * @param params - Paramètres optionnels de la requête
 * @param params.top_n - Nombre de compositions par donjon (0 = toutes)
 * @param params.min_usage - Filtre les compositions avec moins d'utilisations (défaut: 3)
 *
 * @returns Promise<TopTeamCompositionsByDungeon[]> - Compositions par donjon
 *
 */
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

// ========================================
// ANALYSES PAR NIVEAU DE CLÉ
// ========================================

/**
 * Récupère les statistiques des spécialisations groupées par niveau de clé
 * Utilise des brackets prédéfinis : Very High Keys (20+), High Keys (18-19), Mid Keys (16-17)
 *
 * @param params - Paramètres optionnels de la requête
 * @param params.top_n - Limite le nombre de spécialisations par bracket et par rôle (0 = toutes)
 * @param params.min_usage - Filtre les spécialisations avec moins d'utilisations (défaut: 5)
 *
 * @returns Promise<MetaByKeyLevels[]> - Métadonnées par niveau de clé
 *
 */
export const getMetaByKeyLevels = async (
  params?: GetMetaByKeyLevelsParams
): Promise<MetaByKeyLevels[]> => {
  try {
    const { data } = await api.get<MetaByKeyLevels[]>(
      "raiderio/mythicplus/analytics/key-levels",
      { params }
    );
    return data;
  } catch (error) {
    console.error("Error fetching meta by key levels:", error);
    throw error;
  }
};

// ========================================
// ANALYSES PAR RÉGION
// ========================================

/**
 * Récupère les statistiques des spécialisations groupées par région
 * Couvre les régions : US, EU, KR, TW
 *
 * @param params - Paramètres optionnels de la requête
 * @param params.top_n - Limite le nombre de spécialisations par région et par rôle (0 = toutes)
 *
 * @returns Promise<MetaByRegion[]> - Métadonnées par région
 *
 */
export const getMetaByRegion = async (
  params?: GetMetaByRegionParams
): Promise<MetaByRegion[]> => {
  try {
    const { data } = await api.get<MetaByRegion[]>(
      "raiderio/mythicplus/analytics/regions",
      { params }
    );
    return data;
  } catch (error) {
    console.error("Error fetching meta by region:", error);
    throw error;
  }
};

// ========================================
// ANALYSES UTILITAIRES
// ========================================

/**
 * Récupère les statistiques générales du dataset Mythic+
 * Inclut : nombre total de runs, scores moyens, compositions uniques, etc.
 *
 * @returns Promise<OverallStatsData> - Statistiques générales
 *
 */
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

/**
 * Récupère la distribution des runs par niveau de clé Mythic+
 * Utile pour comprendre la répartition des joueurs par difficulté
 *
 * @returns Promise<KeyLevelDistribution[]> - Distribution des niveaux de clé
 *
 */
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
