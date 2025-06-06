import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState, useEffect, useCallback } from "react";
import characterService from "@/libs/characterService";
import {
  EnrichedUserCharacter,
  SyncAndEnrichResult,
  RefreshAndEnrichResult,
  CharacterError,
  CharacterErrorCode,
  RateLimitState,
} from "@/types/character/character";
import {
  calculateRateLimitState,
  isRateLimitActive,
  getRemainingTime,
  formatTimeRemaining,
} from "@/utils/character/character";
import { useBattleNetLink } from "./useBattleNetLink";
import toast from "react-hot-toast";

// ============================================================================
// GESTION CENTRALISÉE DES ERREURS
// ============================================================================

/**
 * Hook pour gérer les erreurs de personnages de manière centralisée
 */
export function useCharacterErrorHandler() {
  return useCallback((error: unknown, defaultMessage: string) => {
    if (error instanceof CharacterError) {
      switch (error.code) {
        case CharacterErrorCode.UNAUTHORIZED:
          toast.error("Please link your Battle.net account first");
          break;
        case CharacterErrorCode.FORBIDDEN:
          toast.error(
            "Access denied. Your Battle.net session may have expired. Please re-link your account."
          );
          console.error("403 Forbidden Error Details:", error);
          break;
        case CharacterErrorCode.NOT_FOUND:
          toast.error("Character not found");
          break;
        case CharacterErrorCode.RATE_LIMIT:
          // Rate limit géré séparément dans le hook principal
          toast.error(error.message);
          break;
        case CharacterErrorCode.SERVER_ERROR:
          toast.error("Server error. Please try again later");
          break;
        default:
          toast.error(error.message || defaultMessage);
      }
    } else {
      toast.error(defaultMessage);
    }
  }, []);
}

// ============================================================================
// HOOK PRINCIPAL POUR LES PERSONNAGES ENRICHIS
// ============================================================================

/**
 * Hook principal pour gérer les personnages enrichis avec rate limiting
 *
 * Features:
 * - Récupération des personnages depuis BDD (instantané)
 * - Sync et enrichissement (première fois)
 * - Refresh et enrichissement (mise à jour)
 * - Enrichissement individuel
 * - Rate limiting avec timer UI
 * - Gestion d'erreurs avec toast
 */
