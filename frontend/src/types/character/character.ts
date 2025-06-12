// ============================================================================
// TYPES POUR LE SYSTÈME DE PERSONNAGES ENRICHI
// ============================================================================

/**
 * Personnage enrichi stocké en base de données
 */
export interface EnrichedUserCharacter {
  id: number;
  user_id: number;
  character_id: number;

  // Infos de base (sync)
  name: string;
  realm: string;
  region: string;
  class: string;
  race: string;
  faction: string;
  level: number;

  // Infos enrichies (enrichisseurs)
  gender?: string;
  active_spec_name?: string;
  active_spec_role?: string;
  active_spec_id?: number;
  achievement_points?: number;
  honorable_kills?: number;
  avatar_url?: string;
  inset_avatar_url?: string;
  main_raw_url?: string;
  profile_url?: string;

  // Données JSON structurées (pour futurs enrichisseurs)
  equipment_json?: any; // Équipements détaillés
  mythic_plus_json?: any; // Données M+
  raids_json?: any; // Progression raids
  stats_json?: any; // Statistiques diverses
  talents_json?: any; // Build talents

  // Métadonnées
  is_displayed: boolean;
  last_api_update: string; // ISO date string
}

/**
 * Résultat des opérations sync et enrichissement
 */
export interface SyncAndEnrichResult {
  message: string;
  result: {
    synced_count: number;
    enriched_count: number;
    updated_count: number;
    errors: string[];
  };
}

/**
 * Résultat des opérations refresh et enrichissement
 */
export interface RefreshAndEnrichResult {
  message: string;
  result: {
    synced_count: number;
    enriched_count: number;
    updated_count: number;
    errors: string[];
  };
}

/**
 * Réponse de l'API GET /characters
 */
export interface GetCharactersResponse {
  characters: EnrichedUserCharacter[];
  count: number;
}

/**
 * Erreur de rate limiting avec détails
 */
export interface RateLimitError {
  error: string;
  type: "rate_limit";
  wait_time?: string; // Ex: "3m" ou "45s"
}

/**
 * Codes d'erreur spécifiques au système de personnages
 */
export enum CharacterErrorCode {
  RATE_LIMIT = "RATE_LIMIT",
  UNAUTHORIZED = "UNAUTHORIZED",
  NOT_FOUND = "NOT_FOUND",
  FORBIDDEN = "FORBIDDEN",
  SERVER_ERROR = "SERVER_ERROR",
  NETWORK_ERROR = "NETWORK_ERROR",
}

/**
 * Erreur personnalisée pour le système de personnages
 */
export class CharacterError extends Error {
  constructor(
    public code: CharacterErrorCode,
    message: string,
    public originalError?: unknown,
    public waitTime?: string
  ) {
    super(message);
    this.name = "CharacterError";
  }
}

/**
 * Configuration pour les timers de rate limiting
 */
export interface RateLimitState {
  isRateLimited: boolean;
  expiryTime: Date | null;
  timeRemaining: number;
  waitTime?: string;
}

/**
 * État d'un personnage (pour UI)
 */
export interface CharacterUIState {
  isRecentlyUpdated: boolean;
  isEnriched: boolean;
  avatarUrl: string;
  needsUpdate: boolean;
}
