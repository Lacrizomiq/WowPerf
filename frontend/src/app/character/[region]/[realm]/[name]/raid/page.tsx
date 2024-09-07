"use client";

import RaidOverview from "@/components/Raids/RaidOverview";

export default function RaidPage({
  params,
}: {
  params: { region: string; realm: string; name: string };
}) {
  const { region, realm, name } = params;
  const initialExpansion = "DF"; // You can change this to the desired initial expansion

  return (
    <RaidOverview
      characterName={name}
      realmSlug={realm}
      region={region}
      namespace={`profile-${region}`}
      locale="en_GB"
      initialExpansion={initialExpansion}
    />
  );
}
