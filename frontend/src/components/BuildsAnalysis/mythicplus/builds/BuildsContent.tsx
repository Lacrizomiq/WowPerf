// BuildsContent.tsx
"use client";

import { useState } from "react";
import {
  useGetStatPriorities,
  useGetTopTalentBuilds,
} from "@/hooks/useBuildsAnalysisApi";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import StatPriorities from "./StatPriorities";
import TopBuilds from "./TopBuilds";

interface BuildsContentProps {
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function BuildsContent({ className, spec }: BuildsContentProps) {
  // Fetch data using hooks
  const {
    data: statData,
    isLoading: isLoadingStats,
    error: statsError,
  } = useGetStatPriorities(className, spec);

  const {
    data: buildsData,
    isLoading: isLoadingBuilds,
    error: buildsError,
  } = useGetTopTalentBuilds(className, spec);

  // Loading states
  if (isLoadingStats || isLoadingBuilds) {
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-600"></div>
      </div>
    );
  }

  // Error states
  if (statsError || buildsError) {
    return (
      <div className="bg-red-900/20 border border-red-500 rounded-md p-4 my-4">
        <h3 className="text-red-500 text-lg font-medium">Error loading data</h3>
        <p className="text-slate-300">
          {statsError?.message ||
            buildsError?.message ||
            "An unknown error occurred"}
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Stat Priorities Section */}
      <div>{statData && <StatPriorities stats={statData} />}</div>
      {/* Top Build Section */}
      <div className="mb-10">
        {buildsData && buildsData.length > 0 && (
          <TopBuilds builds={buildsData} className={className} spec={spec} />
        )}
      </div>
    </div>
  );
}
