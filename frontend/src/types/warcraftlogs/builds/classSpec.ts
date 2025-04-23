// types/warcraftlogs/builds/classSpec.ts

// Valid WoW classes (lowercase as expected by the API)
export type WowClassParam =
  | "warrior"
  | "paladin"
  | "hunter"
  | "rogue"
  | "priest"
  | "deathknight"
  | "shaman"
  | "mage"
  | "warlock"
  | "monk"
  | "druid"
  | "demonhunter"
  | "evoker";

// All possible specs (lowercase as expected by the API)
export type WowSpecParam =
  // Warrior
  | "arms"
  | "fury"
  | "protection"
  // Paladin
  | "holy"
  | "protection"
  | "retribution"
  // Hunter
  | "beastmastery"
  | "marksmanship"
  | "survival"
  // Rogue
  | "assassination"
  | "outlaw"
  | "subtlety"
  // Priest
  | "discipline"
  | "holy"
  | "shadow"
  // Death Knight
  | "blood"
  | "frost"
  | "unholy"
  // Shaman
  | "elemental"
  | "enhancement"
  | "restoration"
  // Mage
  | "arcane"
  | "fire"
  | "frost"
  // Warlock
  | "affliction"
  | "demonology"
  | "destruction"
  // Monk
  | "brewmaster"
  | "mistweaver"
  | "windwalker"
  // Druid
  | "balance"
  | "feral"
  | "guardian"
  | "restoration"
  // Demon Hunter
  | "havoc"
  | "vengeance"
  // Evoker
  | "devastation"
  | "preservation"
  | "augmentation";
