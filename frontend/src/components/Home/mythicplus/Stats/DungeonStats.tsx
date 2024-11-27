"use client";

import React, { useState, useEffect } from "react";
import { useGetDungeonStats } from "@/hooks/useRaiderioApi";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import { DungeonStat } from "@/types/dungeonStats";
import { Dungeon } from "@/types/mythicPlusRuns";
import DungeonSelector from "./Selector/DungeonSelector";
import RegionSelector from "./Selector/RegionSelector";
import { OverallStats } from "./OverallStats";
import { SpecStats } from "./SpecStats";
import { TeamComposition } from "./TeamComposition";

const DungeonStats: React.FC = () => {
  const [season] = useState("season-tww-1");
  const [region, setRegion] = useState("world");
  const [dungeon, setDungeon] = useState("all");
  const [activeTab, setActiveTab] = useState("overall");

  const {
    data: statsData,
    isLoading,
    error,
  } = useGetDungeonStats(season, region);

  const { data: dungeonData } = useGetBlizzardMythicDungeonPerSeason(season);
  const [dungeons, setDungeons] = useState<Dungeon[]>([]);

  useEffect(() => {
    if (dungeonData?.dungeons) {
      setDungeons(dungeonData.dungeons);
    }
  }, [dungeonData]);

  if (isLoading) return <div className="text-white">Loading stats...</div>;
  if (error)
    return (
      <div className="text-red-500">Error loading stats: {error.message}</div>
    );

  const currentDungeonStats =
    statsData?.find((stat: DungeonStat) => stat.dungeon_slug === dungeon) ||
    statsData?.[0];

  if (!currentDungeonStats) {
    return (
      <div className="text-white">No data available for this dungeon.</div>
    );
  }

  const tabs = [
    { id: "overall", label: "Overall Stats" },
    { id: "specs", label: "Spec Distribution" },
    { id: "compositions", label: "Team Compositions" },
  ];

  const getKeyRange = (levelStats: Record<string, number>) => {
    const levels = Object.keys(levelStats).map(Number);
    return `+${Math.min(...levels)} / +${Math.max(...levels)}`;
  };

  const renderContent = () => {
    switch (activeTab) {
      case "overall":
        return <OverallStats stats={currentDungeonStats} />;
      case "specs":
        return <SpecStats stats={currentDungeonStats} />;
      case "compositions":
        return <TeamComposition stats={currentDungeonStats} />;
      default:
        return null;
    }
  };

  return (
    <div className="p-4 bg-[#0a0a0a] bg-opacity-80">
      <h2 className="text-2xl font-bold text-white mb-4">
        Dungeon Statistics for{" "}
        {dungeon
          .split("-")
          .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(" ")}{" "}
        dungeon in region: {region.toUpperCase()}
      </h2>

      <div className="mb-4 flex space-x-4">
        <RegionSelector
          regions={["us", "eu", "kr", "tw", "cn"]}
          onRegionChange={setRegion}
          selectedRegion={region}
        />
        <DungeonSelector
          dungeons={dungeons}
          onDungeonChange={setDungeon}
          selectedDungeon={dungeon}
        />
      </div>

      <div className="space-y-4">
        <div className="p-4">
          <p className="text-white">
            <span className="font-bold">Last update:</span>{" "}
            {new Intl.DateTimeFormat("en-US", {
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
            the class you love !
          </p>
        </div>

        <div className="p-4 bg-deep-blue rounded-lg shadow-2xl">
          <h3 className="text-xl font-bold text-white mb-2">
            Mythic+ Keystones Range
          </h3>
          <p className="text-white text-lg">
            {getKeyRange(currentDungeonStats.LevelStats)}
          </p>
        </div>
      </div>

      <div className="mb-6 mt-6">
        <div className="flex space-x-2 border-b border-gray-700">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-2 text-sm font-medium rounded-t-lg transition-colors ${
                activeTab === tab.id
                  ? "text-white bg-deep-blue border-b-2 border-blue-500"
                  : "text-gray-400 hover:text-white hover:bg-deep-blue/50"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      <div className="mt-6">{renderContent()}</div>
    </div>
  );
};

export default DungeonStats;
