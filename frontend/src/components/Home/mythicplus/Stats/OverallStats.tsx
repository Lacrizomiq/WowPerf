import React from "react";
import StatsBarChart from "../../../Charts/StatsBarChart";
import StatsSectionHeader from "../../../Charts/StatsSectionHeader";
import { ChartData } from "../../../Charts/StatsBarChart";
import { DungeonStat, RoleStats } from "@/types/dungeonStats";

interface StatsProps {
  stats: DungeonStat;
}

export const OverallStats: React.FC<StatsProps> = ({ stats }) => {
  const prepareChartData = (role: keyof RoleStats): ChartData[] => {
    const roleStats = stats?.RoleStats || {};
    const classStats = roleStats[role] || {};

    const total = Object.values(classStats).reduce(
      (sum: number, count: number) => sum + count,
      0
    );

    return Object.entries(classStats)
      .map(([className, count]) => {
        const percentage = Number(((count / total) * 100).toFixed(2));
        return {
          name: className,
          value: percentage, // use the percentage as value
          rawValue: count, // keep the raw count as additional data
          percentage,
          color: `var(--color-${className.toLowerCase().replace(" ", "-")})`,
        };
      })
      .sort((a, b) => b.percentage - a.percentage);
  };

  const roles: Array<keyof RoleStats> = ["tank", "healer", "dps"];

  return (
    <div className="space-y-8">
      {roles.map((role) => {
        const chartData = prepareChartData(role);
        if (chartData.length === 0) return null;

        const totalPlayers = chartData.reduce(
          (sum, entry) => sum + entry.rawValue,
          0
        );

        return (
          <div key={role} className="bg-deep-blue p-4 rounded-lg shadow-2xl">
            <StatsSectionHeader
              title={role.toUpperCase()}
              total={totalPlayers}
            />
            <StatsBarChart
              data={chartData}
              formatter={(value: number) =>
                `${value.toFixed(2)}% (${
                  chartData.find((d) => d.value === value)?.rawValue
                })`
              }
            />
          </div>
        );
      })}
    </div>
  );
};
