// Mapping of dungeons to their encounterID
export const DUNGEON_ENCOUNTER_MAPPING: Record<
  string,
  { id: number; name: string }
> = {
  // Season 2
  "cinderbrew-meadery": {
    id: 12661,
    name: "Cinderbrew Meadery",
  },
  "darkflame-cleft": {
    id: 12651,
    name: "Darkflame Cleft",
  },
  "operation-mechagon-workshop": {
    id: 112098,
    name: "Mechagon Workshop",
  },
  "operation-floodgate": {
    id: 12773,
    name: "Operation: Floodgate",
  },
  "priory-of-the-sacred-flame": {
    id: 12649,
    name: "Priory of the Sacred Flame",
  },
  "the-motherlode": {
    id: 61594,
    name: "The MOTHERLODE!!",
  },
  "the-rookery": {
    id: 12648,
    name: "The Rookery",
  },
  "theater-of-pain": {
    id: 62293,
    name: "Theater of Pain",
  },
  // Season 1
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
  "liberation-of-undermine": {
    warcraftLogsId: 42,
    blizzardId: 1296,
    name: "Liberation of Undermine",
    shortName: "LOU",
    slug: "liberation-of-undermine",
  },
};

// Mapping for raid encounters
export const RAID_ENCOUNTER_MAPPING: Record<string, RaidEncounterMapping> = {
  // Liberation of Undermine
  "vexie-and-the-geargrinders": {
    warcraftLogsId: 3009,
    blizzardId: 2639,
    name: "Vexie and the Geargrinders",
    slug: "vexie-and-the-geargrinders",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_zone_undermine.jpg`,
    icon: "inv_achievement_zone_undermine",
  },
  "cauldron-of-carnage": {
    warcraftLogsId: 3010,
    blizzardId: 2640,
    name: "Cauldron of Carnage",
    slug: "cauldron-of-carnage",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_zone_undermine.jpg`,
    icon: "inv_achievement_zone_undermine",
  },
  "rik-reverb": {
    warcraftLogsId: 3011,
    blizzardId: 2641,
    name: "Rik Reverb",
    slug: "rik-reverb",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_zone_undermine.jpg`,
    icon: "inv_achievement_zone_undermine",
  },
  "stix-bunkjunker": {
    warcraftLogsId: 3012,
    blizzardId: 2642,
    name: "Stix Bunkjunker",
    slug: "stix-bunkjunker",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_zone_undermine.jpg`,
    icon: "inv_achievement_zone_undermine",
  },
  "sprocketmonger-lockenstock": {
    warcraftLogsId: 3013,
    blizzardId: 2653,
    name: "Sprocketmonger Lockenstock",
    slug: "sprocketmonger-lockenstock",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_zone_undermine.jpg`,
    icon: "inv_achievement_zone_undermine",
  },
  "onearmed-bandit": {
    warcraftLogsId: 3014,
    blizzardId: 2644,
    name: "One-Armed Bandit",
    slug: "onearmed-bandit",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_zone_undermine.jpg`,
    icon: "inv_achievement_zone_undermine",
  },
  "mugzee-heads-of-security": {
    warcraftLogsId: 3015,
    blizzardId: 2645,
    name: "Mug'Zee, Heads of Security",
    slug: "mugzee-heads-of-security",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_zone_undermine.jpg`,
    icon: "inv_achievement_zone_undermine",
  },
  "chrome-king-gallywix": {
    warcraftLogsId: 3016,
    blizzardId: 2646,
    name: "Chrome King Gallywix",
    slug: "chrome-king-gallywix",
    iconUrl: `https://wow.zamimg.com/images/wow/icons/large/inv_achievement_zone_undermine.jpg`,
    icon: "inv_achievement_zone_undermine",
  },
  // Nerub'ar Palace
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
