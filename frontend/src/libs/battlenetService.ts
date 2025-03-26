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

// Detailed error codes
export enum BattleNetErrorCode {
  ALREADY_LINKED = "battle_net_already_linked",
  LINK_FAILED = "battle_net_link_failed",
  UNLINK_FAILED = "battle_net_unlink_failed",
  NETWORK_ERROR = "battle_net_network_error",
  UNAUTHORIZED = "battle_net_unauthorized",
  CALLBACK_FAILED = "battle_net_callback_failed",
  INVALID_STATE = "battle_net_invalid_state",
  INVALID_CODE = "battle_net_invalid_code",
  INVALID_GRANT = "battle_net_invalid_grant",
  TOKEN_EXPIRED = "battle_net_token_expired",
  SCOPE_MISSING = "battle_net_scope_missing",
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
  // Initiates the Battle.net OAuth flow
  async initiateLinking(): Promise<{ url: string }> {
    try {
      const response = await api.get<{ url: string }>("/auth/battle-net/link", {
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
