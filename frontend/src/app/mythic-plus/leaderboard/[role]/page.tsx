"use client";

import React from "react";
import { useGetRoleLeaderboard } from "@/hooks/useWarcraftLogsApi";
import { Role } from "@/types/warcraftlogs/globalLeaderboard";
import { LeaderboardTable } from "@/components/MythicPlus/Leaderboard/RoleLeaderboard/LeaderboardTable";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import Link from "next/link";
import { MoveLeft } from "lucide-react";

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
        return "Top Mythic+ DPS";
      case "tank":
        return "Top Mythic+ Tanks";
      case "healer":
        return "Top Mythic+ Healers";
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
    <div className="p-16 max-w-7xl mx-auto bg-black">
      <div className="flex flex-row gap-4 mb-6 items-center">
        <Link
          href="/mythic-plus/leaderboard"
          className="flex flex-row gap-2 items-center hover:underline"
        >
          <MoveLeft className="w-4 h-4" />
          Back to leaderboard
        </Link>
      </div>
      <p className="text-white text-xl font-bold">
        The {role} leaderboard only includes the very best players in terms of
        score
      </p>
      <p className="text-white text-md mb-6">
        The score is calculated based on the top score of each dungeon of the
        player.
      </p>
      <Card className="bg-deep-blue">
        <CardHeader>
          <CardTitle className="text-2xl font-bold text-center">
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
