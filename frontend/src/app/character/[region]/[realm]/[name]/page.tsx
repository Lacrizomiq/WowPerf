"use client";

import CharacterSummary from "@/components/Character/CharacterSummary";

export default function CharacterPage({
  params,
}: {
  params: { region: string; realm: string; name: string };
}) {
  const { region, realm, name } = params;

  return (
    <main className="bg-gradient-dark">
      <CharacterSummary region={region} realm={realm} name={name} />
    </main>
  );
}
