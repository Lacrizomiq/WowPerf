// src/hooks/useGoogleAuth.ts

import { useMutation } from "@tanstack/react-query";
import { googleAuthService } from "@/libs/googleAuthService";
import { AuthError } from "@/libs/authService";
import { useAuth } from "@/providers/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { getOAuthErrorDisplay } from "@/utils/auth/oauthErrors";

// ===== HOOK POUR INITIER GOOGLE LOGIN =====

export const useGoogleLogin = () => {
  return useMutation({
    mutationFn: async () => {
      // Simple appel - le backend redirige automatiquement
      await googleAuthService.initiateGoogleLogin();
    },
    onError: (error) => {
      console.error("Failed to initiate Google login:", error);
      // L'erreur sera gérée par le composant qui utilise le hook
    },
  });
};

// ===== HOOK POUR LA PAGE CALLBACK =====

export interface GoogleCallbackState {
  isProcessing: boolean;
  error: AuthError | null;
  isSuccess: boolean;
  errorDisplay: ReturnType<typeof getOAuthErrorDisplay> | null;
}

export const useGoogleCallback = () => {
  const [state, setState] = useState<GoogleCallbackState>({
    isProcessing: true,
    error: null,
    isSuccess: false,
    errorDisplay: null,
  });

  const { checkAuth } = useAuth();
  const router = useRouter();

  useEffect(() => {
    const processCallback = async () => {
      // 1. Parser les paramètres de callback
      const params = googleAuthService.parseCallbackParams();

      // 2. Vérifier si erreur dans les params (de Google)
      if (googleAuthService.hasCallbackError(params)) {
        const authError = googleAuthService.mapCallbackError(params);
        const errorDisplay = getOAuthErrorDisplay(authError.code);

        setState({
          isProcessing: false,
          error: authError,
          isSuccess: false,
          errorDisplay,
        });

        googleAuthService.cleanupCallbackUrl();
        return;
      }

      // 3. Vérifier si erreur backend dans l'URL
      const backendError = googleAuthService.extractRedirectError();
      if (backendError) {
        const authError = googleAuthService.mapBackendError(
          backendError.code,
          backendError.message
        );
        const errorDisplay = getOAuthErrorDisplay(authError.code);

        setState({
          isProcessing: false,
          error: authError,
          isSuccess: false,
          errorDisplay,
        });

        googleAuthService.cleanupCallbackUrl();
        return;
      }

      // 4. Pas d'erreur visible - vérifier l'auth
      try {
        // ✅ OPTIMISATION : Attendre que les cookies soient bien définis
        await new Promise((resolve) => setTimeout(resolve, 200));

        // Vérifier l'état d'authentification
        await checkAuth();

        setState({
          isProcessing: false,
          error: null,
          isSuccess: true,
          errorDisplay: null,
        });

        googleAuthService.cleanupCallbackUrl();

        // ✅ OPTIMISATION : Délai plus court pour la redirection
        setTimeout(() => {
          // Vérifier s'il y a un query param new_user
          const urlParams = new URLSearchParams(window.location.search);
          const isNewUser = urlParams.get("new_user") === "true";

          // Redirection conditionnelle
          router.push(isNewUser ? "/" : "/profile");
        }, 1000); // 1 seconde pour voir le succès
      } catch (error) {
        console.error("Auth check failed after callback:", error);

        const authError = new AuthError(
          "auth_processing_failed" as any,
          "Failed to verify authentication after Google login"
        );
        const errorDisplay = getOAuthErrorDisplay(authError.code);

        setState({
          isProcessing: false,
          error: authError,
          isSuccess: false,
          errorDisplay,
        });

        googleAuthService.cleanupCallbackUrl();
      }
    };

    processCallback();
  }, [checkAuth, router]);

  return state;
};

// ===== HOOK COMBINÉ (OPTIONNEL) =====

export const useGoogleAuth = () => {
  const loginMutation = useGoogleLogin();

  return {
    // Méthodes
    loginWithGoogle: loginMutation.mutate,

    // États
    isLoading: loginMutation.isPending,
    error: loginMutation.error,

    // Pour réinitialiser après erreur
    reset: loginMutation.reset,
  };
};
