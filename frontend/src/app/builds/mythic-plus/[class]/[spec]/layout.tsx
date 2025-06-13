// BuildLayout.tsx
"use client";

import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import BuildHeader from "@/components/BuildsAnalysis/mythicplus/layout/BuildHeader";
import BuildNav from "@/components/BuildsAnalysis/mythicplus/layout/BuildNav";
import BuildFilters from "@/components/BuildsAnalysis/mythicplus/layout/BuildFilters";
import ContentTypeTabs from "@/components/BuildsAnalysis/layout/ContentTypeTabs";
import { useState, useEffect } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import { Dungeon } from "@/types/mythicPlusRuns";

export default function BuildLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ class: string; spec: string }>;
}) {
  const [className, setClassName] = useState<WowClassParam | null>(null);
  const [spec, setSpec] = useState<WowSpecParam | null>(null);
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();

  // Fetch dungeons data for mapping slugs to encounter IDs
  const season = "season-tww-2";
  const { data: dungeonData } = useGetBlizzardMythicDungeonPerSeason(season);
  const [dungeonId, setDungeonId] = useState<string>("all");

  // Resolve params asynchronously
  useEffect(() => {
    const resolveParams = async () => {
      const resolvedParams = await params;
      setClassName(resolvedParams.class as WowClassParam);
      setSpec(resolvedParams.spec as WowSpecParam);
    };
    resolveParams();
  }, [params]);

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

  // Early return while params are resolving (after all hooks)
  if (!className || !spec) {
    return (
      <div className="w-full text-slate-100 min-h-screen bg-[#1A1D21] flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-2 border-primary border-t-transparent"></div>
      </div>
    );
  }

  // Determine if the active tab is "builds"
  const isBuildsTab =
    !pathname.includes("talents") && !pathname.includes("gear");

  const isTalentsTab = pathname.includes("talents");
  const isGearTab = pathname.includes("gear");

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
    let newPath = `/builds/mythic-plus/${newClass}`;

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
    let newPath = `/builds/mythic-plus/${className}/${newSpec}`;

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
    <div className="w-full text-slate-100 min-h-screen bg-[#1A1D21]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-4">
        {/* Header Section */}
        <div className="border-b border-slate-800 w-full">
          <header className="pt-4 pb-6">
            <h1 className="text-3xl md:text-4xl font-bold mb-2 text-white">
              Builds
            </h1>
            <p className="text-muted-foreground text-base md:text-lg">
              Explore optimal talent builds, gear setups, and stat priorities
              for all classes and specializations.
            </p>
          </header>

          {/* Content Type Tabs */}
          <ContentTypeTabs
            className={className}
            spec={spec}
            activeTab="mythic-plus"
          />
        </div>

        <div className="pt-6">
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
    </div>
  );
}
