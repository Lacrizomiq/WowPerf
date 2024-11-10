import api, { APIError } from "./api";
import axios, { AxiosError } from "axios";

// Response types
interface BattleNetLinkResponse {
  message: string;
  linked: boolean;
  battleTag?: string;
}

interface BattleNetStatusResponse {
  linked: boolean;
  battleTag?: string;
}

// Error handling
export enum BattleNetErrorCode {
  ALREADY_LINKED = "battle_net_already_linked",
  LINK_FAILED = "battle_net_link_failed",
  UNLINK_FAILED = "battle_net_unlink_failed",
  NETWORK_ERROR = "battle_net_network_error",
  UNAUTHORIZED = "battle_net_unauthorized",
}

// Battle.net error
export class BattleNetError extends Error {
  constructor(
    public code: BattleNetErrorCode,
    message: string,
    public originalError?: unknown
  ) {
    super(message);
    this.name = "BattleNetError";
  }
}

// Battle.net service
export const battleNetService = {
  // Initiate Battle.net OAuth flow to link account
  async initiateLinking(): Promise<string> {
    try {
      const response = await api.get<{ url: string }>("/auth/battle-net/link");
      return response.data.url;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        throw new BattleNetError(
          BattleNetErrorCode.LINK_FAILED,
          err.response?.data.error ||
            "Failed to initiate Battle.net OAuth flow",
          err
        );
      }
      throw new BattleNetError(
        BattleNetErrorCode.NETWORK_ERROR,
        "An unexpected network error occurred",
        error
      );
    }
  },

  // Get Battle.net link status
  async getLinkStatus(): Promise<BattleNetStatusResponse> {
    try {
      const response = await api.get<BattleNetStatusResponse>(
        "/auth/battle-net/status"
      );
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        throw new BattleNetError(
          BattleNetErrorCode.UNAUTHORIZED,
          err.response?.data.error || "Failed to get Battle.net link status",
          err
        );
      }
      throw new BattleNetError(
        BattleNetErrorCode.NETWORK_ERROR,
        "An unexpected network error occurred",
        error
      );
    }
  },

  // Unlink Battle.net account
  async unlinkAccount(): Promise<void> {
    try {
      await api.post("/auth/battle-net/unlink");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        throw new BattleNetError(
          BattleNetErrorCode.UNLINK_FAILED,
          err.response?.data.error || "Failed to unlink Battle.net account",
          err
        );
      }
      throw new BattleNetError(
        BattleNetErrorCode.NETWORK_ERROR,
        "An unexpected network error occurred",
        error
      );
    }
  },
};

export type { BattleNetLinkResponse, BattleNetStatusResponse };
