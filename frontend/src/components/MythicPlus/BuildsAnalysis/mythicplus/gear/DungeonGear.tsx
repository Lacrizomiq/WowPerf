// DungeonGear.tsx - Version harmonisée
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
import { useRouter, usePathname, useSearchParams } from "next/navigation";

interface DungeonGearProps {
  className: WowClassParam;
  spec: WowSpecParam;
  encounterId: string;
  selectedSlotId: number | null;
}

export default function DungeonGear({
  className,
  spec,
  encounterId,
  selectedSlotId,
}: DungeonGearProps) {
  // Hooks pour la navigation
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

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

  // Handle slot selection change
  const handleSlotChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newSlotId = e.target.value;
    const params = new URLSearchParams(searchParams.toString());

    if (newSlotId === "all") {
      params.delete("slot_id");
    } else {
      params.set("slot_id", newSlotId);
    }

    router.push(`${pathname}?${params.toString()}`);
  };

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
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-600"></div>
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
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 text-center">
        <p className="text-slate-400">
          No gear data available for this dungeon.
        </p>
      </div>
    );
  }

  // Group items by slot
  const itemsBySlot = groupBySlot(itemsData);

  // Filter slots based on selectedSlotId
  const slotsToDisplay =
    selectedSlotId !== null
      ? [selectedSlotId]
      : ITEM_SLOT_DISPLAY_ORDER.filter(
          (slotId) => itemsBySlot[slotId] && itemsBySlot[slotId].length > 0
        );

  const dungeonName = getDungeonName();

  // Prepare title and description based on filtering
  const title =
    selectedSlotId !== null
      ? `Most Popular ${ITEM_SLOT_NAMES[selectedSlotId]} for ${dungeonName}`
      : `Most Popular Gear for ${dungeonName}`;

  const description =
    selectedSlotId !== null
      ? `These items represent the most frequently used ${ITEM_SLOT_NAMES[
          selectedSlotId
        ].toLowerCase()} 
       for ${dungeonName} with ${className} ${spec} in Mythic+.`
      : `These items represent the most frequently used gear for ${dungeonName} 
       with ${className} ${spec} in Mythic+.`;

  return (
    <div>
      <div className="mb-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4">
          <h2 className="text-2xl font-bold text-white">{title}</h2>

          {/* Slot filter with better UX */}
          <div className="inline-flex items-center bg-slate-800/50 rounded-lg border border-slate-700 hover:border-purple-400 transition-colors duration-150 px-3 py-2 shadow-sm">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              className="text-purple-400 mr-2"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"
              />
            </svg>
            <span className="text-purple-300 mr-3 whitespace-nowrap font-medium">
              Item Slot:
            </span>
            <div className="relative">
              <select
                className="appearance-none bg-slate-800/70 hover:bg-slate-700 text-white pl-3 pr-8 py-1 rounded focus:ring-2 focus:ring-purple-500 focus:outline-none cursor-pointer transition-colors duration-150"
                value={selectedSlotId?.toString() || "all"}
                onChange={handleSlotChange}
                aria-label="Select item slot to filter"
              >
                <option value="all">All Slots</option>
                {ITEM_SLOT_DISPLAY_ORDER.map((id) => (
                  <option key={id} value={id.toString()}>
                    {ITEM_SLOT_NAMES[id]}
                  </option>
                ))}
              </select>
              {/* Custom dropdown arrow */}
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-purple-300">
                <svg
                  className="fill-current h-4 w-4"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                >
                  <path d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" />
                </svg>
              </div>
            </div>
          </div>
        </div>
        <p className="text-slate-400">{description}</p>
      </div>

      {/* Grid layout */}
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
        {slotsToDisplay.map((slotId) => (
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
