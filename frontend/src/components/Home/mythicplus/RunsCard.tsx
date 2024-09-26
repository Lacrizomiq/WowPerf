"use client";

import React, { useMemo } from "react";
import { useGetRaiderioMythicPlusBestRuns } from "@/hooks/useRaiderioApi";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import { Dungeon } from "@/types/mythicPlusRuns";
import Image from "next/image";

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
  const { data: dungeonData } =
    useGetBlizzardMythicDungeonPerSeason("season-tww-1");
  const {
    data: mythicPlusData,
    isLoading,
    error,
  } = useGetRaiderioMythicPlusBestRuns(season, region, dungeon, page);

  const dungeonMap = useMemo(() => {
    if (dungeonData?.dungeons) {
      return dungeonData.dungeons.reduce(
        (acc: Record<string, Dungeon>, dungeon: Dungeon) => {
          acc[dungeon.Slug.toLowerCase()] = dungeon;
          return acc;
        },
        {}
      );
    }
    return {};
  }, [dungeonData]);

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
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-4">
      {mythicPlusData.rankings.map((ranking: any) => {
        const dungeonSlug = ranking.run.dungeon.slug.toLowerCase();
        const dungeonInfo = dungeonMap[dungeonSlug];

        return (
          <div
            key={ranking.run.keystone_run_id}
            className="rounded-xl overflow-hidden shadow-lg h-64"
          >
            <div className="bg-deep-blue bg-opacity-80 h-full flex flex-col justify-between">
              <div
                className="flex justify-between h-full"
                style={{
                  backgroundImage: `url(${dungeonInfo?.MediaURL})`,
                  backgroundSize: "cover",
                  backgroundPosition: "center",
                  backgroundRepeat: "no-repeat",
                }}
              >
                <h2 className="text-white text-lg font-bold p-4">
                  {ranking.run.dungeon.name}
                </h2>
              </div>
              <div className="p-4">
                <p className="text-white text-sm mb-2">Rank: {ranking.rank}</p>
                <p className="text-white text-sm mb-2">
                  Level: {ranking.run.mythic_level}
                </p>
                <p className="text-white text-sm mb-2">
                  Time: {(ranking.run.clear_time_ms / 1000 / 60).toFixed(2)}{" "}
                  minutes
                </p>
                <p className="text-white text-sm">Score: {ranking.score}</p>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default RunsCard;
