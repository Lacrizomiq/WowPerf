// components/MythicPlus/BuildsAnalysis/gear/DungeonGear.tsx
"use client";

import { useEffect } from "react";
import { useGetPopularItems } from "@/hooks/useBuildsAnalysisApi";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import {
  groupBySlot,
  ITEM_SLOT_NAMES,
  ITEM_SLOT_DISPLAY_ORDER,
} from "@/utils/buildsAnalysis/dataTransformer";
import ItemTable from "./ItemTable";
import { Dungeon } from "@/types/mythicPlusRuns";

interface DungeonGearProps {
  className: WowClassParam;
  spec: WowSpecParam;
  encounterId: string;
}

export default function DungeonGear({
  className,
  spec,
  encounterId,
}: DungeonGearProps) {
  // Get popular items data for this specific dungeon
  const {
    data: itemsData,
    isLoading,
    error,
  } = useGetPopularItems(className, spec, encounterId);

  // Get dungeon data to get the dungeon name
  const season = "season-tww-2";
  const { data: dungeonData } = useGetBlizzardMythicDungeonPerSeason(season);

  // Refresh tooltips Wowhead when data is loaded
  useEffect(() => {
    if (itemsData && typeof window !== "undefined" && window.$WowheadPower) {
      window.$WowheadPower.refreshLinks();
    }
  }, [itemsData]);

  // Find the dungeon name from its ID
  const getDungeonName = () => {
    if (!dungeonData?.dungeons) return `Dungeon #${encounterId}`;

    const dungeon = dungeonData.dungeons.find(
      (d: Dungeon) => d.EncounterID?.toString() === encounterId
    );

    return dungeon ? dungeon.Name : `Dungeon #${encounterId}`;
  };

  // Handle loading and error states
  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-900/20 border border-red-500 rounded-md p-4 my-4">
        <h3 className="text-red-500 text-lg font-medium">
          Error loading dungeon gear data
        </h3>
        <p className="text-slate-300">
          {error instanceof Error ? error.message : "An unknown error occurred"}
        </p>
      </div>
    );
  }

  if (!itemsData || itemsData.length === 0) {
    return (
      <div className="bg-slate-800 rounded-lg p-5 text-center">
        <p className="text-slate-400">
          No gear data available for this dungeon.
        </p>
      </div>
    );
  }

  // Group items by slot
  const itemsBySlot = groupBySlot(itemsData);

  // Filter slots that have items
  const filteredSlots = ITEM_SLOT_DISPLAY_ORDER.filter(
    (slotId) => itemsBySlot[slotId] && itemsBySlot[slotId].length > 0
  );

  const dungeonName = getDungeonName();

  return (
    <div>
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-white mb-4">
          Most Popular Gear for {dungeonName}
        </h2>
        <p className="text-slate-400">
          These items represent the most frequently used gear for {dungeonName}{" "}
          with {className} {spec} in Mythic+.
        </p>
      </div>

      {/* Limited grid to 2 columns maximum */}
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
        {filteredSlots.map((slotId) => (
          <ItemTable
            key={slotId}
            slotName={ITEM_SLOT_NAMES[slotId]}
            items={itemsBySlot[slotId]}
          />
        ))}
      </div>
    </div>
  );
}
