// src/hooks/useUserProfile.tsx
"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useEffect, useCallback } from "react";
import {
  userService,
  UserServiceError,
  UserErrorCode,
  type UserProfile,
  type ApiResponse,
} from "@/libs/userService";
import { useAuth } from "@/providers/AuthContext";
import { useRouter } from "next/navigation";
import { csrfService } from "@/libs/csrfService";

// Improved types for mutations
interface MutationState {
  isError: boolean;
  error: Error | null;
  isPending: boolean;
  errorCode?: UserErrorCode;
}

interface MutationResponse {
  success: boolean;
  error?: string;
  code?: UserErrorCode;
}

export function useUserProfile() {
  const queryClient = useQueryClient();
  const router = useRouter();
  const { isAuthenticated, logout } = useAuth();

  // Improved centralized error handler
  const handleError = useCallback(
    async (error: unknown): Promise<MutationResponse> => {
      console.error("User profile error:", error);

      if (error instanceof UserServiceError) {
        switch (error.code) {
          case UserErrorCode.INVALID_CSRF_TOKEN:
            try {
              // Attempt to refresh the CSRF token
              await csrfService.getToken(true);
              return {
                success: false,
                error: "Security token refreshed. Please try again.",
                code: error.code,
              };
            } catch {
              await logout();
              return {
                success: false,
                error: "Security verification failed. Please log in again.",
                code: error.code,
              };
            }

          case UserErrorCode.UNAUTHORIZED:
            csrfService.clearToken();
            await logout();
            return {
              success: false,
              error: "Session expired. Please login again.",
              code: error.code,
            };

          case UserErrorCode.EMAIL_EXISTS:
          case UserErrorCode.USERNAME_EXISTS:
          case UserErrorCode.INVALID_EMAIL:
          case UserErrorCode.INVALID_USERNAME:
          case UserErrorCode.INVALID_PASSWORD:
          case UserErrorCode.USERNAME_CHANGE_LIMIT:
            return {
              success: false,
              error: error.message,
              code: error.code,
            };

          default:
            return {
              success: false,
              error: "An unexpected error occurred",
              code: UserErrorCode.SERVER_ERROR,
            };
        }
      }

      return {
        success: false,
        error: "An unexpected error occurred",
        code: UserErrorCode.SERVER_ERROR,
      };
    },
    [logout]
  );

  // Query for user profile
  const {
    data: profile,
    isLoading,
    error: queryError,
  } = useQuery<UserProfile>({
    queryKey: ["userProfile"],
    queryFn: userService.getProfile,
    enabled: isAuthenticated,
    retry: (failureCount, error) => {
      if (error instanceof UserServiceError) {
        return error.code !== UserErrorCode.UNAUTHORIZED && failureCount < 3;
      }
      return failureCount < 3;
    },
  });

  // Mutations with improved CSRF handling
  const updateEmailMutation = useMutation<ApiResponse, Error, string>({
    mutationFn: userService.updateEmail,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["userProfile"] });
    },
    onError: async (error) => {
      const response = await handleError(error);
      throw new Error(response.error);
    },
  });

  const changePasswordMutation = useMutation<
    ApiResponse,
    Error,
    { currentPassword: string; newPassword: string }
  >({
    mutationFn: ({ currentPassword, newPassword }) =>
      userService.changePassword(currentPassword, newPassword),
    onError: async (error) => {
      const response = await handleError(error);
      throw new Error(response.error);
    },
  });

  const changeUsernameMutation = useMutation<ApiResponse, Error, string>({
    mutationFn: userService.changeUsername,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["userProfile"] });
    },
    onError: async (error) => {
      const response = await handleError(error);
      throw new Error(response.error);
    },
  });

  const deleteAccountMutation = useMutation<ApiResponse, Error, void>({
    mutationFn: userService.deleteAccount,
    onSuccess: async () => {
      await logout();
      queryClient.clear();
      router.push("/login");
    },
    onError: async (error) => {
      const response = await handleError(error);
      throw new Error(response.error);
    },
  });

  // Error handling for query
  useEffect(() => {
    if (queryError) {
      handleError(queryError).then((response) => {
        if (response.code === UserErrorCode.UNAUTHORIZED) {
          router.push("/login");
        }
      });
    }
  }, [queryError, router, handleError]);

  // Wrapper methods with improved error handling
  const updateEmail = async (newEmail: string): Promise<MutationResponse> => {
    try {
      await updateEmailMutation.mutateAsync(newEmail);
      return { success: true };
    } catch (error) {
      return await handleError(error);
    }
  };

  const changePassword = async (
    currentPassword: string,
    newPassword: string
  ): Promise<MutationResponse> => {
    try {
      await changePasswordMutation.mutateAsync({
        currentPassword,
        newPassword,
      });
      return { success: true };
    } catch (error) {
      return await handleError(error);
    }
  };

  const changeUsername = async (
    newUsername: string
  ): Promise<MutationResponse> => {
    try {
      await changeUsernameMutation.mutateAsync(newUsername);
      return { success: true };
    } catch (error) {
      return await handleError(error);
    }
  };

  const deleteAccount = async (): Promise<MutationResponse> => {
    try {
      await deleteAccountMutation.mutateAsync();
      return { success: true };
    } catch (error) {
      return await handleError(error);
    }
  };

  // Mutation states for the UI
  const mutationStates: Record<string, MutationState> = {
    updateEmail: {
      isError: updateEmailMutation.isError,
      error: updateEmailMutation.error as Error | null,
      isPending: updateEmailMutation.isPending,
      errorCode: (updateEmailMutation.error as UserServiceError)?.code,
    },
    changePassword: {
      isError: changePasswordMutation.isError,
      error: changePasswordMutation.error as Error | null,
      isPending: changePasswordMutation.isPending,
      errorCode: (changePasswordMutation.error as UserServiceError)?.code,
    },
    changeUsername: {
      isError: changeUsernameMutation.isError,
      error: changeUsernameMutation.error as Error | null,
      isPending: changeUsernameMutation.isPending,
      errorCode: (changeUsernameMutation.error as UserServiceError)?.code,
    },
    deleteAccount: {
      isError: deleteAccountMutation.isError,
      error: deleteAccountMutation.error as Error | null,
      isPending: deleteAccountMutation.isPending,
      errorCode: (deleteAccountMutation.error as UserServiceError)?.code,
    },
  };

  return {
    profile,
    isLoading,
    error: queryError,
    mutationStates,
    updateEmail,
    changePassword,
    changeUsername,
    deleteAccount,
  };
}
