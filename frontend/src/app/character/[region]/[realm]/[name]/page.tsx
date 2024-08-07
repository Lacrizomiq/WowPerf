"use client";

import Header from "@/components/Header/Header";
import CharacterSummary from "@/components/Character/CharacterSummary";
import CharacterGear from "@/components/Character/CharacterGear";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import CharacterTalent from "@/components/Character/CharacterTalent";
import CharacterMythicPlus from "@/components/Character/CharacterMythicPlus";

export default function CharacterPage({
  params,
}: {
  params: { region: string; realm: string; name: string };
}) {
  const { region, realm, name } = params;

  useWowheadTooltips();

  return (
    <main className="bg-gradient-dark">
      <Header />
      <CharacterSummary region={region} realm={realm} name={name} />
      <CharacterGear region={region} realm={realm} name={name} />
      <CharacterTalent region={region} realm={realm} name={name} />
      <CharacterMythicPlus region={region} realm={realm} name={name} />
    </main>
  );
}
