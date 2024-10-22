import React from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  LabelList,
  Cell,
} from "recharts";
import { DungeonStat } from "@/types/dungeonStats";

interface StatsProps {
  stats: DungeonStat;
}

export const OverallStats: React.FC<StatsProps> = ({ stats }) => {
  const prepareChartData = (role: string) => {
    const roleStats = stats?.RoleStats || {};
    const classStats =
      roleStats[role.toLowerCase() as keyof typeof roleStats] || {};
    const total = Object.values(classStats).reduce(
      (sum: number, count) => sum + (count as number),
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

  return (
    <div className="space-y-8">
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
          <div key={role} className="bg-deep-blue p-4 rounded-lg shadow-2xl">
            <h3 className="text-xl font-bold text-white mb-4 capitalize">
              {role} - Total:{" "}
              {chartData.reduce((sum, entry) => sum + entry.count, 0)} players
            </h3>
            <ResponsiveContainer width="100%" height={400}>
              <BarChart
                data={chartData}
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
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        );
      })}
    </div>
  );
};
