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
      <div className="grid grid-cols-1 xl:grid-cols-3 gap-6 max-w-[1400px] mx-auto">
        <RoleLeaderboardCard
          role="dps"
          title="All Stars (DPS)"
          data={dpsData || []}
          isLoading={isDpsLoading}
          error={dpsError}
        />
        <RoleLeaderboardCard
          role="tank"
          title="All Stars (Tanks)"
          data={tankData || []}
          isLoading={isTankLoading}
          error={tankError}
        />
        <RoleLeaderboardCard
          role="healer"
          title="All Stars (Healers)"
          data={healerData || []}
          isLoading={isHealerLoading}
          error={healerError}
        />
      </div>
    </div>
  );
};
