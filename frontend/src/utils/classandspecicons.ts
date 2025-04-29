// Types for class and spec icons
export interface SpecIcons {
  [key: string]: string;
}

export interface ClassData {
  classIcon: string;
  spec: SpecIcons;
}

export type ClassIconsMapping = {
  [key: string]: ClassData;
};

// Types for WarcraftLogs API data
export interface WowSpec {
  id: number;
  name: string;
  slug: string;
}

export interface WowClass {
  id: number;
  name: string;
  slug: string;
  specs: WowSpec[];
}

// Mapping between WarcraftLogs class IDs and class names
export const CLASS_ID_MAPPING: { [key: number]: string } = {
  1: "DeathKnight",
  2: "Druid",
  3: "Hunter",
  4: "Mage",
  5: "Monk",
  6: "Paladin",
  7: "Priest",
  8: "Rogue",
  9: "Shaman",
  10: "Warlock",
  11: "Warrior",
  12: "DemonHunter",
  13: "Evoker",
};

// constants
export const CLASS_ICONS_MAPPING: ClassIconsMapping = {
  DeathKnight: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/spell_deathknight_classicon.jpg",
    spec: {
      Blood:
        "https://render.worldofwarcraft.com/us/icons/56/spell_deathknight_bloodpresence.jpg",
      Frost:
        "https://render.worldofwarcraft.com/us/icons/56/spell_deathknight_frostpresence.jpg",
      Unholy:
        "https://render.worldofwarcraft.com/us/icons/56/spell_deathknight_unholypresence.jpg",
    },
  },
  DemonHunter: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_demonhunter.jpg",
    spec: {
      Havoc:
        "https://render.worldofwarcraft.com/us/icons/56/ability_demonhunter_specdps.jpg",
      Vengeance:
        "https://render.worldofwarcraft.com/us/icons/56/ability_demonhunter_spectank.jpg",
    },
  },
  Druid: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_druid.jpg",
    spec: {
      Balance:
        "https://render.worldofwarcraft.com/us/icons/56/spell_nature_starfall.jpg",
      Feral:
        "https://render.worldofwarcraft.com/us/icons/56/ability_druid_catform.jpg",
      Guardian:
        "https://render.worldofwarcraft.com/us/icons/56/ability_racial_bearform.jpg",
      Restoration:
        "https://render.worldofwarcraft.com/us/icons/56/spell_nature_healingtouch.jpg",
    },
  },
  Evoker: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_evoker.jpg",
    spec: {
      Devastation:
        "https://render.worldofwarcraft.com/us/icons/56/classicon_evoker_devastation.jpg",
      Preservation:
        "https://render.worldofwarcraft.com/us/icons/56/classicon_evoker_preservation.jpg",
      Augmentation:
        "https://render.worldofwarcraft.com/us/icons/56/classicon_evoker_augmentation.jpg",
    },
  },
  Hunter: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_hunter.jpg",
    spec: {
      BeastMastery:
        "https://render.worldofwarcraft.com/us/icons/56/ability_hunter_bestialdiscipline.jpg",
      Marksmanship:
        "https://render.worldofwarcraft.com/us/icons/56/ability_hunter_focusedaim.jpg",
      Survival:
        "https://render.worldofwarcraft.com/us/icons/56/ability_hunter_camouflage.jpg",
    },
  },
  Mage: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_mage.jpg",
    spec: {
      Arcane:
        "https://render.worldofwarcraft.com/us/icons/56/spell_holy_magicalsentry.jpg",
      Fire: "https://render.worldofwarcraft.com/us/icons/56/spell_fire_firebolt02.jpg",
      Frost:
        "https://render.worldofwarcraft.com/us/icons/56/spell_frost_frostbolt02.jpg",
    },
  },
  Monk: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_monk.jpg",
    spec: {
      Brewmaster:
        "https://render.worldofwarcraft.com/us/icons/56/spell_monk_brewmaster_spec.jpg",
      Windwalker:
        "https://render.worldofwarcraft.com/us/icons/56/spell_monk_windwalker_spec.jpg",
      Mistweaver:
        "https://render.worldofwarcraft.com/us/icons/56/spell_monk_mistweaver_spec.jpg",
    },
  },
  Paladin: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_paladin.jpg",
    spec: {
      Holy: "https://render.worldofwarcraft.com/us/icons/56/spell_holy_holybolt.jpg",
      Protection:
        "https://render.worldofwarcraft.com/us/icons/56/ability_paladin_shieldofthetemplar.jpg",
      Retribution:
        "https://render.worldofwarcraft.com/us/icons/56/spell_holy_auraoflight.jpg",
    },
  },
  Priest: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_priest.jpg",
    spec: {
      Discipline:
        "https://render.worldofwarcraft.com/us/icons/56/spell_holy_powerwordshield.jpg",
      Holy: "https://render.worldofwarcraft.com/us/icons/56/spell_holy_guardianspirit.jpg",
      Shadow:
        "https://render.worldofwarcraft.com/us/icons/56/spell_shadow_shadowwordpain.jpg",
    },
  },
  Rogue: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_rogue.jpg",
    spec: {
      Assassination:
        "https://render.worldofwarcraft.com/us/icons/56/ability_rogue_deadlybrew.jpg",
      Outlaw:
        "https://render.worldofwarcraft.com/us/icons/56/ability_rogue_waylay.jpg",
      Subtlety:
        "https://render.worldofwarcraft.com/us/icons/56/ability_stealth.jpg",
    },
  },
  Shaman: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_shaman.jpg",
    spec: {
      Elemental:
        "https://render.worldofwarcraft.com/us/icons/56/spell_nature_lightning.jpg",
      Enhancement:
        "https://render.worldofwarcraft.com/us/icons/56/spell_shaman_improvedstormstrike.jpg",
      Restoration:
        "https://render.worldofwarcraft.com/us/icons/56/spell_nature_magicimmunity.jpg",
    },
  },
  Warlock: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_warlock.jpg",
    spec: {
      Affliction:
        "https://render.worldofwarcraft.com/us/icons/56/spell_shadow_deathcoil.jpg",
      Demonology:
        "https://render.worldofwarcraft.com/us/icons/56/spell_shadow_metamorphosis.jpg",
      Destruction:
        "https://render.worldofwarcraft.com/us/icons/56/spell_shadow_rainoffire.jpg",
    },
  },
  Warrior: {
    classIcon:
      "https://render.worldofwarcraft.com/us/icons/56/classicon_warrior.jpg",
    spec: {
      Arms: "https://render.worldofwarcraft.com/us/icons/56/ability_warrior_savageblow.jpg",
      Fury: "https://render.worldofwarcraft.com/us/icons/56/ability_warrior_innerrage.jpg",
      Protection:
        "https://render.worldofwarcraft.com/us/icons/56/ability_warrior_defensivestance.jpg",
    },
  },
};

