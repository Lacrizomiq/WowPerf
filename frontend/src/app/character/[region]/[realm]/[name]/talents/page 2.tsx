"use client";

import CharacterTalent from "@/components/Character/CharacterTalent";

export default function CharacterTalentPage({ params }: any) {
  const { region, realm, name } = params;

  return (
    <CharacterTalent
      region={region}
      realm={realm}
      name={name}
      namespace={`profile-${region}`}
      locale="en_GB"
    />
  );
}
