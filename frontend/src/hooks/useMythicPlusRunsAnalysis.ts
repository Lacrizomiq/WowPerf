import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
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

// ========================================
// TYPES POUR LES OPTIONS REACT QUERY
// ========================================

/**
 * Type générique pour les options React Query personnalisées
 * Exclut queryKey et queryFn qui sont gérés automatiquement
 */
type QueryOptions<TData> = Omit<
  UseQueryOptions<TData, Error, TData, (string | number | undefined)[]>,
  "queryKey" | "queryFn"
>;

// ========================================
// CONFIGURATION PAR DÉFAUT
// ========================================

/**
 * Configuration React Query par défaut pour les endpoints Mythic+
 */
const defaultQueryConfig = {
  staleTime: 1000 * 60 * 5, // 5 minutes - données considérées fraîches
  gcTime: 1000 * 60 * 30, // 30 minutes - garde en cache
  retry: 2, // 2 tentatives en cas d'échec réseau
  refetchOnWindowFocus: false, // Évite les refetch intempestifs en développement
};

// ========================================
// ANALYSES GLOBALES
// ========================================

/**
 * Hook pour récupérer les statistiques des spécialisations par rôle
 *
 * @param role - Le rôle à analyser (tank, healer, dps)
 * @param params - Paramètres optionnels de la requête
 * @param params.top_n - Limite le nombre de résultats (0 = tous)
 * @param options - Options React Query personnalisées
 *
 * @returns Résultat de la query avec data, loading, error, etc.
 *
 * @example
 * ```typescript
 * // Récupère toutes les spécialisations tank
 * const { data: tanks, isLoading } = useSpecsByRole('tank');
 *
 * // Top 5 DPS avec options personnalisées
 * const { data: topDPS } = useSpecsByRole('dps', { top_n: 5 });
 * ```
 */
export const useSpecsByRole = (
  role: Role,
  params?: { top_n?: number },
  options?: QueryOptions<SpecByRole[]>
): UseQueryResult<SpecByRole[], Error> => {
  return useQuery({
    queryKey: ["raiderio", "mythicplus", "specs", "role", role, params?.top_n],
    queryFn: () => getSpecByRole(role, params),
    enabled: !!role,
    ...defaultQueryConfig,
    ...options,
  });
};

/**
 * Hook pour récupérer les compositions d'équipes les plus populaires globalement
 *
 * @param params - Paramètres optionnels de la requête
 * @param params.limit - Nombre maximum de compositions (défaut: 20)
 * @param params.min_usage - Filtre les compositions avec moins d'utilisations
 * @param options - Options React Query personnalisées
 *
 * @returns Résultat de la query avec les compositions populaires
 *
 * @example
 * ```typescript
 * // Top 10 compositions avec au moins 5 utilisations
 * const { data: compositions } = useTopTeamCompositionsGlobal({
 *   limit: 10,
 *   min_usage: 5
 * });
 *
 * // Avec gestion d'erreur personnalisée
 * const { data, error, refetch } = useTopTeamCompositionsGlobal(undefined, {
 *   onError: (error) => toast.error(`Erreur: ${error.message}`)
 * });
 * ```
 */
export const useTopTeamCompositionsGlobal = (
  params?: { limit?: number; min_usage?: number },
  options?: QueryOptions<TopTeamCompositionsGlobal[]>
): UseQueryResult<TopTeamCompositionsGlobal[], Error> => {
  return useQuery({
    queryKey: [
      "raiderio",
      "mythicplus",
      "compositions",
      "global",
      params?.limit,
      params?.min_usage,
    ],
    queryFn: () => getTopTeamCompositionsGlobal(params),
    ...defaultQueryConfig,
    ...options,
  });
};

// ========================================
// ANALYSES PAR DONJON
// ========================================

