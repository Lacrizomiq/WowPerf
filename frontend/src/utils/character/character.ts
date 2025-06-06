// frontend/src/utils/character/character.ts

import {
  EnrichedUserCharacter,
  RateLimitState,
} from "@/types/character/character";

// ============================================================================
// UTILITAIRES RATE LIMITING
// ============================================================================

/**
 * Extrait le temps d'attente depuis un message d'erreur rate limit
 * Ex: "Please wait 3m before next sync" → "3m"
 */
export function extractWaitTime(errorMessage: string): string | undefined {
  const match = errorMessage.match(/wait (\d+[ms])/);
  return match?.[1];
}

/**
 * Convertit un temps d'attente en millisecondes
 * Ex: "3m" → 180000, "45s" → 45000
 */
export function parseWaitTime(waitTime: string): number {
  const match = waitTime.match(/^(\d+)([ms])$/);
  if (!match) return 0;

  const value = parseInt(match[1]);
  const unit = match[2];

  switch (unit) {
    case "m":
      return value * 60 * 1000;
    case "s":
      return value * 1000;
    default:
      return 0;
  }
}

/**
 * Formate un temps restant en format "2m 34s"
 */
export function formatTimeRemaining(milliseconds: number): string {
  const totalSeconds = Math.ceil(milliseconds / 1000);
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;

  if (minutes > 0) {
    return `${minutes}m ${seconds}s`;
  }
  return `${seconds}s`;
}

/**
 * Calcule l'état du rate limiting à partir d'un temps d'attente
 */
export function calculateRateLimitState(waitTime?: string): RateLimitState {
  if (!waitTime) {
    return {
      isRateLimited: false,
      expiryTime: null,
      timeRemaining: 0,
    };
  }

  const waitMilliseconds = parseWaitTime(waitTime);
  const expiryTime = new Date(Date.now() + waitMilliseconds);

  return {
    isRateLimited: true,
    expiryTime,
    timeRemaining: waitMilliseconds,
    waitTime,
  };
}

/**
 * Vérifie si le rate limiting est encore actif
 */
export function isRateLimitActive(rateLimitState: RateLimitState): boolean {
  if (!rateLimitState.isRateLimited || !rateLimitState.expiryTime) {
    return false;
  }

  return Date.now() < rateLimitState.expiryTime.getTime();
}

/**
 * Calcule le temps restant avant la fin du rate limiting
 */
export function getRemainingTime(rateLimitState: RateLimitState): number {
  if (!rateLimitState.expiryTime) return 0;

  const remaining = rateLimitState.expiryTime.getTime() - Date.now();
  return Math.max(0, remaining);
}

// ============================================================================
// UTILITAIRES PERSONNAGES
// ============================================================================

/**
 * Vérifie si un personnage a été récemment mis à jour
 */
export function isCharacterRecentlyUpdated(
  character: EnrichedUserCharacter,
  maxAgeHours: number = 24
): boolean {
  if (!character.last_api_update) return false;

  const lastUpdate = new Date(character.last_api_update);
  const maxAge = maxAgeHours * 60 * 60 * 1000; // Convertir en millisecondes

  return Date.now() - lastUpdate.getTime() < maxAge;
}

/**
 * Vérifie si un personnage a des données enrichies
 */
export function isCharacterEnriched(character: EnrichedUserCharacter): boolean {
  return !!(
    character.active_spec_name ||
    character.achievement_points ||
    character.avatar_url
  );
}

/**
 * Obtient l'URL de l'avatar avec fallback
 */
export function getCharacterAvatarUrl(
  character: EnrichedUserCharacter
): string {
  return (
    character.avatar_url ||
    character.inset_avatar_url ||
    "/default-character-avatar.png"
  );
}

/**
 * Vérifie si un personnage a besoin d'une mise à jour
 */
export function characterNeedsUpdate(
  character: EnrichedUserCharacter,
  maxAgeHours: number = 24
): boolean {
  return (
    !isCharacterRecentlyUpdated(character, maxAgeHours) ||
    !isCharacterEnriched(character)
  );
}

/**
 * Trie les personnages par priorité (favoris, niveau, dernière update)
 */
export function sortCharactersByPriority(
  characters: EnrichedUserCharacter[]
): EnrichedUserCharacter[] {
  return [...characters].sort((a, b) => {
    // 1. Favoris en premier (si tu as ce champ)
    // if (a.is_favorite !== b.is_favorite) {
    //   return b.is_favorite ? 1 : -1;
    // }

    // 2. Niveau décroissant
    if (a.level !== b.level) {
      return b.level - a.level;
    }

    // 3. Dernière mise à jour (plus récent en premier)
    const aUpdate = new Date(a.last_api_update || 0).getTime();
    const bUpdate = new Date(b.last_api_update || 0).getTime();

    return bUpdate - aUpdate;
  });
}

/**
 * Filtre les personnages affichables
 */
export function getDisplayableCharacters(
  characters: EnrichedUserCharacter[]
): EnrichedUserCharacter[] {
  return characters.filter((char) => char.is_displayed);
}

/**
 * Groupe les personnages par serveur/realm
 */
export function groupCharactersByRealm(
  characters: EnrichedUserCharacter[]
): Record<string, EnrichedUserCharacter[]> {
  return characters.reduce((acc, character) => {
    const realm = character.realm;
    if (!acc[realm]) {
      acc[realm] = [];
    }
    acc[realm].push(character);
    return acc;
  }, {} as Record<string, EnrichedUserCharacter[]>);
}

/**
 * Calcule des statistiques sur les personnages
 */
export function getCharacterStats(characters: EnrichedUserCharacter[]) {
  const total = characters.length;
  const enriched = characters.filter(isCharacterEnriched).length;
  const needingUpdate = characters.filter((char) =>
    characterNeedsUpdate(char)
  ).length;
  const maxLevel = Math.max(...characters.map((char) => char.level), 0);
  const avgLevel =
    total > 0
      ? Math.round(
          characters.reduce((sum, char) => sum + char.level, 0) / total
        )
      : 0;

  return {
    total,
    enriched,
    needingUpdate,
    maxLevel,
    avgLevel,
    enrichedPercentage: total > 0 ? Math.round((enriched / total) * 100) : 0,
  };
}
