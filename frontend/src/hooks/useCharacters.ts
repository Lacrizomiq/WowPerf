import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState, useEffect, useCallback } from "react";
import characterService from "@/libs/characterService";
import { battleNetService } from "@/libs/battlenetService";
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

export function useCharacterErrorHandler() {
  return useCallback((error: unknown, defaultMessage: string) => {
    if (error instanceof CharacterError) {
      switch (error.code) {
        case CharacterErrorCode.UNAUTHORIZED:
          toast.error("Battle.net connection required");
          break;
        case CharacterErrorCode.FORBIDDEN:
          toast.error("Access denied. Please re-link your Battle.net account.");
          break;
        case CharacterErrorCode.NOT_FOUND:
          toast.error("Character not found");
          break;
        case CharacterErrorCode.RATE_LIMIT:
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
// HOOK PRINCIPAL SIMPLIFIÉ
// ============================================================================

export function useCharacters(region: string = "eu") {
  const queryClient = useQueryClient();
  const { linkStatus } = useBattleNetLink();
  const handleError = useCharacterErrorHandler();

  // Rate limiting state
  const [rateLimitState, setRateLimitState] = useState<RateLimitState>({
    isRateLimited: false,
    expiryTime: null,
    timeRemaining: 0,
  });

  const [currentTimeRemaining, setCurrentTimeRemaining] = useState(0);

  // Timer pour le rate limiting
  useEffect(() => {
    if (!isRateLimitActive(rateLimitState)) {
      setCurrentTimeRemaining(0);
      return;
    }

    const interval = setInterval(() => {
      const remaining = getRemainingTime(rateLimitState);
      setCurrentTimeRemaining(remaining);

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
  // 🔥 QUERY SIMPLIFIÉE - TOUJOURS RÉCUPÉRER LES PERSONNAGES BDD
  // ============================================================================

  const charactersQuery = useQuery<EnrichedUserCharacter[], Error>({
    queryKey: ["characters"],
    queryFn: () => characterService.getCharacters(),
    // 🔥 TOUJOURS ACTIVÉ - ne dépend plus du linkStatus
    enabled: true,
    retry: (failureCount, error) => {
      // Ne pas retry sur 401, c'est normal si pas de personnages
      if (
        error instanceof CharacterError &&
        error.code === CharacterErrorCode.UNAUTHORIZED
      ) {
        return false;
      }
      return failureCount < 2;
    },
  });

  // ============================================================================
  // 🔥 AUTO-RELINK + SYNC LOGIC
  // ============================================================================

  /**
   * Fonction principale : Lance sync OU auto-relink si nécessaire
   */
  const smartSyncAndEnrich = useCallback(async () => {
    try {
      // Essayer la sync directement
      return await characterService.syncAndEnrich(region);
    } catch (error) {
      // Si 401 (token expiré), lancer le relink automatiquement
      if (
        error instanceof CharacterError &&
        error.code === CharacterErrorCode.UNAUTHORIZED
      ) {
        toast.loading("Battle.net token expired. Re-linking your account...", {
          duration: 3000,
        });

        // 🔥 Déclencher le relink avec flag auto_relink
        const { url } = await battleNetService.initiateLinking(true);

        // Ouvrir la fenêtre d'auth (ou rediriger)
        window.location.href = url;

        throw new Error("Redirecting to Battle.net authentication...");
      }

      // Autres erreurs, les laisser passer
      throw error;
    }
  }, [region]);

  // ============================================================================
  // MUTATIONS
  // ============================================================================

  const syncAndEnrichMutation = useMutation<SyncAndEnrichResult, Error>({
    mutationFn: smartSyncAndEnrich,
    onSuccess: (data) => {
      const { result } = data;

      if (result.enriched_count > 0) {
        toast.success(
          `${result.synced_count} characters synchronized, ${result.enriched_count} enriched!`,
          { duration: 4000 }
        );
      } else if (result.synced_count > 0) {
        toast.success(`${result.synced_count} characters synchronized!`, {
          duration: 3000,
        });
      }

      if (result.errors && result.errors.length > 0) {
        toast.error(`Some issues occurred. ${result.errors.length} errors.`);
      }

      // Actualiser les données
      queryClient.invalidateQueries({ queryKey: ["characters"] });
      queryClient.invalidateQueries({ queryKey: ["battleNetLinkStatus"] });

      // Reset rate limiting
      setRateLimitState({
        isRateLimited: false,
        expiryTime: null,
        timeRemaining: 0,
      });
    },
    onError: (error: Error) => {
      // Ignorer l'erreur de redirection vers Battle.net
      if (error.message.includes("Redirecting to Battle.net")) {
        return;
      }

      if (
        error instanceof CharacterError &&
        error.code === CharacterErrorCode.RATE_LIMIT
      ) {
        const newRateLimitState = calculateRateLimitState(error.waitTime);
        setRateLimitState(newRateLimitState);
        setCurrentTimeRemaining(newRateLimitState.timeRemaining);
      } else {
        handleError(error, "Failed to sync characters");
      }
    },
  });

  const refreshAndEnrichMutation = useMutation<RefreshAndEnrichResult, Error>({
    mutationFn: () => characterService.refreshAndEnrich(region),
    onSuccess: (data) => {
      const { result } = data;
      toast.success(`${result.enriched_count} characters refreshed!`);
      queryClient.invalidateQueries({ queryKey: ["characters"] });
    },
    onError: (error: Error) => {
      if (
        error instanceof CharacterError &&
        error.code === CharacterErrorCode.RATE_LIMIT
      ) {
        const newRateLimitState = calculateRateLimitState(error.waitTime);
        setRateLimitState(newRateLimitState);
        setCurrentTimeRemaining(newRateLimitState.timeRemaining);
      } else {
        handleError(error, "Failed to refresh characters");
      }
    },
  });

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
  // COMPUTED VALUES
  // ============================================================================

  const characters: EnrichedUserCharacter[] = charactersQuery.data || [];
  const hasCharacters = characters.length > 0;
  const isRateLimited = isRateLimitActive(rateLimitState);

  const rateLimitMessage =
    isRateLimited && currentTimeRemaining > 0
      ? `Please wait ${formatTimeRemaining(
          currentTimeRemaining
        )} before next sync`
      : null;

  const isDisabled = {
    sync: isRateLimited || syncAndEnrichMutation.isPending,
    refresh: isRateLimited || refreshAndEnrichMutation.isPending,
    individual: enrichCharacterMutation.isPending,
  };

  const refetchCharacters = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: ["characters"] });
  }, [queryClient]);

  // ============================================================================
  // ACTIONS
  // ============================================================================

  const actions = {
    syncAndEnrich: useCallback(() => {
      if (!isRateLimited) {
        syncAndEnrichMutation.mutate();
      }
    }, [isRateLimited, syncAndEnrichMutation]),

    refreshAndEnrich: useCallback(() => {
      if (!isRateLimited && linkStatus?.linked) {
        refreshAndEnrichMutation.mutate();
      } else if (!linkStatus?.linked) {
        toast.error("Please link your Battle.net account first");
      }
    }, [isRateLimited, linkStatus, refreshAndEnrichMutation]),

    enrichCharacter: useCallback(
      (characterId: number) => {
        enrichCharacterMutation.mutate(characterId);
      },
      [enrichCharacterMutation]
    ),

    refetchCharacters,
  };

  // ============================================================================
  // RETURN
  // ============================================================================

  return {
    // 🔥 Données toujours disponibles si en BDD
    characters,
    hasCharacters,
    isLoadingCharacters: charactersQuery.isLoading,
    charactersError: charactersQuery.error,

    // Actions
    actions,

    // Loading states
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

    // UI states
    ui: {
      isDisabled,
      canSync: !isDisabled.sync, // 🔥 Toujours possible (auto-relink si nécessaire)
      canRefresh: linkStatus?.linked === true && !isDisabled.refresh,
      showRateLimit: isRateLimited,
    },

    // Meta
    region,
    isAuthenticated: linkStatus?.linked === true,
  };
}
