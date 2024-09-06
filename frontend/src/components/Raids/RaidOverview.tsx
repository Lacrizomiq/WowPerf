import React, { useState } from "react";
import { Raid } from "@/types/raids";
import { useGetBlizzardRaidsByExpansion } from "@/hooks/useBlizzardApi";
import StaticRaidsList from "./StaticRaidsList";
import ExpansionSelector from "./ExpansionSelector";

interface RaidOverviewProps {
  initialExpansion: string;
}

const RaidOverview: React.FC<RaidOverviewProps> = ({ initialExpansion }) => {
  const {
    data: raids,
    isLoading,
    error,
  } = useGetBlizzardRaidsByExpansion(initialExpansion);

  const [selectedExpansion, setSelectedExpansion] = useState(initialExpansion);

  const handleExpansionChange = (expansion: string) => {
    setSelectedExpansion(expansion);
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <ExpansionSelector
        raids={raids || []}
        onExpansionChange={handleExpansionChange}
        selectedExpansion={selectedExpansion}
      />
      <StaticRaidsList raids={raids || []} />
    </div>
  );
};

export default RaidOverview;
