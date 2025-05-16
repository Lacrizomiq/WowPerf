// components/performance/mythicplus/MythicPlusContent.tsx
"use client";

import React, { useState, useMemo } from "react";
import FilterSection from "./FilterSection";
import SpecCard from "./SpecCard";
import { useGetSpecAverageGlobalScore } from "@/hooks/useWarcraftLogsApi";
import { SpecAverageGlobalScore } from "@/types/warcraftlogs/globalLeaderboardAnalysis";

export default function MythicPlusContent() {
  type RoleType = "Tank" | "Healer" | "DPS" | "ALL";

  const [selectedRole, setSelectedRole] = useState<RoleType>("ALL");
  const [selectedClass, setSelectedClass] = useState<string | null>(null);
  const [selectedDungeon, setSelectedDungeon] = useState<string>("all");

  const { data, isLoading, error } = useGetSpecAverageGlobalScore();

  const specs = useMemo(() => {
    if (!data) return [] as SpecAverageGlobalScore[];
    return Array.isArray(data) ? data : [];
  }, [data]);

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

    // TODO : Implement this to filter the score per dungeon instead of global when the user select a dungeon
    // if (selectedDungeon !== "all") {
    // filtered = filtered.filter(
    // (spec) => spec.performanceByDungeon?.[selectedDungeon]
    // );
    // }

    return filtered;
  }, [sortedSpecs, selectedRole, selectedClass]);

  // Gestionnaires d'événements
  const handleRoleChange = (role: RoleType) => {
    setSelectedRole(role);
  };

  const handleClassChange = (className: string | null) => {
    setSelectedClass(className);
  };

  const handleDungeonChange = (dungeon: string) => {
    setSelectedDungeon(dungeon);
  };

  const resetFilters = () => {
    setSelectedRole("ALL");
    setSelectedClass(null);
    setSelectedDungeon("all");
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

  const isFiltered =
    selectedRole !== "ALL" ||
    selectedClass !== null ||
    selectedDungeon !== "all";

  return (
    <>
      {/* Filter Section */}
      <FilterSection
        selectedRole={selectedRole}
        selectedClass={selectedClass}
        selectedDungeon={selectedDungeon}
        availableClasses={availableClasses}
        onRoleChange={handleRoleChange}
        onClassChange={handleClassChange}
        onDungeonChange={handleDungeonChange}
        onResetFilters={resetFilters}
        isFiltered={isFiltered}
      />

      {/* Data Description */}
      <div className="mt-6">
        <p className="text-gray-400">
          Showing average global scores for the very best players of each spec.
        </p>
        <p className="text-gray-400">
          Data is based on high-level Mythic+ runs from the current season.
        </p>
        <p className="text-gray-400 mb-4">
          Note: Chinese players are not included in the data.
        </p>
      </div>

      {/* Specialization Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
        {filteredSpecs.map((spec, index) => (
          <SpecCard
            key={`${spec.class}-${spec.spec}-${index}`}
            specData={spec}
            selectedRole={selectedRole}
          />
        ))}
      </div>

      {/* Message when no results */}
      {filteredSpecs.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-400">No specs match the selected filters.</p>
          <button
            onClick={resetFilters}
            className="mt-4 px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700"
          >
            Reset Filters
          </button>
        </div>
      )}
    </>
  );
}
