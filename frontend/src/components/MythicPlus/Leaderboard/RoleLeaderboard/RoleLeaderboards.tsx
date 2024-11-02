"use client";

import React from "react";
import { useGetRoleLeaderboard } from "@/hooks/useWarcraftLogsApi";
import { RoleLeaderboardCard } from "./RoleLeaderboardCard";

export const RoleLeaderboards: React.FC = () => {
  const {
    data: dpsData,
    isLoading: isDpsLoading,
    error: dpsError,
  } = useGetRoleLeaderboard("dps", 10);
  const {
    data: healerData,
    isLoading: isHealerLoading,
    error: healerError,
  } = useGetRoleLeaderboard("healer", 10);
  const {
    data: tankData,
    isLoading: isTankLoading,
    error: tankError,
  } = useGetRoleLeaderboard("tank", 10);

  return (
    <div className=" mx-auto px-8 py-6 bg-black">
      <div className="grid grid-cols-1 xl:grid-cols-3 gap-6 max-w-[1400px] mx-auto">
        <RoleLeaderboardCard
          role="dps"
          title="Top Mythic+ DPS"
          data={dpsData || []}
          isLoading={isDpsLoading}
          error={dpsError}
        />
        <RoleLeaderboardCard
          role="healer"
          title="Top Mythic+ Healers"
          data={healerData || []}
          isLoading={isHealerLoading}
          error={healerError}
        />
        <RoleLeaderboardCard
          role="tank"
          title="Top Mythic+ Tanks"
          data={tankData || []}
          isLoading={isTankLoading}
          error={tankError}
        />
      </div>

      <div className="mt-6">
        <h2 className="text-left text-xl font-bold">
          Leaderboard for each dungeon
        </h2>
      </div>
    </div>
  );
};
