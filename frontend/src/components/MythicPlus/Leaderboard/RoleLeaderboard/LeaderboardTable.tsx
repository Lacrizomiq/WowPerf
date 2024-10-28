// components/MythicPlus/Leaderboard/LeaderboardTable.tsx

import React from "react";
import { RoleLeaderboardEntry } from "@/types/warcraftlogs/globalLeaderboard";

interface LeaderboardTableProps {
  entries: RoleLeaderboardEntry[];
}

export const LeaderboardTable: React.FC<LeaderboardTableProps> = ({
  entries,
}) => {
  // Helper function to get the proper class color
  const getClassColor = (className: string) => {
    // Convert class names like "DeathKnight" to "death-knight" for CSS classes
    const formattedClass = className
      .replace(/([A-Z])/g, "-$1")
      .toLowerCase()
      .replace(/^-/, "");
    return `class-color--${formattedClass}`;
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
                <span className={getClassColor(entry.class)}>{entry.name}</span>

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
