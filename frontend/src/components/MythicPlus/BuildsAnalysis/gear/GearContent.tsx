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
  const slotId = searchParams.get("slot_id");

  // Convert slotId to number if it exists, otherwise null
  const selectedSlotId = slotId !== null ? parseInt(slotId) : null;

  // Initialize Wowhead tooltips
  useWowheadTooltips();

  // Refresh tooltips when parameters change
  useEffect(() => {
    if (typeof window !== "undefined" && window.$WowheadPower) {
      window.$WowheadPower.refreshLinks();
    }
  }, [encounterId, slotId]);

  return (
    <div className="space-y-6">
      {/* Conditional content based on encounter and slot */}
      {encounterId ? (
        <DungeonGear
          className={className}
          spec={spec}
          encounterId={encounterId}
          selectedSlotId={selectedSlotId}
        />
      ) : (
        <AllDungeonsGear
          className={className}
          spec={spec}
          selectedSlotId={selectedSlotId}
        />
      )}
    </div>
  );
}
