// components/MythicPlus/BuildsAnalysis/enchants-gems/BestGemsOverview.tsx
"use client";

import Image from "next/image";
import { GemUsage } from "@/types/warcraftlogs/builds/buildsAnalysis";
import { getGemIconUrl } from "@/utils/buildsAnalysis/gemsData";

interface BestGemsOverviewProps {
  gemsData: GemUsage[];
}

// Function to generate a unique key for a gem combination
const getGemCombinationKey = (gem: GemUsage): string => {
  return `${gem.gem_ids_array.sort().join("-")}-${gem.gems_count}`;
};

export default function BestGemsOverview({ gemsData }: BestGemsOverviewProps) {
  if (!gemsData || gemsData.length === 0) {
    return (
      <div className="bg-slate-800 rounded-lg p-5 text-center">
        <p className="text-slate-400">No gem data available.</p>
      </div>
    );
  }

  // Create a map to track combined usage of gem combinations across all slots
  const gemCombinationsMap = new Map<
    string,
    {
      key: string;
      gemIds: number[];
      gemsCount: number;
      totalUsage: number;
      avgKeyLevel: number;
      slots: Set<number>;
    }
  >();

  // Aggregate gem usage across all slots
  gemsData.forEach((gem) => {
    const key = getGemCombinationKey(gem);
    if (!gemCombinationsMap.has(key)) {
      gemCombinationsMap.set(key, {
        key,
        gemIds: [...gem.gem_ids_array],
        gemsCount: gem.gems_count,
        totalUsage: gem.usage_count,
        avgKeyLevel: gem.avg_keystone_level * gem.usage_count, // Weighted avg, will divide later
        slots: new Set([gem.item_slot]),
      });
    } else {
      const existingCombination = gemCombinationsMap.get(key)!;
      existingCombination.totalUsage += gem.usage_count;
      existingCombination.avgKeyLevel +=
        gem.avg_keystone_level * gem.usage_count;
      existingCombination.slots.add(gem.item_slot);
    }
  });

  // Convert to array and calculate final averages
  const gemCombinations = Array.from(gemCombinationsMap.values()).map(
    (comb) => {
      return {
        ...comb,
        avgKeyLevel: comb.avgKeyLevel / comb.totalUsage,
      };
    }
  );

  // Sort by total usage (most popular first)
  const sortedCombinations = gemCombinations.sort(
    (a, b) => b.totalUsage - a.totalUsage
  );

  // Take the top 8 combinations
  const topCombinations = sortedCombinations.slice(0, 8);

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {topCombinations.map((combination) => (
        <div
          key={combination.key}
          className="bg-slate-900 rounded-lg overflow-hidden"
        >
          {/* Card Header */}
          <div className="bg-slate-800 px-4 py-2 flex justify-between items-center">
            <h3 className="text-white font-semibold text-base">
              {combination.gemsCount}{" "}
              {combination.gemsCount === 1 ? "Gem" : "Gems"}
            </h3>
            <span className="text-xs text-slate-300">
              Used in {combination.slots.size}{" "}
              {combination.slots.size === 1 ? "slot" : "slots"}
            </span>
          </div>

          {/* Card Content */}
          <div className="p-4">
            {/* Gem Icons - using our utility function */}
            <div className="flex justify-center gap-2 mb-3">
              {combination.gemIds.map((gemId, idx) => (
                <a
                  key={`${combination.key}-gem-${idx}`}
                  href={`https://www.wowhead.com/item=${gemId}`}
                  data-wowhead={`item=${gemId}`}
                  className="block"
                >
                  <div className="border-2 border-slate-700 bg-slate-800 rounded w-10 h-10 overflow-hidden">
                    <Image
                      src={getGemIconUrl(gemId)}
                      alt={`Gem ${gemId}`}
                      width={40}
                      height={40}
                      className="rounded"
                      unoptimized
                    />
                  </div>
                </a>
              ))}
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 gap-2">
              <div className="bg-slate-800 p-2 rounded">
                <div className="text-slate-400 text-xs">Avg Key</div>
                <div className="text-white font-medium">
                  +{Math.round(combination.avgKeyLevel * 10) / 10}
                </div>
              </div>
              <div className="bg-slate-800 p-2 rounded">
                <div className="text-slate-400 text-xs">Usage</div>
                <div className="text-white font-medium">
                  {combination.totalUsage}
                </div>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
