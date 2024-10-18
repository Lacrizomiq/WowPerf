"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import { userService } from "@/libs/userService";
import { authService } from "@/libs/authService";
import { useRouter } from "next/navigation";
import { AxiosError } from "axios";

export function useUserProfile() {
  const queryClient = useQueryClient();
  const router = useRouter();

  const {
    data: profile,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["userProfile"],
    queryFn: async () => {
      console.log("Fetching user profile...");
      const data = await userService.getProfile();
      console.log("Profile data received:", data);
      return data;
    },
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
    updateEmail: (newEmail: string) => updateEmailMutation.mutate(newEmail),
    changePassword: (currentPassword: string, newPassword: string) =>
      changePasswordMutation.mutate({ currentPassword, newPassword }),
    changeUsername: (newUsername: string) =>
      changeUsernameMutation.mutate(newUsername),
    deleteAccount: () => deleteAccountMutation.mutate(),
    isUpdatingEmail: updateEmailMutation.isPending,
    isChangingPassword: changePasswordMutation.isPending,
    isChangingUsername: changeUsernameMutation.isPending,
    isDeletingAccount: deleteAccountMutation.isPending,
  };
}
