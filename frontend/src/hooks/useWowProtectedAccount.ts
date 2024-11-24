import { useQuery } from "@tanstack/react-query";
import {
  wowService,
  WoWProfile,
  WoWError,
  WoWErrorCode,
} from "@/libs/wowProtectedAccountService";
import { useBattleNetLink } from "./useBattleNetLink";
import { useCallback } from "react";
import toast from "react-hot-toast";

export function useWoWProfile() {
  const { linkStatus } = useBattleNetLink();

  const {
    data: wowProfile,
    isLoading,
    error,
    refetch,
  } = useQuery<WoWProfile>({
    queryKey: ["wowProfile"],
    queryFn: async () => {
      const profile = await wowService.getWoWProfile();
      return profile;
    },
    enabled: linkStatus?.linked === true,
    retry: (failureCount, error) => {
      if (error instanceof WoWError) {
        return error.code !== WoWErrorCode.UNAUTHORIZED && failureCount < 3;
      }
      return failureCount < 3;
    },
  });

  const handleError = useCallback((error: unknown) => {
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
        default:
          toast.error("Failed to fetch WoW profile");
      }
    }
  }, []);

  return {
    wowProfile,
    isLoading,
    error,
    refetch,
  };
}
