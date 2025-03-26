// app/mythic-plus/statistics/page.tsx
"use client";

import { useStats } from "@/providers/StatsContext";
import { useGetDungeonStats } from "@/hooks/useRaiderioApi";
import { OverallStats } from "@/components/Home/mythicplus/Stats/OverallStats";
import { DungeonStat } from "@/types/dungeonStats";

export default function StatisticsPage() {
  const { season, region, dungeon } = useStats();
  const { data: statsData } = useGetDungeonStats("season-tww-2", region);

  const currentDungeonStats =
    statsData?.find((stat: DungeonStat) => stat.dungeon_slug === dungeon) ||
    statsData?.[0];

  return <OverallStats stats={currentDungeonStats} />;
}
