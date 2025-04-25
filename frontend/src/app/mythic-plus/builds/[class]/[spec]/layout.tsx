"use client";

import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import BuildHeader from "@/components/MythicPlus/BuildsAnalysis/layout/BuildHeader";
import BuildNav from "@/components/MythicPlus/BuildsAnalysis/layout/BuildNav";
import BuildFilters from "@/components/MythicPlus/BuildsAnalysis/layout/BuildFilters";
import { useState } from "react";

export default function BuildLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { class: string; spec: string };
}) {
  const className = params.class as WowClassParam;
  const spec = params.spec as WowSpecParam;

  const [dungeonId, setDungeonId] = useState<string>("all");

  return (
    <div className="container mx-auto p-4 bg-slate-900 text-slate-100 min-h-screen">
      {/* Header Section */}
      <BuildHeader className={className} spec={spec} />

      {/* Main Navigation */}
      <BuildNav defaultTab="builds">
        {/* Filters Section */}
        <BuildFilters
          className={className}
          spec={spec}
          onDungeonChange={(value) => setDungeonId(value)}
        />

        {/* Content will be injected here */}
        {children}
      </BuildNav>
    </div>
  );
}
