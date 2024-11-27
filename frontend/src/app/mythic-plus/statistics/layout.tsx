// app/mythic-plus/statistics/layout.tsx
"use client";

import React from "react";
import Link from "next/link";
import { StatsProvider } from "@/providers/StatsContext";
import { useStats } from "@/providers/StatsContext";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import { useGetDungeonStats } from "@/hooks/useRaiderioApi";
import DungeonSelector from "@/components/Home/mythicplus/Stats/Selector/DungeonSelector";
import RegionSelector from "@/components/Home/mythicplus/Stats/Selector/RegionSelector";
import { DungeonStat } from "@/types/dungeonStats";

function StatisticsHeader() {
  const { region, dungeon, season } = useStats();
  const { data: statsData } = useGetDungeonStats(season, region);

  const currentDungeonStats =
    statsData?.find((stat: DungeonStat) => stat.dungeon_slug === dungeon) ||
    statsData?.[0];

  return (
    <>
      <h2 className="text-2xl font-bold text-white mb-4 mt-8">
        Dungeon Statistics for{" "}
        {dungeon
          .split("-")
          .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(" ")}{" "}
        dungeon in region: {region.toUpperCase()}
      </h2>

      <div className="space-y-4">
        <div className="p-4">
          <p className="text-white">
            <span className="font-bold">Last update:</span>{" "}
            {statsData &&
              statsData[0]?.updated_at &&
              new Intl.DateTimeFormat("en-US", {
                weekday: "long",
                day: "2-digit",
                month: "long",
                year: "numeric",
              }).format(new Date(statsData[0].updated_at))}
          </p>
          <p className="text-white mt-4">
            The data is updated every week on Tuesday, coming from{" "}
            <a
              href="https://raider.io"
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-500"
            >
              Raider.io
            </a>{" "}
            best mythic + runs.
          </p>
          <p className="text-white">
            As it aggregates data from the best runs and the very top teams, it
            may not be 100% accurate to determine the best class / spec or team
            composition as some players / teams are highlighted many times and
            can skew the data.
          </p>
          <p className="text-white mt-4">
            Remember to play the game and enjoy the journey with your friends on
            the class you love!
          </p>
        </div>
      </div>
    </>
  );
}

function StatisticsNav() {
  const { region, dungeon, season, setRegion, setDungeon } = useStats();

  const { data: dungeonData } =
    useGetBlizzardMythicDungeonPerSeason("season-tww-1");

  const { data: statsData } = useGetDungeonStats(season, region);

  const currentDungeonStats =
    statsData?.find((stat: DungeonStat) => stat.dungeon_slug === dungeon) ||
    statsData?.[0];

  const getKeyRange = (levelStats: Record<string, number>) => {
    const levels = Object.keys(levelStats).map(Number);
    return `+${Math.min(...levels)} / +${Math.max(...levels)}`;
  };

  return (
    <div className="mt-4">
      <div className="mb-4 flex space-x-4">
        <RegionSelector
          regions={["us", "eu", "kr", "tw", "cn"]}
          onRegionChange={setRegion}
          selectedRegion={region}
        />
        <DungeonSelector
          dungeons={dungeonData?.dungeons || []}
          onDungeonChange={setDungeon}
          selectedDungeon={dungeon}
        />
      </div>

      {currentDungeonStats && (
        <div className="p-4 bg-deep-blue rounded-lg shadow-2xl mb-4 mt-4">
          <h3 className="text-xl font-bold text-white mb-2">
            Mythic+ Keystones Range
          </h3>
          <p className="text-white text-lg">
            {getKeyRange(currentDungeonStats.LevelStats)}
          </p>
        </div>
      )}

      <nav className="mt-6 mb-6 border-b border-gray-700">
        <div className="flex space-x-4">
          <Link
            href="/mythic-plus/statistics"
            className="px-4 py-2 text-sm font-medium text-white hover:bg-deep-blue/50 transition-colors hover:bg-blue-200"
          >
            Overall Stats
          </Link>
          <Link
            href="/mythic-plus/statistics/spec"
            className="px-4 py-2 text-sm font-medium text-white hover:bg-deep-blue/50 transition-colors"
          >
            Spec Distribution
          </Link>
          <Link
            href="/mythic-plus/statistics/team"
            className="px-4 py-2 text-sm font-medium text-white hover:bg-deep-blue/50 transition-colors"
          >
            Team Compositions
          </Link>
        </div>
      </nav>
    </div>
  );
}

export default function StatisticsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="p-4 bg-[#0a0a0a] bg-opacity-80">
      <StatsProvider>
        <StatisticsHeader />
        <StatisticsNav />
        {children}
      </StatsProvider>
    </div>
  );
}
