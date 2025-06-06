// src/libs/googleAuthService.ts

import api from "./api";
import { AuthError, AuthErrorCode } from "./authService";
import {
  GoogleOAuthCallbackParams,
  OAuthErrorResponse,
} from "@/types/auth/oauth.types";
import axios, { AxiosError } from "axios";

// Service to handle Google OAuth authentication flows
export const googleAuthService = {
  // Initiate the Google OAuth login flow
  // This will redirect the user to Google OAuth consent screen
  async initiateGoogleLogin(): Promise<void> {
    try {
      // fonctionne en dev/prod automatiquement
      window.location.href = "/api/auth/google/login";
    } catch (error) {
      // Gestion d'erreurs
      console.log("Failed to initiate Google login", error);

      throw new AuthError(
        AuthErrorCode.OAUTH_FAILED,
        "Failed to start Google sign-in"
      );
    }
  },

  // Parses the OAuth callback parameters from the URL
  // Used on the callback page to extract code, state, and error params
  parseCallbackParams(): GoogleOAuthCallbackParams {
    const params = new URLSearchParams(window.location.search);

    return {
      code: params.get("code") || undefined,
      state: params.get("state") || undefined,
      error: params.get("error") || undefined,
      error_description: params.get("error_description") || undefined,
      scope: params.get("scope") || undefined,
      authuser: params.get("authuser") || undefined,
      prompt: params.get("prompt") || undefined,
    };
  },

  // Validates if the callback contains an error
  hasCallbackError(params: GoogleOAuthCallbackParams): boolean {
    return !!params.error;
  },

  // Maps OAuth callback errors to our AuthError system
  mapCallbackError(params: GoogleOAuthCallbackParams): AuthError {
    const { error, error_description } = params;

    // Map common OAuth errors to our error codes
    switch (error) {
      case "access_denied":
        return new AuthError(
          AuthErrorCode.OAUTH_CANCELLED,
          "You cancelled the Google sign-in process"
        );

      case "invalid_request":
      case "invalid_grant":
      case "unsupported_grant_type":
        return new AuthError(
          AuthErrorCode.OAUTH_INVALID_CALLBACK,
          error_description || "Invalid authentication request"
        );

      case "server_error":
      case "temporarily_unavailable":
        return new AuthError(
          AuthErrorCode.OAUTH_FAILED,
          "Google authentication service is temporarily unavailable"
        );

      default:
        return new AuthError(
          AuthErrorCode.OAUTH_FAILED,
          error_description || "Google sign-in failed"
        );
    }
  },

  // Extracts error information from the redirect URL
  // This is for errors that come from our backend after processing
  extractRedirectError(): { code?: string; message?: string } | null {
    const params = new URLSearchParams(window.location.search);
    const errorCode = params.get("error");
    const errorMessage = params.get("message");

    if (errorCode || errorMessage) {
      return {
        code: errorCode || undefined,
        message: errorMessage || undefined,
      };
    }

    return null;
  },

  // Maps backend OAuth errors to AuthError
  mapBackendError(code?: string, message?: string): AuthError {
    if (!code) {
      return new AuthError(
        AuthErrorCode.OAUTH_FAILED,
        message || "Authentication failed"
      );
    }

    // Map backend error codes to our frontend error codes
    const errorMap: Record<string, AuthErrorCode> = {
      auth_cancelled: AuthErrorCode.OAUTH_CANCELLED,
      auth_failed: AuthErrorCode.OAUTH_FAILED,
      auth_processing_failed: AuthErrorCode.OAUTH_PROCESSING_FAILED,
      invalid_callback: AuthErrorCode.OAUTH_INVALID_CALLBACK,
      token_exchange_failed: AuthErrorCode.OAUTH_TOKEN_EXCHANGE_FAILED,
      user_info_failed: AuthErrorCode.OAUTH_USER_INFO_FAILED,
      email_already_linked: AuthErrorCode.EMAIL_ALREADY_LINKED,
      state_mismatch: AuthErrorCode.OAUTH_STATE_MISMATCH,
    };

    const mappedCode = errorMap[code] || AuthErrorCode.OAUTH_FAILED;

    return new AuthError(mappedCode, message || "Authentication failed");
  },

  // Clears OAuth-related data from the URL without page reload
  // Useful after processing callback params to clean up the URL
  cleanupCallbackUrl(): void {
    const url = new URL(window.location.href);

    // Remove all OAuth-related query parameters
    const paramsToRemove = [
      "code",
      "state",
      "error",
      "error_description",
      "scope",
      "authuser",
      "prompt",
      "message",
    ];

    paramsToRemove.forEach((param) => url.searchParams.delete(param));

    // Update URL without reload
    window.history.replaceState({}, document.title, url.toString());
  },
};
