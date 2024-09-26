"use client";

import React from "react";
import { useGetRaiderioMythicPlusBestRuns } from "@/hooks/useRaiderioApi";

interface RunsCardProps {
  season: string;
  region: string;
  dungeon: string;
  page: number;
}

const RunsCard: React.FC<RunsCardProps> = ({
  season,
  region,
  dungeon,
  page,
}) => {
  const {
    data: mythicPlusData,
    isLoading,
    error,
  } = useGetRaiderioMythicPlusBestRuns(season, region, dungeon, page);

  if (isLoading)
    return <div className="text-white text-center p-4">Loading...</div>;
  if (error)
    return (
      <div className="text-red-500 text-center p-4">
        Error: {(error as Error).message}
      </div>
    );
  if (!mythicPlusData || !mythicPlusData.rankings)
    return (
      <div className="text-yellow-500 text-center p-4">No data available</div>
    );

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {mythicPlusData.rankings.map((ranking: any) => (
        <div
          key={ranking.run.keystone_run_id}
          className="bg-deep-blue p-4 rounded-lg shadow-lg"
        >
          <h2 className="text-white text-lg font-bold mb-2">
            {ranking.run.dungeon.name}
          </h2>
          <p className="text-white text-sm mb-2">Rank: {ranking.rank}</p>
          <p className="text-white text-sm mb-2">
            Level: {ranking.run.mythic_level}
          </p>
          <p className="text-white text-sm mb-2">
            Time: {(ranking.run.clear_time_ms / 1000 / 60).toFixed(2)} minutes
          </p>
          <p className="text-white text-sm">Score: {ranking.score}</p>
        </div>
      ))}
    </div>
  );
};

export default RunsCard;
