import {
  TalentBuildByDungeon,
  StatPriority,
  OptimalBuild,
  OptimalItemInfo,
} from "@/types/warcraftlogs/builds/buildsAnalysis";

// Group items by item slot (items, enchantments, gems)
interface ItemWithSlot {
  item_slot: number;
  [key: string]: any;
}

// Group items by item slot (items, enchantments, gems)
export function groupBySlot<T extends ItemWithSlot>(
  items: T[]
): Record<number, T[]> {
  return items.reduce((acc, item) => {
    if (!acc[item.item_slot]) {
      acc[item.item_slot] = [];
    }
    acc[item.item_slot].push(item);
    return acc;
  }, {} as Record<number, T[]>);
}

// Mapping of item slot IDs to item slot names
export const ITEM_SLOT_NAMES: Record<number, string> = {
  0: "Head",
  1: "Neck",
  2: "Shoulders",
  3: "Shirt",
  4: "Chest",
  5: "Waist",
  6: "Legs",
  7: "Feet",
  8: "Wrists",
  9: "Hands",
  10: "Ring 1",
  11: "Ring 2",
  12: "Trinket 1",
  13: "Trinket 2",
  14: "Back",
  15: "Main Hand",
  16: "Off Hand",
  17: "Tabard",
};

// Order of display for item slots
export const ITEM_SLOT_DISPLAY_ORDER: number[] = [
  0, // Head
  1, // Neck
  2, // Shoulders
  14, // Back
  4, // Chest
  8, // Wrists
  9, // Hands
  5, // Waist
  6, // Legs
  7, // Feet
  10, // Ring 1
  11, // Ring 2
  12, // Trinket 1
  13, // Trinket 2
  15, // Main Hand
  16, // Off Hand
  // 3, 17 - Shirt and Tabard excluded because less important for performance
];

// Get the URL for an item icon
export function getItemIconUrl(iconName: string): string {
  return `https://wow.zamimg.com/images/wow/icons/large/${iconName}`;
}

// Get the quality class for an item
export function getItemQualityClass(quality: number): string {
  switch (quality) {
    case 4:
      return "epic"; // Purple
    case 3:
      return "rare"; // Blue
    case 2:
      return "uncommon"; // Green
    default:
      return "common"; // White
  }
}

// Group talents by dungeon ID and dungeon name
interface TalentsByDungeon {
  [dungeonId: string]: {
    dungeonName: string;
    talents: TalentBuildByDungeon[];
  };
}

// Group talents by dungeon ID and dungeon name
export function groupTalentsByDungeon(
  talents: TalentBuildByDungeon[]
): TalentsByDungeon {
  return talents.reduce((acc, talent) => {
    const id = talent.encounter_id.toString();
    if (!acc[id]) {
      acc[id] = {
        dungeonName: talent.dungeon_name,
        talents: [],
      };
    }
    acc[id].talents.push(talent);
    return acc;
  }, {} as TalentsByDungeon);
}

// Group stats by category (minor and secondary)
interface StatsByCategory {
  minor: StatPriority[];
  secondary: StatPriority[];
}

export function groupStatsByCategory(stats: StatPriority[]): StatsByCategory {
  return stats.reduce(
    (acc, stat) => {
      if (!acc[stat.stat_category]) {
        acc[stat.stat_category] = [];
      }
      acc[stat.stat_category].push(stat);
      return acc;
    },
    { minor: [], secondary: [] } as StatsByCategory
  );
}

// Convert optimal items to array of slot IDs and items
export interface SlotWithItem {
  slotId: number;
  item: OptimalItemInfo;
}

export function convertOptimalItemsToArray(
  build: OptimalBuild
): SlotWithItem[] {
  return Object.entries(build.top_items).map(([slotId, item]) => ({
    slotId: parseInt(slotId, 10),
    item,
  }));
}
