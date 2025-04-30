"use client";

import { useEffect } from "react";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import {
  useGetEnchantUsage,
  useGetGemUsage,
} from "@/hooks/useBuildsAnalysisApi";
import {
  groupBySlot,
  ITEM_SLOT_NAMES,
  ITEM_SLOT_DISPLAY_ORDER,
} from "@/utils/buildsAnalysis/dataTransformer";
import {
  classNameToPascalCase,
  specNameToPascalCase,
} from "@/utils/classandspecicons";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import BestEnchantsOverview from "./BestEnchantsOverview";
import BestGemsOverview from "./BestGemsOverview";
import EnchantTable from "./EnchantTable";
import GemTable from "./GemTable";
import {
  EnchantUsage,
  GemUsage,
} from "@/types/warcraftlogs/builds/buildsAnalysis";

interface EnchantGemsContentProps {
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function EnchantGemsContent({
  className,
  spec,
}: EnchantGemsContentProps) {
  // Fetch enchants and gems data
  const {
    data: enchantsData,
    isLoading: enchantsLoading,
    error: enchantsError,
  } = useGetEnchantUsage(className, spec);

  const {
    data: gemsData,
    isLoading: gemsLoading,
    error: gemsError,
  } = useGetGemUsage(className, spec);

  // Initialize Wowhead tooltips - still needed for gem tooltips
  useWowheadTooltips();

  // Refresh tooltips when gems data changes
  useEffect(() => {
    if (gemsData && typeof window !== "undefined" && window.$WowheadPower) {
      window.$WowheadPower.refreshLinks();
    }
  }, [gemsData]);

  // Handle loading states
  if (enchantsLoading || gemsLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

  // Handle error states
  if (enchantsError || gemsError) {
    return (
      <div className="bg-red-900/20 border border-red-500 rounded-md p-4 my-4">
        <h3 className="text-red-500 text-lg font-medium">
          Error loading enchants and gems data
        </h3>
        <p className="text-slate-300">
          {enchantsError instanceof Error
            ? enchantsError.message
            : gemsError instanceof Error
            ? gemsError.message
            : "An unknown error occurred"}
        </p>
      </div>
    );
  }

  // Check if data exists
  if (
    (!enchantsData || enchantsData.length === 0) &&
    (!gemsData || gemsData.length === 0)
  ) {
    return (
      <div className="bg-slate-800 rounded-lg p-5 text-center">
        <p className="text-slate-400">No enchants and gems data available.</p>
      </div>
    );
  }

  // Group enchants and gems by slot
  const enchantsBySlot = enchantsData
    ? groupBySlot<EnchantUsage>(enchantsData)
    : {};
  const gemsBySlot = gemsData ? groupBySlot<GemUsage>(gemsData) : {};

  // Filter for slots that have enchants or gems
  const enchantSlots = enchantsData
    ? ITEM_SLOT_DISPLAY_ORDER.filter(
        (slotId) => enchantsBySlot[slotId] && enchantsBySlot[slotId].length > 0
      )
    : [];

  const gemSlots = gemsData
    ? ITEM_SLOT_DISPLAY_ORDER.filter(
        (slotId) => gemsBySlot[slotId] && gemsBySlot[slotId].length > 0
      )
    : [];

  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-white mb-2">
          {classNameToPascalCase(className)} {specNameToPascalCase(spec)}{" "}
          Enchants & Gems
        </h1>
        <p className="text-slate-400">
          Most popular enchants and gems used by top{" "}
          {classNameToPascalCase(className)} {specNameToPascalCase(spec)}{" "}
          players in Mythic+.
        </p>
      </div>

      {/* Best Enchants Overview Section */}
      {enchantsData && enchantsData.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-white mb-4">
            Best Enchants by Slot
          </h2>
          <BestEnchantsOverview
            enchantsBySlot={enchantsBySlot}
            slotIds={enchantSlots}
          />
        </div>
      )}

      {/* Best Gems Overview Section */}
      {gemsData && gemsData.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-white mb-4">
            Most Popular Gem Combinations
          </h2>
          <BestGemsOverview gemsData={gemsData} />
        </div>
      )}

      {/* Enchants Tables Section */}
      {enchantsData && enchantsData.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-white mb-4">
            Enchants by Slot
          </h2>
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
            {enchantSlots.map((slotId) => (
              <EnchantTable
                key={`enchant-${slotId}`}
                slotName={ITEM_SLOT_NAMES[slotId]}
                enchants={enchantsBySlot[slotId]}
              />
            ))}
          </div>
        </div>
      )}

      {/* Gems Tables Section */}
      {gemsData && gemsData.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-white mb-4">Gems by Slot</h2>
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
            {gemSlots.map((slotId) => (
              <GemTable
                key={`gem-${slotId}`}
                slotName={ITEM_SLOT_NAMES[slotId]}
                gems={gemsBySlot[slotId]}
              />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
