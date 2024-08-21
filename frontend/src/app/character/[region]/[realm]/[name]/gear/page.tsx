"use client";

import CharacterGear from "@/components/Character/CharacterGear";

export default function CharacterGearPage({ params }: any) {
  const { region, realm, name } = params;

  return (
    <CharacterGear
      region={region}
      realm={realm}
      name={name}
      namespace={`profile-${region}`}
      locale="en_GB"
    />
  );
}
