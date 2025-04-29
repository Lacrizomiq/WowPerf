"use client";

import { useEffect } from "react";
import { useSearchParams } from "next/navigation";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import AllDungeonsGear from "./AllDungeonsGear";
import DungeonGear from "./DungeonGear";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";

interface GearContentProps {
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function GearContent({ className, spec }: GearContentProps) {
  const searchParams = useSearchParams();
  const encounterId = searchParams.get("encounter_id");

  // Initialize Wowhead tooltips
  useWowheadTooltips();

  // Refresh tooltips when parameters change
  useEffect(() => {
    if (typeof window !== "undefined" && window.$WowheadPower) {
      window.$WowheadPower.refreshLinks();
    }
  }, [encounterId]);

  return (
    <div className="space-y-6">
      {/* Conditional content exactly like in TalentsContent */}
      {encounterId ? (
        <DungeonGear
          className={className}
          spec={spec}
          encounterId={encounterId}
        />
      ) : (
        <AllDungeonsGear className={className} spec={spec} />
      )}
    </div>
  );
}
