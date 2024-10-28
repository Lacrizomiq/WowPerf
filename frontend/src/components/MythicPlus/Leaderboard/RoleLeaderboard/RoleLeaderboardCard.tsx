import React from "react";
import Link from "next/link";
import {
  RoleLeaderboardEntry,
  Role,
} from "@/types/warcraftlogs/globalLeaderboard";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { LeaderboardTable } from "./LeaderboardTable";

interface RoleLeaderboardCardProps {
  role: Role;
  title: string;
  data: RoleLeaderboardEntry[];
  isLoading: boolean;
  error: Error | null;
}

export const RoleLeaderboardCard: React.FC<RoleLeaderboardCardProps> = ({
  role,
  title,
  data,
  isLoading,
  error,
}) => {
  if (error) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <Card className="w-full h-full bg-deep-blue rounded-lg glow-effect">
      <CardHeader>
        <CardTitle className="text-xl font-bold">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div>Loading...</div>
        ) : (
          <>
            <LeaderboardTable entries={data} />
            <div className="mt-4 text-center border-t border-gray-700 pt-4">
              <Link
                href={`/mythic-plus/leaderboard/${role}`}
                className="text-blue-400 hover:text-blue-300 transition-colors"
              >
                View More...
              </Link>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
};
