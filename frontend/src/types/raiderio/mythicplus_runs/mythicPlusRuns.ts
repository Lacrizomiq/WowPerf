// frontend/src/types/raiderio/mythicplus_runs/mythicPlusRuns.ts

/**
 * Énumération des rôles disponibles dans WoW Mythic+
 */
export enum Role {
  TANK = "tank",
  HEALER = "healer",
  DPS = "dps",
}

/**
 * Représente les statistiques d'une spécialisation par rôle
 * Utilisé pour les endpoints :
 * - /api/raiderio/mythicplus/analytics/specs/{role} - tank, healer, dps
 */

export interface SpecByRole {
  class: string;
  spec: string;
  display: string;
  usage_count: number;
  percentage: number;
  rank: number;
}

/**
 * Représente les compositions de team les plus utilisées
 * Utilisé pour l'endpoint :
 * - /api/raiderio/mythicplus/analytics/compositions
 */
export interface TopTeamCompositionsGlobal {
  tank: string;
  healer: string;
  dps1: string;
  dps2: string;
  dps3: string;
  usage_count: number;
  percentage: number;
  rank: number;
  avg_score: number;
}

/**
 * Représente les statistiques d'une spécialisation par donjon et rôle
 * Utilisé pour les endpoints :
 * - /api/raiderio/mythicplus/analytics/dungeons/{dungeon_slug}/specs/{role}
 */
export interface SpecByDungeonAndRole {
  class: string;
  spec: string;
  display: string;
  dungeon_slug: string;
  dungeon_name: string;
  usage_count: number;
  percentage: number;
  rank_in_dungeon: number;
}

/**
 * Représente les compositions de team les plus utilisées par donjon
 * Utilisé pour l'endpoint :
 * - api/raiderio/mythicplus/analytics/dungeons/compositions?top_n=5&min_usage=10
 */
export interface TopTeamCompositionsByDungeon {
  tank: string;
  healer: string;
  dps1: string;
  dps2: string;
  dps3: string;
  dungeon_slug: string;
  dungeon_name: string;
  usage_count: number;
  percentage: number;
  rank_in_dungeon: number;
  avg_score: number;
}

/**
 * Représente les statistiques de spécialisations par niveau de clé
 * Utilisé pour l'endpoint :
 * - /api/raiderio/mythicplus/analytics/key-levels
 */
export interface MetaByKeyLevels {
  role: string;
  key_level_bracket: string;
  class: string;
  spec: string;
  display: string;
  usage_count: number;
  percentage: number;
  rank: number;
  avg_score: number;
}

/**
 * Représente les statistiques de spécialisations par région
 * Utilisé pour l'endpoint :
 * - /api/raiderio/mythicplus/analytics/regions
 */
export interface MetaByRegion {
  role: string;
  region: string;
  class: string;
  spec: string;
  display: string;
  usage_count: number;
  percentage_in_region: number;
  rank_in_region: number;
}

/**
 * Représente les statistiques générales du dataset
 * Utilisé pour l'endpoint :
 * - /api/raiderio/mythicplus/analytics/stats/overall
 */
export interface OverallStatsData {
  total_runs: number;
  runs_with_score: number;
  unique_compositions: number;
  unique_dungeons: number;
  unique_regions: number;
  oldest_run: string;
  newest_run: string;
  avg_score: number;
  avg_key_level: number;
}

/**
 * Représente la distribution des niveaux de clés
 * Utilisé pour l'endpoint :
 * - /api/raiderio/mythicplus/analytics/stats/key-levels
 */
export interface KeyLevelDistribution {
  mythic_level: number;
  count: number;
  percentage: number;
  rank: number;
  avg_score: number;
  min_score: number;
  max_score: number;
}
