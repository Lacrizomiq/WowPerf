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

function handleApiError(error: unknown, defaultMessage: string): never {
  if (axios.isAxiosError(error)) {
    const err = error as AxiosError<APIError | RateLimitError>;

    switch (err.response?.status) {
      case 401:
        throw new CharacterError(
          CharacterErrorCode.UNAUTHORIZED,
          "Authentication required. Please link your Battle.net account."
        );

      case 403:
        // 🔥 NOUVEAU: Détecter spécifiquement "account_not_linked"
        const errorData = err.response?.data as any;
        if (errorData?.code === "account_not_linked") {
          throw new CharacterError(
            CharacterErrorCode.UNAUTHORIZED, // 🔥 Mapper vers UNAUTHORIZED pour déclencher auto-relink
            "Battle.net account not linked. Please connect your account."
          );
        }

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

    if (err.response?.data?.error) {
      throw new CharacterError(
        CharacterErrorCode.NETWORK_ERROR,
        err.response.data.error,
        error
      );
    }
  }

  throw new CharacterError(
    CharacterErrorCode.NETWORK_ERROR,
    defaultMessage,
    error
  );
}

// ============================================================================
// SERVICE PRINCIPAL POUR LES PERSONNAGES ENRICHIS
// ============================================================================

export const characterService = {
  /**
   * 🔥 SIMPLIFIÉ: Récupère TOUJOURS les personnages depuis la BDD
   * Fonctionne même si le token Battle.net est expiré
   * Le backend doit retourner les personnages stockés en base
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
      // Si erreur, retourner tableau vide au lieu de throw
      // L'utilisateur verra "aucun personnage" au lieu d'une erreur
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        return []; // Pas encore de personnages synchronisés
      }
      throw handleApiError(error, "Failed to get characters");
    }
  },

  /**
   * Synchronise et enrichit tous les personnages d'un compte Battle.net
   * Usage: Première utilisation après liaison OAuth OU re-sync après token expiré
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
   * Usage: Bouton "Refresh" pour mises à jour régulières
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
   * Enrichit un personnage spécifique
   * Usage: Bouton "Update" sur un personnage individuel
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
