"use client";

import Header from "@/components/Header/Header";
import CharacterSummary from "@/components/Character/CharacterSummary";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import CharacterTalent from "@/components/Character/CharacterTalent";
import CharacterMythicPlus from "@/components/Character/CharacterMythicPlus";
import CharacterGear from "@/components/Character/CharacterGear";
import { useState } from "react";

export default function CharacterLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { region: string; realm: string; name: string };
}) {
  const { region, realm, name } = params;
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
        return <CharacterTalent region={region} realm={realm} name={name} />;
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
      <nav className="flex space-x-4 bg-black p-4 text-white">
        <button onClick={() => setSelectedTab("gear")}>Gear</button>
        <button onClick={() => setSelectedTab("talents")}>Talents</button>
      </nav>

      <div>{renderContent()}</div>
    </main>
  );
}
