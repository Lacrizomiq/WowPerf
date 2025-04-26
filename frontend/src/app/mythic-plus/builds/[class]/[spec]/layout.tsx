"use client";

import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import BuildHeader from "@/components/MythicPlus/BuildsAnalysis/layout/BuildHeader";
import BuildNav from "@/components/MythicPlus/BuildsAnalysis/layout/BuildNav";
import BuildFilters from "@/components/MythicPlus/BuildsAnalysis/layout/BuildFilters";
import { useState } from "react";
import { usePathname } from "next/navigation";

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

  // Determine if the active tab is "builds"
  const isBuildsTab =
    !pathname.includes("talents") &&
    !pathname.includes("gear") &&
    !pathname.includes("enchants-gems");

  const [dungeonId, setDungeonId] = useState<string>("all");

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
            onDungeonChange={(value) => setDungeonId(value)}
            showDungeonSelector={!isBuildsTab}
          />

          {/* Content will be injected here */}
          {children}
        </BuildNav>
      </div>
    </div>
  );
}
