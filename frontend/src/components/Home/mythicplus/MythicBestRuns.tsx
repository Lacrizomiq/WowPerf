import React, { useState, useEffect } from "react";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import RunsCard from "./RunsCard";
import DungeonSelector from "./Selector/DungeonSelector";
import RegionSelector from "./Selector/RegionSelector";
import { Dungeon } from "@/types/mythicPlusRuns";
import Image from "next/image";

const MythicPlusBestRuns: React.FC = () => {
  const season = "season-tww-2";
  const [region, setRegion] = useState("world");
  const [dungeon, setDungeon] = useState("all");
  const [page, setPage] = useState(0);

  const { data: dungeonData } = useGetBlizzardMythicDungeonPerSeason(season);

  const [dungeons, setDungeons] = useState<Dungeon[]>([]);

  useEffect(() => {
    if (dungeonData && dungeonData.dungeons) {
      setDungeons(dungeonData.dungeons);
    }
  }, [dungeonData]);

  const handleDungeonChange = (selectedDungeonSlug: string) => {
    console.log("selectedDungeonSlug", selectedDungeonSlug);
    setDungeon(selectedDungeonSlug);
  };

  return (
    <div className="relative w-full h-full bg-black">
      <div className="relative z-10 h-full overflow-auto">
        <div className="relative z-10 h-full overflow-auto">
          <div className="max-w-7xl mx-auto p-6">
            <h2 className="text-2xl font-bold text-white mb-6">
              Mythic+ Best Runs Leaderboard for TWW Season 1
            </h2>

            <div className="flex space-x-4 mb-6">
              <RegionSelector
                regions={["US", "EU", "TW", "KR", "CN"]}
                onRegionChange={setRegion}
                selectedRegion={region}
              />
              <DungeonSelector
                dungeons={dungeons}
                onDungeonChange={handleDungeonChange}
                selectedDungeon={dungeon}
              />
            </div>

            <RunsCard
              season={season}
              region={region}
              dungeon={dungeon}
              page={page}
            />

            {/* Pagination to be added here if needed */}
          </div>
        </div>
      </div>
    </div>
  );
};

export default MythicPlusBestRuns;
