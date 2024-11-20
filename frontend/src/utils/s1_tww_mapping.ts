// Mapping of dungeons to their encounterID
export const DUNGEON_ENCOUNTER_MAPPING: Record<
  string,
  { id: number; name: string }
> = {
  "arakara-city-of-echoes": {
    id: 12660,
    name: "Ara-Kara, City of Echoes",
  },
  "city-of-threads": {
    id: 12669,
    name: "City of Threads",
  },
  "grim-batol": {
    id: 60670,
    name: "Grim Batol",
  },
  "mists-of-tirna-scithe": {
    id: 62290,
    name: "Mists of Tirna Scithe",
  },
  "siege-of-boralus": {
    id: 61822,
    name: "Siege of Boralus",
  },
  "the-dawnbreaker": {
    id: 12662,
    name: "The Dawnbreaker",
  },
  "the-necrotic-wake": {
    id: 62286,
    name: "The Necrotic Wake",
  },
  "the-stonevault": {
    id: 12652,
    name: "The Stonevault",
  },
};

// Mapping between WarcraftLogs encounter IDs and Blizzard encounter IDs
export interface RaidEncounterMapping {
  warcraftLogsId: number; // WarcraftLogs encounter ID
  blizzardId: number; // Blizzard encounter ID
  name: string; // Official encounter name
  slug: string; // Encounter slug for URLs
  iconUrl?: string; // Optional icon override
  icon?: string; // Optional icon override
}

export interface RaidZoneMapping {
  warcraftLogsId: number; // WarcraftLogs zone ID
  blizzardId: number; // Blizzard raid ID
  name: string; // Raid name
  shortName: string; // Short version of raid name
  slug: string; // Raid slug for URLs
}

// Mapping for specific raids and their encounters
export const RAID_ZONE_MAPPING: Record<string, RaidZoneMapping> = {
  "nerubar-palace": {
    warcraftLogsId: 38,
    blizzardId: 1273,
    name: "Nerub'ar Palace",
    shortName: "NP",
    slug: "nerubar-palace",
  },
  // Add other raids as they become available
};

// Mapping for raid encounters
export const RAID_ENCOUNTER_MAPPING: Record<string, RaidEncounterMapping> = {
  "ulgrax-the-devourer": {
    warcraftLogsId: 2902,
    blizzardId: 2607,
    name: "Ulgrax the Devourer",
    slug: "ulgrax-the-devourer",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_raidnerubian_nerubianhulk.jpg`,
    icon: "inv_achievement_raidnerubian_nerubianhulk",
  },
  "the-bloodbound-horror": {
    warcraftLogsId: 2917,
    blizzardId: 2611,
    name: "The Bloodbound Horror",
    slug: "the-bloodbound-horror",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_raidnerubian_blackblood.jpg`,
    icon: "inv_achievement_raidnerubian_blackblood",
  },
  sikran: {
    warcraftLogsId: 2898,
    blizzardId: 2599,
    name: "Sikran, Captain of the Sureki",
    slug: "sikran",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_raidnerubian_nerubianevolved.jpg`,
    icon: "inv_achievement_raidnerubian_nerubianevolved",
  },
  rashanan: {
    warcraftLogsId: 2918,
    blizzardId: 2609,
    name: "Rasha'nan",
    slug: "rashanan",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_raidnerubian_flyingnerubianevolved.jpg`,
    icon: "inv_achievement_raidnerubian_flyingnerubianevolved",
  },
  "broodtwister-ovinax": {
    warcraftLogsId: 2919,
    blizzardId: 2612,
    name: "Broodtwister Ovi'nax",
    slug: "broodtwister-ovinax",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_raidnerubian_swarmmother.jpg`,
    icon: "inv_achievement_raidnerubian_swarmmother",
  },
  "nexus-princess-kyveza": {
    warcraftLogsId: 2920,
    blizzardId: 2601,
    name: "Nexus-Princess Ky'veza",
    slug: "nexus-princess-kyveza",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_raidnerubian_etherealassasin.jpg`,
    icon: "inv_achievement_raidnerubian_etherealassasin",
  },
  "the-silken-court": {
    warcraftLogsId: 2921,
    blizzardId: 2608,
    name: "The Silken Court",
    slug: "the-silken-court",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_raidnerubian_council.jpg`,
    icon: "inv_achievement_raidnerubian_council",
  },
  "queen-ansurek": {
    warcraftLogsId: 2922,
    blizzardId: 2602,
    name: "Queen Ansurek",
    slug: "queen-ansurek",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_raidnerubian_queenansurek.jpg`,
    icon: "inv_achievement_raidnerubian_queenansurek",
  },
};

// Helper functions to get encounter information
export const getEncounterByWarcraftLogsId = (
  warcraftLogsId: number
): RaidEncounterMapping | undefined => {
  return Object.values(RAID_ENCOUNTER_MAPPING).find(
    (encounter) => encounter.warcraftLogsId === warcraftLogsId
  );
};

export const getEncounterByBlizzardId = (
  blizzardId: number
): RaidEncounterMapping | undefined => {
  return Object.values(RAID_ENCOUNTER_MAPPING).find(
    (encounter) => encounter.blizzardId === blizzardId
  );
};

export const getZoneByWarcraftLogsId = (
  warcraftLogsId: number
): RaidZoneMapping | undefined => {
  return Object.values(RAID_ZONE_MAPPING).find(
    (zone) => zone.warcraftLogsId === warcraftLogsId
  );
};
