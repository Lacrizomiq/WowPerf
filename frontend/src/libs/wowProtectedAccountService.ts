import api, { APIError } from "./api";
import axios, { AxiosError } from "axios";
import {
  WoWProfile,
  CharacterBasicInfo,
  UserCharacter,
  SyncResult,
  RefreshResult,
  CharacterDisplayResponse,
  FavoriteCharacterResponse,
  WoWError,
  WoWErrorCode,
} from "@/types/userCharacter/userCharacter";

/**
 * Service to interact with the Battle.net API and endpoints
 * related to WoW characters
 */
export const wowService = {
  /**
   * Fetch the WoW profile of the user
   */
  async getWoWProfile(region: string = "eu"): Promise<WoWProfile> {
    try {
      const response = await api.get<WoWProfile>("/wow/profile", {
        headers: {
          Region: region,
          Accept: "application/json",
        },
        withCredentials: true,
      });

      // Extract the region from the URL in the response
      const extractedRegion =
        response.data._links.self.href.match(
          /https:\/\/(\w+)\.api\.blizzard\.com/
        )?.[1] || region;

      // Filter to keep only characters level 80+
      const filteredAccounts = response.data.wow_accounts.map((account) => ({
        ...account,
        characters: account.characters.filter(
          (character) => character.level >= 80
        ),
      }));

      return {
        ...response.data,
        region: extractedRegion,
        wow_accounts: filteredAccounts,
      };
    } catch (error) {
      return this.handleApiError(error, "Failed to fetch WoW profile");
    }
  },

  /**
   * Fetch the details of a protected character
   */
  async getProtectedCharacter(
    realmId: number,
    characterId: number,
    region: string = "eu",
    locale: string = "en_GB"
  ): Promise<any> {
    try {
      const response = await api.get(`/wow/profile/protected-character`, {
        headers: {
          Region: region,
          Accept: "application/json",
        },
        params: {
          realmId,
          characterId,
          locale,
        },
        withCredentials: true,
      });
      return response.data;
    } catch (error) {
      return this.handleApiError(
        error,
        "Failed to fetch protected character profile"
      );
    }
  },

  /**
   * List all available characters of the Battle.net account
   */
  async listAccountCharacters(
    region: string = "eu"
  ): Promise<CharacterBasicInfo[]> {
    try {
      const response = await api.get<CharacterBasicInfo[]>(
        "/wow/profile/characters",
        {
          headers: {
            Region: region,
            Accept: "application/json",
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      return this.handleApiError(error, "Failed to list account characters");
    }
  },

  /**
   * Synchronize all level 80+ characters of the account
   */
  async syncAllAccountCharacters(region: string = "eu"): Promise<SyncResult> {
    try {
      const response = await api.post<SyncResult>(
        "/wow/profile/characters/sync",
        {},
        {
          headers: {
            Region: region,
            Accept: "application/json",
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      return this.handleApiError(error, "Failed to synchronize characters");
    }
  },

  /**
   * Refresh the existing characters and add the new ones
   */
  async refreshUserCharacters(region: string = "eu"): Promise<RefreshResult> {
    try {
      const response = await api.post<RefreshResult>(
        "/wow/profile/characters/refresh",
        {},
        {
          headers: {
            Region: region,
            Accept: "application/json",
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      return this.handleApiError(error, "Failed to refresh characters");
    }
  },

  /**
   * Fetch the synchronized characters of the user
   */
  async getUserCharacters(): Promise<UserCharacter[]> {
    try {
      const response = await api.get<UserCharacter[]>(
        "/wow/profile/user/characters",
        {
          headers: {
            Accept: "application/json",
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      return this.handleApiError(error, "Failed to get user characters");
    }
  },

  /**
   * Set a character as favorite
   */
  async setFavoriteCharacter(
    characterId: number
  ): Promise<{ message: string }> {
    try {
      const response = await api.put<{ message: string }>(
        `/wow/profile/characters/${characterId}/favorite`,
        {},
        {
          headers: {
            Accept: "application/json",
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      return this.handleApiError(error, "Failed to set favorite character");
    }
  },

  /**
   * Enable or disable the display of a character
   */
  async toggleCharacterDisplay(
    characterId: number,
    display: boolean
  ): Promise<{ message: string }> {
    try {
      const response = await api.put<{ message: string }>(
        `/wow/profile/characters/${characterId}/display`,
        { display },
        {
          headers: {
            Accept: "application/json",
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error) {
      return this.handleApiError(error, "Failed to toggle character display");
    }
  },

  /**
   * Centralized error handling
   */
  handleApiError(error: unknown, defaultMessage: string): never {
    // Axios errors (HTTP errors)
    if (axios.isAxiosError(error)) {
      const err = error as AxiosError<APIError>;

      switch (err.response?.status) {
        case 401:
          throw new WoWError(
            WoWErrorCode.UNAUTHORIZED,
            "Authentication required. Please link your Battle.net account."
          );
        case 403:
          throw new WoWError(
            WoWErrorCode.TOKEN_EXPIRED,
            "Your session has expired. Please re-link your Battle.net account."
          );
        case 404:
          throw new WoWError(
            WoWErrorCode.NOT_FOUND,
            "The requested resource was not found."
          );
        case 500:
        case 502:
        case 503:
        case 504:
          throw new WoWError(
            WoWErrorCode.SERVER_ERROR,
            "A server error occurred. Please try again later."
          );
      }

      // For other HTTP error codes
      if (err.response?.data?.error) {
        throw new WoWError(
          WoWErrorCode.NETWORK_ERROR,
          err.response.data.error,
          error
        );
      }
    }

    // Non-Axios errors (JS errors, network errors, etc.)
    throw new WoWError(WoWErrorCode.NETWORK_ERROR, defaultMessage, error);
  },
};
