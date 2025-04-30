"use client";

import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import BuildHeader from "@/components/MythicPlus/BuildsAnalysis/layout/BuildHeader";
import BuildNav from "@/components/MythicPlus/BuildsAnalysis/layout/BuildNav";
import BuildFilters from "@/components/MythicPlus/BuildsAnalysis/layout/BuildFilters";
import { useState, useEffect } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import { Dungeon } from "@/types/mythicPlusRuns";

export default function BuildLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { class: string; spec: string };
}) {
  const className = params.class as WowClassParam;
  const spec = params.spec as WowSpecParam;
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();

  // Determine if the active tab is "builds"
  const isBuildsTab =
    !pathname.includes("talents") && !pathname.includes("gear");

  const isTalentsTab = pathname.includes("talents");
  const isGearTab = pathname.includes("gear");

  // Fetch dungeons data for mapping slugs to encounter IDs
  const season = "season-tww-2";
  const { data: dungeonData } = useGetBlizzardMythicDungeonPerSeason(season);
  const [dungeonId, setDungeonId] = useState<string>("all");

  // Update dungeonId from URL on initial load and when URL changes
  useEffect(() => {
    if (!dungeonData?.dungeons) return;

    const encounterId = searchParams.get("encounter_id");

    if (encounterId) {
      // Search for the dungeon by EncounterID (case-sensitive)
      const dungeon = dungeonData.dungeons.find(
        (d: Dungeon) => d.EncounterID?.toString() === encounterId
      );

      if (dungeon) {
        setDungeonId(dungeon.Slug);
      }
    } else {
      setDungeonId("all");
    }
  }, [searchParams, dungeonData]);

  // Handle dungeon change
  const handleDungeonChange = (value: string) => {
    setDungeonId(value);

    // Update: apply for talents AND gear
    if (isTalentsTab || isGearTab) {
      const params = new URLSearchParams(searchParams.toString());

      if (value === "all") {
        params.delete("encounter_id");
        router.push(`${pathname}?${params.toString()}`);
      } else if (dungeonData?.dungeons) {
        // Find the EncounterID corresponding to the selected slug
        const selectedDungeon = dungeonData.dungeons.find(
          (d: Dungeon) => d.Slug === value
        );

        if (selectedDungeon?.EncounterID) {
          params.set("encounter_id", selectedDungeon.EncounterID.toString());
          router.push(`${pathname}?${params.toString()}`);
        }
      }
    }
  };

  // Handle class change
  const handleClassChange = (newClass: WowClassParam) => {
    // Build the new path while keeping the same tab
    let newPath = `/mythic-plus/builds/${newClass}`;

    // Determine the first spec of the new class
    // This is already handled by ClassSelector, which redirects to the first spec,
    // but we need to add the correct URL suffixes

    if (pathname.includes("/talents")) {
      newPath += `/${spec}/talents`;

      // Keep the URL parameters for talents if necessary
      if (
        (isTalentsTab || isGearTab) && // Update: include gear for parameter conservation
        dungeonId !== "all" &&
        searchParams.has("encounter_id")
      ) {
        newPath += `?${searchParams.toString()}`;
      }
    } else if (pathname.includes("/gear")) {
      newPath += `/${spec}/gear`;

      // Add: Conservation des paramètres pour gear
      if (
        isGearTab &&
        dungeonId !== "all" &&
        searchParams.has("encounter_id")
      ) {
        newPath += `?${searchParams.toString()}`;
      }
    } else if (pathname.includes("/enchants-gems")) {
      newPath += `/${spec}/enchants-gems`;
    }

    // Let ClassSelector handle the navigation, we don't need to router.push here
  };

  // Handle spec change
  const handleSpecChange = (newSpec: WowSpecParam) => {
    // Build the new path while keeping the same tab
    let newPath = `/mythic-plus/builds/${className}/${newSpec}`;

    if (pathname.includes("/talents")) {
      newPath += "/talents";

      // Keep the URL parameters for talents if necessary
      if ((isTalentsTab || isGearTab) && dungeonId !== "all") {
        // Mise à jour: inclure gear
        const params = new URLSearchParams(searchParams.toString());
        router.push(`${newPath}?${params.toString()}`);
        return;
      }
    } else if (pathname.includes("/gear")) {
      newPath += "/gear";

      // Ajout: Conservation des paramètres pour gear
      if (isGearTab && dungeonId !== "all") {
        const params = new URLSearchParams(searchParams.toString());
        router.push(`${newPath}?${params.toString()}`);
        return;
      }
    } else if (pathname.includes("/enchants-gems")) {
      newPath += "/enchants-gems";
    }

    router.push(newPath);
  };

  return (
    <div className="w-full bg-black text-slate-100 min-h-screen">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-4">
        {/* Header Section */}
        <BuildHeader className={className} spec={spec} />

        {/* Main Navigation */}
        <BuildNav defaultTab="builds" className={className} spec={spec}>
          {/* Filters Section */}
          <BuildFilters
            className={className}
            spec={spec}
            onDungeonChange={handleDungeonChange}
            onClassChange={handleClassChange}
            onSpecChange={handleSpecChange}
            showDungeonSelector={!isBuildsTab}
            selectedDungeon={dungeonId}
          />

          {/* Content will be injected here */}
          {children}
        </BuildNav>
      </div>
    </div>
  );
}
