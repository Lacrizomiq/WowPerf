import React from "react";
import Link from "next/link";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { buildCharacterUrl, realmService } from "@/utils/realmMappingUtility";
import type { RegionType } from "@/utils/realmMappingUtility";
import type { BestTenPlayerPerSpec } from "@/types/warcraftlogs/globalLeaderboardAnalysis";

interface TopPlayersTableProps {
  players: BestTenPlayerPerSpec[];
  className: string;
}

const TopPlayersTable: React.FC<TopPlayersTableProps> = ({
  players,
  className,
}) => {
  // Helper function to format class name for CSS
  const formatClassNameForCSS = (className: string): string => {
    return className.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
  };

  // Build character URL from player data
  const handleCharacterUrl = (player: BestTenPlayerPerSpec): string => {
    try {
      const realm = realmService.getRealmByName(player.server_name);

      if (!realm) {
        console.warn(
          `Unable to find realm for server name ${player.server_name}`,
          "Falling back to basic URL structure"
        );
        const region = player.server_region.toLowerCase();
        const name = player.name.toLowerCase();
        return `/character/${region}/${player.server_name.toLowerCase()}/${name}`;
      }

      const url = buildCharacterUrl(
        player.name,
        realm.id,
        player.server_region as RegionType
      );

      if (!url) {
        throw new Error("Failed to build character URL");
      }

      return url;
    } catch (error) {
      console.error("Error building character URL:", error);
      return "#";
    }
  };

  const visiblePlayers = players.slice(0, 10);

  return (
    <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 overflow-x-auto">
      <Table>
        <TableHeader>
          <TableRow className="border-slate-700">
            <TableHead className="w-16">Rank</TableHead>
            <TableHead>Player</TableHead>
            <TableHead>Realm</TableHead>
            <TableHead className="text-right">M+ Score</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {visiblePlayers.map((player, index) => (
            <TableRow
              key={`${player.name}-${index}`}
              className="border-slate-700"
            >
              <TableCell className="font-medium">
                <span className={index < 3 ? "text-purple-400" : ""}>
                  {player.rank}
                </span>
              </TableCell>
              <TableCell>
                <Link
                  href={handleCharacterUrl(player)}
                  className="hover:underline"
                  style={{
                    color: `var(--color-${formatClassNameForCSS(className)})`,
                  }}
                >
                  {player.name}
                </Link>
              </TableCell>
              <TableCell>
                {player.server_region}-{player.server_name}
              </TableCell>
              <TableCell className="text-right">
                {Math.round(player.total_score).toLocaleString()}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      {players.length > 10 && (
        <div className="mt-4 text-center">
          <Button
            variant="outline"
            className="border-slate-700 hover:bg-slate-700"
          >
            View More Players
          </Button>
        </div>
      )}
    </div>
  );
};

export default TopPlayersTable;
