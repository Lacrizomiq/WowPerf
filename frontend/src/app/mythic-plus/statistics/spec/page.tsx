// app/mythic-plus/statistics/spec/page.tsx
"use client";

import { useStats } from "@/providers/StatsContext";
import { useGetDungeonStats } from "@/hooks/useRaiderioApi";
import { SpecStats } from "@/components/Home/mythicplus/Stats/SpecStats";
import { DungeonStat } from "@/types/dungeonStats";

export default function SpecDistributionPage() {
  const { season, region, dungeon } = useStats();
  const { data: statsData } = useGetDungeonStats(season, region);

  const currentDungeonStats =
    statsData?.find((stat: DungeonStat) => stat.dungeon_slug === dungeon) ||
    statsData?.[0];

  return <SpecStats stats={currentDungeonStats} />;
}
