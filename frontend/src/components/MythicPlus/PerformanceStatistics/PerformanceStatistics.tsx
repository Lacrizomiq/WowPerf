// components/MythicPlus/PerformanceStatistics/PerformanceStatistics.tsx
"use client";

import React, { useMemo, useState } from "react";
import SpecScoreCard from "./SpecScoreCard";
import { SpecAverageGlobalScore } from "@/types/warcraftlogs/globalLeaderboardAnalysis";
import { specMapping } from "@/utils/specmapping";
import { useGetSpecAverageGlobalScore } from "@/hooks/useWarcraftLogsApi";

const PerformanceStatistics: React.FC = () => {
  type RoleType = "TANK" | "HEALER" | "DPS" | "ALL";

  const [selectedRole, setSelectedRole] = useState<RoleType>("ALL");

  const { data, isLoading, error } = useGetSpecAverageGlobalScore();

  const specs = useMemo(() => {
    if (!data) return [] as SpecAverageGlobalScore[];
    return Array.isArray(data) ? data : [];
  }, [data]);

  const sortedSpecs = useMemo(() => {
    if (!specs || specs.length === 0) return [] as SpecAverageGlobalScore[];
    return [...specs].sort((a, b) => b.avg_global_score - a.avg_global_score);
  }, [specs]);

  const filteredSpecs = useMemo(() => {
    if (selectedRole === "ALL") {
      return sortedSpecs;
    }

    return sortedSpecs.filter((spec) => {
      if (!spec.class || !spec.spec) return false;
      const specInfo = specMapping[spec.class]?.[spec.spec];
      return specInfo && specInfo.role === selectedRole;
    });
  }, [sortedSpecs, selectedRole]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-white text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-white mx-auto"></div>
          <p className="mt-4">Loading spec performance data...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-white text-center">
          <div className="bg-red-500 p-4 rounded-md inline-block">
            <p>Error loading data. Please try again later.</p>
          </div>
        </div>
      </div>
    );
  }

  if (!specs || specs.length === 0) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-white text-center">
          <p>No specification data available.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between mb-6 mt-12">
        <h1 className="text-3xl font-bold text-white mb-4 md:mb-0">
          Class & Spec Performance
        </h1>

        <div className="flex space-x-2">
          <button
            className={`px-4 py-2 rounded-md text-sm ${
              selectedRole === "ALL"
                ? "bg-blue-600 text-white"
                : "bg-gray-800 text-gray-300"
            }`}
            onClick={() => setSelectedRole("ALL")}
          >
            All Roles
          </button>
          <button
            className={`px-4 py-2 rounded-md text-sm ${
              selectedRole === "TANK"
                ? "bg-blue-600 text-white"
                : "bg-gray-800 text-gray-300"
            }`}
            onClick={() => setSelectedRole("TANK")}
          >
            Tank
          </button>
          <button
            className={`px-4 py-2 rounded-md text-sm ${
              selectedRole === "HEALER"
                ? "bg-blue-600 text-white"
                : "bg-gray-800 text-gray-300"
            }`}
            onClick={() => setSelectedRole("HEALER")}
          >
            Healer
          </button>
          <button
            className={`px-4 py-2 rounded-md text-sm ${
              selectedRole === "DPS"
                ? "bg-blue-600 text-white"
                : "bg-gray-800 text-gray-300"
            }`}
            onClick={() => setSelectedRole("DPS")}
          >
            DPS
          </button>
        </div>
      </div>

      <p className="text-gray-400">
        Showing average global scores for the very best players of each spec.
      </p>
      <p className="text-gray-400 mb-8">
        Data is based on high-level Mythic+ runs from the current season.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-4">
        {filteredSpecs.map((spec, index) => (
          <SpecScoreCard
            key={`${spec.class}-${spec.spec}-${index}`}
            specData={spec}
          />
        ))}
      </div>
    </div>
  );
};

export default PerformanceStatistics;
