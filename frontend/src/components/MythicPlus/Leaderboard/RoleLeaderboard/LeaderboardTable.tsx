import React from "react";
import { RoleLeaderboardEntry } from "@/types/warcraftlogs/globalLeaderboard";
import Link from "next/link";
import { buildCharacterUrl } from "@/utils/realmMappingUtility";
import type { RegionType } from "@/utils/realmMappingUtility";
import { realmService } from "@/utils/realmMappingUtility";

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

  const handleCharacterUrl = (entry: RoleLeaderboardEntry): string => {
    try {
      // Trouver le realm par son nom
      const realm = realmService.getRealmByName(entry.server_name);

      if (!realm) {
        console.warn(
          `Unable to find realm for server name ${entry.server_name}`,
          "Falling back to basic URL structure"
        );
        // Fallback à une URL basique si on ne trouve pas le realm
        const region = entry.server_region.toLowerCase();
        const name = entry.name.toLowerCase();
        return `/character/${region}/${entry.server_name.toLowerCase()}/${name}`;
      }

      // Utiliser buildCharacterUrl avec l'ID trouvé
      const url = buildCharacterUrl(
        entry.name,
        realm.id,
        entry.server_region as RegionType
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
                    href={handleCharacterUrl(entry)}
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
