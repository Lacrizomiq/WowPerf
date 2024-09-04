import React from "react";
import { MythicPlusRuns } from "@/types/mythicPlusRuns";

interface DungeonDetailsProps {
  run: MythicPlusRuns;
}

const DungeonDetails: React.FC<DungeonDetailsProps> = ({ run }) => {
  return (
    <div className="mt-8 p-6 bg-deep-blue rounded-lg">
      <h2 className="text-2xl font-bold mb-4">
        Detailed run for {run.Dungeon.Name}
      </h2>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <p>Level: {run.KeystoneLevel}</p>
          <p>Score: {run.MythicRating.toFixed(2)}</p>
          <p>Time: {(run.Duration / 1000 / 60).toFixed(2)} minutes</p>
          <p>Completed: {new Date(run.CompletedTimestamp).toLocaleString()}</p>
        </div>
        <div>
          <p>Affixes:</p>
          <ul>
            {run.Affixes.map((affix) => (
              <li key={affix.ID}>{affix.Name}</li>
            ))}
          </ul>
        </div>
      </div>
      <div className="mt-4">
        <h3 className="text-xl font-bold mb-2">Party Members:</h3>
        <ul>
          {run.Members.map((member) => (
            <li key={member.CharacterID}>
              {member.CharacterName} - {member.Specialization} {member.RaceName}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default DungeonDetails;
