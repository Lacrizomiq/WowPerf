import React from "react";
import { StaticRaid, RaidProgressionData, Raid, RaidMode } from "@/types/raids";

interface DynamicRaidProgressProps {
  staticRaid: StaticRaid;
  raidProgressionData: RaidProgressionData | null;
  isLoading: boolean;
  error: any;
}

const DynamicRaidProgress: React.FC<DynamicRaidProgressProps> = ({
  staticRaid,
  raidProgressionData,
  isLoading,
  error,
}) => {
  if (isLoading) return <div>Loading raid progress...</div>;
  if (error) return <div>Error loading raid progress: {error.message}</div>;
  if (!raidProgressionData) return <div>No raid progress data available.</div>;

  const findRaidProgress = (): Raid | undefined => {
    return raidProgressionData.expansions
      .find((exp) => exp.name === staticRaid.Expansion)
      ?.raids.find((raid) => raid.id === staticRaid.ID);
  };

  const raidProgress = findRaidProgress();

  if (!raidProgress) {
    return <div>No progression data available for this raid.</div>;
  }

  return (
    <div className="mt-8">
      <h2 className="text-2xl font-bold mb-4">{staticRaid.Name} Progress</h2>
      {raidProgress.modes.map((mode: RaidMode) => (
        <div key={mode.difficulty} className="mb-4">
          <h3 className="text-xl font-semibold">{mode.difficulty}</h3>
          <p>Status: {mode.status}</p>
          <p>
            Progress: {mode.progress.completed_count}/
            {mode.progress.total_count}
          </p>
          <ul className="mt-2">
            {mode.progress.encounters.map((encounter) => (
              <li
                key={encounter.id}
                className="flex justify-between items-center"
              >
                <span>{encounter.name}</span>
                <span>Kills: {encounter.completed_count}</span>
              </li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
};

export default DynamicRaidProgress;
