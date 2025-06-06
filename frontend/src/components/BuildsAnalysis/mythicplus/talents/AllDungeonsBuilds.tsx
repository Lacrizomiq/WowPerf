// AllDungeonsBuilds.tsx - Version harmonisée
"use client";

import { useGetTopTalentBuilds } from "@/hooks/useBuildsAnalysisApi";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import TalentBuildCard from "./TalentBuildCard";

interface AllDungeonsBuildsProps {
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function AllDungeonsBuilds({
  className,
  spec,
}: AllDungeonsBuildsProps) {
  // Fetch top talent builds data
  const {
    data: buildsData,
    isLoading,
    error,
  } = useGetTopTalentBuilds(className, spec);

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
          Error loading builds
        </h3>
        <p className="text-slate-300">
          {error.message || "An unknown error occurred"}
        </p>
      </div>
    );
  }

  // No data state
  if (!buildsData || buildsData.length === 0) {
    return (
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 text-center">
        <p className="text-slate-400">No builds available yet.</p>
      </div>
    );
  }

  // Assign ranks to builds
  const rankLabels = ["1st", "2nd", "3rd"];

  return (
    <div>
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-white mb-4">
          Most Popular Talents Builds Across All Dungeons
        </h2>
        <p className="text-slate-400">
          These builds represent the most frequently used talent combinations
          across all Mythic+ dungeons for {className} {spec}.
        </p>
      </div>

      <div className="space-y-8">
        {buildsData.slice(0, 3).map((build, index) => (
          <TalentBuildCard
            key={build.talent_import}
            build={build}
            className={className}
            spec={spec}
            rank={rankLabels[index]}
          />
        ))}
      </div>
    </div>
  );
}
