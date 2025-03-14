import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { wowService } from "@/libs/wowProtectedAccountService";
import React, { useEffect } from "react";
import {
  WoWProfile,
  CharacterBasicInfo,
  UserCharacter,
  SyncResult,
  RefreshResult,
  CharacterDisplayResponse,
  FavoriteCharacterResponse,
  WoWError,
  WoWErrorCode,
} from "@/types/userCharacter/userCharacter";
import { useBattleNetLink } from "./useBattleNetLink";
import { useCallback } from "react";
import toast from "react-hot-toast";

/**
 * Hook to handle WoW errors in a centralized way
 */
export function useWoWErrorHandler() {
  return useCallback((error: unknown, defaultMessage: string) => {
    if (error instanceof WoWError) {
      switch (error.code) {
        case WoWErrorCode.UNAUTHORIZED:
          toast.error("Please link your Battle.net account first");
          break;
        case WoWErrorCode.TOKEN_EXPIRED:
          toast.error(
            "Session expired. Please re-link your Battle.net account"
          );
          break;
        case WoWErrorCode.NOT_FOUND:
          toast.error("The requested character was not found");
          break;
        case WoWErrorCode.SERVER_ERROR:
          toast.error("A server error occurred. Please try again later");
          break;
        default:
          toast.error(defaultMessage);
      }
    } else {
      toast.error(defaultMessage);
    }
  }, []);
}

/**
 * Main hook to fetch the WoW profile of the user
 */
export function useWoWProfile(region: string = "eu") {
  const { linkStatus } = useBattleNetLink();
  const handleError = useWoWErrorHandler();

  const result = useQuery<WoWProfile, Error>({
    queryKey: ["wowProfile", region],
    queryFn: async () => wowService.getWoWProfile(region),
    enabled: linkStatus?.linked === true,
    retry: (failureCount, error) => {
      if (error instanceof WoWError) {
        return error.code !== WoWErrorCode.UNAUTHORIZED && failureCount < 3;
      }
      return failureCount < 3;
    },
  });

  // Gérer les erreurs avec useEffect
  useEffect(() => {
    if (result.error) {
      handleError(result.error, "Failed to fetch WoW profile");
    }
  }, [result.error, handleError]);

  return result;
}

/**
 * Hook to manage WoW characters
 */
export function useWoWCharacters() {
  const queryClient = useQueryClient();
  const { linkStatus } = useBattleNetLink();
  const handleError = useWoWErrorHandler();

  // Fetch the WoW profile (to get the region)
  const profileQuery = useWoWProfile();
  const wowProfile = profileQuery.data;
  const region = wowProfile?.region || "eu";

  // List the available characters of the Battle.net account
  const accountCharactersQuery = useQuery<CharacterBasicInfo[], Error>({
    queryKey: ["accountCharacters", region],
    queryFn: () => wowService.listAccountCharacters(region),
    enabled: linkStatus?.linked === true && !profileQuery.isLoading,
    retry: (failureCount, error) => {
      if (error instanceof WoWError) {
        return error.code !== WoWErrorCode.UNAUTHORIZED && failureCount < 3;
      }
      return failureCount < 3;
    },
  });

  // Handle errors for accountCharactersQuery
  React.useEffect(() => {
    if (accountCharactersQuery.error) {
      handleError(
        accountCharactersQuery.error,
        "Failed to list account characters"
      );
    }
  }, [accountCharactersQuery.error, handleError]);

  // Fetch the synchronized characters of the user
  const userCharactersQuery = useQuery<UserCharacter[], Error>({
    queryKey: ["userCharacters"],
    queryFn: () => wowService.getUserCharacters(),
    enabled: linkStatus?.linked === true,
    retry: (failureCount, error) => {
      if (error instanceof WoWError) {
        return error.code !== WoWErrorCode.UNAUTHORIZED && failureCount < 3;
      }
      return failureCount < 3;
    },
  });

  // Handle errors for userCharactersQuery
  React.useEffect(() => {
    if (userCharactersQuery.error) {
      handleError(userCharactersQuery.error, "Failed to get user characters");
    }
  }, [userCharactersQuery.error, handleError]);

  // Mutation to synchronize all characters
  const syncMutation = useMutation<SyncResult, Error>({
    mutationFn: () => wowService.syncAllAccountCharacters(region),
    onSuccess: (data) => {
      toast.success(`${data.count} characters synchronized successfully`);
      queryClient.invalidateQueries({ queryKey: ["userCharacters"] });
    },
    onError: (error: Error) => {
      handleError(error, "Failed to synchronize characters");
    },
  });

  // Mutation to refresh the characters
  const refreshMutation = useMutation<RefreshResult, Error>({
    mutationFn: () => wowService.refreshUserCharacters(region),
    onSuccess: (data) => {
      toast.success(
        `${data.new_characters} new characters added, ${data.updated_characters} characters updated`
      );
      queryClient.invalidateQueries({ queryKey: ["userCharacters"] });
    },
    onError: (error: Error) => {
      handleError(error, "Failed to refresh characters");
    },
  });

  // Mutation to set a character as favorite
  const setFavoriteMutation = useMutation<{ message: string }, Error, number>({
    mutationFn: (characterId: number) =>
      wowService.setFavoriteCharacter(characterId),
    onSuccess: () => {
      toast.success("Character set as favorite successfully");

      // Invalidate both the characters and the user profile queries
      queryClient.invalidateQueries({ queryKey: ["userCharacters"] });
      queryClient.invalidateQueries({ queryKey: ["userProfile"] }); // Important to update favorite_character_id
    },
    onError: (error: Error) => {
      handleError(error, "Failed to set favorite character");
    },
  });

  // Mutation to enable/disable the display of a character
  const toggleDisplayMutation = useMutation<
    { message: string },
    Error,
    { characterId: number; display: boolean }
  >({
    mutationFn: ({ characterId, display }) =>
      wowService.toggleCharacterDisplay(characterId, display),
    onSuccess: () => {
      toast.success("Character display updated successfully");
      queryClient.invalidateQueries({ queryKey: ["userCharacters"] });
    },
    onError: (error: Error) => {
      handleError(error, "Failed to update character display");
    },
  });

  return {
    // WoW profile
    wowProfile,
    isLoadingProfile: profileQuery.isLoading,
    profileError: profileQuery.error,

    // Battle.net characters
    accountCharacters: accountCharactersQuery.data,
    isLoadingAccountCharacters: accountCharactersQuery.isLoading,
    accountCharactersError: accountCharactersQuery.error,
    refetchAccountCharacters: accountCharactersQuery.refetch,

    // Synchronized characters
    userCharacters: userCharactersQuery.data,
    isLoadingUserCharacters: userCharactersQuery.isLoading,
    userCharactersError: userCharactersQuery.error,
    refetchUserCharacters: userCharactersQuery.refetch,

    // Actions
    syncCharacters: syncMutation.mutate,
    isSyncing: syncMutation.isPending,
    syncError: syncMutation.error,

    refreshCharacters: refreshMutation.mutate,
    isRefreshing: refreshMutation.isPending,
    refreshError: refreshMutation.error,

    setFavoriteCharacter: setFavoriteMutation.mutate,
    isSettingFavorite: setFavoriteMutation.isPending,
    setFavoriteError: setFavoriteMutation.error,

    toggleCharacterDisplay: toggleDisplayMutation.mutate,
    isTogglingDisplay: toggleDisplayMutation.isPending,
    toggleDisplayError: toggleDisplayMutation.error,

    // Current region
    region,
  };
}
