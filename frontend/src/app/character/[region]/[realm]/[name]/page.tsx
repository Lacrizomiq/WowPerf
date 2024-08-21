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
  params: {
    region: string;
    realm: string;
    name: string;
    namespace: string;
    locale: string;
  };
}) {
  const { region, realm, name, namespace, locale } = params;

  useWowheadTooltips();

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
      <CharacterGear
        region={region}
        realm={realm}
        name={name}
        namespace={`profile-${region}`}
        locale="en_GB"
      />
      <CharacterTalent region={region} realm={realm} name={name} />
      <CharacterMythicPlus region={region} realm={realm} name={name} />
    </main>
  );
}