/**
 * Hook pour récupérer les spécialisations pour un donjon et rôle spécifiques
 *
 * @param dungeonSlug - Identifiant unique du donjon
 * @param role - Le rôle à analyser (tank, healer, dps)
 * @param params - Paramètres optionnels de la requête
 * @param params.top_n - Nombre de spécialisations par donjon (0 = toutes)
 * @param options - Options React Query personnalisées
 *
 * @returns Résultat de la query avec les spécialisations par donjon
 *
 * @example
 * ```typescript
 * // Top 3 tanks pour un donjon spécifique
 * const { data: tankSpecs } = useSpecsByDungeonAndRole(
 *   'mists-of-tirna-scithe',
 *   'tank',
 *   { top_n: 3 }
 * );
 *
 * // Avec query conditionnelle
 * const { data } = useSpecsByDungeonAndRole(
 *   selectedDungeon,
 *   'healer',
 *   undefined,
 *   { enabled: !!selectedDungeon } // Ne s'exécute que si un donjon est sélectionné
 * );
 * ```
 */
export const useSpecsByDungeonAndRole = (
  dungeonSlug: string,
  role: Role,
  params?: { top_n?: number },
  options?: QueryOptions<SpecByDungeonAndRole[]>
): UseQueryResult<SpecByDungeonAndRole[], Error> => {
  return useQuery({
    queryKey: [
      "raiderio",
      "mythicplus",
      "specs",
      "dungeon",
      dungeonSlug,
      role,
      params?.top_n,
    ],
    queryFn: () => getSpecByDungeonAndRole(dungeonSlug, role, params),
    enabled: !!dungeonSlug && !!role,
    ...defaultQueryConfig,
    ...options,
  });
};

/**
 * Hook pour récupérer les compositions d'équipes groupées par donjon
 *
 * @param params - Paramètres optionnels de la requête
 * @param params.top_n - Nombre de compositions par donjon (0 = toutes)
 * @param params.min_usage - Filtre les compositions avec moins d'utilisations
 * @param options - Options React Query personnalisées
 *
 * @returns Résultat de la query avec les compositions par donjon
 *
 * @example
 * ```typescript
 * // Top 5 compositions par donjon
 * const { data: dungeonComps } = useTopTeamCompositionsByDungeon({
 *   top_n: 5,
 *   min_usage: 10
 * });
 *
 * // Filtrer pour un donjon spécifique après récupération
 * const mistComps = dungeonComps?.filter(
 *   comp => comp.dungeon_slug === 'mists-of-tirna-scithe'
 * );
 * ```
 */
export const useTopTeamCompositionsByDungeon = (
  params?: { top_n?: number; min_usage?: number },
  options?: QueryOptions<TopTeamCompositionsByDungeon[]>
): UseQueryResult<TopTeamCompositionsByDungeon[], Error> => {
  return useQuery({
    queryKey: [
      "raiderio",
      "mythicplus",
      "compositions",
      "dungeon",
      params?.top_n,
      params?.min_usage,
    ],
    queryFn: () => getTopTeamCompositionsByDungeon(params),
    ...defaultQueryConfig,
    ...options,
  });
};

// ========================================
// ANALYSES PAR NIVEAU DE CLÉ
// ========================================

/**
 * Hook pour récupérer les métadonnées par niveau de clé
 * Organise les données par brackets : Very High Keys (20+), High Keys (18-19), Mid Keys (16-17)
 *
 * @param params - Paramètres optionnels de la requête
 * @param params.min_usage - Filtre les spécialisations avec moins d'utilisations
 * @param options - Options React Query personnalisées
 *
 * @returns Résultat de la query avec les métadonnées par niveau
 *
 * @example
 * ```typescript
 * // Récupère la méta high keys
 * const { data: keyMeta } = useMetaByKeyLevels({ min_usage: 5 });
 *
 * // Filtrer pour les very high keys uniquement
 * const veryHighKeys = keyMeta?.filter(
 *   spec => spec.key_level_bracket === 'Very High Keys (20+)'
 * );
 *
 * // Grouper par rôle
 * const metaByRole = keyMeta?.reduce((acc, spec) => {
 *   if (!acc[spec.role]) acc[spec.role] = [];
 *   acc[spec.role].push(spec);
 *   return acc;
 * }, {} as Record<string, typeof keyMeta>);
 * ```
 */
