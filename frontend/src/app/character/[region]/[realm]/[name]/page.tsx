"use client";

import CharacterSummary from "@/components/Character/CharacterSummary";
import CharacterGear from "@/components/Character/CharacterGear";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";

export default function CharacterPage({
  params,
}: {
  params: { region: string; realm: string; name: string };
}) {
  const { region, realm, name } = params;

  useWowheadTooltips();

  return (
    <main className="bg-gradient-dark">
      <CharacterSummary region={region} realm={realm} name={name} />
      <CharacterGear region={region} realm={realm} name={name} />
    </main>
  );
}
