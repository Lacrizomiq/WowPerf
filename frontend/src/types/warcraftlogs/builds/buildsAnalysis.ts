// Types for builds analysis

// Popular items analysis
export interface PopularItem {
  encounter_id: number;
  item_slot: number;
  item_id: number;
  item_name: string;
  item_icon: string;
  item_quality: number; // 2: Common, 3: Rare, 4: Epic
  item_level: number;
  usage_count: number;
  usage_percentage: number;
  avg_keystone_level: number;
  rank: number;
}

// Enchant usage analysis
export interface EnchantUsage {
  item_slot: number;
  permanent_enchant_id: number;
  permanent_enchant_name: string;
  usage_count: number;
  avg_keystone_level: number;
  avg_item_level: number;
  max_keystone_level: number;
  rank: number;
}

// Gem usage analysis
export interface GemUsage {
  item_slot: number;
  gems_count: number;
  gem_ids_array: number[];
  gem_icons_array: string[];
  gem_levels_array: number[];
  usage_count: number;
  avg_keystone_level: number;
  avg_item_level: number;
  rank: number;
}

// Top talent builds analysis
export interface TopTalentBuild {
  talent_import: string;
  total_usage: number;
  avg_usage_percentage: number;
  avg_keystone_level: number;
}

// Talent builds by dungeon analysis
export interface TalentBuildByDungeon {
  class: string;
  spec: string;
  encounter_id: number;
  dungeon_name: string;
  talent_import: string;
  total_usage: number;
  avg_usage_percentage: number;
  avg_keystone_level: number;
}

// Stat priority analysis
export interface StatPriority {
  stat_name: string;
  stat_category: "minor" | "secondary";
  avg_value: number;
  min_value: number;
  max_value: number;
  total_samples: number;
  avg_keystone_level: number;
  priority_rank: number;
}

// Optimal item info
export interface OptimalItemInfo {
  name: string;
  icon: string;
  quality: number;
  usage_count: number;
}

// Optimal build
export interface OptimalBuild {
  top_talent_import: string;
  stat_priority: string;
  top_items: {
    [slotId: string]: OptimalItemInfo;
  };
}

// Class and spec summary
export interface ClassSpecSummary {
  avg_keystone_level: number;
  max_keystone_level: number;
  avg_item_level: number;
  top_talent_import: string;
  stat_priority: string;
  dungeons_count: number;
}

// Spec comparison for a given class
export interface SpecComparison {
  spec: string;
  avg_keystone_level: number;
  max_keystone_level: number;
  avg_item_level: number;
  dungeons_count: number;
  stat_priority: string;
}

// Response types
export type PopularItemsResponse = PopularItem[];
export type EnchantUsageResponse = EnchantUsage[];
export type GemUsageResponse = GemUsage[];
export type TopTalentBuildsResponse = TopTalentBuild[];
export type TalentBuildsByDungeonResponse = TalentBuildByDungeon[];
export type StatPrioritiesResponse = StatPriority[];
export type OptimalBuildResponse = OptimalBuild;
export type ClassSpecSummaryResponse = ClassSpecSummary;
export type SpecComparisonResponse = SpecComparison[];
