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

export const TeamComposition: React.FC<StatsProps> = ({ stats }) => {
  const prepareTeamData = () => {
    const teamComps = stats.TeamComp;
    return Object.entries(teamComps)
      .map(([key, value]) => ({
        id: key,
        count: value.count,
        composition: value.composition,
        percentage: (
          (value.count /
            Object.values(teamComps).reduce(
              (sum, comp) => sum + comp.count,
              0
            )) *
          100
        ).toFixed(2),
      }))
      .sort((a, b) => b.count - a.count);
  };

  const formatRole = (role: string): string => {
    if (role === "tank") return "Tank";
    if (role === "healer") return "Healer";
    if (role.startsWith("dps")) return "Dps";
    return role;
  };

  const teamData = prepareTeamData();

  return (
    <div className="bg-deep-blue p-4 rounded-lg shadow-2xl">
      <h3 className="text-xl font-bold text-white mb-4">
        Popular Team Compositions
      </h3>
      <div className="space-y-4">
        {teamData.map((team) => (
          <div key={team.id} className="border border-gray-700 rounded-lg p-4">
            <div className="flex justify-between items-center">
              <span className="text-white font-bold">
                {team.count} teams ({team.percentage}%)
              </span>
            </div>
            <div className="grid grid-cols-5 gap-4 mt-2">
              {Object.entries(team.composition).map(([role, info]) => (
                <div key={role} className="text-center">
                  <div className="text-gray-400 text-sm">
                    {formatRole(role)}
                  </div>
                  <div className="text-white">{info.class}</div>
                  <div className="text-gray-300 text-sm">{info.spec}</div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
