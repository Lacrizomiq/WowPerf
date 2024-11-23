import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useCallback, useState, useEffect } from "react";
import toast from "react-hot-toast";
import {
  battleNetService,
  BattleNetError,
  BattleNetErrorCode,
  BattleNetStatusResponse,
} from "@/libs/battlenetService";
import { useAuth } from "@/providers/AuthContext";

export interface BattleNetLinkState {
  linkStatus: BattleNetStatusResponse | null;
  isLoading: boolean;
  error: Error | null;
  isLinking: boolean;
  isUnlinking: boolean;
  initiateLink: () => Promise<{ success: boolean; error?: string }>;
  unlinkAccount: () => Promise<{ success: boolean; error?: string }>;
}

export function useBattleNetLink(): BattleNetLinkState {
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuth();
  const [isLinking, setIsLinking] = useState(false);

  // Query pour le statut de la liaison
  const {
    data: linkStatus,
    isLoading,
    error: queryError,
  } = useQuery({
    queryKey: ["battleNetLinkStatus"],
    queryFn: battleNetService.getLinkStatus,
    enabled: isAuthenticated,
    retry: (failureCount, error: unknown) => {
      if (error instanceof BattleNetError) {
        return (
          error.code !== BattleNetErrorCode.UNAUTHORIZED && failureCount < 3
        );
      }
      return failureCount < 3;
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
    refetchOnWindowFocus: true,
  });

  // Gérer les erreurs via useEffect
  useEffect(() => {
    if (
      queryError instanceof BattleNetError &&
      queryError.code !== BattleNetErrorCode.UNAUTHORIZED
    ) {
      toast.error(queryError.message);
    }
  }, [queryError]);

  // Mutation pour le unlinking
  const unlinkMutation = useMutation({
    mutationFn: battleNetService.unlinkAccount,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["battleNetLinkStatus"] });
      toast.success("Battle.net account unlinked successfully");
    },
    onError: (error: unknown) => {
      const message =
        error instanceof BattleNetError
          ? error.message
          : "Failed to unlink Battle.net account";
      toast.error(message);
    },
  });

  // Initier la liaison
  const initiateLink = useCallback(async () => {
    setIsLinking(true);
    try {
      const { url } = await battleNetService.initiateLinking();
      window.location.href = url;
      return { success: true };
    } catch (error) {
      const message =
        error instanceof BattleNetError
          ? error.message
          : "Failed to initiate Battle.net linking";
      toast.error(message);
      return { success: false, error: message };
    } finally {
      setIsLinking(false);
    }
  }, []);

  // Délier le compte
  const unlinkAccount = async () => {
    const toastId = toast.loading("Unlinking Battle.net account...");
    try {
      await unlinkMutation.mutateAsync();
      toast.success("Account unlinked successfully", { id: toastId });
      return { success: true };
    } catch (error) {
      const message =
        error instanceof BattleNetError
          ? error.message
          : "Failed to unlink Battle.net account";
      toast.error(message, { id: toastId });
      return { success: false, error: message };
    }
  };

  return {
    linkStatus: linkStatus || null,
    isLoading,
    error: queryError as Error | null,
    initiateLink,
    unlinkAccount,
    isUnlinking: unlinkMutation.isPending,
    isLinking,
  };
}
