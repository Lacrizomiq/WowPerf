export interface Season {
  slug: string;
  name: string;
  shortName: string;
  id: number;
  Dungeons: null | Dungeon[];
}

export interface Dungeon {
  ID: number;
  ChallengeModeID: number;
  EncounterID: number;
  Slug: string;
  Name: string;
  ShortName: string;
  MediaURL: string;
  Icon: string;
  KeyStoneUpgrades: KeyStoneUpgrade[];
  Seasons: Season[];
}

export interface Affix {
  ID: number;
  Name: string;
  Icon: string;
  WowheadURL: string;
}

export interface Member {
  CharacterID: number;
  CharacterName: string;
  RealmID: number;
  RealmSlug: string;
  EquippedItemLevel: number;
  RaceID: number;
  RaceName: string;
  SpecializationID: number;
  Specialization: string;
}

export interface MythicPlusSeasonInfo {
  CharacterName: string;
  RealmSlug: string;
  SeasonID: number;
  OverallMythicRating: number;
  OverallMythicRatingHex: string;
  BestRuns: MythicPlusRuns[];
}

export interface MythicPlusRuns {
  CompletedTimestamp: string;
  DungeonID: number;
  Dungeon: Dungeon;
  ShortName: string;
  Duration: number;
  IsCompletedWithinTime: boolean;
  KeyStoneUpgrades: number;
  KeystoneLevel: number;
  MythicRating: number;
  SeasonID: number;
  Season: Season;
  Affixes: Affix[];
  Members: Member[];
}

export interface KeyStoneUpgrade {
  ChallengeModeID: number;
  QualifyingDuration: number;
  UpgradeLevel: number;
}

export interface MythicDungeonProps {
  characterName: string;
  realmSlug: string;
  namespace: string;
  locale: string;
  region: string;
  seasonSlug: string;
}
