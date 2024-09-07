import React from "react";
import Image from "next/image";
import { StaticRaid, RaidProgressionData } from "@/types/raids";

interface StaticRaidsListProps {
  raids: StaticRaid[];
  raidProgressionData: RaidProgressionData | null;
  onRaidSelect: (raid: StaticRaid) => void;
}

const StaticRaidsList: React.FC<StaticRaidsListProps> = ({
  raids,
  raidProgressionData,
  onRaidSelect,
}) => {
  const getDifficultyInfo = (raidId: number) => {
    if (!raidProgressionData) return "";
    const raid = raidProgressionData.expansions
      .flatMap((exp) => exp.raids)
      .find((r) => r.id === raidId);

    if (!raid) return "";

    const difficultyOrder = ["Mythic", "Heroic", "Normal"];
    const difficultyMap: { [key: string]: string } = {
      Mythic: "M",
      Heroic: "H",
      Normal: "N",
    };

    return difficultyOrder
      .map((diff) => {
        const mode = raid.modes.find((m) => m.difficulty === diff);
        if (mode) {
          const { completed_count, total_count } = mode.progress;
          return `${completed_count}/${total_count}${difficultyMap[diff]}`;
        }
        return null;
      })
      .filter(Boolean)
      .join(", ");
  };

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 mt-4">
      {raids.map((raid) => (
        <div
          key={raid.ID}
          className="bg-deep-blue p-4 rounded-lg cursor-pointer hover:bg-blue-700 transition-colors duration-200 flex flex-col items-center"
          onClick={() => onRaidSelect(raid)}
        >
          <Image
            src={raid.MediaURL}
            alt={raid.Name}
            width={200}
            height={200}
            className="rounded-md mb-2"
          />
          <h3 className="font-bold text-lg mb-2 text-center">{raid.Name}</h3>
          <p className="text-sm text-center">{getDifficultyInfo(raid.ID)}</p>
        </div>
      ))}
    </div>
  );
};

export default StaticRaidsList;
