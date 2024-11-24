import api, { APIError } from "./api";
import axios, { AxiosError } from "axios";

export interface WoWProfile {
  sub: string;
  id: number;
  region?: string;
  _links: {
    self: {
      href: string;
    };
  };
  wow_accounts: Array<{
    characters: Array<{
      name: string;
      realm: {
        name: string;
        slug: string;
      };
      playable_class: {
        id: number;
        name: string;
      };
      level: number;
    }>;
  }>;
}

export enum WoWErrorCode {
  UNAUTHORIZED = "wow_unauthorized",
  TOKEN_EXPIRED = "wow_token_expired",
  NETWORK_ERROR = "wow_network_error",
  BATTLE_NET_NOT_LINKED = "battle_net_not_linked",
}

export class WoWError extends Error {
  constructor(
    public code: WoWErrorCode,
    message: string,
    public originalError?: unknown
  ) {
    super(message);
    this.name = "WoWError";
  }
}

export const wowService = {
  async getWoWProfile(): Promise<WoWProfile> {
    try {
      const response = await api.get<WoWProfile>("/wow/profile", {
        headers: {
          Region: "eu",
          Accept: "application/json",
        },
        withCredentials: true,
      });

      // Extract the region from the _links.self.href
      const region =
        response.data._links.self.href.match(
          /https:\/\/(\w+)\.api\.blizzard\.com/
        )?.[1] || "eu";

      const filteredAccounts = response.data.wow_accounts.map((account) => ({
        ...account,
        characters: account.characters.filter(
          (character) => character.level >= 80
        ),
      }));

      return {
        ...response.data,
        region,
        wow_accounts: filteredAccounts,
      };
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        if (err.response?.status === 401) {
          throw new WoWError(
            WoWErrorCode.UNAUTHORIZED,
            "Battle.net authentication required"
          );
        }
      }
      throw new WoWError(
        WoWErrorCode.NETWORK_ERROR,
        "Failed to fetch WoW profile"
      );
    }
  },

  async getProtectedCharacter(
    realmId: number,
    characterId: number,
    region: string = "eu",
    locale: string = "en_GB"
  ): Promise<any> {
    try {
      const response = await api.get(
        `/wow/protected-character/${realmId}-${characterId}`,
        {
          headers: {
            Region: region,
            Accept: "application/json",
          },
          params: {
            namespace: `profile-${region}`,
            locale,
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        if (err.response?.status === 401) {
          throw new WoWError(
            WoWErrorCode.UNAUTHORIZED,
            "Battle.net authentication required"
          );
        }
      }
      throw new WoWError(
        WoWErrorCode.NETWORK_ERROR,
        "Failed to fetch protected character profile"
      );
    }
  },
};
