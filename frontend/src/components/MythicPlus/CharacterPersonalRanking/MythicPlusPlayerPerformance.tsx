import React, { useMemo } from "react";
import Image from "next/image";
import { getSpecIconById } from "@/utils/classandspecicons";
import { DUNGEON_ENCOUNTER_MAPPING } from "@/utils/s1_tww_mapping";
import {
  getPerformanceColorClass,
  getPerformanceColor,
} from "@/utils/rankingColor";
import { Trophy, Star, Award, Globe2 } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

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
  dungeonData?: { dungeons: Array<{ Slug: string; Icon: string }> };
}

const MythicPlusPlayerPerformance: React.FC<
  MythicPlusPlayerPerformanceProps
> = ({ playerData, dungeonData }) => {
  const encounterToSlugMapping = useMemo(() => {
    const mapping: Record<number, string> = {};
    Object.entries(DUNGEON_ENCOUNTER_MAPPING).forEach(([slug, data]) => {
      mapping[data.id] = slug;
    });
    return mapping;
  }, []);

  const getDungeonInfo = (encounterId: number) => {
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
      {/* Stats Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <Card className="bg-gray-900/50 border-gray-800">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="bg-blue-900/30 p-3 rounded-lg">
                <Trophy className="w-6 h-6 text-blue-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Best Perf. Avg</p>
                <p
                  className="text-2xl font-bold"
                  style={{
                    color: getPerformanceColor(
                      playerData.zoneRankings.bestPerformanceAverage
                    ),
                  }}
                >
                  {Math.floor(playerData.zoneRankings.bestPerformanceAverage)}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gray-900/50 border-gray-800">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="bg-blue-900/30 p-3 rounded-lg">
                <Globe2 className="w-6 h-6 text-blue-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400">World Rank</p>
                <p
                  className="text-2xl font-bold"
                  style={{
                    color: getPerformanceColor(allStarsData.rankPercent),
                  }}
                >
                  #{allStarsData.rank}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gray-900/50 border-gray-800">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="bg-blue-900/30 p-3 rounded-lg">
                <Star className="w-6 h-6 text-blue-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Region Rank</p>
                <p
                  className="text-2xl font-bold"
                  style={{
                    color: getPerformanceColor(allStarsData.rankPercent),
                  }}
                >
                  #{allStarsData.regionRank}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gray-900/50 border-gray-800">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="bg-blue-900/30 p-3 rounded-lg">
                <Award className="w-6 h-6 text-blue-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Server Rank</p>
                <p
                  className="text-2xl font-bold"
                  style={{
                    color: getPerformanceColor(allStarsData.rankPercent),
                  }}
                >
                  #{allStarsData.serverRank}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Additional Info */}
      <div className="text-center mb-6">
        <h2 className="text-2xl font-bold mb-4">
          Player rankings for {allStarsData.spec} Spec
        </h2>
        <div className="flex justify-center items-center space-x-6">
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
            Top %:{" "}
            <span
              style={{
                color: getPerformanceColor(allStarsData.rankPercent),
              }}
            >
              {formatTopPercent(allStarsData.rankPercent)}
            </span>
          </p>
          <p className="text-sm text-gray-400">
            Kills Logged:{" "}
            <span className="text-white">
              {playerData.zoneRankings.rankings.reduce(
                (acc, curr) => acc + curr.totalKills,
                0
              )}
            </span>
          </p>
        </div>
      </div>

      {/* Modern Table */}
      <div className="bg-gray-900/50 rounded-lg border border-gray-800">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-800">
                <th className="py-4 px-4 text-left text-gray-400 font-medium">
                  Dungeon
                </th>
                <th className="py-4 px-4 text-center text-gray-400 font-medium">
                  Best %
                </th>
                <th className="py-4 px-4 text-center text-gray-400 font-medium">
                  Score
                </th>
                <th className="py-4 px-4 text-center text-gray-400 font-medium">
                  Runs
                </th>
                <th className="py-4 px-4 text-center text-gray-400 font-medium">
                  Median %
                </th>
                <th className="py-4 px-4 text-right text-gray-400 font-medium">
                  All Stars
                </th>
              </tr>
            </thead>
            <tbody>
              {playerData.zoneRankings.rankings.map((dungeon) => {
                const dungeonInfo = getDungeonInfo(dungeon.encounter.id);

                return (
                  <tr
                    key={dungeon.encounter.id}
                    className={`border-b border-gray-800 hover:bg-gray-800/50 transition-colors 
                      ${dungeon.rankPercent === 0 ? "opacity-50" : ""}`}
                  >
                    <td className="py-4 px-4">
                      <div className="flex items-center gap-3">
                        {dungeonInfo && (
                          <Image
                            src={`https://wow.zamimg.com/images/wow/icons/large/${dungeonInfo.icon}.jpg`}
                            alt={dungeonInfo.name}
                            width={32}
                            height={32}
                            className="rounded"
                            unoptimized
                          />
                        )}
                        <span className="font-medium">
                          {dungeon.encounter.name}
                        </span>
                      </div>
                    </td>
                    <td className="py-4 px-4">
                      <div className="flex justify-center items-center gap-2">
                        {dungeon.spec && (
                          <Image
                            src={
                              getSpecIconById(
                                playerData.classID,
                                dungeon.spec
                              ) || ""
                            }
                            alt={dungeon.spec}
                            width={20}
                            height={20}
                            className="rounded"
                          />
                        )}
                        <div className="flex items-center gap-2">
                          <div className="relative w-24 bg-gray-800 h-2 rounded-full overflow-hidden">
                            <div
                              className="absolute left-0 top-0 h-full transition-all"
                              style={{
                                width: `${dungeon.rankPercent}%`,
                                backgroundColor: getPerformanceColor(
                                  dungeon.rankPercent
                                ),
                                opacity: 0.5,
                              }}
                            />
                          </div>
                          <span
                            style={{
                              color: getPerformanceColor(dungeon.rankPercent),
                            }}
                          >
                            {Math.floor(dungeon.rankPercent)}%
                          </span>
                        </div>
                      </div>
                    </td>
                    <td
                      className="py-4 px-4 text-center font-medium"
                      style={{
                        color: getPerformanceColor(dungeon.rankPercent),
                      }}
                    >
                      {dungeon.allStars.points.toFixed(0)}
                    </td>
                    <td className="py-4 px-4 text-center text-gray-300">
                      {dungeon.totalKills}
                    </td>
                    <td
                      className="py-4 px-4 text-center font-medium"
                      style={{
                        color: getPerformanceColor(dungeon.medianPercent),
                      }}
                    >
                      {Math.floor(dungeon.medianPercent)}%
                    </td>
                    <td
                      className="py-4 px-4 text-right font-medium"
                      style={{
                        color: getPerformanceColor(dungeon.rankPercent),
                      }}
                    >
                      {dungeon.allStars.rank}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default MythicPlusPlayerPerformance;
