"use client";

import React, { useState, useEffect } from "react";
import Header from "@/components/Header/Header";
import CharacterSummary from "@/components/Character/CharacterSummary";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import CharacterTalent from "@/components/Character/CharacterTalent";
import CharacterGear from "@/components/Character/CharacterGear";
import MythicDungeonOverview from "@/components/MythicPlus/MythicOverview";
import RaidOverview from "@/components/Raids/RaidOverview";
import { Shield, ScrollText, Sword, Hourglass } from "lucide-react";
import { useGetBlizzardCharacterProfile } from "@/hooks/useBlizzardApi";
import "@/app/globals.css";
import Sidebar from "@/components/Header/Sidebar";
export default function CharacterLayout({
  params,
}: {
  params: {
    region: string;
    realm: string;
    name: string;
    seasonSlug: string;
    expansion?: string;
    namespace: string;
    locale: string;
  };
}) {
  const { region, realm, name, seasonSlug, expansion } = params;
  const [selectedTab, setSelectedTab] = useState<string>("gear");
  const [mainMargin, setMainMargin] = useState(64);
  const {
    data: characterProfile,
    isLoading,
    error,
  } = useGetBlizzardCharacterProfile(
    region,
    realm,
    name,
    `profile-${region}`,
    "en_GB"
  );

  useWowheadTooltips();

  console.log("Character Profile Loading:", isLoading);
  console.log("Character Profile Error:", error);
  console.log("Character Profile Data:", characterProfile);

  const renderContent = () => {
    return (
      <>
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
    backgroundSize: "cover",
    backgroundPosition: "center",
    backgroundAttachment: "fixed",
  };

  const defaultBackgroundClass = "bg-deep-blue";

  if (isLoading) {
    return (
      <div className="min-h-screen p-1 bg-[#090909] text-white">Loading...</div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen p-1 bg-[#090909] text-white">
        Error: {error.message}
      </div>
    );
  }

  return (
    <div className="flex min-h-screen bg-[#0a0a0a] text-white">
      <Sidebar setMainMargin={setMainMargin} />
      <div
        className="flex-1 transition-all duration-300"
        style={{ marginLeft: `${mainMargin}px` }}
      >
        <div className="max-w-7xl mx-auto p-5">
          <CharacterSummary
            region={region}
            realm={realm}
            name={name}
            namespace={`profile-${region}`}
            locale="en_GB"
          />

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
                ${index !== 0 ? "border-l border-[#003660]" : ""}`}
                >
                  {tab.icon}
                  <span>{tab.name}</span>
                </button>
              ))}
            </nav>
          </div>
          <div
            className={`max-w-7xl mx-auto rounded-2xl shadow-2xl ${
              characterProfile?.spec_id
                ? `bg-spec-${characterProfile.spec_id}`
                : defaultBackgroundClass
            }`}
            style={backgroundStyle}
          >
            {renderContent()}
          </div>
        </div>
      </div>
    </div>
  );
}
