// BestGemsOverview.tsx - Version harmonisée
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
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 text-center">
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

  // Calculate total usage for percentage calculation
  const totalUsageAll = topCombinations.reduce(
    (acc, comb) => acc + comb.totalUsage,
    0
  );

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {topCombinations.map((combination) => {
        // Calculate usage percentage for consistency with enchants display
        const usagePercentage = (combination.totalUsage / totalUsageAll) * 100;

        return (
          <div
            key={combination.key}
            className="bg-slate-800/30 rounded-lg overflow-hidden border border-slate-700"
          >
            {/* Card Header - Styled like enchants header */}
            <div className="bg-slate-800/50 px-3 py-2">
              <h3 className="text-white font-medium text-base">
                {combination.gemsCount}{" "}
                {combination.gemsCount === 1 ? "Gem" : "Gems"}
                <span className="text-xs text-slate-300 ml-2">
                  Used in {combination.slots.size}{" "}
                  {combination.slots.size === 1 ? "slot" : "slots"}
                </span>
              </h3>
            </div>

            {/* Card Content */}
            <div className="p-3">
              {/* Gem Icons */}
              <div className="flex mb-2">
                {combination.gemIds.map((gemId, idx) => (
                  <a
                    key={`${combination.key}-gem-${idx}`}
                    href={`https://www.wowhead.com/item=${gemId}`}
                    data-wowhead={`item=${gemId}`}
                    className={`block ${idx > 0 ? "-ml-1" : ""}`}
                  >
                    <div className="border border-slate-600 bg-slate-900 rounded-full w-8 h-8 overflow-hidden">
                      <Image
                        src={getGemIconUrl(gemId)}
                        alt={`Gem ${gemId}`}
                        width={32}
                        height={32}
                        className="rounded-full"
                        unoptimized
                      />
                    </div>
                  </a>
                ))}
              </div>

              {/* Stats - Styled like enchants stats */}
              <div className="grid grid-cols-2">
                <div>
                  <span className="text-slate-400 text-xs">Avg Key</span>
                  <div className="text-white font-medium">
                    +{Math.round(combination.avgKeyLevel * 10) / 10}
                  </div>
                </div>
                <div className="text-right">
                  <span className="text-slate-400 text-xs">Usage</span>
                  <div className="text-white font-medium">
                    {Math.round(usagePercentage)}
                  </div>
                </div>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}
