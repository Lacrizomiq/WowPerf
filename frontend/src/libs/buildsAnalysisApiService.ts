import api from "./api";
import {
  WowClassParam,
  WowSpecParam,
} from "../types/warcraftlogs/builds/classSpec";
import {
  PopularItemsResponse,
  EnchantUsageResponse,
  GemUsageResponse,
  TopTalentBuildsResponse,
  TalentBuildsByDungeonResponse,
  StatPrioritiesResponse,
  OptimalBuildResponse,
  ClassSpecSummaryResponse,
  SpecComparisonResponse,
} from "../types/warcraftlogs/builds/buildsAnalysis";

// Get popular items by class and spec
export const getPopularItems = async (
  className: WowClassParam,
  spec: WowSpecParam
) => {
  try {
    const { data } = await api.get<PopularItemsResponse>(
      "warcraftlogs/mythicplus/builds/analysis/items",
      { params: { class: className, spec } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching popular items:", error);
    throw error;
  }
};

// Get enchant usage by class and spec
export const getEnchantUsage = async (
  className: WowClassParam,
  spec: WowSpecParam
) => {
  try {
    const { data } = await api.get<EnchantUsageResponse>(
      "/warcraftlogs/mythicplus/builds/analysis/enchants",
      { params: { class: className, spec } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching enchant usage:", error);
    throw error;
  }
};

// Get gem usage by class and spec
export const getGemUsage = async (
  className: WowClassParam,
  spec: WowSpecParam
) => {
  try {
    const { data } = await api.get<GemUsageResponse>(
      "/warcraftlogs/mythicplus/builds/analysis/gems",
      { params: { class: className, spec } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching gem usage:", error);
    throw error;
  }
};

// Get top talent builds by class and spec
export const getTopTalentBuilds = async (
  className: WowClassParam,
  spec: WowSpecParam
) => {
  try {
    const { data } = await api.get<TopTalentBuildsResponse>(
      "/warcraftlogs/mythicplus/builds/analysis/talents/top",
      { params: { class: className, spec } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching top talent builds:", error);
    throw error;
  }
};

// Get talent builds by dungeon for a specific class and spec
export const getTalentBuildsByDungeon = async (
  className: WowClassParam,
  spec: WowSpecParam
) => {
  try {
    const { data } = await api.get<TalentBuildsByDungeonResponse>(
      "/warcraftlogs/mythicplus/builds/analysis/talents/dungeons",
      { params: { class: className, spec } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching talent builds by dungeon:", error);
    throw error;
  }
};

// Get stat priorities by class and spec
export const getStatPriorities = async (
  className: WowClassParam,
  spec: WowSpecParam
) => {
  try {
    const { data } = await api.get<StatPrioritiesResponse>(
      "/warcraftlogs/mythicplus/builds/analysis/stats",
      { params: { class: className, spec } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching stat priorities:", error);
    throw error;
  }
};

// Get optimal build by class and spec
export const getOptimalBuild = async (
  className: WowClassParam,
  spec: WowSpecParam
) => {
  try {
    const { data } = await api.get<OptimalBuildResponse>(
      "/warcraftlogs/mythicplus/builds/analysis/optimal",
      { params: { class: className, spec } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching optimal build:", error);
    throw error;
  }
};

// Get class and spec summary
export const getClassSpecSummary = async (
  className: WowClassParam,
  spec: WowSpecParam
) => {
  try {
    const { data } = await api.get<ClassSpecSummaryResponse>(
      "/warcraftlogs/mythicplus/builds/analysis/summary",
      { params: { class: className, spec } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching class/spec summary:", error);
    throw error;
  }
};

// Get spec comparison for a given class
export const getSpecComparison = async (className: WowClassParam) => {
  try {
    const { data } = await api.get<SpecComparisonResponse>(
      "/warcraftlogs/mythicplus/builds/analysis/specs/comparison",
      { params: { class: className } }
    );
    return data;
  } catch (error) {
    console.error("Error fetching spec comparison:", error);
    throw error;
  }
};