// Function to get class name from WarcraftLogs class ID
export const getClassNameById = (classId: number): string => {
  const className = CLASS_ID_MAPPING[classId];
  if (!className) {
    console.warn(`Class ID not found: ${classId}`);
    return "";
  }
  return className;
};

// Function to get class icon using WarcraftLogs class ID
export const getClassIconById = (classId: number): string => {
  const className = getClassNameById(classId);
  if (!className) return "";

  return getClassIcon(className);
};

// Function to get spec icon using WarcraftLogs class ID and spec name
export const getSpecIconById = (classId: number, specName: string): string => {
  const className = getClassNameById(classId);
  if (!className) return "";

  const normalizedSpecName = normalizeWowName(specName);
  return getSpecIcon(className, normalizedSpecName);
};

// utils functions to get spec icons
export const getSpecIcon = (className: string, specName: string): string => {
  const classData = CLASS_ICONS_MAPPING[className];
  if (!classData) {
    console.warn(`Class not found: ${className}`);
    return "";
  }

  const specIcon = classData.spec[specName];
  if (!specIcon) {
    console.warn(`Spec not found: ${specName} for class ${className}`);
    return "";
  }

  return specIcon;
};

// utils functions to get class icons
export const getClassIcon = (className: string): string => {
  const classData = CLASS_ICONS_MAPPING[className];
  if (!classData) {
    console.warn(`Class not found: ${className}`);
    return "";
  }

  return classData.classIcon;
};

// utils functions to normalize wow names
export const normalizeWowName = (name: string): string => {
  // Gestion des cas spéciaux
  if (name.toLowerCase() === "beastmastery") return "BeastMastery";

  // for other cases, capitalize the first letter of each word
  return name
    .split(/(?=[A-Z])/)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join("");
};

// Function to get spec data from WarcraftLogs spec name and class ID
export const getSpecData = (
  classId: number,
  specName: string
): { icon: string; name: string } => {
  const className = getClassNameById(classId);
  const normalizedSpecName = normalizeWowName(specName);

  return {
    icon: getSpecIcon(className, normalizedSpecName),
    name: specName,
  };
};

// Function to format the display name of a class
export const formatDisplayClassName = (className: string): string => {
  if (className === "deathknight") return "Death Knight";
  if (className === "demonhunter") return "Demon Hunter";
  return className.charAt(0).toUpperCase() + className.slice(1);
};

// Function to format the display name of a spec
export const formatDisplaySpecName = (spec: string): string => {
  if (spec === "beastmastery") return "Beast Mastery";
  return spec.charAt(0).toUpperCase() + spec.slice(1);
};

// Function to convert a class name to PascalCase for icons
export const classNameToPascalCase = (name: string): string => {
  if (name === "deathknight") return "DeathKnight";
  if (name === "demonhunter") return "DemonHunter";
  return name.charAt(0).toUpperCase() + name.slice(1);
};

// Function to convert a spec name to PascalCase for icons
export const specNameToPascalCase = (spec: string): string => {
  if (spec === "beastmastery") return "BeastMastery";
  return spec.charAt(0).toUpperCase() + spec.slice(1);
};
