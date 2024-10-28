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
    <div className="container mx-auto px-4 p-6">
      <p className="text-left text-xl font-bold">
        Data are updated every 24 hours.
      </p>
      <p className="text-left text-sm">
        Due to technical limitations, the data is not updated in real-time and
        only includes the very best players of each role
      </p>
      <p className="text-left text-sm mb-6">
        The data is provided by{" "}
        <a
          href="https://www.warcraftlogs.com/about"
          target="_blank"
          className="text-blue-400 hover:text-blue-300 transition-colors"
        >
          Warcraft Logs
        </a>{" "}
        , check them out for more detailed data.
      </p>
      <div className="grid grid-cols-1 xl:grid-cols-3 gap-6 max-w-[1400px] mx-auto">
        <RoleLeaderboardCard
          role="dps"
          title="Top Mythic+ DPS"
          data={dpsData || []}
          isLoading={isDpsLoading}
          error={dpsError}
        />
        <RoleLeaderboardCard
          role="tank"
          title="Top Mythic+ Tanks"
          data={tankData || []}
          isLoading={isTankLoading}
          error={tankError}
        />
        <RoleLeaderboardCard
          role="healer"
          title="Top Mythic+ Healers"
          data={healerData || []}
          isLoading={isHealerLoading}
          error={healerError}
        />
      </div>
    </div>
  );
};
