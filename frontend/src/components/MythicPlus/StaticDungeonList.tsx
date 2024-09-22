import React from "react";
import Image from "next/image";
import { Dungeon, MythicPlusRuns } from "@/types/mythicPlusRuns";

interface StaticDungeonListProps {
  dungeons: Dungeon[];
  mythicPlusRuns: MythicPlusRuns[];
  onDungeonClick: (dungeon: Dungeon) => void;
  selectedDungeon: Dungeon | null;
}

const StaticDungeonList: React.FC<StaticDungeonListProps> = ({
  dungeons,
  mythicPlusRuns,
  onDungeonClick,
  selectedDungeon,
}) => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-4">
      {dungeons.map((dungeon) => {
        const run = mythicPlusRuns.find((r) => r.Dungeon.ID === dungeon.ID);
        return (
          <div
            key={dungeon.ID}
            className={`rounded-xl overflow-hidden bg-deep-blue shadow-lg cursor-pointer transition-all duration-300 ${
              selectedDungeon?.ID === dungeon.ID
                ? "ring-2 ring-blue-500 shadow-2xl scale-105"
                : "hover:shadow-xl hover:scale-105"
            }`}
            onClick={() => onDungeonClick(dungeon)}
          >
            <div className="h-36 bg-gray-300 relative">
              <Image
                src={dungeon.MediaURL}
                alt={dungeon.Name}
                layout="fill"
                objectFit="cover"
                className="mb-2"
              />
            </div>
            <div className="p-4">
              <h3 className="font-bold text-lg mb-2 text-center text-white">
                {dungeon.Name}
              </h3>
              {run ? (
                <div className="text-center">
                  <p className="text-blue-100">Level: {run.KeystoneLevel}</p>
                </div>
              ) : (
                <p className="text-center text-blue-300">
                  No run data available
                </p>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default StaticDungeonList;
