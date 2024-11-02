"use client";

import React, { useState } from "react";
import { RoleLeaderboards } from "./RoleLeaderboard/RoleLeaderboards";
import DungeonLeaderboard from "./DungeonsLeaderboard/DungeonLeaderboard";
const LeaderboardTabs = () => {
  const [activeTab, setActiveTab] = useState("global");

  const tabs = [
    { id: "global", label: "Global Rankings" },
    { id: "dungeons", label: "Dungeon Rankings" },
  ];

  const renderContent = () => {
    switch (activeTab) {
      case "global":
        return <RoleLeaderboards />;
      case "dungeons":
        return <DungeonLeaderboard />;
      default:
        return null;
    }
  };

  return (
    <div className="mx-auto px-8 py-6 bg-black w-full h-full">
      <div className="mb-6">
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

      <div>{renderContent()}</div>
    </div>
  );
};

export default LeaderboardTabs;
