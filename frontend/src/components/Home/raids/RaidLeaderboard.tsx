import React, { useState, useEffect } from "react";
import { useGetBlizzardRaidsByExpansion } from "@/hooks/useBlizzardApi";
import RegionSelector from "./Selector/RegionSelector";
import RaidsSelector from "./Selector/RaidsSelector";
import LeaderBoardCards from "./LeaderboardCards";
import { StaticRaid } from "@/types/raids";
import Image from "next/image";

const RaidLeaderboard: React.FC = () => {
  const [region, setRegion] = useState("world");
  const [selectedRaid, setSelectedRaid] = useState<StaticRaid | null>(null);
  const [isMounted, setIsMounted] = useState(false);

  const { data: raidsData, isLoading: isRaidsLoading } =
    useGetBlizzardRaidsByExpansion("TWW");

  useEffect(() => {
    setIsMounted(true);
  }, []);

  useEffect(() => {
    if (raidsData && !isRaidsLoading) {
      const defaultRaid = raidsData.find(
        (raid) => raid.Slug === "nerubar-palace"
      );
      if (defaultRaid) {
        setSelectedRaid(defaultRaid);
      }
    }
  }, [raidsData, isRaidsLoading]);

  const handleRaidChange = (raid: StaticRaid) => {
    setSelectedRaid(raid);
  };

  return (
    <div className="relative w-full h-full">
      {isMounted && (
        <div className="fixed h-full w-full">
          <Image
            src="/tww.png"
            alt="World of Warcraft The War Within"
            layout="fill"
            objectFit="cover"
            quality={100}
            priority
            className="filter brightness-50"
          />
        </div>
      )}
      <div className="relative z-10 h-full overflow-auto">
        <div className="max-w-7xl mx-auto p-6">
          <h2 className="text-2xl font-bold text-white mb-6">
            Raid Leaderboard for The War Within
          </h2>

          <div className="flex space-x-4 mb-6">
            <RegionSelector
              regions={["US", "EU", "TW", "KR", "CN"]}
              onRegionChange={setRegion}
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
            raid="nerubar-palace"
            difficulty="mythic"
            region={region}
          />
        </div>
      </div>
    </div>
  );
};

export default RaidLeaderboard;
