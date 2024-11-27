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

export const SpecStats: React.FC<StatsProps> = ({ stats }) => {
  const prepareSpecData = () => {
    if (!stats?.SpecStats) return [];

    const specStats = stats.SpecStats;
    const formattedData = [];

    for (const [className, specs] of Object.entries(specStats)) {
      for (const [specName, count] of Object.entries(specs)) {
        formattedData.push({
          name: `${specName} ${className}`,
          count,
          color: `var(--color-${className.toLowerCase().replace(" ", "-")})`,
        });
      }
    }

    return formattedData.sort((a, b) => b.count - a.count);
  };

  const specData = prepareSpecData();
  const totalPlayers = specData.reduce((sum, entry) => sum + entry.count, 0);

  return (
    <div className="bg-deep-blue p-4 rounded-lg shadow-2xl">
      <h3 className="text-xl font-bold text-white mb-4">
        Specialization Distribution - Total: {totalPlayers} players
      </h3>
      <ResponsiveContainer
        width="100%"
        height={Math.max(600, specData.length * 35)}
      >
        <BarChart
          data={specData}
          layout="vertical"
          margin={{ top: 20, right: 50, left: 20, bottom: 5 }}
        >
          <XAxis type="number" stroke="white" />
          <YAxis
            type="category"
            dataKey="name"
            width={180}
            stroke="white"
            tick={{ fill: "white" }}
          />
          <Tooltip
            contentStyle={{
              backgroundColor: "#fff",
              border: "1px solid #333",
              borderRadius: "4px",
              color: "black",
            }}
            cursor={{ fill: "transparent" }}
          />
          <Bar dataKey="count" name="Players">
            {specData.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={entry.color} />
            ))}
            <LabelList
              dataKey="count"
              position="right"
              fill="white"
              formatter={(value: number) => `${value}`}
            />
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};
