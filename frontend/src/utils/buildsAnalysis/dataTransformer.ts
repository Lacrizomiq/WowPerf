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
