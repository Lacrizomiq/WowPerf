export interface MythicPlusRun {
  clear_time_ms: number;
  completed_at: string;
  deleted_at: string | null;
  dungeon: Dungeon;
  faction: string;
  keystone_run_id: number;
  keystone_team_id: number;
  keystone_time_ms: number;
  loggedSources: string[];
  logged_details: string | null;
  logged_run_id: string | null;
  mythic_level: number;
  num_chests: number;
  num_modifiers_active: number;
  roster: Roster[];
  score: number;
  season: string;
  status: string;
  time_remaining_ms: number;
  weekly_modifiers: Affix[];
}

export interface Dungeon {
  expansion_id: number;
  group_finder_activity_ids: number[];
  icon_url: string;
  id: number;
  keystone_timer_ms: number;
  map_challenge_mode_id: number;
  name: string;
  num_bosses: number;
  patch: string;
  short_name: string;
  slug: string;
  type: string;
  wowInstanceId: number;
}

export interface Roster {
  character: Character;
  guild: Guild | null;
  isTransfer: boolean;
  items: Items;
  oldCharacter: string | null;
  ranks: Ranks;
  role: string;
}

export interface Character {
  class: Class;
  faction: string;
  id: number;
  level: number;
  name: string;
  path: string;
  persona_id: number;
  race: Race;
  realm: Realm;
  recruitmentProfiles: string[];
  region: Region;
  spec: Spec;
  talentLoadout: TalentLoadout;
}

export interface Class {
  id: number;
  name: string;
  slug: string;
}

export interface Race {
  faction: string;
  id: number;
  name: string;
  slug: string;
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
  wowRealmId: number;
}

export interface Region {
  name: string;
  short_name: string;
  slug: string;
}

export interface Spec {
  class_id: number;
  id: number;
  is_melee: boolean;
  name: string;
  patch: string;
  role: string;
  slug: string;
}

export interface TalentLoadout {
  heroSubTreeId: number;
  loadoutText: string;
  specId: number;
}

export interface Guild {
  faction: string;
  id: number;
  name: string;
  path: string;
  realm: Realm;
  region: Region;
}

export interface Items {
  item_level_equipped: number;
  item_level_total: number;
  items: { [key: string]: EquippedItem };
  updated_at: string;
}

export interface EquippedItem {
  bonuses: number[];
  enchant?: number;
  gems?: number[];
  icon: string;
  is_legendary: boolean;
  item_id: number;
  item_level: number;
  item_quality: number;
  name: string;
  tier?: string;
}

export interface Ranks {
  realm: number;
  region: number;
  score: number;
  world: number;
}

export interface Affix {
  description: string;
  icon: string;
  id: number;
  name: string;
  slug: string;
}
