export enum AttackType {
  MELEE = "MELEE",
  RANGED = "RANGED",
  SPELL = "SPELL",
}

interface SpecializationInfo {
  id: number;
  attackType: AttackType;
  role: "TANK" | "HEALER" | "DPS";
}

interface ClassSpecializations {
  [spec: string]: SpecializationInfo;
}

interface WoWClassSpecs {
  [className: string]: ClassSpecializations;
}

export const specMapping: WoWClassSpecs = {
  Hunter: {
    BeastMastery: { id: 253, attackType: AttackType.RANGED, role: "DPS" },
    Marksmanship: { id: 254, attackType: AttackType.RANGED, role: "DPS" },
    Survival: { id: 255, attackType: AttackType.MELEE, role: "DPS" },
  },
  Shaman: {
    Elemental: { id: 262, attackType: AttackType.SPELL, role: "DPS" },
    Enhancement: { id: 263, attackType: AttackType.MELEE, role: "DPS" },
    Restoration: { id: 264, attackType: AttackType.SPELL, role: "HEALER" },
  },
  Druid: {
    Balance: { id: 102, attackType: AttackType.SPELL, role: "DPS" },
    Feral: { id: 103, attackType: AttackType.MELEE, role: "DPS" },
    Guardian: { id: 104, attackType: AttackType.MELEE, role: "TANK" },
    Restoration: { id: 105, attackType: AttackType.SPELL, role: "HEALER" },
  },
  Evoker: {
    Devastation: { id: 1467, attackType: AttackType.SPELL, role: "DPS" },
    Preservation: { id: 1468, attackType: AttackType.SPELL, role: "HEALER" },
    Augmentation: { id: 1473, attackType: AttackType.SPELL, role: "DPS" },
  },
  Warrior: {
    Arms: { id: 71, attackType: AttackType.MELEE, role: "DPS" },
    Fury: { id: 72, attackType: AttackType.MELEE, role: "DPS" },
    Protection: { id: 73, attackType: AttackType.MELEE, role: "TANK" },
  },
  DeathKnight: {
    Blood: { id: 250, attackType: AttackType.MELEE, role: "TANK" },
    Frost: { id: 251, attackType: AttackType.MELEE, role: "DPS" },
    Unholy: { id: 252, attackType: AttackType.MELEE, role: "DPS" },
  },
  Paladin: {
    Holy: { id: 65, attackType: AttackType.SPELL, role: "HEALER" },
    Protection: { id: 66, attackType: AttackType.MELEE, role: "TANK" },
    Retribution: { id: 70, attackType: AttackType.MELEE, role: "DPS" },
  },
  Priest: {
    Discipline: { id: 256, attackType: AttackType.SPELL, role: "HEALER" },
    Holy: { id: 257, attackType: AttackType.SPELL, role: "HEALER" },
    Shadow: { id: 258, attackType: AttackType.SPELL, role: "DPS" },
  },
  Monk: {
    Brewmaster: { id: 268, attackType: AttackType.MELEE, role: "TANK" },
    Mistweaver: { id: 270, attackType: AttackType.SPELL, role: "HEALER" },
    Windwalker: { id: 269, attackType: AttackType.MELEE, role: "DPS" },
  },
  Mage: {
    Arcane: { id: 62, attackType: AttackType.SPELL, role: "DPS" },
    Fire: { id: 63, attackType: AttackType.SPELL, role: "DPS" },
    Frost: { id: 64, attackType: AttackType.SPELL, role: "DPS" },
  },
  Rogue: {
    Assassination: { id: 259, attackType: AttackType.MELEE, role: "DPS" },
    Outlaw: { id: 260, attackType: AttackType.MELEE, role: "DPS" },
    Subtlety: { id: 261, attackType: AttackType.MELEE, role: "DPS" },
  },
  DemonHunter: {
    Havoc: { id: 577, attackType: AttackType.MELEE, role: "DPS" },
    Vengeance: { id: 581, attackType: AttackType.MELEE, role: "TANK" },
  },
  Warlock: {
    Affliction: { id: 265, attackType: AttackType.SPELL, role: "DPS" },
    Demonology: { id: 266, attackType: AttackType.SPELL, role: "DPS" },
    Destruction: { id: 267, attackType: AttackType.SPELL, role: "DPS" },
  },
};

// Helper functions
export const getSpecInfoById = (
  specId: number
):
  | { className: string; specName: string; info: SpecializationInfo }
  | undefined => {
  for (const [className, specs] of Object.entries(specMapping)) {
    for (const [specName, info] of Object.entries(specs)) {
      if (info.id === specId) {
        return { className, specName, info };
      }
    }
  }
  return undefined;
};

export const getAttackTypeForSpec = (
  specId: number
): AttackType | undefined => {
  const specInfo = getSpecInfoById(specId);
  return specInfo?.info.attackType;
};
