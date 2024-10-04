import React, { useState, useEffect } from "react";
import { useGetDungeonStats } from "@/hooks/useRaiderioApi";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  LabelList,
  Cell,
} from "recharts";
import { ClassStats, DungeonStat, RoleStats } from "@/types/dungeonStats";
import DungeonSelector from "./Selector/DungeonSelector";
import RegionSelector from "./Selector/RegionSelector";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import { Dungeon } from "@/types/mythicPlusRuns";

const DungeonStats: React.FC = () => {
  const [season] = useState("season-tww-1");
  const [region, setRegion] = useState("world");
  const [dungeon, setDungeon] = useState("all");

  const {
    data: statsData,
    isLoading,
    error,
  } = useGetDungeonStats(season, region);

  console.log("stats data ", statsData);

  const { data: dungeonData } = useGetBlizzardMythicDungeonPerSeason(season);

  const [dungeons, setDungeons] = useState<Dungeon[]>([]);

  useEffect(() => {
    if (dungeonData && dungeonData.dungeons) {
      setDungeons(dungeonData.dungeons);
    }
  }, [dungeonData]);

  if (isLoading) return <div className="text-white">Loading stats...</div>;
  if (error)
    return (
      <div className="text-red-500">Error loading stats: {error.message}</div>
    );

  const currentDungeonStats =
    statsData?.find((stat) => stat.dungeon_slug === dungeon) || statsData?.[0];

  if (!currentDungeonStats)
    return (
      <div className="text-white">No data available for this dungeon.</div>
    );

  const getLevelRange = (levelStats: Record<string, number>) => {
    const levels = Object.keys(levelStats).map(Number);
    const minLevel = Math.min(...levels);
    const maxLevel = Math.max(...levels);
    return `+${minLevel} / +${maxLevel}`;
  };

  const prepareChartData = (role: string) => {
    const roleStats = currentDungeonStats?.RoleStats || {};
    const classStats =
      roleStats[role.toLowerCase() as keyof typeof roleStats] || {};
    const total = Object.values(classStats).reduce(
      (sum, count) => sum + (count as number),
      0
    );
    return Object.entries(classStats)
      .map(([className, count]) => ({
        className,
        count: count as number,
        percentage: Number((((count as number) / total) * 100).toFixed(2)),
        color: `var(--color-${className.toLowerCase().replace(" ", "-")})`,
      }))
      .sort((a, b) => b.percentage - a.percentage);
  };

  const roles = ["tank", "healer", "dps"];

  console.log("statsData", statsData);

  return (
    <div className="p-4 bg-[#0a0a0a] bg-opacity-80">
      <h2 className="text-2xl font-bold text-white mb-4">
        Dungeon Statistics for{" "}
        {dungeon
          .split("-")
          .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(" ")}{" "}
        dungeon in region: {region.toUpperCase()}
      </h2>

      <div className="mb-4 flex space-x-4">
        <RegionSelector
          regions={["us", "eu", "kr", "tw", "cn"]}
          onRegionChange={setRegion}
          selectedRegion={region}
        />
        <DungeonSelector
          dungeons={dungeons}
          onDungeonChange={setDungeon}
          selectedDungeon={dungeon}
        />
      </div>

      <div className="p-4">
        <p>
          Last update:{" "}
          {new Intl.DateTimeFormat("en-US", {
            weekday: "long",
            day: "2-digit",
            month: "long",
            year: "numeric",
          }).format(new Date(statsData[0].updated_at))}
        </p>
      </div>

      <div className="p-4 bg-deep-blue rounded-lg mb-4">
        <h3 className="text-xl font-bold text-white mb-2">
          Mythic+ KeystoneLevel Range
        </h3>
        <p className="text-white text-lg">
          {getLevelRange(currentDungeonStats.LevelStats)}
        </p>
      </div>

      <div className="space-y-8 pt-4 ">
        {roles.map((role) => {
          const chartData = prepareChartData(role);
          if (chartData.length === 0) {
            return (
              <div key={role} className="p-4 rounded-lg">
                <h3 className="text-xl font-bold text-white mb-2 capitalize">
                  {role}
                </h3>
                <p className="text-white">No data available for this role.</p>
              </div>
            );
          }
          return (
            <div key={role} className="bg-deep-blue p-4 rounded-lg">
              <h3 className="text-xl font-bold text-white mb-4 capitalize">
                {role} - Total:{" "}
                {chartData.reduce((sum, entry) => sum + entry.count, 0)} players
              </h3>
              <ResponsiveContainer width="100%" height={400}>
                <BarChart
                  data={chartData}
                  layout="horizontal"
                  margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
                >
                  <XAxis type="category" dataKey="className" />
                  <YAxis type="number" domain={[0, 100]} />
                  <Tooltip
                    formatter={(value: number, name: string, props: any) => [
                      `${
                        props.payload.count
                      } players (${props.payload.percentage.toFixed(2)}%)`,
                      props.payload.className,
                    ]}
                    contentStyle={{ backgroundColor: "#000", border: "none" }}
                    cursor={{ fill: "transparent" }}
                  />
                  <Bar dataKey="percentage" name="Percentage">
                    {chartData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                    <LabelList
                      dataKey="percentage"
                      position="top"
                      formatter={(value: number) =>
                        value > 5 ? `${value.toFixed(2)}%` : ""
                      }
                      style={{
                        fill: "white",
                        fontWeight: "bold",
                        textShadow: "1px 1px 1px #000",
                      }}
                    />
                    <LabelList
                      dataKey="percentage"
                      position="top"
                      formatter={(value: number) =>
                        value <= 5 ? `${value.toFixed(2)}%` : ""
                      }
                      style={{
                        fill: "white",
                        fontSize: "12px",
                        textShadow: "1px 1px 1px #000",
                      }}
                    />
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default DungeonStats;
