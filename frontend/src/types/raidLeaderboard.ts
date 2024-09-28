export interface RaidLeaderboard {
  progression: ProgressionItem[];
}

export interface ProgressionItem {
  guilds: GuildProgression[];
  progress: number;
  totalGuilds: number;
}

export interface GuildProgression {
  defeatedAt: string;
  guild: Guild;
  recruitmentProfiles: RecruitmentProfile[];
  streamers: Streamers;
}

export interface Guild {
  faction: string;
  id: number;
  name: string;
  path: string;
  realm: Realm;
  region: Region;
  color?: string;
  logo?: string;
  alt_name?: string;
}

export interface Realm {
  altName: string | null;
  altSlug: string;
  connectedRealmId: number;
  id: number;
  isConnected: boolean;
  locale: string;
  name: string;
  realmType: string;
  slug: string;
  wowConnectedRealmId: number;
  wowRealmId: number | null;
}

export interface Region {
  name: string;
  short_name: string;
  slug: string;
}

export interface RecruitmentProfile {
  activity_type: string;
  entity_type: string;
  recruitment_profile_id: number;
}

export interface Streamers {
  count: number;
  stream: Stream | null;
}

export interface Stream {
  community_ids: any[];
  game_id: string;
  id: string;
  language: string;
  name: string;
  started_at: string;
  thumbnail_url: string;
  title: string;
  type: string;
  user_id: string;
  viewer_count: number;
}
