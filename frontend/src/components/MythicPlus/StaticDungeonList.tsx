import React from "react";
import Image from "next/image";
import { Dungeon, MythicPlusRuns } from "@/types/mythicPlusRuns";

interface StaticDungeonListProps {
  dungeons: Dungeon[];
  mythicPlusRuns: MythicPlusRuns[];
  onDungeonClick: (dungeon: Dungeon) => void;
}

const StaticDungeonList: React.FC<StaticDungeonListProps> = ({
  dungeons,
  mythicPlusRuns,
  onDungeonClick,
}) => {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4 pt-4">
      {dungeons.map((dungeon) => {
        const run = mythicPlusRuns.find((r) => r.Dungeon.ID === dungeon.ID);
        return (
          <div
            key={dungeon.ID}
            className="flex flex-col items-center p-4 bg-deep-blue rounded-lg cursor-pointer hover:bg-blue-900 transition-colors duration-200"
            onClick={() => onDungeonClick(dungeon)}
          >
            <Image
              src={dungeon.MediaURL}
              alt={dungeon.Name}
              width={100}
              height={40}
              className="rounded-md border-2 border-gray-700"
            />
            <span className="font-bold mt-2">{dungeon.Name}</span>
            <span className="text-sm text-gray-400">{dungeon.ShortName}</span>
            {run && (
              <div className="mt-2 text-center">
                <p className="text-sm"> +{run.KeystoneLevel}</p>
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
};

export default StaticDungeonList;
