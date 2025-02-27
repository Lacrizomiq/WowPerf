// components/MythicPlus/SpecDetail/SpecDetailView.tsx
"use client";

import React, { useMemo } from "react";
import Link from "next/link";
import Image from "next/image";
import { useGetSpecAverageGlobalScore } from "@/hooks/useWarcraftLogsApi";
import { useGetBestTenPlayerPerSpec } from "@/hooks/useWarcraftLogsApi";
import { getSpecIcon, normalizeWowName } from "@/utils/classandspecicons";
import { specMapping } from "@/utils/specmapping";
import DungeonPerformance from "./MaxKeyLevelSpecDungeon";
import { getSpecBackground } from "@/utils/classandspecbackgrounds";
// Import realm mapping utilities for character URL generation
import { buildCharacterUrl, realmService } from "@/utils/realmMappingUtility";
import type { RegionType } from "@/utils/realmMappingUtility";

interface SpecDetailViewProps {
  slug: string;
}

const SpecDetailView: React.FC<SpecDetailViewProps> = ({ slug }) => {
  // Parse the slug to get className and specName
  const [className, specName] = useMemo(() => {
    if (!slug) return [null, null];
    const parts = slug.split("-");
    if (parts.length === 2) {
      return [
        parts[0].charAt(0).toUpperCase() + parts[0].slice(1),
        parts[1].charAt(0).toUpperCase() + parts[1].slice(1),
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

  // Build character URL from player data
  const handleCharacterUrl = (player: any): string => {
    try {
      // Find the realm by its name
      const realm = realmService.getRealmByName(player.server_name);

      if (!realm) {
        console.warn(
          `Unable to find realm for server name ${player.server_name}`,
          "Falling back to basic URL structure"
        );
        // Fallback to basic URL structure if the realm is not found
        const region = player.server_region.toLowerCase();
        const name = player.name.toLowerCase();
        return `/character/${region}/${player.server_name.toLowerCase()}/${name}`;
      }

      // Use buildCharacterUrl with the found ID
      const url = buildCharacterUrl(
        player.name,
        realm.id,
        player.server_region as RegionType
      );

      if (!url) {
        throw new Error("Failed to build character URL");
      }

      return url;
    } catch (error) {
      console.error("Error building character URL:", error);
      return "#";
    }
  };

  // Loading state
  if (isLoadingSpecs || isLoadingPlayers) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-white text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-white mx-auto"></div>
          <p className="mt-4">Loading spec data...</p>
        </div>
      </div>
    );
  }

  // Error state
  if (!currentSpecData || !className || !specName) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-white text-center">
          <p>No data available for this specialization.</p>
          <Link href="/mythic-plus/analysis">
            <button className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
              Back to All Specs
            </button>
          </Link>
        </div>
      </div>
    );
  }

  // Helper function to format class name for CSS
  const formatClassNameForCSS = (className: string): string => {
    return className.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
  };

  // Get spec icon
  const specIconUrl = getSpecIcon(className, normalizeWowName(specName));

  // Get spec background class
  const backgroundClass = getSpecBackground(className, specName);

  return (
    <div
      style={{
        backgroundColor: "black",
        color: "#e6f1ff",
        minHeight: "100vh",
      }}
    >
      {/* Header section with spec background */}
      <div
        className={`header-bg py-12 ${backgroundClass} mx-auto`}
        style={{
          backgroundSize: "cover",
          backgroundPosition: "top",
          backgroundRepeat: "no-repeat",
        }}
      >
        <div className="container mx-auto px-4">
          <div className="flex items-center">
            <Link href="/mythic-plus/analysis" className="text-blue-300 mr-4">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-5 w-5"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z"
                  clipRule="evenodd"
                />
              </svg>
            </Link>
            <Image
              src={specIconUrl}
              alt={`${specName} ${className}`}
              className="w-12 h-12 rounded-md mr-4"
              width={48}
              height={48}
              unoptimized
            />
            <div>
              <h1
                className="text-3xl font-bold"
                style={{
                  color: `var(--color-${formatClassNameForCSS(className)})`,
                }}
              >
                {specName} {className}
              </h1>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Stat Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          {/* Global Score Card */}
          <div
            className="p-4 rounded-md"
            style={{ backgroundColor: "#1a365d" }}
          >
            <p className="text-sm text-gray-400 mb-1">Global Score</p>
            <p className="text-3xl font-bold">
              {Math.round(currentSpecData.avg_global_score).toLocaleString()}
            </p>
            <p className="text-xs text-gray-400">
              Average score of the very best player
            </p>
          </div>

          {/* Spec Ranking Card */}
          <div
            className="p-4 rounded-md"
            style={{ backgroundColor: "#1a365d" }}
          >
            <p className="text-sm text-gray-400 mb-1">Spec Ranking</p>
            <p className="text-3xl font-bold">
              {currentSpecData ? "#" + currentSpecData.overall_rank : "N/A"}
            </p>
            <p className="text-xs text-gray-400">Spec rank in average score</p>
          </div>
        </div>

        {/* Main Content Area */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Top Players Section */}
          <div className="lg:col-span-2">
            <h2 className="text-xl font-bold mb-4">Top Players</h2>

            <div className="space-y-4">
              {topPlayers.map((player, index) => (
                <div
                  key={`${player.name}-${index}`}
                  className="p-4 rounded-md transition-all hover:transform hover:-translate-y-1 hover:shadow-lg"
                  style={{ backgroundColor: "#112240" }}
                >
                  <div className="flex items-center">
                    {/* Rank Badge */}
                    <div
                      className={`font-bold text-lg flex items-center justify-center w-7 h-7 rounded-full mr-3`}
                      style={{
                        backgroundColor:
                          player.rank === 1
                            ? "#d6b656" // Gold
                            : player.rank === 2
                            ? "#a6a6a6" // Silver
                            : player.rank === 3
                            ? "#ad8a56" // Bronze
                            : "#1a365d", // Default blue for other ranks
                        color: player.rank <= 3 ? "#000" : "inherit",
                      }}
                    >
                      {player.rank}
                    </div>
                    <div className="flex-grow">
                      <div className="flex items-center">
                        {/* Character name with link to profile */}
                        <Link
                          href={handleCharacterUrl(player)}
                          className={`font-bold text-lg hover:underline`}
                          style={{
                            color: `var(--color-${formatClassNameForCSS(
                              className
                            )})`,
                          }}
                        >
                          {player.name}
                        </Link>
                      </div>
                      <p className="text-sm text-gray-400">
                        {player.server_region} - {player.server_name}
                      </p>
                    </div>
                    {/* Player Score */}
                    <div className="text-right">
                      <p className="text-2xl font-bold">
                        {Math.round(player.total_score).toLocaleString()}
                      </p>
                      <p className="text-xs text-gray-400">Score</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Dungeon Performance Section */}
          <div>
            <h2 className="text-xl font-bold mb-4">Max Key Level</h2>
            <DungeonPerformance className={className} specName={specName} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default SpecDetailView;
