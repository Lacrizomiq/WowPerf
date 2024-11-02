import React, { useState } from "react";
import { useGetDungeonLeaderboard } from "@/hooks/useWarcraftLogsApi";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import DungeonLeaderboardSelector from "./Selector/DungeonLeaderboardSelector";
import ClassSpecSelector from "./Selector/ClassSpecSelector";
import type { WowClass } from "@/types/warcraftlogs/dungeonRankings";
import { CLASS_ICONS_MAPPING } from "@/utils/classandspecicons";
import DungeonLeaderboardTable from "./DungeonsTable";

const DungeonLeaderboard = () => {
  // States for filters
  const [selectedDungeon, setSelectedDungeon] = useState<number | null>(null);
  const [selectedRegion, setSelectedRegion] = useState<string | null>(null);
  const [selectedClass, setSelectedClass] = useState<WowClass | null>(null);
  const [selectedSpec, setSelectedSpec] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);

  // Get dungeons data
  const { data: dungeonData } =
    useGetBlizzardMythicDungeonPerSeason("season-tww-1");

  // Build an options object for the leaderboard
  const buildOptions = () => {
    const options: any = {};

    if (selectedRegion && selectedRegion !== "world") {
      options.serverRegion = selectedRegion;
    }

    if (selectedClass) {
      options.className = selectedClass;
    }

    if (selectedSpec) {
      options.specName = selectedSpec;
    }

    return Object.keys(options).length > 0 ? options : undefined;
  };

  // Get leaderboard data
  const { data: leaderboardData } = useGetDungeonLeaderboard(
    selectedDungeon || 0,
    currentPage,
    buildOptions()
  );

  const regions = ["US", "EU", "KR", "TW", "CN"];

  // Reset filters
  const handleReset = () => {
    setSelectedDungeon(null);
    setSelectedRegion(null);
    setSelectedClass(null);
    setSelectedSpec(null);
    setCurrentPage(1);
  };

  // Reset spec when class changes
  const handleClassChange = (newClass: WowClass | null) => {
    setSelectedClass(newClass);
    setSelectedSpec(null);
  };

  return (
    <div className="p-4 bg-black w-full h-full mb-12">
      <div className="flex flex-wrap gap-4 mb-6">
        <DungeonLeaderboardSelector
          dungeons={dungeonData?.dungeons || []}
          onDungeonChange={(id) => setSelectedDungeon(id)}
          selectedDungeonId={selectedDungeon}
        />

        <ClassSpecSelector
          selectedClass={selectedClass}
          selectedSpec={selectedSpec}
          onClassChange={handleClassChange}
          onSpecChange={setSelectedSpec}
          classMapping={CLASS_ICONS_MAPPING}
        />

        <button
          onClick={handleReset}
          className="px-4 py-2 bg-gradient-blue hover:bg-blue-700 rounded-lg text-white transition-colors shadow-2xl border-none"
        >
          Reset Filters
        </button>
      </div>

      {leaderboardData ? (
        <DungeonLeaderboardTable rankings={leaderboardData.rankings} />
      ) : (
        <div className="text-white text-center p-4 flex justify-center items-center mt-12">
          Please select a dungeon to see the leaderboard
        </div>
      )}
    </div>
  );
};

export default DungeonLeaderboard;
