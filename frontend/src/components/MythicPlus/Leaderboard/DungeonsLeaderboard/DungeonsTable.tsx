import React from "react";
import Image from "next/image";
import Link from "next/link";
import { buildCharacterUrl, RegionType } from "@/utils/realmMappingUtility";
import { getSpecIcon, normalizeWowName } from "@/utils/classandspecicons";
import type { DungeonRanking } from "@/types/warcraftlogs/dungeonRankings";

interface DungeonLeaderboardTableProps {
  rankings: DungeonRanking[];
}

const DungeonLeaderboardTable: React.FC<DungeonLeaderboardTableProps> = ({
  rankings,
}) => {
  const formatTime = (duration: number) => {
    const totalSeconds = Math.floor(duration / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
  };

  const formatDate = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString();
  };

  const getClassColor = (className: string) => {
    const formattedClass = className
      .replace(/([A-Z])/g, "-$1")
      .toLowerCase()
      .replace(/^-/, "");
    return `class-color--${formattedClass}`;
  };

  const getClassHoverStyles = (className: string) => {
    const baseClass = getClassColor(className);
    return `${baseClass} inline relative no-underline transition-all duration-200 hover:after:content-[''] hover:after:absolute hover:after:left-0 hover:after:bottom-[-1px] hover:after:h-[1px] hover:after:w-[100%] hover:after:bg-current`;
  };

  const getAffixIcon = (affixId: number) => {
    // TODO: Implémenter la logique pour obtenir l'icône d'affix si nécessaire
    return `https://wow.zamimg.com/images/wow/icons/large/achievement_boss_archaedas.jpg`;
  };

  return (
    <table className="w-full">
      <thead>
        <tr className="text-gray-300">
          <th className="text-center p-2">Rank</th>
          <th className="text-left p-2">Player</th>
          <th className="text-center p-2">Key Level</th>
          <th className="text-center p-2">Time</th>
          <th className="text-center p-2">Date</th>
          <th className="text-center p-2">Affixes</th>
          <th className="text-right p-2">Score</th>
        </tr>
      </thead>
      <tbody>
        {rankings.map((ranking, index) => (
          <tr
            key={`${ranking.server.id}-${ranking.startTime}-${index}`}
            className="border-t border-gray-700"
          >
            <td className="p-2 text-center">{index + 1}</td>
            <td className="py-2">
              <div className="flex items-center gap-2">
                {ranking.spec && ranking.class && (
                  <div className="w-6 h-6">
                    <Image
                      src={getSpecIcon(
                        normalizeWowName(ranking.class),
                        normalizeWowName(ranking.spec)
                      )}
                      alt={`${ranking.spec} ${ranking.class}`}
                      width={24}
                      height={24}
                      className="rounded-sm"
                      unoptimized
                    />
                  </div>
                )}
                <div className="flex flex-col">
                  <Link
                    href={buildCharacterUrl(
                      ranking.name,
                      ranking.server.id,
                      ranking.server.region as RegionType
                    )}
                    className={getClassHoverStyles(ranking.class)}
                  >
                    {ranking.name}
                  </Link>
                  <span className="text-xs text-gray-400">
                    {ranking.server.name} ({ranking.server.region})
                  </span>
                </div>
              </div>
            </td>
            <td className="text-center">+{ranking.bracketData}</td>
            <td
              className={`text-center ${
                ranking.medal === "gold"
                  ? "text-yellow-400"
                  : ranking.medal === "silver"
                  ? "text-gray-300"
                  : "text-amber-600"
              }`}
            >
              {formatTime(ranking.duration)}
            </td>
            <td className="text-center">{formatDate(ranking.startTime)}</td>
            <td className="text-center">
              <div className="flex justify-center gap-1">
                {ranking.affixes.map((affixId, idx) => (
                  <div key={affixId} className="w-6 h-6">
                    <Image
                      src={getAffixIcon(affixId)}
                      alt={`Affix ${affixId}`}
                      width={24}
                      height={24}
                      className="rounded-sm"
                      unoptimized
                    />
                  </div>
                ))}
              </div>
            </td>
            <td className="text-right">{ranking.score.toFixed(1)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default DungeonLeaderboardTable;
