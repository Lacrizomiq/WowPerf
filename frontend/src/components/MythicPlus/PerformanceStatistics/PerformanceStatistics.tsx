// components/MythicPlus/PerformanceStatistics/PerformanceStatistics.tsx
"use client";

import React, { useMemo, useState } from "react";
import SpecScoreCard from "./SpecScoreCard";
import ClassSelector from "./ClassSelector";
import { SpecAverageGlobalScore } from "@/types/warcraftlogs/globalLeaderboardAnalysis";
import { RefreshCw } from "lucide-react";
import { useGetSpecAverageGlobalScore } from "@/hooks/useWarcraftLogsApi";
import { spec } from "node:test/reporters";

const PerformanceStatistics: React.FC = () => {
  type RoleType = "Tank" | "Healer" | "DPS" | "ALL";

  const [selectedRole, setSelectedRole] = useState<RoleType>("ALL");
  const [selectedClass, setSelectedClass] = useState<string | null>(null);

  const { data, isLoading, error } = useGetSpecAverageGlobalScore();

  const specs = useMemo(() => {
    if (!data) return [] as SpecAverageGlobalScore[];
    return Array.isArray(data) ? data : [];
  }, [data]);

  // Extract all unique class names from the specs data
  const availableClasses = useMemo(() => {
    if (!specs || specs.length === 0) return [] as string[];
    const classSet = new Set(specs.map((spec) => spec.class));
    return Array.from(classSet);
  }, [specs]);

  const sortedSpecs = useMemo(() => {
    if (!specs || specs.length === 0) return [] as SpecAverageGlobalScore[];
    return [...specs].sort((a, b) => b.avg_global_score - a.avg_global_score);
  }, [specs]);

  const filteredSpecs = useMemo(() => {
    let filtered = sortedSpecs;

    // Filter by role if not ALL
    if (selectedRole !== "ALL") {
      filtered = filtered.filter((spec) => spec.role === selectedRole);
    }

    // Filter by class if one is selected
    if (selectedClass) {
      filtered = filtered.filter((spec) => spec.class === selectedClass);
    }

    return filtered;
  }, [sortedSpecs, selectedRole, selectedClass]);

  // Reset all filters to default values
  const resetFilters = () => {
    setSelectedRole("ALL");
    setSelectedClass(null);
  };

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

  const isFiltered = selectedRole !== "ALL" || selectedClass !== null;

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between mb-2 mt-12">
        <h1 className="text-3xl font-bold text-white mb-6 md:mb-0">
          Class & Spec Performance
        </h1>

        <div className="flex flex-col space-y-4">
          {/* Role filters */}
          <div className="flex space-x-2 justify-end">
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
                selectedRole === "Tank"
                  ? "bg-blue-600 text-white"
                  : "bg-gray-800 text-gray-300"
              }`}
              onClick={() => setSelectedRole("Tank")}
            >
              Tank
            </button>
            <button
              className={`px-4 py-2 rounded-md text-sm ${
                selectedRole === "Healer"
                  ? "bg-blue-600 text-white"
                  : "bg-gray-800 text-gray-300"
              }`}
              onClick={() => setSelectedRole("Healer")}
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

          {/* Class selector and reset button */}
          <div className="flex space-x-2 items-center justify-end">
            <ClassSelector
              selectedClass={selectedClass}
              onClassChange={setSelectedClass}
              availableClasses={availableClasses}
            />

            <button
              onClick={resetFilters}
              className={`flex items-center justify-center px-4 py-2 rounded-md text-sm ${
                isFiltered
                  ? "bg-gray-700 text-white hover:bg-gray-600"
                  : "bg-gray-800 text-gray-500 cursor-not-allowed"
              }`}
              disabled={!isFiltered}
              title={isFiltered ? "Reset all filters" : "No filters applied"}
            >
              <RefreshCw className="w-4 h-4 mr-2" />
              Reset
            </button>
          </div>
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
            selectedRole={selectedRole}
          />
        ))}
      </div>

      {filteredSpecs.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-400">No specs match the selected filters.</p>
          <button
            onClick={resetFilters}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            Reset Filters
          </button>
        </div>
      )}
    </div>
  );
};

export default PerformanceStatistics;
