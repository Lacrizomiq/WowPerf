import React, { useState, useEffect } from "react";
import { useGetBlizzardRaidsByExpansion } from "@/hooks/useBlizzardApi";
import RegionSelector from "./Selector/RegionSelector";
import RaidsSelector from "./Selector/RaidsSelector";
import LeaderBoardCards from "./LeaderboardCards";
import { StaticRaid } from "@/types/raids";
import Image from "next/image";

const regions = [
  { value: "us", label: "US" },
  { value: "eu", label: "EU" },
  { value: "tw", label: "TW" },
  { value: "kr", label: "KR" },
  { value: "cn", label: "CN" },
];

const RaidLeaderboard: React.FC = () => {
  const [region, setRegion] = useState("world");
  const [selectedRaid, setSelectedRaid] = useState<StaticRaid | null>(null);
  const [currentPage, setCurrentPage] = useState(0);

  const { data: raidsData, isLoading: isRaidsLoading } =
    useGetBlizzardRaidsByExpansion("TWW");

  useEffect(() => {
    if (raidsData && !isRaidsLoading) {
      const defaultRaid = raidsData.find(
        (raid) => raid.Slug === "liberation-of-undermine"
      );
      if (defaultRaid) {
        setSelectedRaid(defaultRaid);
      }
    }
  }, [raidsData, isRaidsLoading]);

  const handleRaidChange = (raid: StaticRaid) => {
    setSelectedRaid(raid);
  };

  const handleRegionChange = (newRegion: string) => {
    setRegion(newRegion);
  };

  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
  };

  return (
    <div className="relative w-full h-full">
      <div className="relative z-10 h-full overflow-auto">
        <h2 className="text-2xl font-bold text-white mb-6">
          Raid Leaderboard for The War Within
        </h2>

        <div className="flex space-x-4 mb-6">
          <RegionSelector
            regions={regions}
            onRegionChange={handleRegionChange}
            selectedRegion={region}
          />
          {raidsData && (
            <RaidsSelector
              raids={raidsData}
              onRaidChange={handleRaidChange}
              selectedRaid={selectedRaid}
            />
          )}
        </div>

        {selectedRaid && (
          <div className="text-white mb-6">
            <h3 className="text-xl font-semibold">
              Selected Raid: {selectedRaid.Name}
            </h3>
          </div>
        )}

        <LeaderBoardCards
          raid={selectedRaid?.Slug || "liberation-of-undermine"}
          difficulty="mythic"
          region={region}
          limit={20}
          page={currentPage}
          onPageChange={handlePageChange}
        />
      </div>
    </div>
  );
};

export default RaidLeaderboard;
