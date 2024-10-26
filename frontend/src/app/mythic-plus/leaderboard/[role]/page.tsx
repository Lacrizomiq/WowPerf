"use client";

import React from "react";
import { useGetRoleLeaderboard } from "@/hooks/useWarcraftLogsApi";
import { Role } from "@/types/warcraftlogs/globalLeaderboard";
import { LeaderboardTable } from "@/components/MythicPlus/Leaderboard/RoleLeaderboard/LeaderboardTable";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

interface RolePageProps {
  params: {
    role: Role;
  };
}

const RolePage: React.FC<RolePageProps> = ({ params }) => {
  const { role } = params;
  const { data, isLoading, error } = useGetRoleLeaderboard(role, 100);

  const getTitleByRole = (role: Role) => {
    switch (role) {
      case "dps":
        return "All Stars (DPS)";
      case "tank":
        return "All Stars (Tanks)";
      case "healer":
        return "All Stars (Healers)";
      default:
        return "Leaderboard";
    }
  };

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <Card className="bg-deep-blue">
        <CardHeader>
          <CardTitle className="text-2xl font-bold">
            {getTitleByRole(role)}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <LeaderboardTable entries={data || []} />
        </CardContent>
      </Card>
    </div>
  );
};

export default RolePage;
