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
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-4">
      {dungeons.map((dungeon) => {
        const run = mythicPlusRuns.find((r) => r.Dungeon.ID === dungeon.ID);
        return (
          <div
            key={dungeon.ID}
            className="bg-deep-blue p-4 rounded-lg cursor-pointer hover:bg-blue-700 transition-colors duration-200 flex flex-col items-center"
            onClick={() => onDungeonClick(dungeon)}
          >
            <Image
              src={dungeon.MediaURL}
              alt={dungeon.Name}
              width={100}
              height={100}
              className="rounded-md mb-2"
            />
            <h3 className="font-bold text-lg mb-2 text-center">
              {dungeon.Name}
            </h3>
            {run ? (
              <>
                <p>Level: {run.KeystoneLevel}</p>
                <p>Score: {run.MythicRating.toFixed(2)}</p>
              </>
            ) : (
              <p>No run data available</p>
            )}
          </div>
        );
      })}
    </div>
  );
};

export default StaticDungeonList;
