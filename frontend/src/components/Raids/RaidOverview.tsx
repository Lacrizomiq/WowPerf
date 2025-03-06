import React, { useState } from "react";
import {
  useGetBlizzardRaidsByExpansion,
  useGetBlizzardCharacterEncounterRaid,
} from "@/hooks/useBlizzardApi";
import { useGetPlayerRaidRankings } from "@/hooks/useWarcraftLogsApi";
import ExpansionSelector from "./ExpansionSelector";
import StaticRaidsList from "./StaticRaidsList";
import RaidDetails from "./RaidDetails";
import RaidsPlayerPerformance from "./CharacterPersonnalRanking/RaidsPlayerPerformance";
import { StaticRaid } from "@/types/raids";
import { RAID_ZONE_MAPPING } from "@/utils/s1_tww_mapping";

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

  // Add WarcraftLogs data fetch
  const { data: raidRankings, isLoading: isRankingsLoading } =
    useGetPlayerRaidRankings(characterName, realmSlug, region, 42);

  const handleExpansionChange = (newExpansion: string) => {
    setSelectedExpansion(newExpansion);
    setSelectedRaid(null);
  };

  const handleRaidSelect = (raid: StaticRaid) => {
    setSelectedRaid(raid);
  };

  if (isStaticLoading || isProgressionLoading) {
    return (
      <div className="text-white text-center p-4">Loading raid data...</div>
    );
  }

  return (
    <div className="p-6 rounded-xl shadow-lg m-4">
      {raidRankings && (
        <div className="mt-6">
          <RaidsPlayerPerformance playerData={raidRankings} />
        </div>
      )}
      <div className="flex justify-between items-center mb-6 mt-6">
        <h2 className="text-2xl font-bold text-white">Raids Progression</h2>
        <ExpansionSelector
          currentExpansion={selectedExpansion}
          onExpansionChange={handleExpansionChange}
        />
      </div>

      <StaticRaidsList
        raids={staticRaids || []}
        raidProgressionData={raidProgressionData}
        onRaidSelect={handleRaidSelect}
        selectedRaid={selectedRaid}
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
