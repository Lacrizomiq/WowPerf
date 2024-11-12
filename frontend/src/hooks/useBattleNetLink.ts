import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useCallback, useState } from "react";
import {
  battleNetService,
  BattleNetError,
  BattleNetErrorCode,
} from "@/libs/battlenetService";
import { useAuth } from "@/providers/AuthContext";

export function useBattleNetLink() {
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuth();
  const [isLinking, setIsLinking] = useState(false);
  // Query for Battle.net link status
  const {
    data: linkStatus,
    isLoading,
    error: queryError,
  } = useQuery({
    queryKey: ["battleNetLinkStatus"],
    queryFn: battleNetService.getLinkStatus,
    enabled: isAuthenticated,
    retry: (failureCount, error) => {
      if (error instanceof BattleNetError) {
        return (
          error.code !== BattleNetErrorCode.UNAUTHORIZED && failureCount < 3
        );
      }
      return failureCount < 3;
    },
  });

  // Mutation for unlinking account
  const unlinkMutation = useMutation({
    mutationFn: battleNetService.unlinkAccount,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["battleNetLinkStatus"] });
    },
  });

  // Initiate linking
  const initiateLink = useCallback(async () => {
    try {
      await battleNetService.initiateLinking();
      return { success: true as const };
    } catch (error) {
      return {
        success: false as const,
        error:
          error instanceof Error ? error.message : "Failed to initiate linking",
      };
    }
  }, []);

  // Unlink account
  const unlinkAccount = async () => {
    try {
      await unlinkMutation.mutateAsync();
      return { success: true as const };
    } catch (error) {
      return {
        success: false as const,
        error:
          error instanceof Error ? error.message : "Failed to unlink account",
      };
    }
  };

  return {
    linkStatus,
    isLoading,
    error: queryError,
    initiateLink,
    unlinkAccount,
    isUnlinking: unlinkMutation.isPending,
    unlinkError: unlinkMutation.error,
    isLinking,
  };
}
