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
    return (
      <div className="bg-deep-blue p-4 rounded-lg mt-8">
        <h2 className="text-2xl font-bold mb-4 text-white">
          {staticRaid.Name}
        </h2>
        <p>No progression data available for this raid.</p>
      </div>
    );
  }

  return (
    <div className="bg-deep-blue p-4 rounded-lg mt-8">
      <h2 className="text-2xl font-bold mb-4 text-white">{staticRaid.Name}</h2>
      {raidProgress.modes.map((mode: RaidMode) => (
        <div key={mode.difficulty} className="mb-4">
          <div className="flex items-center text-center justify-between">
            <h3 className="text-xl font-bold pb-2">{mode.difficulty}</h3>
            <p>
              Progress: {mode.progress.completed_count}/
              {mode.progress.total_count}
            </p>
          </div>
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
