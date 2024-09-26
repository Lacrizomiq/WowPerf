"use client";

import React, { useState, useEffect } from "react";
import { useGetRaiderioMythicPlusBestRuns } from "@/hooks/useRaiderioApi";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import RunsCard from "./RunsCard";
import DungeonSelector from "./Selector/DungeonSelector";
import RegionSelector from "./Selector/RegionSelector";
import { Dungeon } from "@/types/mythicPlusRuns";

const MythicPlusBestRuns: React.FC = () => {
  const season = "season-tww-1"; // Fixed season
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
    setDungeon(selectedDungeonSlug);
  };

  return (
    <div className="p-6 rounded-xl shadow-lg m-4">
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

      <RunsCard season={season} region={region} dungeon={dungeon} page={page} />

      {/* Pagination pourrait être ajoutée ici si nécessaire */}
    </div>
  );
};

export default MythicPlusBestRuns;