export function useCharacters(region: string = "eu") {
  const queryClient = useQueryClient();
  const { linkStatus } = useBattleNetLink();
  const handleError = useCharacterErrorHandler();

  // ============================================================================
  // ÉTATS RATE LIMITING
  // ============================================================================

  const [rateLimitState, setRateLimitState] = useState<RateLimitState>({
    isRateLimited: false,
    expiryTime: null,
    timeRemaining: 0,
  });

  const [currentTimeRemaining, setCurrentTimeRemaining] = useState(0);

  // Timer pour mettre à jour le temps restant chaque seconde
  useEffect(() => {
    if (!isRateLimitActive(rateLimitState)) {
      setCurrentTimeRemaining(0);
      return;
    }

    const interval = setInterval(() => {
      const remaining = getRemainingTime(rateLimitState);
      setCurrentTimeRemaining(remaining);

      // Rate limit expiré
      if (remaining <= 0) {
        setRateLimitState({
          isRateLimited: false,
          expiryTime: null,
          timeRemaining: 0,
        });
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [rateLimitState]);

  // ============================================================================
  // QUERIES - RÉCUPÉRATION DES DONNÉES
  // ============================================================================

  /**
   * Récupère tous les personnages enrichis depuis la BDD
   * Données instantanées, pas d'appel API Blizzard
   */
  const charactersQuery = useQuery<EnrichedUserCharacter[], Error>({
    queryKey: ["characters"],
    queryFn: () => characterService.getCharacters(),
    enabled: linkStatus?.linked === true,
    retry: (failureCount, error) => {
      if (error instanceof CharacterError) {
        return (
          error.code !== CharacterErrorCode.UNAUTHORIZED && failureCount < 3
        );
      }
      return failureCount < 3;
    },
  });

  // ============================================================================
  // MUTATIONS - OPÉRATIONS D'ÉCRITURE
  // ============================================================================

  /**
   * Synchronisation et enrichissement complet (première fois)
   * Usage: Modal onboarding après liaison OAuth
   */
  const syncAndEnrichMutation = useMutation<SyncAndEnrichResult, Error>({
    mutationFn: () => characterService.syncAndEnrich(region),
    onSuccess: (data) => {
      const { result } = data;

      if (result.enriched_count > 0) {
        toast.success(
          `${result.synced_count} characters synchronized, ${result.enriched_count} enriched!`,
          { duration: 4000 }
        );
      } else if (result.synced_count > 0) {
        toast.success(
          `${result.synced_count} characters synchronized! Enrichment may need a refresh.`,
          { duration: 5000 }
        );
      }

      // Show errors if any
      if (result.errors && result.errors.length > 0) {
        toast.error(
          `Some issues occurred during sync. ${result.errors.length} errors.`,
          { duration: 4000 }
        );
      }

      // Invalider la cache pour recharger les données
      queryClient.invalidateQueries({ queryKey: ["characters"] });

      // Réinitialiser le rate limiting en cas de succès
      setRateLimitState({
        isRateLimited: false,
        expiryTime: null,
        timeRemaining: 0,
      });
    },
    onError: (error: Error) => {
      if (
        error instanceof CharacterError &&
        error.code === CharacterErrorCode.RATE_LIMIT
      ) {
        // Gestion spécifique du rate limiting
        const newRateLimitState = calculateRateLimitState(error.waitTime);
        setRateLimitState(newRateLimitState);
        setCurrentTimeRemaining(newRateLimitState.timeRemaining);
      } else if (
        error instanceof CharacterError &&
        error.code === CharacterErrorCode.SERVER_ERROR
      ) {
        // Server error might be rate limiting in disguise
        const possibleRateLimit = calculateRateLimitState("5m"); // Assume 5min wait
        setRateLimitState(possibleRateLimit);
        setCurrentTimeRemaining(possibleRateLimit.timeRemaining);
        toast.error(
          "Server busy. Please wait a few minutes before trying again."
        );
      } else {
        handleError(error, "Failed to sync and enrich characters");
      }
    },
  });

  /**
   * Rafraîchissement et enrichissement (mise à jour régulière)
   * Usage: Bouton "Refresh" dans l'interface
   */
  const refreshAndEnrichMutation = useMutation<RefreshAndEnrichResult, Error>({
    mutationFn: () => characterService.refreshAndEnrich(region),
    onSuccess: (data) => {
      const { result } = data;

      if (result.enriched_count > 0) {
        toast.success(
          `${result.enriched_count} characters refreshed successfully!`,
          { duration: 4000 }
        );
      } else {
        toast.success(
          `Characters refreshed! Some enrichment data may be unavailable.`,
          { duration: 4000 }
        );
      }

      // Show errors if any
      if (result.errors && result.errors.length > 0) {
        toast.error(
          `Some issues occurred during refresh. ${result.errors.length} errors.`,
          { duration: 4000 }
        );
      }

      queryClient.invalidateQueries({ queryKey: ["characters"] });

      setRateLimitState({
        isRateLimited: false,
        expiryTime: null,
        timeRemaining: 0,
      });
    },
    onError: (error: Error) => {
      if (
        error instanceof CharacterError &&
        error.code === CharacterErrorCode.RATE_LIMIT
      ) {
        const newRateLimitState = calculateRateLimitState(error.waitTime);
        setRateLimitState(newRateLimitState);
        setCurrentTimeRemaining(newRateLimitState.timeRemaining);
      } else if (
        error instanceof CharacterError &&
        error.code === CharacterErrorCode.SERVER_ERROR
      ) {
        // Server error might be rate limiting in disguise
        const possibleRateLimit = calculateRateLimitState("5m"); // Assume 5min wait
        setRateLimitState(possibleRateLimit);
        setCurrentTimeRemaining(possibleRateLimit.timeRemaining);
        toast.error(
          "Server busy. Please wait a few minutes before trying again."
        );
      } else {
        handleError(error, "Failed to refresh characters");
      }
    },
  });

  /**
   * Enrichissement d'un personnage individuel
   * Usage: Bouton "Update" sur un personnage spécifique
   */
  const enrichCharacterMutation = useMutation<
    { message: string },
    Error,
    number
  >({
    mutationFn: (characterId: number) =>
      characterService.enrichCharacter(characterId),
    onSuccess: () => {
      toast.success("Character updated successfully!");
      queryClient.invalidateQueries({ queryKey: ["characters"] });
    },
    onError: (error: Error) => {
      handleError(error, "Failed to update character");
    },
  });

  // ============================================================================
  // HELPERS ET ÉTATS CALCULÉS
  // ============================================================================

  /**
   * Vérifie si une opération est actuellement bloquée par le rate limiting
   */
  const isRateLimited = isRateLimitActive(rateLimitState);

  /**
   * Message du rate limiting pour l'UI
   */
  const rateLimitMessage =
    isRateLimited && currentTimeRemaining > 0
      ? `Please wait ${formatTimeRemaining(
          currentTimeRemaining
        )} before next sync`
      : null;

  /**
   * Vérifie si les boutons doivent être désactivés
   */
  const isDisabled = {
    sync: isRateLimited || syncAndEnrichMutation.isPending,
    refresh: isRateLimited || refreshAndEnrichMutation.isPending,
    individual: enrichCharacterMutation.isPending,
  };

  /**
   * Force l'actualisation des personnages depuis la cache
   */
  const refetchCharacters = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: ["characters"] });
  }, [queryClient]);

  /**
   * Actions disponibles avec vérification du rate limiting
   */
  const actions = {
    syncAndEnrich: useCallback(() => {
      if (!isRateLimited) {
        syncAndEnrichMutation.mutate();
      }
    }, [isRateLimited, syncAndEnrichMutation]),

    refreshAndEnrich: useCallback(() => {
      if (!isRateLimited) {
        refreshAndEnrichMutation.mutate();
      }
    }, [isRateLimited, refreshAndEnrichMutation]),

    enrichCharacter: useCallback(
      (characterId: number) => {
        enrichCharacterMutation.mutate(characterId);
      },
      [enrichCharacterMutation]
    ),

    refetchCharacters,
  };

  // ============================================================================
  // RETURN - INTERFACE PUBLIQUE DU HOOK
  // ============================================================================

  return {
    // Données des personnages
    characters: charactersQuery.data,
    isLoadingCharacters: charactersQuery.isLoading,
    charactersError: charactersQuery.error,

    // Actions disponibles
    actions,

    // États des opérations
    isLoading: {
      sync: syncAndEnrichMutation.isPending,
      refresh: refreshAndEnrichMutation.isPending,
      individual: enrichCharacterMutation.isPending,
    },

    // Rate limiting
    rateLimitState: {
      isRateLimited,
      timeRemaining: currentTimeRemaining,
      message: rateLimitMessage,
      formattedTime:
        currentTimeRemaining > 0
          ? formatTimeRemaining(currentTimeRemaining)
          : null,
    },

    // États UI
    ui: {
      isDisabled,
      canSync: linkStatus?.linked === true && !isDisabled.sync,
      canRefresh: linkStatus?.linked === true && !isDisabled.refresh,
      showRateLimit: isRateLimited,
    },

    // Métadonnées
    region,
    isAuthenticated: linkStatus?.linked === true,
  };
}
