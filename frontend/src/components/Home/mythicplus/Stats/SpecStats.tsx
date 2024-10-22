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
    const specStats = stats.SpecStats;
    const formattedData = [];

    for (const [className, specs] of Object.entries(specStats)) {
      for (const [specName, count] of Object.entries(specs)) {
        formattedData.push({
          className,
          specName,
          fullName: `${specName} ${className}`,
          count,
          color: `var(--color-${className.toLowerCase().replace(" ", "-")})`,
        });
      }
    }

    return formattedData.sort((a, b) => b.count - a.count);
  };

  const specData = prepareSpecData();

  return (
    <div className="bg-deep-blue p-4 rounded-lg shadow-2xl">
      <h3 className="text-xl font-bold text-white mb-4">
        Specialization Distribution
      </h3>
      <ResponsiveContainer
        width="100%"
        height={Math.max(600, specData.length * 35)}
      >
        <BarChart
          data={specData}
          layout="vertical"
          margin={{ top: 20, right: 50, left: 200, bottom: 5 }}
        >
          <XAxis type="number" domain={[0, "dataMax"]} />
          <YAxis
            type="category"
            dataKey="fullName"
            width={180}
            tick={({ x, y, payload }) => (
              <g transform={`translate(${x},${y})`}>
                <text
                  x={-5}
                  y={0}
                  dy={4}
                  textAnchor="end"
                  fill="white"
                  className="text-sm"
                >
                  {payload.value}
                </text>
              </g>
            )}
          />
          <Tooltip
            formatter={(value: number, name: string, props: any) => [
              `${value} players`,
              props.payload.fullName,
            ]}
            contentStyle={{
              backgroundColor: "#1a1a1a",
              border: "1px solid #333",
              borderRadius: "4px",
              color: "white",
            }}
            labelStyle={{ color: "white" }}
            itemStyle={{ color: "white" }}
          />
          <Bar dataKey="count" name="Players">
            {specData.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={entry.color} />
            ))}
            <LabelList
              dataKey="count"
              position="insideRight"
              content={({ x, y, width, value }) => (
                <text
                  x={Number(x) + Number(width) + 10}
                  y={Number(y) + 4}
                  fill="white"
                  textAnchor="start"
                  className="text-sm"
                >
                  {value}
                </text>
              )}
            />
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};
