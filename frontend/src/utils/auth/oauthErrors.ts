// src/utils/auth/oauthErrors.ts

import { AuthErrorCode } from "@/libs/authService";
import { OAuthErrorDisplay } from "@/types/auth/oauth.types";

/**
 * Mapping of OAuth error codes to user-friendly display information
 */
export const OAUTH_ERROR_DISPLAYS: Record<string, OAuthErrorDisplay> = {
  [AuthErrorCode.OAUTH_CANCELLED]: {
    code: AuthErrorCode.OAUTH_CANCELLED,
    title: "Authentication Cancelled",
    message: "You cancelled the Google sign-in process.",
    actions: [
      {
        label: "Try Again",
        href: "/login",
      },
    ],
    recoverable: true,
  },
  [AuthErrorCode.OAUTH_FAILED]: {
    code: AuthErrorCode.OAUTH_FAILED,
    title: "Authentication Failed",
    message: "Google sign-in failed. Please try again.",
    actions: [
      {
        label: "Try Again",
        href: "/login",
      },
      {
        label: "Use Email/Password",
        href: "/login",
      },
    ],
    recoverable: true,
  },
  [AuthErrorCode.OAUTH_PROCESSING_FAILED]: {
    code: AuthErrorCode.OAUTH_PROCESSING_FAILED,
    title: "Processing Error",
    message: "Failed to process authentication. Please try again.",
    actions: [
      {
        label: "Try Again",
        href: "/login",
      },
    ],
    recoverable: true,
  },
  [AuthErrorCode.EMAIL_ALREADY_LINKED]: {
    code: AuthErrorCode.EMAIL_ALREADY_LINKED,
    title: "Email Already in Use",
    message:
      "This email is already associated with another account. Please sign in with your existing account.",
    actions: [
      {
        label: "Sign In",
        href: "/login",
      },
      {
        label: "Forgot Password?",
        href: "/forgot-password",
      },
    ],
    recoverable: true,
  },
  [AuthErrorCode.OAUTH_STATE_MISMATCH]: {
    code: AuthErrorCode.OAUTH_STATE_MISMATCH,
    title: "Security Error",
    message:
      "A security check failed during authentication. This might happen if you took too long to sign in.",
    actions: [
      {
        label: "Try Again",
        href: "/login",
      },
    ],
    recoverable: true,
  },
  // Fallback for unknown errors
  default: {
    code: "unknown",
    title: "Authentication Error",
    message: "An unexpected error occurred during sign-in.",
    actions: [
      {
        label: "Try Again",
        href: "/login",
      },
      {
        label: "Contact Support",
        href: "/support",
      },
    ],
    recoverable: false,
  },
};

/**
 * Get display information for an OAuth error
 */
export function getOAuthErrorDisplay(errorCode: string): OAuthErrorDisplay {
  return OAUTH_ERROR_DISPLAYS[errorCode] || OAUTH_ERROR_DISPLAYS.default;
}
