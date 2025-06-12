"use client";

import React, { useMemo } from "react";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import {
  useGetSpecAverageGlobalScore,
  useGetBestTenPlayerPerSpec,
} from "@/hooks/useWarcraftLogsApi";
import {
  normalizeWowName,
  classNameToPascalCase,
  specNameToPascalCase,
} from "@/utils/classandspecicons";
import { specMapping } from "@/utils/specmapping";
import SpecHeader from "./SpecHeader";
import ScoreEvolution from "./ScoreEvolution";
import TopPlayersTable from "./TopPlayersTable";
import DungeonPerformanceGrid from "./DungeonPerformanceGrid";
import InfoTooltip from "@/components/Shared/InfoTooltip";

interface SpecDetailViewProps {
  slug: string;
}

const SpecDetailView: React.FC<SpecDetailViewProps> = ({ slug }) => {
  // Parse the slug to get className and specName
  const [className, specName] = useMemo(() => {
    if (!slug) return [null, null];
    const parts = slug.split("-");
    if (parts.length === 2) {
      // Use the utility functions to properly format class and spec names
      const rawClassName = parts[0];
      const rawSpecName = parts[1];

      return [
        classNameToPascalCase(rawClassName),
        specNameToPascalCase(rawSpecName),
      ];
    }
    return [null, null];
  }, [slug]);

  // Fetch data for spec averages and top players
  const { data: allSpecsData, isLoading: isLoadingSpecs } =
    useGetSpecAverageGlobalScore();
  const { data: allPlayersData, isLoading: isLoadingPlayers } =
    useGetBestTenPlayerPerSpec();

  // Find current spec data
  const currentSpecData = useMemo(() => {
    if (!allSpecsData || !className || !specName) return null;
    return allSpecsData.find(
      (spec) =>
        spec.class.toLowerCase() === className.toLowerCase() &&
        spec.spec.toLowerCase() === specName.toLowerCase()
    );
  }, [allSpecsData, className, specName]);

  // Find top players for this spec
  const topPlayers = useMemo(() => {
    if (!allPlayersData || !className || !specName) return [];
    return allPlayersData
      .filter(
        (player) =>
          player.class.toLowerCase() === className.toLowerCase() &&
          player.spec.toLowerCase() === specName.toLowerCase()
      )
      .sort((a, b) => a.rank - b.rank);
  }, [allPlayersData, className, specName]);

  // Determine role from spec mapping
  const role = useMemo(() => {
    if (!className || !specName) return "";
    return specMapping[className]?.[specName]?.role || "DPS";
  }, [className, specName]);

  // Loading state
  if (isLoadingSpecs || isLoadingPlayers) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-600 mx-auto"></div>
          <p className="mt-4 text-slate-300">Loading spec data...</p>
        </div>
      </div>
    );
  }

  // Error state
  if (!currentSpecData || !className || !specName) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <p className="text-slate-300">
            No data available for this specialization.
          </p>
          <Link href="/performance-analysis/mythic-plus">
            <button className="mt-4 px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700">
              Back to All Specs
            </button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen text-slate-100">
      {/* Header with back button */}
      <header className="pt-8 pb-6 px-4 md:px-8 border-b border-slate-800">
        <div className="container mx-auto">
          <Link
            href="/performance-analysis"
            className="inline-flex items-center text-slate-300 hover:text-purple-400 transition-colors"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Performance Analysis
          </Link>

          <SpecHeader
            className={className}
            specName={specName}
            currentSpecData={currentSpecData}
            role={role}
          />
        </div>
      </header>

      {/* Main content */}
      <main className="flex-1 container mx-auto px-4 md:px-8 py-8">
        <div className="space-y-10">
          {/* Score Evolution Section */}
          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center">
              Mythic+ Score Evolution
              <InfoTooltip
                content={`This section shows the evolution of the score for ${specName} ${className} specialization over a period of time.`}
                className="ml-2"
                size="lg"
              />
            </h2>
            <ScoreEvolution />
          </section>

          {/* Players Section */}
          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center">
              Top {specName} {className} Players in Mythic+
              <InfoTooltip
                content={`This section shows the top 10 players for ${specName} ${className} specialization in Mythic+.`}
                className="ml-2"
                size="lg"
              />
            </h2>
            <TopPlayersTable players={topPlayers} className={className} />
          </section>

          {/* Dungeon Performance Section */}
          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center">
              Max Key Performance by Dungeon
              <InfoTooltip
                content={`This section shows the max key performance for ${specName} ${className} specialization in Mythic+.`}
                className="ml-2"
                size="lg"
              />
            </h2>
            <DungeonPerformanceGrid className={className} specName={specName} />
          </section>
        </div>
      </main>
    </div>
  );
};

export default SpecDetailView;
