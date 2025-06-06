// BestEnchantsOverview.tsx - Version harmonisée
"use client";

import { EnchantUsage } from "@/types/warcraftlogs/builds/buildsAnalysis";
import { ITEM_SLOT_NAMES } from "@/utils/buildsAnalysis/dataTransformer";

interface BestEnchantsOverviewProps {
  enchantsBySlot: Record<number, EnchantUsage[]>;
  slotIds: number[];
}

export default function BestEnchantsOverview({
  enchantsBySlot,
  slotIds,
}: BestEnchantsOverviewProps) {
  // No enchants data
  if (slotIds.length === 0) {
    return (
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 text-center">
        <p className="text-slate-400">No enchant data available.</p>
      </div>
    );
  }

  // Find the best enchant for each slot (rank 1)
  const bestEnchants = slotIds
    .map((slotId) => {
      const slotEnchants = enchantsBySlot[slotId];
      // Find rank 1 enchant for this slot
      return slotEnchants.find((enchant) => enchant.rank === 1);
    })
    .filter(Boolean) as EnchantUsage[];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {bestEnchants.map((enchant) => {
        // Calculate usage percentage (we have count but not percentage in the data)
        const totalCount = enchantsBySlot[enchant.item_slot].reduce(
          (acc, item) => acc + item.usage_count,
          0
        );
        const usagePercentage = (enchant.usage_count / totalCount) * 100;

        return (
          <div
            key={`${enchant.item_slot}-${enchant.permanent_enchant_id}`}
            className="bg-slate-800/30 rounded-lg overflow-hidden border border-slate-700"
          >
            {/* Slot Name Header */}
            <div className="bg-slate-800/50 px-3 py-2">
              <h3 className="text-white font-medium text-base">
                {ITEM_SLOT_NAMES[enchant.item_slot]}
              </h3>
            </div>

            {/* Enchant Name */}
            <div className="px-3 pt-3 pb-2">
              <div className="text-purple-400 font-medium text-base">
                {enchant.permanent_enchant_name}
              </div>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 px-3 pb-3">
              <div>
                <span className="text-slate-400 text-xs">Avg Key</span>
                <div className="text-white font-medium">
                  +{Math.round(enchant.avg_keystone_level * 10) / 10}
                </div>
              </div>
              <div className="text-right">
                <span className="text-slate-400 text-xs">Usage</span>
                <div className="text-white font-medium">
                  {Math.round(usagePercentage)}%
                </div>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}
