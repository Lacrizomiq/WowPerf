"use client";

import Header from "@/components/Header/Header";
import CharacterSummary from "@/components/Character/CharacterSummary";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import CharacterTalent from "@/components/Character/CharacterTalent";
import CharacterGear from "@/components/Character/CharacterGear";
import MythicDungeonOverview from "@/components/MythicPlus/MythicOverview";
import { useState } from "react";
import RaidOverview from "@/components/Raids/RaidOverview";
export default function CharacterLayout({
  params,
}: {
  params: { region: string; realm: string; name: string; seasonSlug: string };
}) {
  const { region, realm, name, seasonSlug } = params;
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
        return <RaidOverview initialExpansion="DF" />;
      default:
        return null;
    }
  };

  return (
    <main className="bg-gradient-dark">
      <Header />
      <CharacterSummary
        region={region}
        realm={realm}
        name={name}
        namespace={`profile-${region}`}
        locale="en_GB"
      />
      <nav className="flex space-x-4 items-center justify-center font-bold  p-4 text-white">
        <button onClick={() => setSelectedTab("gear")}>Gear</button>
        <button onClick={() => setSelectedTab("talents")}>Talents</button>
        <button onClick={() => setSelectedTab("mythic-plus")}>Mythic+</button>
        <button onClick={() => setSelectedTab("raid-progression")}>
          Raids
        </button>
      </nav>

      <div>{renderContent()}</div>
    </main>
  );
}
