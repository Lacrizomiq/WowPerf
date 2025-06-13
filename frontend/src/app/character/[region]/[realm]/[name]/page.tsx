"use client";

import React, { useState, use } from "react";
import CharacterSummary from "@/components/Character/CharacterSummary";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import CharacterTalent from "@/components/Character/CharacterTalent";
import CharacterGear from "@/components/Character/CharacterGear";
import MythicDungeonOverview from "@/components/MythicPlus/MythicOverview";
import RaidOverview from "@/components/Raids/RaidOverview";
import { Shield, ScrollText, Sword, Hourglass } from "lucide-react";
import { useGetBlizzardCharacterProfile } from "@/hooks/useBlizzardApi";
import "@/app/globals.css";

export default function CharacterPage({
  params,
}: {
  params: Promise<{
    region: string;
    realm: string;
    name: string;
    seasonSlug: string;
    expansion?: string;
    namespace: string;
    locale: string;
  }>;
}) {
  // Utilisation deuse() pour résoudre la Promise dans un Client Component
  const resolvedParams = use(params);
  const { region, realm, name, seasonSlug, expansion } = resolvedParams;

  const [selectedTab, setSelectedTab] = useState<string>("gear");
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

  const renderContent = () => (
    <div className="rounded-xl w-full">
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
                seasonSlug={seasonSlug || "season-tww-2"}
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
  );

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
    <div className="min-h-screen bg-[#0a0a0a] text-white">
      <div className="w-full max-w-7xl mx-auto px-2 sm:px-4 md:px-6">
        <div className="mb-4">
          <CharacterSummary
            region={region}
            realm={realm}
            name={name}
            namespace={`profile-${region}`}
            locale="en_GB"
          />
        </div>

        {/* Navigation */}
        <div className="flex justify-center w-full mb-4">
          <nav className="flex flex-col md:flex-row md:inline-flex justify-center bg-[#002440] rounded-2xl md:rounded-full border-2 border-[#003660] overflow-hidden">
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
                className={`
          flex items-center justify-center gap-2 
          px-6 py-3
          transition-all bg-[#002440]
          ${
            selectedTab === tab.key
              ? "bg-[#003660]"
              : "hover:bg-[#003660] hover:bg-opacity-50"
          }
          ${
            index !== array.length - 1
              ? "md:border-r border-b md:border-b-0 border-[#003660]"
              : ""
          }
          ${index === 0 ? "md:rounded-l-full" : ""}
          ${index === array.length - 1 ? "md:rounded-r-full" : ""}
        `}
              >
                {tab.icon}
                <span className="whitespace-nowrap">{tab.name}</span>
              </button>
            ))}
          </nav>
        </div>

        <div className="w-full">{renderContent()}</div>
      </div>
    </div>
  );
}
