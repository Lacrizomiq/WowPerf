import React from "react";
import Image from "next/image";
import { DungeonStat } from "@/types/dungeonStats";

interface StatsProps {
  stats: DungeonStat;
}

export const TeamComposition: React.FC<StatsProps> = ({ stats }) => {
  const roleIcons = {
    tank: "https://cdn.raiderio.net/assets/img/role_tank-6cee7610058306ba277e82c392987134.png",
    healer:
      "https://cdn.raiderio.net/assets/img/role_healer-984e5e9867d6508a714a9c878d87441b.png",
    dps: "https://cdn.raiderio.net/assets/img/role_dps-eb25989187d4d3ac866d609dc009f090.png",
  };

  const formatClassName = (className: string): string => {
    const classMap: { [key: string]: string } = {
      "Death Knight": "death-knight",
      "Demon Hunter": "demon-hunter",
    };

    return classMap[className] || className.toLowerCase();
  };

  const getRoleIcon = (role: string) => {
    if (role === "tank") return roleIcons.tank;
    if (role === "healer") return roleIcons.healer;
    if (role.startsWith("dps")) return roleIcons.dps;
    return "";
  };

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

  const teamData = prepareTeamData();

  return (
    <div className="bg-deep-blue p-4 rounded-lg shadow-2xl">
      <h3 className="text-xl font-bold text-white mb-4">
        Popular Team Compositions - Number of teams:{" "}
        {teamData.reduce((sum, entry) => sum + entry.count, 0)}
      </h3>
      <div className="space-y-4">
        {teamData.map((team) => (
          <div key={team.id} className="border border-gray-700 rounded-lg p-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-white font-bold">
                {team.count} teams ({team.percentage}%)
              </span>
            </div>
            <div className="grid grid-cols-5 gap-4 mt-4">
              {Object.entries(team.composition).map(([role, info]) => (
                <div key={role} className="text-center">
                  <div className="rounded-full  bg-opacity-50 flex items-center justify-center relative mb-2">
                    <Image
                      src={getRoleIcon(role)}
                      alt={role}
                      width={24}
                      height={24}
                    />
                  </div>
                  <div
                    className={`text-white font-bold class-color--${formatClassName(
                      info.class
                    )}`}
                  >
                    {info.class}
                  </div>
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
