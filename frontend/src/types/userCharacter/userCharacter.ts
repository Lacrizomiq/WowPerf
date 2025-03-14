// types/userCharacter/userCharacter.ts
export interface UserCharacter {
  id: number;
  user_id: number;

  // Character identifiers
  character_id: number;
  name: string;
  realm: string;
  region: string;

  // Basic information
  class: string;
  race: string;
  gender: string;
  faction: string;
  active_spec_name: string;
  active_spec_id: number;
  active_spec_role: string;

  // Level and main stats
  level: number;
  item_level: number;
  mythic_plus_rating: number;
  mythic_plus_rating_color: string;
  achievement_points: number;
  honorable_kills: number;

  // Image URLs
  avatar_url: string;
  inset_avatar_url: string;
  main_raw_url: string;
  profile_url: string;

  // JSON structured data (optional for now)
  equipment_json?: Record<string, any>;
  stats_json?: Record<string, any>;
  talents_json?: Record<string, any>;
  mythic_plus_json?: Record<string, any>;
  raids_json?: Record<string, any>;

  // Metadata
  is_displayed: boolean;
  last_api_update: string;
  created_at: string;
  updated_at: string;
}

// CharacterBasicInfo contains basic informations about a character
export interface CharacterBasicInfo {
  character_id: number;
  name: string;
  realm: string;
  region: string;
  class: string;
  race: string;
  level: number;
  faction: string;
}

// SyncResult represents the result of a sync
export interface SyncResult {
  message: string;
  count: number;
}

// RefreshResult represents the result of a refresh
export interface RefreshResult {
  message: string;
  new_characters: number;
  updated_characters: number;
}

// WoWProfile is the profile of the user
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
    id: number;
    characters: Array<{
      id: number;
      name: string;
      realm: {
        id: number;
        slug: string;
        name: string;
      };
      playable_class: {
        id: number;
        name: string;
      };
      playable_race: {
        id: number;
        name: string;
      };
      level: number;
      faction: {
        type: string;
        name: string;
      };
    }>;
  }>;
}

// CharacterDisplayResponse is the response for the character display endpoints
export interface CharacterDisplayResponse {
  message: string;
}

export interface FavoriteCharacterResponse {
  message: string;
}

// Enum centralised for WoW error codes
export enum WoWErrorCode {
  UNAUTHORIZED = "wow_unauthorized",
  TOKEN_EXPIRED = "wow_token_expired",
  NETWORK_ERROR = "wow_network_error",
  BATTLE_NET_NOT_LINKED = "battle_net_not_linked",
  NOT_FOUND = "wow_not_found",
  SERVER_ERROR = "wow_server_error",
}

// Centralised error class
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
