import React from "react";
import { RoleLeaderboardEntry } from "@/types/warcraftlogs/globalLeaderboard";
import Link from "next/link";
import { normalizeServerName } from "@/utils/serverNameUtils";

interface LeaderboardTableProps {
  entries: RoleLeaderboardEntry[];
}

export const LeaderboardTable: React.FC<LeaderboardTableProps> = ({
  entries,
}) => {
  const getClassColor = (className: string) => {
    const formattedClass = className
      .replace(/([A-Z])/g, "-$1")
      .toLowerCase()
      .replace(/^-/, "");
    return `class-color--${formattedClass}`;
  };

  const getClassHoverStyles = (className: string) => {
    const baseClass = getClassColor(className);
    return `${baseClass} inline relative no-underline transition-all duration-200 hover:after:content-[''] hover:after:absolute hover:after:left-0 hover:after:bottom-[-2px] hover:after:h-[2px] hover:after:w-[100%] hover:after:bg-current`;
  };

  const buildCharacterUrl = (entry: RoleLeaderboardEntry) => {
    const region = entry.server_region.toLowerCase();
    const realm = normalizeServerName(entry.server_name);
    const name = entry.name.toLowerCase();
    return `/character/${region}/${realm}/${name}`;
  };

  return (
    <table className="w-full">
      <thead>
        <tr>
          <th className="text-center p-2">Rank</th>
          <th className="text-center p-2">Name</th>
          <th className="text-right p-2">Score</th>
        </tr>
      </thead>
      <tbody>
        {entries.map((entry) => (
          <tr key={entry.player_id} className="border-t border-gray-700">
            <td className="p-2 text-center">{entry.rank}</td>
            <td className="py-2">
              <div className="flex flex-col">
                <div>
                  <Link
                    href={buildCharacterUrl(entry)}
                    className={getClassHoverStyles(entry.class)}
                  >
                    {entry.name}
                  </Link>
                </div>
                <span className="text-xs text-gray-400">
                  {entry.server_name} ({entry.server_region})
                </span>
              </div>
            </td>
            <td className="text-right">{entry.total_score.toFixed(2)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};
