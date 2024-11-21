import React, { useMemo } from "react";
import Image from "next/image";
import { getSpecIconById } from "@/utils/classandspecicons";
import {
  RAID_ENCOUNTER_MAPPING,
  getEncounterByWarcraftLogsId,
} from "@/utils/s1_tww_mapping";
import { getPerformanceColor } from "@/utils/rankingColor";
import { StaticRaid } from "@/types/raids";
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
              Total Kills:{" "}
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
              Raid Performance - {allStarsData.spec} Spec
            </h2>
          </div>
          <div className="flex items-center gap-2 justify-center space-x-4">
            <div className="flex items-center gap-2">
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

      {/* Boss Performance Table */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="text-left border-b border-gray-700">
              <th className="py-2 px-4">Boss</th>
              <th className="py-2 px-4">Best %</th>
              <th className="py-2 px-4">Points</th>
              <th className="py-2 px-4">Kills</th>
              <th className="py-2 px-4">Median %</th>
              <th className="py-2 px-4 text-right">Rank</th>
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
                  className="border-b border-gray-800 hover:bg-gray-800/30"
                >
                  <td className="py-2 px-4">
                    <div className="flex items-center gap-2">
                      {bossInfo && (
                        <Image
                          src={`https://wow.zamimg.com/images/wow/icons/large/${
                            bossInfo.icon || "ability_rogue_findweakness"
                          }.jpg`}
                          alt={bossInfo.name}
                          width={24}
                          height={24}
                          className="rounded"
                          unoptimized
                        />
                      )}
                      <span>{encounter.encounter.name}</span>
                    </div>
                  </td>

                  <td className="py-2 px-4">
                    <div className="flex items-center gap-2">
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
                      <span
                        style={{
                          color: getPerformanceColor(encounter.rankPercent),
                        }}
                      >
                        {Math.floor(encounter.rankPercent)}
                      </span>
                    </div>
                  </td>

                  <td className="py-2 px-4 text-green-600">
                    {encounter.allStars.points.toFixed(0)}
                  </td>

                  <td className="py-2 px-4">{encounter.totalKills}</td>

                  <td
                    className="py-2 px-4"
                    style={{
                      color: getPerformanceColor(encounter.medianPercent),
                    }}
                  >
                    {Math.round(Math.floor(encounter.medianPercent))}
                  </td>

                  <td className="py-2 px-4 text-right">
                    <span
                      style={{
                        color: getPerformanceColor(encounter.rankPercent),
                      }}
                    >
                      #{encounter.allStars.rank}
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

export default RaidsPlayerPerformance;
