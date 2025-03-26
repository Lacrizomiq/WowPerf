// utils/specSlugs.ts
import { SpecAverageGlobalScore } from "@/types/warcraftlogs/globalLeaderboardAnalysis";

export const generateSpecSlugs = (
  specData: SpecAverageGlobalScore[]
): string[] => {
  return specData.map((item) => {
    // Convert class and spec to lowercase and replace spaces with hyphens
    const classSlug = item.class.toLowerCase().replace(/ /g, "-");
    const specSlug = item.spec.toLowerCase();
    return `${classSlug}-${specSlug}`;
  });
};

// Example usage to get a specific slug
export const getSpecSlug = (className: string, specName: string): string => {
  return `${className
    .toLowerCase()
    .replace(/ /g, "-")}-${specName.toLowerCase()}`;
};
