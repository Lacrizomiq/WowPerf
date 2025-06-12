// utils/gemsData.ts

// Interface for gem data
export interface GemData {
  id: number;
  name: string;
  icon: string;
  quality: number; // 1: normal, 2: rare, 3: epic, etc.
  color?: string; // For colored gems like "red", "blue", etc.
}

// Correspondance between gem ID and their data
export const gemDataMap: Record<number, GemData> = {
  // Blasphemite Gems
  213740: {
    id: 213740,
    name: "Elusive Blasphemite",
    icon: "inv_misc_metagem_b.jpg",
    quality: 3,
  },
  213743: {
    id: 213743,
    name: "Deadly Blasphemite",
    icon: "item_cutmetagemb.jpg",
    quality: 3,
  },
  213746: {
    id: 213746,
    name: "Radiant Blasphemite",
    icon: "inv_misc_gem_x4_metagem_cut.jpg",
    quality: 3,
  },

  // Hybrid gems - Series 1
  213458: {
    id: 213458,
    name: "Inscribed Elixir Crystal",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color4_1.jpg",
    quality: 3,
  },
  213461: {
    id: 213461,
    name: "Embellished Elixir Crystal",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color4_2.jpg",
    quality: 3,
  },
  213455: {
    id: 213455,
    name: "Jagged Elixir Crystal",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color4_3.jpg",
    quality: 3,
  },

  // Hybrid gems - Series 2
  213467: {
    id: 213467,
    name: "Inscribed Dreadstone",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color5_2.jpg",
    quality: 3,
  },
  213470: {
    id: 213470,
    name: "Embellished Dreadstone",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color5_3.jpg",
    quality: 3,
  },
  213473: {
    id: 213473,
    name: "Jagged Dreadstone",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color5_1.jpg",
    quality: 3,
  },

  // Hybrid gems - Series 3
  213479: {
    id: 213479,
    name: "Inscribed Amethyst",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color1_2.jpg",
    quality: 3,
  },
  213482: {
    id: 213482,
    name: "Embellished Amethyst",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color1_3.jpg",
    quality: 3,
  },
  213485: {
    id: 213485,
    name: "Jagged Amethyst",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color1_1.jpg",
    quality: 3,
  },

  // Hybrid gems - Series 4
  213491: {
    id: 213491,
    name: "Inscribed Emerald",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color2_2.jpg",
    quality: 3,
  },
  213494: {
    id: 213494,
    name: "Embellished Emerald",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color2_3.jpg",
    quality: 3,
  },
  213497: {
    id: 213497,
    name: "Jagged Emerald",
    icon: "inv_jewelcrafting_cut-standart-gem-hybrid_color2_1.jpg",
    quality: 3,
  },

  // Citrine Gems - Special Series
  228638: {
    id: 228638,
    name: "Sea-Runed Citrine (Red)",
    icon: "inv_siren_isle_searuned_citrine_red.jpg",
    quality: 3,
  },
  228639: {
    id: 228639,
    name: "Sea-Runed Citrine (Blue)",
    icon: "inv_siren_isle_searuned_citrine_blue.jpg",
    quality: 3,
  },
  228640: {
    id: 228640,
    name: "Sea-Runed Citrine (Pink)",
    icon: "inv_siren_isle_searuned_citrine_pink.jpg",
    quality: 3,
  },
  228634: {
    id: 228634,
    name: "Stormcharged Citrine (Blue)",
    icon: "inv_siren_isle_stormcharged_citrine_blue.jpg",
    quality: 3,
  },
  228636: {
    id: 228636,
    name: "Stormcharged Citrine (Green)",
    icon: "inv_siren_isle_stormcharged_citrine_green.jpg",
    quality: 3,
  },
};

// Function to get gem data from ID
export function getGemData(gemId: number): GemData | undefined {
  return gemDataMap[gemId];
}

// Function to get the gem icon URL
export function getGemIconUrl(gemId: number): string {
  const gemData = getGemData(gemId);
  if (!gemData) {
    // Fallback icon if gem ID is not found
    return "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg";
  }
  return `https://wow.zamimg.com/images/wow/icons/large/${gemData.icon}`;
}

// Function to get color class based on gem quality
export function getGemQualityClass(quality: number): string {
  switch (quality) {
    case 4:
      return "text-purple-400"; // Epic
    case 3:
      return "text-blue-400"; // Rare
    case 2:
      return "text-green-400"; // Uncommon
    default:
      return "text-white"; // Common
  }
}
