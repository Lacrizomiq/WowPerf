import React from "react";
import Image from "next/image";
import { useGetRaiderioRaidLeaderboard } from "@/hooks/useRaiderioApi";
import {
  RaidRankings,
  RaidRanking,
  EncounterDefeated,
  EncounterPulled,
} from "@/types/raidLeaderboard";

interface LeaderBoardCardsProps {
  raid: string;
  difficulty: string;
  region: string;
  limit: number;
  page: number;
}

const LeaderBoardCards: React.FC<LeaderBoardCardsProps> = ({
  raid,
  difficulty,
  region,
  limit,
  page,
}) => {
  const { data, isLoading, error } = useGetRaiderioRaidLeaderboard(
    raid,
    difficulty,
    region,
    limit,
    page
  );

  if (isLoading)
    return <div className="text-white">Loading leaderboard data...</div>;
  if (error)
    return <div className="text-red-500">Error loading leaderboard data.</div>;
  if (!data || !data.raidRankings)
    return <div className="text-white">No leaderboard data available.</div>;

  return (
    <div className="space-y-4">
      {data.raidRankings.map((ranking: RaidRanking, index: number) => (
        <div
          key={index}
          className="bg-deep-blue bg-opacity-80 rounded-2xl overflow-hidden shadow-2xl glow-effect px-6 py-4"
        >
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center">
              {ranking.guild.logo && (
                <Image
                  src={ranking.guild.logo}
                  alt={`${ranking.guild.name} logo`}
                  width={60}
                  height={60}
                  className="rounded-full mr-4"
                />
              )}
              <div>
                <h3 className="text-2xl font-bold text-white">
                  {ranking.guild.name}
                </h3>
                <p className="text-gray-300">
                  {ranking.guild.realm.name} - {ranking.guild.region.name}
                </p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-xl font-bold text-white">
                {ranking.encountersDefeated.length}/8M
              </p>
              <p className="text-xl font-bold text-white">
                Rank: {ranking.rank}
              </p>
              <p className="text-gray-300">Region Rank: {ranking.regionRank}</p>
            </div>
          </div>
          {/*
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h4 className="text-lg font-semibold text-white mb-2">
                Encounters Defeated
              </h4>
              <ul className="space-y-2">
                {ranking.encountersDefeated.map(
                  (encounter: EncounterDefeated, eIndex: number) => (
                    <li key={eIndex} className="text-gray-300">
                      {encounter.slug}:{" "}
                      {new Date(encounter.firstDefeated).toLocaleString()}
                    </li>
                  )
                )}
              </ul>
            </div>
            <div>
              <h4 className="text-lg font-semibold text-white mb-2">
                Encounters Pulled
              </h4>
              <ul className="space-y-2">
                {ranking.encountersPulled.map(
                  (encounter: EncounterPulled, eIndex: number) => (
                    <li key={eIndex} className="text-gray-300">
                      {encounter.slug}: {encounter.numPulls} pulls
                      {encounter.isDefeated
                        ? " (Defeated)"
                        : ` (Best: ${encounter.bestPercent}%)`}
                    </li>
                  )
                )}
              </ul>
            </div>
          </div>
          */}
        </div>
      ))}
    </div>
  );
};

export default LeaderBoardCards;
