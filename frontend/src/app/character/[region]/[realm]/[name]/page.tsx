"use client";

import React, { useState } from "react";
import Header from "@/components/Header/Header";
import CharacterSummary from "@/components/Character/CharacterSummary";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import CharacterTalent from "@/components/Character/CharacterTalent";
import CharacterGear from "@/components/Character/CharacterGear";
import MythicDungeonOverview from "@/components/MythicPlus/MythicOverview";
import RaidOverview from "@/components/Raids/RaidOverview";
import { Shield, ScrollText, Sword, Hourglass } from "lucide-react";

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
        <div className="bg-[#002440] flex justify-center p-5 mt-5 rounded-t-xl">
          <nav className="flex justify-center  bg-[#002440] border-2 border-[#003660] w-fit">
            {[
              { name: "Gear", icon: <Shield size={20} />, key: "gear" },
              {
                name: "Talents",
                icon: <ScrollText size={20} />,
                key: "talents",
              },
              {
                name: "Mythic+",
                icon: <Hourglass size={20} />,
                key: "mythic-plus",
              },
              {
                name: "Raids",
                icon: <Sword size={20} />,
                key: "raid-progression",
              },
            ].map((tab) => (
              <button
                key={tab.key}
                onClick={() => setSelectedTab(tab.key)}
                className={`flex items-center space-x-2 px-6 py-3  transition-all bg-[#002440] justify-center border-2 border-[#003660]
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
        </div>
        <div className="bg-[#002440] rounded-b-xl p-5">{renderContent()}</div>
      </div>
    </div>
  );
}
