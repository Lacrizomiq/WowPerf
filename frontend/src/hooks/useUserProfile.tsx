// src/hooks/useUserProfile.tsx
"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useEffect, useCallback } from "react";
import {
  userService,
  UserServiceError,
  UserErrorCode,
} from "@/libs/userService";
import { useAuth } from "@/providers/AuthContext";
import { useRouter } from "next/navigation";
import { AxiosError } from "axios";
import { csrfService } from "@/libs/csrfService";
import { APIError } from "@/libs/api";

// Mutation states
interface MutationState {
  isError: boolean;
  error: Error | null;
  isPending: boolean;
}

export function useUserProfile() {
  const queryClient = useQueryClient();
  const router = useRouter();
  const { isAuthenticated, logout } = useAuth();

  // Centralized error handler
  const handleError = useCallback(
    async (error: unknown) => {
      console.error("User profile error:", error);

      if (error instanceof UserServiceError) {
        switch (error.code) {
          case UserErrorCode.UNAUTHORIZED:
            csrfService.clearToken();
            await logout();
            return "Session expired. Please login again.";
          case UserErrorCode.UNAUTHORIZED:
            // Proper logout if unauthorized
            await logout();
            return "Session expired. Please login again.";

          case UserErrorCode.PROFILE_NOT_FOUND:
            return "Profile not found";

          case UserErrorCode.EMAIL_EXISTS:
            return "This email is already in use";

          case UserErrorCode.INVALID_EMAIL:
            return "Invalid email format";

          case UserErrorCode.USERNAME_EXISTS:
            return "This username is already taken";

          case UserErrorCode.USERNAME_CHANGE_LIMIT:
            return "Username can only be changed once every 30 days";

          case UserErrorCode.INVALID_PASSWORD:
            return "Invalid password format";

          case UserErrorCode.NETWORK_ERROR:
            return "Network error. Please try again";

          default:
            return error.message;
        }
      }

      if (error instanceof AxiosError) {
        const err = error as AxiosError<APIError>;
        if (err.response?.status === 401) {
          csrfService.clearToken();
          await logout();
          return "Session expired. Please login again.";
        }
      }

      return "An unexpected error occurred";
    },
    [logout]
  );

  // Query for user profile
  const {
    data: profile,
    isLoading,
    error: queryError,
  } = useQuery({
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

  // Mutation for email update
  const updateEmailMutation = useMutation({
    mutationFn: userService.updateEmail,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["userProfile"] });
    },
    onError: async (error) => {
      const errorMessage = await handleError(error);
      throw new Error(errorMessage);
    },
  });

  // Mutation for password change
  const changePasswordMutation = useMutation({
    mutationFn: ({
      currentPassword,
      newPassword,
    }: {
      currentPassword: string;
      newPassword: string;
    }) => userService.changePassword(currentPassword, newPassword),
    onError: async (error) => {
      const errorMessage = await handleError(error);
      throw new Error(errorMessage);
    },
  });

  // Mutation for username change
  const changeUsernameMutation = useMutation({
    mutationFn: userService.changeUsername,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["userProfile"] });
    },
    onError: async (error) => {
      const errorMessage = await handleError(error);
      throw new Error(errorMessage);
    },
  });

  // Mutation for account deletion
  const deleteAccountMutation = useMutation({
    mutationFn: userService.deleteAccount,
    onSuccess: async () => {
      await logout();
      queryClient.clear();
    },
    onError: async (error) => {
      const errorMessage = await handleError(error);
      throw new Error(errorMessage);
    },
  });

  // Query error handling
  useEffect(() => {
    const handleQueryError = async () => {
      if (queryError) {
        const errorMessage = await handleError(queryError);
        console.error(errorMessage);

        if (
          queryError instanceof UserServiceError &&
          queryError.code === UserErrorCode.UNAUTHORIZED
        ) {
          router.push("/login");
        }
      }
    };

    handleQueryError();
  }, [queryError, router, handleError]);

  // Wrapped methods with error handling
  const updateEmail = async (newEmail: string) => {
    try {
      await updateEmailMutation.mutateAsync(newEmail);
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error ? error.message : "Failed to update email",
      };
    }
  };

  const changePassword = async (
    currentPassword: string,
    newPassword: string
  ) => {
    try {
      await changePasswordMutation.mutateAsync({
        currentPassword,
        newPassword,
      });
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error ? error.message : "Failed to change password",
      };
    }
  };

  const changeUsername = async (newUsername: string) => {
    try {
      await changeUsernameMutation.mutateAsync(newUsername);
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error ? error.message : "Failed to change username",
      };
    }
  };

  const deleteAccount = async () => {
    try {
      await deleteAccountMutation.mutateAsync();
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error ? error.message : "Failed to delete account",
      };
    }
  };

  // Mutation states for UI
  const mutationStates: Record<string, MutationState> = {
    updateEmail: {
      isError: updateEmailMutation.isError,
      error: updateEmailMutation.error as Error | null,
      isPending: updateEmailMutation.isPending,
    },
    changePassword: {
      isError: changePasswordMutation.isError,
      error: changePasswordMutation.error as Error | null,
      isPending: changePasswordMutation.isPending,
    },
    changeUsername: {
      isError: changeUsernameMutation.isError,
      error: changeUsernameMutation.error as Error | null,
      isPending: changeUsernameMutation.isPending,
    },
    deleteAccount: {
      isError: deleteAccountMutation.isError,
      error: deleteAccountMutation.error as Error | null,
      isPending: deleteAccountMutation.isPending,
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
