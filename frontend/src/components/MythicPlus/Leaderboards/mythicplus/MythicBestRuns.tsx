// MythicPlusBestRuns.tsx
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
  const [currentPage, setCurrentPage] = useState(0); // Renommé pour plus de clarté

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
    setCurrentPage(0); // Réinitialiser la page lors du changement de donjon
  };

  const handleRegionChange = (newRegion: string) => {
    setRegion(newRegion);
    setCurrentPage(0); // Réinitialiser la page lors du changement de région
  };

  // Ajout de cette fonction pour gérer le changement de page
  const handlePageChange = (newPage: number) => {
    console.log("Changing page to:", newPage);
    setCurrentPage(newPage);
  };

  return (
    <div className="relative w-full h-full">
      <div className="relative z-10 h-full overflow-auto">
        <div className="relative z-10 h-full overflow-auto">
          <div className="max-w-7xl mx-auto">
            <h2 className="text-2xl font-bold text-white mb-6">
              Mythic+ Best Runs Leaderboard for TWW Season 1
            </h2>

            <div className="flex flex-wrap gap-3 mb-6">
              <RegionSelector
                regions={["US", "EU", "TW", "KR", "CN"]}
                onRegionChange={handleRegionChange}
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
              page={currentPage}
              onPageChange={handlePageChange}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default MythicPlusBestRuns;
