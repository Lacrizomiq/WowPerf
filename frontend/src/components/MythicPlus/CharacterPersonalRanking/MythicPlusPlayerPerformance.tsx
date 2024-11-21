import React, { useMemo } from "react";
import Image from "next/image";
import { getSpecIconById } from "@/utils/classandspecicons";
import { DUNGEON_ENCOUNTER_MAPPING } from "@/utils/s1_tww_mapping";
import {
  getPerformanceColorClass,
  getPerformanceColor,
} from "@/utils/rankingColor";

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

interface DungeonInfo {
  slug: string;
  icon: string;
  name: string;
}

const MythicPlusPlayerPerformance: React.FC<
  MythicPlusPlayerPerformanceProps & {
    dungeonData?: { dungeons: Array<{ Slug: string; Icon: string }> };
  }
> = ({ playerData, dungeonData }) => {
  // Create an inverse mapping to find the slug from the ID
  const encounterToSlugMapping = useMemo(() => {
    const mapping: Record<number, string> = {};
    Object.entries(DUNGEON_ENCOUNTER_MAPPING).forEach(([slug, data]) => {
      mapping[data.id] = slug;
    });
    return mapping;
  }, []);

  // Function to get the dungeon info from the ID
  const getDungeonInfo = (encounterId: number): DungeonInfo | null => {
    const slug = encounterToSlugMapping[encounterId];
    if (!slug || !dungeonData) return null;

    const dungeonInfo = dungeonData.dungeons.find((d) => d.Slug === slug);
    if (!dungeonInfo) return null;

    return {
      slug,
      icon: dungeonInfo.Icon,
      name: DUNGEON_ENCOUNTER_MAPPING[slug].name,
    };
  };

  // Check if allStars data is available
  if (!playerData?.zoneRankings?.allStars?.[0]) {
    return (
      <div className="bg-deep-blue rounded-lg p-4 text-white shadow-2xl">
        <p className="text-center">No performance data available</p>
      </div>
    );
  }

  const allStarsData = playerData.zoneRankings.allStars[0];
  const specIcon = getSpecIconById(playerData.classID, allStarsData.spec);

  const formatTopPercent = (rankPercent: number): string => {
    const topPercent = 100 - rankPercent;
    if (topPercent < 0.01) return "< 0.01%";
    return `${topPercent.toFixed(3)}%`;
  };

  return (
    <div className="bg-deep-blue rounded-lg p-4 text-white shadow-2xl">
      {/* Header Section */}
      <div className="mb-6 text-center">
        <div className="flex justify-between items-center mb-4">
          <div>
            <h3 className="text-lg text-gray-400">Best Perf. Avg</h3>
            <p
              className="text-4xl font-bold"
              style={{
                color: getPerformanceColor(
                  playerData.zoneRankings.bestPerformanceAverage
                ),
              }}
            >
              {Math.floor(
                playerData.zoneRankings.bestPerformanceAverage
              ).toFixed(0)}
            </p>
          </div>
          <div className="text-right">
            <p className="text-sm text-gray-400">
              Median Perf. Avg:{" "}
              <span
                style={{
                  color: getPerformanceColor(
                    playerData.zoneRankings.medianPerformanceAverage
                  ),
                }}
              >
                {playerData.zoneRankings.medianPerformanceAverage.toFixed(1)}
              </span>
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
        <div className="flex flex-col items-center justify-center gap-4">
          <div>
            <h2 className="text-2xl font-bold">
              Player rankings for {allStarsData.spec} Spec
            </h2>
          </div>
          <div className="flex items-center gap-2 justify-center space-x-4">
            {/* World Rank */}
            <div className="flex items-center gap-2 ">
              <div>
                <p className="text-xs text-gray-400">World</p>
                <p
                  className="font-bold text-white"
                  style={{
                    color: getPerformanceColor(allStarsData.rankPercent),
                  }}
                >
                  #{allStarsData.rank}
                </p>
              </div>
            </div>

            {/* Region Rank */}
            <div className="flex items-center gap-2">
              <div>
                <p className="text-xs text-gray-400">Region</p>
                <p
                  className="font-bold text-white"
                  style={{
                    color: getPerformanceColor(allStarsData.rankPercent),
                  }}
                >
                  #{allStarsData.regionRank}
                </p>
              </div>
            </div>

            {/* Server Rank */}
            <div className="flex items-center gap-2">
              <div>
                <p className="text-xs text-gray-400">Realm</p>
                <p
                  className="font-bold text-white"
                  style={{
                    color: getPerformanceColor(allStarsData.rankPercent),
                  }}
                >
                  #{allStarsData.serverRank}
                </p>
              </div>
            </div>

            {/* Rank Percent */}
            <div className="flex items-center gap-2">
              <div>
                <p className="text-xs text-gray-400">Top %</p>
                <p
                  className="font-bold text-white"
                  style={{
                    color: getPerformanceColor(allStarsData.rankPercent),
                  }}
                >
                  {formatTopPercent(allStarsData.rankPercent)}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="text-left border-b border-gray-700">
              <th className="py-2 px-4">Dungeon</th>
              <th className="py-2 px-4">Best %</th>
              <th className="py-2 px-4">Score</th>
              <th className="py-2 px-4">Number of Runs</th>
              <th className="py-2 px-4">Median %</th>
              <th className="py-2 px-4 text-right">All Stars</th>
            </tr>
          </thead>
          <tbody>
            {playerData.zoneRankings.rankings.map((dungeon) => {
              const dungeonInfo = getDungeonInfo(dungeon.encounter.id);

              return (
                <tr
                  key={dungeon.encounter.id}
                  className="border-b border-gray-800 hover:bg-gray-800/30"
                >
                  {/* Dungeon name with icon */}
                  <td className="py-2 px-4">
                    <div className="flex items-center gap-2">
                      {dungeonInfo && (
                        <Image
                          src={`https://wow.zamimg.com/images/wow/icons/large/${dungeonInfo.icon}.jpg`}
                          alt={dungeonInfo.name}
                          width={24}
                          height={24}
                          className="rounded"
                          unoptimized
                        />
                      )}
                      <span>{dungeon.encounter.name}</span>
                    </div>
                  </td>

                  {/* Best % (points verts) */}
                  <td className="py-2 px-4 ">
                    <div className="flex items-center gap-2">
                      {dungeon.spec && (
                        <Image
                          src={
                            getSpecIconById(playerData.classID, dungeon.spec) ||
                            ""
                          }
                          alt={dungeon.spec}
                          width={20}
                          height={20}
                          className="rounded"
                        />
                      )}
                      <span
                        style={{
                          color: getPerformanceColor(dungeon.rankPercent),
                        }}
                      >
                        {Math.floor(dungeon.rankPercent)}
                      </span>
                    </div>
                  </td>

                  {/* Score */}
                  <td className="py-2 px-4 text-green-600">
                    {dungeon.allStars.points.toFixed(0)}
                  </td>

                  {/* Number of Runs */}
                  <td className="py-2 px-4">{dungeon.totalKills}</td>

                  {/* Median % */}
                  <td
                    className="py-2 px-4"
                    style={{
                      color: getPerformanceColor(dungeon.medianPercent),
                    }}
                  >
                    {Math.round(Math.floor(dungeon.medianPercent))}
                  </td>

                  {/* All Stars Rank */}
                  <td className="py-2 px-4 text-right">
                    <span
                      className="text-yellow-300 font-medium"
                      style={{
                        color: getPerformanceColor(dungeon.rankPercent),
                      }}
                    >
                      {dungeon.allStars.rank}
                    </span>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default MythicPlusPlayerPerformance;
