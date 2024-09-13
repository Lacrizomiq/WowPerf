"use client";

import React, { useState } from "react";
import Header from "@/components/Header/Header";
import CharacterSummary from "@/components/Character/CharacterSummary";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import CharacterTalent from "@/components/Character/CharacterTalent";
import CharacterGear from "@/components/Character/CharacterGear";
import MythicDungeonOverview from "@/components/MythicPlus/MythicOverview";
import RaidOverview from "@/components/Raids/RaidOverview";
import { Shield, Book, Sword, Activity } from "lucide-react";

export default function CharacterLayout({
  params,
}: {
  params: {
    region: string;
    realm: string;
    name: string;
    seasonSlug: string;
    expansion?: string;
  };
}) {
  const { region, realm, name, seasonSlug, expansion } = params;
  const [selectedTab, setSelectedTab] = useState<string>("gear");

  useWowheadTooltips();

  const renderContent = () => {
    switch (selectedTab) {
      case "gear":
        return (
          <CharacterGear
            region={region}
            realm={realm}
            name={name}
            namespace={`profile-${region}`}
            locale="en_GB"
          />
        );
      case "talents":
        return (
          <CharacterTalent
            region={region}
            realm={realm}
            name={name}
            namespace={`profile-${region}`}
            locale="en_GB"
          />
        );
      case "mythic-plus":
        return (
          <MythicDungeonOverview
            characterName={name}
            realmSlug={realm}
            region={region}
            namespace={`profile-${region}`}
            locale="en_GB"
            seasonSlug={seasonSlug || "season-tww-1"}
          />
        );
      case "raid-progression":
        return (
          <RaidOverview
            characterName={name}
            realmSlug={realm}
            region={region}
            namespace={`profile-${region}`}
            locale="en_GB"
            expansion={expansion || "TWW"}
          />
        );
      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen p-5 bg-[#000c1a] text-white">
      <div className="max-w-7xl mx-auto">
        <Header />
        <CharacterSummary
          region={region}
          realm={realm}
          name={name}
          namespace={`profile-${region}`}
          locale="en_GB"
        />
        <nav className="flex justify-center mb-5 space-x-4">
          {[
            { name: "Gear", icon: <Shield size={20} />, key: "gear" },
            { name: "Talents", icon: <Book size={20} />, key: "talents" },
            { name: "Mythic+", icon: <Sword size={20} />, key: "mythic-plus" },
            {
              name: "Raids",
              icon: <Activity size={20} />,
              key: "raid-progression",
            },
          ].map((tab) => (
            <button
              key={tab.key}
              onClick={() => setSelectedTab(tab.key)}
              className={`flex items-center space-x-2 px-6 py-3 rounded-lg transition-all bg-[#002440] justify-center
                ${
                  selectedTab === tab.key
                    ? "bg-[#003660]"
                    : "hover:bg-[#003660] hover:bg-opacity-50"
                }`}
            >
              {tab.icon}
              <span>{tab.name}</span>
            </button>
          ))}
        </nav>
        <div>{renderContent()}</div>
      </div>
    </div>
  );
}
