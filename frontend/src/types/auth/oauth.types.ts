// src/types/auth/oauth.types.ts
import { AuthMethod } from "@/libs/authService";

// Google OAuth callback parameters from the OAuth flow
export interface GoogleOAuthCallbackParams {
  /** Authorization code from Google */
  code?: string;
  /** State parameter for CSRF protection */
  state?: string;
  /** Error code if the OAuth flow failed */
  error?: string;
  /** Error description providing more details */
  error_description?: string;
  /** The scope that was granted (on success) */
  scope?: string;
  /** The user who authorized (on success) */
  authuser?: string;
  /** The prompt parameter that was used */
  prompt?: string;
}

// Response from initiating Google OAuth login
export interface GoogleOAuthInitResponse {
  /** The URL to redirect the user to Google */
  auth_url: string;
  /** State parameter for CSRF validation */
  state: string;
}

// OAuth provider information
export interface OAuthProvider {
  /** Unique identifier for the provider */
  id: string;
  /** Display name of the provider */
  name: string;
  /** Icon component or URL */
  icon?: string;
  /** Whether this provider is enabled */
  enabled: boolean;
}

// OAuth account link status
export interface OAuthLinkStatus {
  /** The OAuth provider (e.g., "google") */
  provider: string;
  /** Whether the account is linked */
  linked: boolean;
  /** The email associated with the OAuth account */
  email?: string;
  /** When the account was linked */
  linked_at?: string;
}

// Result of OAuth authentication
export interface OAuthAuthResult {
  /** Whether this is a new user registration */
  is_new_user: boolean;
  /** The authentication method used */
  method: "login" | "signup" | "link";
  /** User information */
  user: {
    id: number;
    username: string;
    email: string;
    auth_methods: AuthMethod[];
  };
  /** Optional message for the user */
  message?: string;
}

// OAuth error response structure
export interface OAuthErrorResponse {
  /** Error code matching AuthErrorCode enum */
  code: string;
  /** Human-readable error message */
  message: string;
  /** Additional error details */
  details?: string;
  /** The OAuth provider that failed */
  provider?: string;
}

// Configuration for OAuth providers
export interface OAuthConfig {
  /** Google OAuth configuration */
  google: {
    enabled: boolean;
    clientId?: string;
    /** Scopes requested from Google */
    scopes: string[];
  };
  // Future providers
  // discord?: { ... }
}

// User authentication state with OAuth information
export interface AuthStateWithOAuth {
  /** Whether the user is authenticated */
  isAuthenticated: boolean;
  /** Loading state */
  isLoading: boolean;
  /** User information */
  user: {
    username: string;
    email?: string;
    /** Authentication methods available to the user */
    authMethods?: AuthMethod[];
    /** OAuth accounts linked to this user */
    oauthAccounts?: OAuthLinkStatus[];
  } | null;
}

// Parameters for handling OAuth errors in the UI
export interface OAuthErrorDisplay {
  /** The error code */
  code: string;
  /** User-friendly title */
  title: string;
  /** Detailed error message */
  message: string;
  /** Suggested actions for the user */
  actions?: Array<{
    label: string;
    href?: string;
    onClick?: () => void;
  }>;
  /** Whether the error is recoverable */
  recoverable: boolean;
}

// OAuth state stored during the flow
export interface OAuthFlowState {
  /** The OAuth provider being used */
  provider: string;
  /** Timestamp when the flow started */
  initiated_at: number;
  /** Optional redirect URL after successful auth */
  redirect_to?: string;
  /** Any additional state data */
  metadata?: Record<string, unknown>;
}
