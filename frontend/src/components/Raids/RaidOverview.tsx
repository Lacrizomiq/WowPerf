import React, { useState, useEffect } from "react";
import {
  useGetBlizzardRaidsByExpansion,
  useGetBlizzardCharacterEncounterRaid,
} from "@/hooks/useBlizzardApi";
import ExpansionSelector from "./ExpansionSelector";
import StaticRaidsList from "./StaticRaidsList";
import RaidDetails from "./RaidDetails";
import { StaticRaid, RaidProgressionData } from "@/types/raids";

interface RaidOverviewProps {
  characterName: string;
  realmSlug: string;
  region: string;
  namespace: string;
  locale: string;
  expansion: string;
}

const RaidOverview: React.FC<RaidOverviewProps> = ({
  characterName,
  realmSlug,
  region,
  namespace,
  locale,
  expansion,
}) => {
  const [selectedExpansion, setSelectedExpansion] = useState(expansion);
  const [selectedRaid, setSelectedRaid] = useState<StaticRaid | null>(null);

  const { data: staticRaids, isLoading: isStaticLoading } =
    useGetBlizzardRaidsByExpansion(selectedExpansion);
  const { data: raidProgressionData, isLoading: isProgressionLoading } =
    useGetBlizzardCharacterEncounterRaid(
      region,
      realmSlug,
      characterName,
      namespace,
      locale
    );

  const handleExpansionChange = (newExpansion: string) => {
    setSelectedExpansion(newExpansion);
    setSelectedRaid(null);
  };

  const handleRaidSelect = (raid: StaticRaid) => {
    setSelectedRaid(raid);
  };

  if (isStaticLoading || isProgressionLoading) {
    return <div>Loading raid data...</div>;
  }

  return (
    <div className="p-6 bg-gradient-dark shadow-lg rounded-lg glow-effect m-12 max-w-6xl mx-auto">
      <ExpansionSelector
        currentExpansion={selectedExpansion}
        onExpansionChange={handleExpansionChange}
      />
      <StaticRaidsList
        raids={staticRaids || []}
        raidProgressionData={raidProgressionData}
        onRaidSelect={handleRaidSelect}
      />
      {selectedRaid && (
        <RaidDetails
          staticRaid={selectedRaid}
          raidProgressionData={raidProgressionData}
        />
      )}
    </div>
  );
};

export default RaidOverview;