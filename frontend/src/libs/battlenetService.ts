import api, { APIError } from "./api";
import axios, { AxiosError } from "axios";

// More detailed response types
export interface BattleNetLinkResponse {
  linked: boolean;
  battleTag?: string;
  message: string;
  code: string;
  expiresIn?: number;
  scope?: string;
}

export interface BattleNetStatusResponse {
  linked: boolean;
  battleTag?: string;
  error?: string;
  code?: string;
}

// Detailed error codes - UPDATED with new backend error codes
export enum BattleNetErrorCode {
  ALREADY_LINKED = "battle_net_already_linked", // Battle.net account already linked
  LINK_FAILED = "battle_net_link_failed", // Failed to initiate Battle.net linking
  UNLINK_FAILED = "battle_net_unlink_failed", // Failed to unlink Battle.net account
  NETWORK_ERROR = "battle_net_network_error", // Network error during Battle.net linking
  UNAUTHORIZED = "battle_net_unauthorized", // Unauthorized access
  CALLBACK_FAILED = "battle_net_callback_failed", // Failed to handle OAuth callback
  INVALID_STATE = "battle_net_invalid_state", // Invalid OAuth state parameter
  INVALID_CODE = "battle_net_invalid_code", // Invalid OAuth code parameter
  INVALID_GRANT = "battle_net_invalid_grant", // Invalid OAuth grant parameter
  TOKEN_EXPIRED = "battle_net_token_expired", // Token expired
  SCOPE_MISSING = "battle_net_scope_missing", // Missing OAuth scope

  STATE_USER_MISMATCH = "state_user_mismatch", // Security: OAuth state user != authenticated user
  USER_NOT_AUTHENTICATED = "user_not_authenticated", // User session not found in callback
  INVALID_OAUTH_PARAMS = "invalid_oauth_params", // Missing code or state parameters
  TOKEN_EXCHANGE_FAILED = "token_exchange_failed", // Failed to exchange code for token
  AUTH_INITIATION_FAILED = "auth_initiation_failed", // Failed to initiate OAuth flow
  STATUS_CHECK_FAILED = "status_check_failed", // Failed to check Battle.net status
  ACCOUNT_NOT_LINKED = "account_not_linked", // Battle.net account not linked
  TOKEN_NOT_FOUND = "token_not_found", // Token not found in context
  INVALID_TOKEN_FORMAT = "invalid_token_format", // Invalid token format
  PROFILE_FETCH_FAILED = "profile_fetch_failed", // Failed to fetch Battle.net profile
  AUTH_VALIDATION_FAILED = "auth_validation_failed", // Failed to validate authentication
}

export class BattleNetError extends Error {
  constructor(
    public code: BattleNetErrorCode,
    message: string,
    public originalError?: unknown,
    public details?: any
  ) {
    super(message);
    this.name = "BattleNetError";
  }
}

export const battleNetService = {
  // Support du paramètre autoRelink
  async initiateLinking(autoRelink: boolean = false): Promise<{ url: string }> {
    try {
      const params: Record<string, string> = {};
      if (autoRelink) {
        params.auto_relink = "true";
      }

      const response = await api.get<{ url: string }>("/auth/battle-net/link", {
        params,
        withCredentials: true,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        const errorCode =
          err.response?.data?.code || BattleNetErrorCode.LINK_FAILED;
        const errorMessage =
          err.response?.data?.error || "Failed to initiate Battle.net linking";

        throw new BattleNetError(
          errorCode as BattleNetErrorCode,
          errorMessage,
          err,
          err.response?.data
        );
      }
      throw new BattleNetError(
        BattleNetErrorCode.NETWORK_ERROR,
        "Network error during Battle.net linking",
        error
      );
    }
  },

  // Gets the status of the Battle.net link
  async getLinkStatus(): Promise<BattleNetStatusResponse> {
    try {
      const response = await api.get<BattleNetStatusResponse>(
        "/auth/battle-net/status"
      );
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        if (err.response?.status === 401) {
          throw new BattleNetError(
            BattleNetErrorCode.UNAUTHORIZED,
            "Unauthorized access",
            err
          );
        }
        throw new BattleNetError(
          BattleNetErrorCode.NETWORK_ERROR,
          err.response?.data.error || "Failed to get Battle.net status",
          err,
          err.response?.data
        );
      }
      throw new BattleNetError(
        BattleNetErrorCode.NETWORK_ERROR,
        "Network error while getting status",
        error
      );
    }
  },

  // Unlinks the Battle.net account
  async unlinkAccount(): Promise<void> {
    try {
      await api.post("/auth/battle-net/unlink");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        throw new BattleNetError(
          BattleNetErrorCode.UNLINK_FAILED,
          err.response?.data.error || "Failed to unlink Battle.net account",
          err,
          err.response?.data
        );
      }
      throw new BattleNetError(
        BattleNetErrorCode.NETWORK_ERROR,
        "Network error while unlinking account",
        error
      );
    }
  },
};
