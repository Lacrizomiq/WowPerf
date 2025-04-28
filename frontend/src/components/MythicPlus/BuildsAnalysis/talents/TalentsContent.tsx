"use client";

import { useEffect } from "react";
import { useSearchParams } from "next/navigation";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import AllDungeonsBuilds from "./AllDungeonsBuilds";
import DungeonBuilds from "./DungeonBuilds";

interface TalentsContentProps {
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function TalentsContent({
  className,
  spec,
}: TalentsContentProps) {
  const searchParams = useSearchParams();
  const encounterId = searchParams.get("encounter_id");

  return (
    <div className="space-y-6">
      {/* Content */}
      {encounterId ? (
        <DungeonBuilds
          className={className}
          spec={spec}
          encounterId={encounterId}
        />
      ) : (
        <AllDungeonsBuilds className={className} spec={spec} />
      )}
    </div>
  );
}