export const useMetaByKeyLevels = (
  params?: { min_usage?: number },
  options?: QueryOptions<MetaByKeyLevels[]>
): UseQueryResult<MetaByKeyLevels[], Error> => {
  return useQuery({
    queryKey: [
      "raiderio",
      "mythicplus",
      "meta",
      "key-levels",
      params?.min_usage,
    ],
    queryFn: () => getMetaByKeyLevels(params),
    ...defaultQueryConfig,
    ...options,
  });
};

// ========================================
// ANALYSES PAR RÉGION
// ========================================

/**
 * Hook pour récupérer les métadonnées par région
 * Couvre les régions : US, EU, KR, TW
 *
 * @param options - Options React Query personnalisées
 *
 * @returns Résultat de la query avec les métadonnées par région
 *
 * @example
 * ```typescript
 * // Récupère la méta par région
 * const { data: regionMeta } = useMetaByRegion();
 *
 * // Filtrer pour l'Europe
 * const euMeta = regionMeta?.filter(spec => spec.region === 'EU');
 *
 * // Comparer les régions pour un rôle spécifique
 * const tankByRegion = regionMeta?.filter(spec => spec.role === 'Tank');
 * ```
 */
export const useMetaByRegion = (
  options?: QueryOptions<MetaByRegion[]>
): UseQueryResult<MetaByRegion[], Error> => {
  return useQuery({
    queryKey: ["raiderio", "mythicplus", "meta", "region"],
    queryFn: () => getMetaByRegion(),
    ...defaultQueryConfig,
    ...options,
  });
};

// ========================================
// ANALYSES UTILITAIRES
// ========================================

/**
 * Hook pour récupérer les statistiques générales du dataset
 * Inclut : total runs, scores moyens, compositions uniques, périodes couvertes
 *
 * @param options - Options React Query personnalisées
 *
 * @returns Résultat de la query avec les statistiques générales
 *
 * @example
 * ```typescript
 * // Statistiques générales pour un dashboard
 * const { data: stats, isLoading } = useOverallStats();
 *
 * if (!isLoading && stats) {
 *   console.log(`Analyse de ${stats.total_runs.toLocaleString()} runs`);
 *   console.log(`Score moyen: ${stats.avg_score}`);
 *   console.log(`${stats.unique_compositions} compositions uniques`);
 * }
 *
 * // Avec refresh automatique
 * const { data } = useOverallStats({
 *   refetchInterval: 60000 // Actualise chaque minute
 * });
 * ```
 */
export const useOverallStats = (
  options?: QueryOptions<OverallStatsData>
): UseQueryResult<OverallStatsData, Error> => {
  return useQuery({
    queryKey: ["raiderio", "mythicplus", "stats", "overall"],
    queryFn: () => getOverallStats(),
    ...defaultQueryConfig,
    ...options,
  });
};

/**
 * Hook pour récupérer la distribution des runs par niveau de clé
 * Utile pour les graphiques et analyses de popularité par difficulté
 *
 * @param options - Options React Query personnalisées
 *
 * @returns Résultat de la query avec la distribution des niveaux
 *
 * @example
 * ```typescript
 * // Distribution pour un graphique
 * const { data: distribution } = useKeyLevelDistribution();
 *
 * // Trouver le niveau le plus populaire
 * const mostPopular = distribution?.reduce((prev, current) =>
 *   current.count > prev.count ? current : prev
 * );
 *
 * // Préparer les données pour Chart.js
 * const chartData = distribution?.map(level => ({
 *   x: level.mythic_level,
 *   y: level.count,
 *   label: `+${level.mythic_level} (${level.percentage}%)`
 * }));
 * ```
 */
export const useKeyLevelDistribution = (
  options?: QueryOptions<KeyLevelDistribution[]>
): UseQueryResult<KeyLevelDistribution[], Error> => {
  return useQuery({
    queryKey: ["raiderio", "mythicplus", "stats", "key-level-distribution"],
    queryFn: () => getKeyLevelDistribution(),
    ...defaultQueryConfig,
    ...options,
  });
};
