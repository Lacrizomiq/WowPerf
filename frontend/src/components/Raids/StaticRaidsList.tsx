import React from "react";
import Image from "next/image";
import { StaticRaid, RaidProgressionData } from "@/types/raids";

interface StaticRaidsListProps {
  raids: StaticRaid[];
  raidProgressionData: RaidProgressionData | null;
  onRaidSelect: (raid: StaticRaid) => void;
  selectedRaid: StaticRaid | null;
}

const StaticRaidsList: React.FC<StaticRaidsListProps> = ({
  raids,
  raidProgressionData,
  onRaidSelect,
  selectedRaid,
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
          className={`rounded-xl overflow-hidden bg-deep-blue shadow-lg cursor-pointer transition-all duration-300 ${
            selectedRaid?.ID === raid.ID
              ? "ring-2 ring-[#001830] shadow-2xl scale-105 glow-effect"
              : "hover:shadow-xl"
          }`}
          onClick={() => onRaidSelect(raid)}
        >
          <div className="h-48 bg-gray-300 relative">
            <Image
              src={raid.MediaURL}
              alt={raid.Name}
              layout="fill"
              className="object-cover mb-2"
            />
          </div>
          <div className="p-4">
            <h3 className="font-bold text-lg mb-2 text-center text-white">
              {raid.Name}
            </h3>
            <p className="text-sm text-center text-white">
              {getDifficultyInfo(raid.ID)}
            </p>
          </div>
        </div>
      ))}
    </div>
  );
};

export default StaticRaidsList;
