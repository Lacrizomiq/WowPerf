import React from "react";
import { StaticRaid, RaidProgressionData, RaidMode } from "@/types/raids";

interface RaidDetailsProps {
  staticRaid: StaticRaid;
  raidProgressionData: RaidProgressionData | null;
}

const RaidDetails: React.FC<RaidDetailsProps> = ({
  staticRaid,
  raidProgressionData,
}) => {
  const raidProgress = raidProgressionData?.expansions
    .flatMap((exp) => exp.raids)
    .find((raid) => raid.id === staticRaid.ID);

  if (!raidProgress) {
    return <div>No progression data available for this raid.</div>;
  }

  return (
    <div className="bg-deep-blue p-4 rounded-lg mt-8">
      <h2 className="text-2xl font-bold mb-4">{staticRaid.Name}</h2>
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

export default RaidDetails;
