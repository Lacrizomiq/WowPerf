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

export interface ChartData {
  name: string;
  value: number;
  percentage: number;
  color: string;
  rawValue: number;
}

interface StatsBarChartProps {
  data: ChartData[];
  vertical?: boolean;
  height?: number;
  leftMargin?: number;
  formatter?: (value: number) => string;
}

interface TooltipProps {
  payload?: Array<{
    payload: ChartData;
  }>;
}

interface YAxisProps {
  x: number;
  y: number;
  payload: {
    value: string;
  };
}

// Mise à jour de ChartData

const StatsBarChart: React.FC<StatsBarChartProps> = ({
  data,
  vertical = false,
  height = 400,
  leftMargin = 50,
  formatter = (value) => `${value}%`,
}) => {
  return (
    <ResponsiveContainer width="100%" height={height}>
      <BarChart
        data={data}
        layout={vertical ? "vertical" : "horizontal"}
        margin={{ top: 20, right: 30, left: leftMargin, bottom: 5 }}
      >
        {vertical ? (
          <>
            <XAxis type="number" domain={[0, "dataMax"]} />
            <YAxis
              type="category"
              dataKey="name"
              width={180}
              tick={({ x, y, payload }) => (
                <text
                  x={x - 5}
                  y={y}
                  dy={4}
                  textAnchor="end"
                  fill="white"
                  className="text-sm"
                >
                  {payload.value}
                </text>
              )}
            />
          </>
        ) : (
          <>
            <XAxis type="category" dataKey="name" />
            <YAxis type="number" domain={[0, 100]} />
          </>
        )}
        <Tooltip
          formatter={(value: number, name: string, props: any) => [
            formatter(value),
            props.payload.name,
          ]}
          contentStyle={{
            backgroundColor: "#fff",
            border: "none",
            color: "black",
          }}
          cursor={{ fill: "transparent" }}
        />
        <Bar dataKey="value">
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={entry.color} />
          ))}
          <LabelList
            dataKey="value"
            position={vertical ? "right" : "top"}
            formatter={formatter}
            style={{
              fill: "white",
              fontWeight: "bold",
              textShadow: "1px 1px 1px #000",
            }}
          />
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
};

export default StatsBarChart;
