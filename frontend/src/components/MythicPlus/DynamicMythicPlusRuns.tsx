import React from "react";
import { useGetBlizzardCharacterMythicPlusBestRuns } from "@/hooks/useBlizzardApi";
import { MythicPlusRuns } from "@/types/mythicPlusRuns";

interface DynamicMythicPlusRunsProps {
  characterName: string;
  realmSlug: string;
  region: string;
  namespace: string;
  locale: string;
  seasonId: string;
}

const DynamicMythicPlusRuns: React.FC<DynamicMythicPlusRunsProps> = ({
  characterName,
  realmSlug,
  region,
  namespace,
  locale,
  seasonId,
}) => {
  const {
    data: mythicPlusRuns,
    isLoading,
    error,
  } = useGetBlizzardCharacterMythicPlusBestRuns(
    region,
    realmSlug,
    characterName,
    namespace,
    locale,
    seasonId
  );

  console.log("Mythic+ Runs Data:", mythicPlusRuns);

  if (isLoading) return <div>Loading Mythic+ runs...</div>;
  if (error) return <div>Error loading Mythic+ runs</div>;
  if (!mythicPlusRuns)
    return <div>No Mythic+ data available for this season</div>;
  if (mythicPlusRuns.length === 0)
    return <div>No Mythic+ runs found for this season</div>;

  return (
    <div className="w-full mt-6">
      <h2 className="text-2xl font-bold mb-4">Best Mythic+ Runs</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
        {mythicPlusRuns.map((run: MythicPlusRuns) => (
          <div
            key={`${run.Dungeon.ID}-${run.CompletedTimestamp}`}
            className="bg-deep-blue p-4 rounded-lg"
          >
            <h3 className="font-bold">{run.Dungeon.Name}</h3>
            <p>Level: {run.KeystoneLevel}</p>
            <p>Score: {run.MythicRating.toFixed(2)}</p>
            <p>Time: {(run.Duration / 1000 / 60).toFixed(2)} minutes</p>
            <p>
              Completed: {new Date(run.CompletedTimestamp).toLocaleString()}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default DynamicMythicPlusRuns;
