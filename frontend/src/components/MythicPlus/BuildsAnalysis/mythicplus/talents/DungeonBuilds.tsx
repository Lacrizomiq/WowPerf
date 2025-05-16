// DungeonBuilds.tsx - Version harmonisée
"use client";

import { useGetTalentBuildsByDungeon } from "@/hooks/useBuildsAnalysisApi";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import { groupTalentsByDungeon } from "@/utils/buildsAnalysis/dataTransformer";
import TalentBuildCard from "./TalentBuildCard";

interface DungeonBuildsProps {
  className: WowClassParam;
  spec: WowSpecParam;
  encounterId: string;
}

export default function DungeonBuilds({
  className,
  spec,
  encounterId,
}: DungeonBuildsProps) {
  // Fetch talent builds by dungeon data
  const {
    data: dungeonBuildsData,
    isLoading,
    error,
  } = useGetTalentBuildsByDungeon(className, spec);

  // Loading state
  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-600"></div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="bg-red-900/20 border border-red-500 rounded-md p-4 my-4">
        <h3 className="text-red-500 text-lg font-medium">
          Error loading dungeon builds
        </h3>
        <p className="text-slate-300">
          {error.message || "An unknown error occurred"}
        </p>
      </div>
    );
  }

  // No data state
  if (!dungeonBuildsData || dungeonBuildsData.length === 0) {
    return (
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 text-center">
        <p className="text-slate-400">
          No builds available for this dungeon yet.
        </p>
      </div>
    );
  }

  // Find builds that match the encounterId
  const matchingBuilds = dungeonBuildsData.filter(
    (build) => build.encounter_id.toString() === encounterId
  );

  // If no matching builds found
  if (matchingBuilds.length === 0) {
    return (
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 text-center">
        <p className="text-slate-400">
          No builds available for the selected dungeon.
        </p>
      </div>
    );
  }

  // Get the dungeon name from the first matching build
  const dungeonName = matchingBuilds[0].dungeon_name;

  // Sort builds by usage percentage (highest first)
  const sortedBuilds = [...matchingBuilds].sort(
    (a, b) => b.avg_usage_percentage - a.avg_usage_percentage
  );

  // Assign ranks to builds if there are multiple
  const rankLabels = ["1st", "2nd", "3rd"];

  return (
    <div>
      <div className="space-y-8">
        {sortedBuilds.slice(0, 3).map((build, index) => (
          <TalentBuildCard
            key={`${build.talent_import}-${index}`}
            build={build}
            className={className}
            spec={spec}
            rank={sortedBuilds.length > 1 ? rankLabels[index] : undefined}
            dungeonName={dungeonName}
          />
        ))}
      </div>
    </div>
  );
}
