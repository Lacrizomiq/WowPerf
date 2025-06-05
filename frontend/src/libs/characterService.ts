import api, { APIError } from "./api";
import axios, { AxiosError } from "axios";
import {
  EnrichedUserCharacter,
  SyncAndEnrichResult,
  RefreshAndEnrichResult,
  RateLimitError,
  CharacterError,
  CharacterErrorCode,
  GetCharactersResponse,
} from "@/types/character/character";
import { extractWaitTime } from "@/utils/character/character";

// ============================================================================
// GESTION D'ERREURS
// ============================================================================

/**
 * Gestion centralisée des erreurs avec support spécialisé rate limiting
 *
 * Erreurs gérées:
 * - HTTP 429: Rate limit avec extraction du temps d'attente
 * - HTTP 401: Authentication requise
 * - HTTP 403: Ownership (personnage n'appartient pas à l'utilisateur)
 * - HTTP 404: Personnage non trouvé
 * - HTTP 5xx: Erreurs serveur
 *
 * @param error - Erreur originale (Axios ou autre)
 * @param defaultMessage - Message par défaut si erreur non reconnue
 * @throws CharacterError avec code spécialisé et waitTime pour rate limiting
 */
function handleApiError(error: unknown, defaultMessage: string): never {
  // Erreurs Axios (HTTP)
  if (axios.isAxiosError(error)) {
    const err = error as AxiosError<APIError | RateLimitError>;

    switch (err.response?.status) {
      case 401:
        throw new CharacterError(
          CharacterErrorCode.UNAUTHORIZED,
          "Authentication required. Please link your Battle.net account."
        );

      case 403:
        throw new CharacterError(
          CharacterErrorCode.FORBIDDEN,
          "Access denied. Character may not belong to this user."
        );

      case 404:
        throw new CharacterError(
          CharacterErrorCode.NOT_FOUND,
          "Character not found."
        );

      case 429:
        // Rate limit spécifique avec extraction du temps d'attente
        const rateLimitData = err.response?.data as RateLimitError;
        const waitTime = rateLimitData?.error
          ? extractWaitTime(rateLimitData.error)
          : undefined;

        throw new CharacterError(
          CharacterErrorCode.RATE_LIMIT,
          rateLimitData?.error ||
            "Rate limit exceeded. Please try again later.",
          error,
          waitTime
        );

      case 500:
      case 502:
      case 503:
      case 504:
        throw new CharacterError(
          CharacterErrorCode.SERVER_ERROR,
          "A server error occurred. Please try again later."
        );
    }

    // Autres codes HTTP avec message personnalisé
    if (err.response?.data?.error) {
      throw new CharacterError(
        CharacterErrorCode.NETWORK_ERROR,
        err.response.data.error,
        error
      );
    }
  }

  // Erreurs non-Axios (réseau, JavaScript, etc.)
  throw new CharacterError(
    CharacterErrorCode.NETWORK_ERROR,
    defaultMessage,
    error
  );
}

// ============================================================================
// SERVICE PRINCIPAL POUR LES PERSONNAGES ENRICHIS
// ============================================================================

/**
 * Service pour gérer les personnages enrichis via le nouveau système orchestré
 *
 * Architecture:
 * - syncAndEnrich: Première sync après OAuth (modal onboarding)
 * - refreshAndEnrich: Bouton refresh quotidien
 * - getCharacters: Affichage instantané depuis BDD
 * - enrichCharacter: Mise à jour individuelle
 */
export const characterService = {
  /**
   * Synchronise et enrichit tous les personnages d'un compte Battle.net
   *
   * Usage: Première utilisation après liaison OAuth
   * Flux: Battle.net API → Sync BDD → Enrichissement → Sauvegarde
   *
   * @param region - Région Battle.net (eu, us, kr, tw)
   * @returns Résultat avec compteurs et erreurs éventuelles
   */
  async syncAndEnrich(region: string = "eu"): Promise<SyncAndEnrichResult> {
    try {
      const response = await api.post<SyncAndEnrichResult>(
        "/characters/sync-and-enrich",
        {},
        {
          headers: {
            Region: region,
            Accept: "application/json",
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      throw handleApiError(error, "Failed to sync and enrich characters");
    }
  },

  /**
   * Rafraîchit et enrichit les personnages existants
   *
   * Usage: Bouton "Refresh" pour mises à jour régulières
   * Flux: BDD existante → API Blizzard → Enrichissement → Update BDD
   *
   * @param region - Région Battle.net (eu, us, kr, tw)
   * @returns Résultat avec compteurs et erreurs éventuelles
   */
  async refreshAndEnrich(
    region: string = "eu"
  ): Promise<RefreshAndEnrichResult> {
    try {
      const response = await api.post<RefreshAndEnrichResult>(
        "/characters/refresh-and-enrich",
        {},
        {
          headers: {
            Region: region,
            Accept: "application/json",
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      throw handleApiError(error, "Failed to refresh and enrich characters");
    }
  },

  /**
   * Récupère tous les personnages enrichis de l'utilisateur
   *
   * Usage: Affichage principal - données instantanées depuis BDD
   * Avantages: Pas d'appel API Blizzard, données enrichies disponibles
   *
   * @returns Liste des personnages avec toutes les données enrichies
   */
  async getCharacters(): Promise<EnrichedUserCharacter[]> {
    try {
      const response = await api.get<GetCharactersResponse>("/characters", {
        headers: {
          Accept: "application/json",
        },
        withCredentials: true,
      });
      return response.data.characters;
    } catch (error) {
      throw handleApiError(error, "Failed to get characters");
    }
  },

  /**
   * Enrichit un personnage spécifique
   *
   * Usage: Bouton "Update" sur un personnage individuel
   * Flux: Character ID → API Blizzard → Enrichissement → Update BDD
   *
   * @param characterId - ID du personnage à enrichir
   * @returns Message de confirmation
   */
  async enrichCharacter(characterId: number): Promise<{ message: string }> {
    try {
      const response = await api.post<{ message: string }>(
        `/characters/${characterId}/enrich`,
        {},
        {
          headers: {
            Accept: "application/json",
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      throw handleApiError(error, `Failed to enrich character ${characterId}`);
    }
  },
};

export default characterService;
