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
    return (
      <>
        <div className="flex justify-center p-5 mt-5 rounded-xl">
          <nav className="flex justify-center mt-5 bg-[#002440] overflow-hidden rounded-full border-2 border-[#003660]">
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
            ].map((tab, index, array) => (
              <button
                key={tab.key}
                onClick={() => setSelectedTab(tab.key)}
                className={`flex items-center space-x-2 px-6 py-3 transition-all bg-[#002440] justify-center
                ${
                  selectedTab === tab.key
                    ? "bg-[#003660]"
                    : "hover:bg-[#003660] hover:bg-opacity-50"
                }
                ${index === 0 ? "rounded-l-full" : ""}
                ${index === array.length - 1 ? "rounded-r-full" : ""}
                ${index !== 0 ? "border-l border-[#003660]" : ""}
      `}
              >
                {tab.icon}
                <span>{tab.name}</span>
              </button>
            ))}
          </nav>
        </div>
        <div className="rounded-xl mt-5">
          {(() => {
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
          })()}
        </div>
      </>
    );
  };

  const backgroundStyle = {
    backgroundImage: `linear-gradient(rgba(0, 0, 0, 0.4), rgba(0, 0, 0, 0.4)), url('https://wow.zamimg.com/images/tools/dragonflight-talent-calc/blizzard/talentbg-warlock-affliction.jpg')`,
    backgroundSize: "cover",
    backgroundPosition: "center",
    backgroundAttachment: "fixed",
  };

  return (
    <div className="min-h-screen p-1 bg-[#090909] text-white ">
      <div className="max-w-7xl mx-auto p-5" style={backgroundStyle}>
        <Header />

        <CharacterSummary
          region={region}
          realm={realm}
          name={name}
          namespace={`profile-${region}`}
          locale="en_GB"
        />
        <div className="rounded-xl mt-5 shadow-2xl">{renderContent()}</div>
      </div>
    </div>
  );
}
