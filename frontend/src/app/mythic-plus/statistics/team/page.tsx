// app/mythic-plus/statistics/team/page.tsx
"use client";

import { useStats } from "@/providers/StatsContext";
import { useGetDungeonStats } from "@/hooks/useRaiderioApi";
import { TeamComposition } from "@/components/Home/mythicplus/Stats/TeamComposition";
import { DungeonStat } from "@/types/dungeonStats";
export default function TeamCompositionPage() {
  const { season, region, dungeon } = useStats();
  const { data: statsData } = useGetDungeonStats(season, region);

  const currentDungeonStats =
    statsData?.find((stat: DungeonStat) => stat.dungeon_slug === dungeon) ||
    statsData?.[0];

  return <TeamComposition stats={currentDungeonStats} />;
}
