"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import { userService } from "@/libs/userService";
import { authService } from "@/libs/authService";
import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import { AxiosError } from "axios";

export function useUserProfile() {
  const queryClient = useQueryClient();
  const router = useRouter();

  const { isAuthenticated } = useAuth();

  const {
    data: profile,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["userProfile"],
    queryFn: userService.getProfile,
    enabled: isAuthenticated, // Only run the query if the user is authenticated
  });

  useEffect(() => {
    if (error instanceof AxiosError && error.response?.status === 401) {
      router.push("/login");
    }
  }, [error, router]);

  const updateEmailMutation = useMutation({
    mutationFn: userService.updateEmail,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["userProfile"] });
    },
  });

  const changePasswordMutation = useMutation({
    mutationFn: ({
      currentPassword,
      newPassword,
    }: {
      currentPassword: string;
      newPassword: string;
    }) => userService.changePassword(currentPassword, newPassword),
  });

  const changeUsernameMutation = useMutation({
    mutationFn: userService.changeUsername,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["userProfile"] });
    },
  });

  const deleteAccountMutation = useMutation({
    mutationFn: userService.deleteAccount,
    onSuccess: async () => {
      await authService.logout();
      queryClient.clear();
      router.push("/login");
    },
  });

  return {
    profile,
    isLoading,
    error,
    updateEmail: (newEmail: string): Promise<void> =>
      updateEmailMutation.mutateAsync(newEmail),
    changePassword: (
      currentPassword: string,
      newPassword: string
    ): Promise<void> =>
      changePasswordMutation.mutateAsync({ currentPassword, newPassword }),
    changeUsername: (newUsername: string): Promise<void> =>
      changeUsernameMutation.mutateAsync(newUsername),
    deleteAccount: (): Promise<void> => deleteAccountMutation.mutateAsync(),
    isUpdatingEmail: updateEmailMutation.isPending,
    isChangingPassword: changePasswordMutation.isPending,
    isChangingUsername: changeUsernameMutation.isPending,
    isDeletingAccount: deleteAccountMutation.isPending,
  };
}
