import React from "react";
import Image from "next/image";
import { getSpecIconById } from "@/utils/classandspecicons";
import {
  RAID_ENCOUNTER_MAPPING,
  getEncounterByWarcraftLogsId,
} from "@/utils/s1_tww_mapping";
import { getPerformanceColor } from "@/utils/rankingColor";
import { StaticRaid } from "@/types/raids";
import { Card, CardContent } from "@/components/ui/card";
import { Trophy, Star, Award, Globe2 } from "lucide-react";

interface RaidEncounter {
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

interface RaidPlayerPerformanceProps {
  playerData: {
    classID: number;
    zoneRankings: {
      bestPerformanceAverage: number;
      medianPerformanceAverage: number;
      rankings: RaidEncounter[];
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
  staticRaid?: StaticRaid;
}

const RaidsPlayerPerformance: React.FC<RaidPlayerPerformanceProps> = ({
  playerData,
  staticRaid,
}) => {
  if (!playerData?.zoneRankings?.allStars?.[0]) {
    return (
      <div className="bg-deep-blue rounded-lg p-4 text-white shadow-2xl">
        <p className="text-center">No raid performance data available</p>
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
          Raid Performance - {allStarsData.spec} Spec
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
            Total Kills:{" "}
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
                  Boss
                </th>
                <th className="py-4 px-4 text-center text-gray-400 font-medium">
                  Best %
                </th>
                <th className="py-4 px-4 text-center text-gray-400 font-medium">
                  Score
                </th>
                <th className="py-4 px-4 text-center text-gray-400 font-medium">
                  Kills
                </th>
                <th className="py-4 px-4 text-center text-gray-400 font-medium">
                  Median %
                </th>
                <th className="py-4 px-4 text-right text-gray-400 font-medium">
                  Rank
                </th>
              </tr>
            </thead>
            <tbody>
              {playerData.zoneRankings.rankings.map((encounter) => {
                const bossInfo = getEncounterByWarcraftLogsId(
                  encounter.encounter.id
                );

                return (
                  <tr
                    key={encounter.encounter.id}
                    className={`border-b border-gray-800 hover:bg-gray-800/50 transition-colors
                      ${encounter.rankPercent === 0 ? "opacity-50" : ""}`}
                  >
                    <td className="py-4 px-4">
                      <div className="flex items-center gap-3">
                        {bossInfo && (
                          <Image
                            src={`https://wow.zamimg.com/images/wow/icons/large/${
                              bossInfo.icon || "ability_rogue_findweakness"
                            }.jpg`}
                            alt={bossInfo.name}
                            width={32}
                            height={32}
                            className="rounded"
                            unoptimized
                          />
                        )}
                        <span className="font-medium">
                          {encounter.encounter.name}
                        </span>
                      </div>
                    </td>
                    <td className="py-4 px-4">
                      <div className="flex justify-center items-center gap-2">
                        {encounter.spec && (
                          <Image
                            src={
                              getSpecIconById(
                                playerData.classID,
                                encounter.spec
                              ) || ""
                            }
                            alt={encounter.spec}
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
                                width: `${encounter.rankPercent}%`,
                                backgroundColor: getPerformanceColor(
                                  encounter.rankPercent
                                ),
                                opacity: 0.5,
                              }}
                            />
                          </div>
                          <span
                            style={{
                              color: getPerformanceColor(encounter.rankPercent),
                            }}
                          >
                            {Math.floor(encounter.rankPercent)}%
                          </span>
                        </div>
                      </div>
                    </td>
                    <td
                      className="py-4 px-4 text-center font-medium"
                      style={{
                        color: getPerformanceColor(encounter.rankPercent),
                      }}
                    >
                      {encounter.allStars.points.toFixed(0)}
                    </td>
                    <td className="py-4 px-4 text-center text-gray-300">
                      {encounter.totalKills}
                    </td>
                    <td
                      className="py-4 px-4 text-center font-medium"
                      style={{
                        color: getPerformanceColor(encounter.medianPercent),
                      }}
                    >
                      {Math.floor(encounter.medianPercent)}%
                    </td>
                    <td
                      className="py-4 px-4 text-right font-medium"
                      style={{
                        color: getPerformanceColor(encounter.rankPercent),
                      }}
                    >
                      #{encounter.allStars.rank}
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

export default RaidsPlayerPerformance;
