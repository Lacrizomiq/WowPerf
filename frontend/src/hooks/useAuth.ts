import { useMutation, useQueryClient } from "@tanstack/react-query";
import { authService } from "@/libs/authService";
import { useRouter } from "next/navigation";

// Hook to signup a user
export const useSignup = () => {
  const router = useRouter();

  return useMutation({
    mutationFn: authService.signup,
    onSuccess: () => {
      router.push("/login");
    },
  });
};

// Hook to login a user
export const useLogin = () => {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: authService.login,
    onSuccess: (token) => {
      queryClient.clear();
      router.push("/");
    },
  });
};

// Hook to logout a user
export const useLogout = () => {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: authService.logout,
    onSuccess: () => {
      queryClient.clear();
      router.push("/");
    },
  });
};

// Hook to check if the user is authenticated
export const useIsAuthenticated = () => {
  return authService.isAuthenticated();
};
