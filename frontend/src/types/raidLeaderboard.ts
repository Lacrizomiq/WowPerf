// Raid Rankings Interfaces

export interface RaidRankings {
  raidRankings: RaidRanking[];
}

export interface RaidRanking {
  rank: number;
  regionRank: number;
  guild: Guild;
  encountersDefeated: EncounterDefeated[];
  encountersPulled: EncounterPulled[];
}

export interface Guild {
  id: number;
  name: string;
  faction: string;
  realm: Realm;
  region: Region;
  path: string;
  logo: string;
  color: string;
}

export interface Realm {
  id: number;
  connectedRealmId: number;
  wowRealmId: number;
  wowConnectedRealmId: number;
  name: string;
  altName: string | null;
  slug: string;
  altSlug: string;
  locale: string;
  isConnected: boolean;
  realmType: string;
}

export interface Region {
  name: string;
  slug: string;
  short_name: string;
}

export interface EncounterDefeated {
  slug: string;
  lastDefeated: string;
  firstDefeated: string;
}

export interface EncounterPulled {
  id: number;
  slug: string;
  numPulls: number;
  pullStartedAt: string;
  bestPercent: number;
  isDefeated: boolean;
}
