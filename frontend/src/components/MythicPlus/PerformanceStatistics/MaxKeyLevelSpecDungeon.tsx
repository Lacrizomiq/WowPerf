// components/MythicPlus/DungeonPerformance/DungeonPerformance.tsx
"use client";

import React, { useMemo } from "react";
import { useGetMaxKeyLevelPerSpecAndDungeon } from "@/hooks/useWarcraftLogsApi";
import { MaxKeyLevelsPerSpecAndDungeon } from "@/types/warcraftlogs/globalLeaderboardAnalysis";

interface DungeonPerformanceProps {
  className: string;
  specName: string;
}

const DungeonPerformance: React.FC<DungeonPerformanceProps> = ({
  className,
  specName,
}) => {
  const {
    data: dungeonData,
    isLoading,
    error,
  } = useGetMaxKeyLevelPerSpecAndDungeon();

  // Filter dungeon data for the specific class and spec
  const specDungeonData = useMemo(() => {
    if (!dungeonData) return [];
    return dungeonData.filter(
      (dungeon) =>
        dungeon.class.toLowerCase() === className.toLowerCase() &&
        dungeon.spec.toLowerCase() === specName.toLowerCase()
    );
  }, [dungeonData, className, specName]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[200px]">
        <div className="text-white text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-white mx-auto"></div>
          <p className="mt-4">Loading dungeon performance...</p>
        </div>
      </div>
    );
  }

  if (error || !specDungeonData.length) {
    return (
      <div className="flex items-center justify-center min-h-[200px]">
        <p className="text-white text-center">
          No dungeon performance data available.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {specDungeonData.map((dungeon, index) => (
        <div
          key={index}
          className="card p-4"
          style={{ backgroundColor: "#112240" }}
        >
          <div className="flex justify-between items-center mb-2">
            <h3 className="font-bold">{dungeon.dungeon_name}</h3>
            <span className="text-lg font-bold">{dungeon.max_key_level}</span>
          </div>
          <p className="text-xs text-gray-400 mb-2">Max Key Level</p>
        </div>
      ))}
    </div>
  );
};

export default DungeonPerformance;
