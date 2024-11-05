import React from "react";
import Image from "next/image";
import { getSpecIconById } from "@/utils/classandspecicons";
import { Trophy, Globe2, Server } from "lucide-react"; // Utilisation des icônes Lucide

interface DungeonPerformance {
  encounter: {
    id: number;
    name: string;
  };
  rankPercent: number;
  medianPercent: number;
  totalKills: number;
  fastestKill: number;
  spec: string;
  allStars: {
    points: number;
    rank: number;
  };
}

interface MythicPlusPlayerPerformanceProps {
  playerData: {
    classID: number;
    zoneRankings: {
      bestPerformanceAverage: number;
      medianPerformanceAverage: number;
      rankings: DungeonPerformance[];
      allStars: Array<{
        points: number;
        spec: string;
        rank: number;
        regionRank: number;
        serverRank: number;
        rankPercent: number;
      }>;
    };
  };
}

const MythicPlusPlayerPerformance: React.FC<
  MythicPlusPlayerPerformanceProps
> = ({ playerData }) => {
  // Get spec icon using classID and spec from allStars
  const specIcon = getSpecIconById(
    playerData.classID,
    playerData.zoneRankings.allStars[0].spec
  );

  return (
    <div className="bg-deep-blue rounded-lg p-4 text-white shadow-2xl">
      {/* Header Section */}
      <div className="mb-6 text-center">
        <div className="flex justify-between items-center mb-4">
          <div>
            <h3 className="text-lg text-gray-400">Best Perf. Avg</h3>
            <p className="text-4xl font-bold text-yellow-400">
              {playerData.zoneRankings.bestPerformanceAverage.toFixed(1)}
            </p>
          </div>
          <div className="text-right">
            <p className="text-sm text-gray-400">
              Median Perf. Avg:{" "}
              {playerData.zoneRankings.medianPerformanceAverage.toFixed(1)}
            </p>
            <p className="text-sm text-gray-400">
              Kills Logged:{" "}
              {playerData.zoneRankings.rankings.reduce(
                (acc, curr) => acc + curr.totalKills,
                0
              )}
            </p>
          </div>
        </div>
        <div className="flex flex-col items-center justify-center gap-2">
          <div>
            <h2 className="text-2xl font-bold">
              Player rankings for {playerData.zoneRankings.allStars[0].spec}{" "}
              Spec
            </h2>
          </div>
          <div className="flex items-center gap-2">
            {/* World Rank */}
            <div className="flex items-center gap-2">
              <div>
                <p className="text-xs text-gray-400">World</p>
                <p className="font-bold text-white">
                  #{playerData.zoneRankings.allStars[0].rank}
                </p>
              </div>
            </div>

            {/* Region Rank */}
            <div className="flex items-center gap-2">
              <div>
                <p className="text-xs text-gray-400">Region</p>
                <p className="font-bold text-white">
                  #{playerData.zoneRankings.allStars[0].regionRank}
                </p>
              </div>
            </div>

            {/* Server Rank */}
            <div className="flex items-center gap-2">
              <div>
                <p className="text-xs text-gray-400">Realm</p>
                <p className="font-bold text-white">
                  #{playerData.zoneRankings.allStars[0].serverRank}
                </p>
              </div>
            </div>

            {/* Rank Percent */}
            <div className="flex items-center gap-2">
              <div>
                <p className="text-xs text-gray-400">Top %</p>
                <p className="font-bold text-white">
                  {(
                    100 - playerData.zoneRankings.allStars[0].rankPercent
                  ).toFixed(3)}
                  %
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Table Section */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="text-left border-b border-gray-700">
              <th className="py-2 px-4">Dungeon</th>
              <th className="py-2 px-4">Best %</th>
              <th className="py-2 px-4">Points</th>
              <th className="py-2 px-4">Number of Runs</th>
              <th className="py-2 px-4">Median %</th>
              <th className="py-2 px-4 text-right">All Stars</th>
            </tr>
          </thead>
          <tbody>
            {playerData.zoneRankings.rankings.map((dungeon) => (
              <tr
                key={dungeon.encounter.id}
                className="border-b border-gray-800 hover:bg-gray-800/30"
              >
                <td className="py-2 px-4 ">
                  <span>{dungeon.encounter.name}</span>
                </td>
                <td className="py-2 px-4 flex items-center gap-2">
                  <span>{dungeon.rankPercent.toFixed()}</span>
                  {getSpecIconById(playerData.classID, dungeon.spec) && (
                    <Image
                      src={getSpecIconById(playerData.classID, dungeon.spec)}
                      alt={dungeon.spec}
                      width={20}
                      height={20}
                      className="rounded"
                    />
                  )}
                </td>
                <td className="py-2 px-4 text-green-400">
                  {dungeon.allStars.points.toFixed(1)}
                </td>
                <td className="py-2 px-4">{dungeon.totalKills}</td>
                <td className="py-2 px-4">
                  {Math.round(Math.floor(dungeon.medianPercent))}
                </td>
                <td className="py-2 px-4 text-right">
                  <span className="text-yellow-300 font-medium">
                    #{dungeon.allStars.rank}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default MythicPlusPlayerPerformance;
