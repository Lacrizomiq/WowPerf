// utils/classandspecbackgrounds.ts
export const getSpecBackground = (
  className: string,
  specName: string
): string => {
  // Normalize the class and spec names to match the API data format
  const normalizedClass = className
    .toLowerCase()
    .replace(/knight|hunter/gi, (match) => match.toLowerCase()); // Handle "DeathKnight" -> "deathknight", "DemonHunter" -> "demonhunter"
  const normalizedSpec = specName.toLowerCase().replace(/mastery/gi, "mastery"); // Handle "BeastMastery" -> "beastmastery"

  // Map class and spec to background class, using the normalized API data format
  const backgroundMap: Record<string, Record<string, string>> = {
    warrior: {
      arms: "bg-spec-71",
      fury: "bg-spec-72",
      protection: "bg-spec-73",
    },
    paladin: {
      holy: "bg-spec-65",
      protection: "bg-spec-66",
      retribution: "bg-spec-70",
    },
    hunter: {
      beastmastery: "bg-spec-253",
      marksmanship: "bg-spec-254",
      survival: "bg-spec-255",
    },
    rogue: {
      assassination: "bg-spec-259",
      outlaw: "bg-spec-260",
      subtlety: "bg-spec-261",
    },
    priest: {
      discipline: "bg-spec-256",
      holy: "bg-spec-257",
      shadow: "bg-spec-258",
    },
    deathknight: {
      blood: "bg-spec-250",
      frost: "bg-spec-251",
      unholy: "bg-spec-252",
    },
    shaman: {
      elemental: "bg-spec-262",
      enhancement: "bg-spec-263",
      restoration: "bg-spec-264",
    },
    mage: {
      arcane: "bg-spec-62",
      fire: "bg-spec-63",
      frost: "bg-spec-64",
    },
    warlock: {
      affliction: "bg-spec-265",
      demonology: "bg-spec-266",
      destruction: "bg-spec-267",
    },
    monk: {
      brewmaster: "bg-spec-268",
      mistweaver: "bg-spec-270",
      windwalker: "bg-spec-269",
    },
    druid: {
      balance: "bg-spec-102",
      feral: "bg-spec-103",
      guardian: "bg-spec-104",
      restoration: "bg-spec-105",
    },
    demonhunter: {
      havoc: "bg-spec-577",
      vengeance: "bg-spec-581",
    },
    evoker: {
      devastation: "bg-spec-1467",
      preservation: "bg-spec-1468",
      augmentation: "bg-spec-1473",
    },
  };

  const classKey = normalizedClass;
  const specKey = normalizedSpec;
  return backgroundMap[classKey]?.[specKey] || "bg-deep-blue";
};
