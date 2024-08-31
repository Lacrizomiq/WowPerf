import React from "react";
import Image from "next/image";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";

interface MythicDungeonProps {
  seasonSlug: string;
}

const MythicDungeon: React.FC<MythicDungeonProps> = ({ seasonSlug }) => {
  const {
    data: dungeonData,
    isLoading,
    error,
  } = useGetBlizzardMythicDungeonPerSeason(seasonSlug);

  if (isLoading) return <div>Loading dungeon data...</div>;
  if (error) return <div>Error loading dungeon data: {error.message}</div>;
  if (!dungeonData) return <div>No dungeon data found</div>;

  console.log("Dungeon Data:", dungeonData);

  return (
    <div className="p-4 bg-gradient-dark shadow-lg rounded-lg glow-effect m-12 max-w-6xl mx-auto">
      <div className="flex flex-col justify-between items-center mb-4">
        {dungeonData.dungeons.map((dungeon) => {
          return (
            <div
              key={dungeon.id}
              className="flex flex-col justify-center items-center gap-2"
            >
              <Image
                src={dungeon.MediaURL}
                alt={dungeon.Name}
                width={100}
                height={40}
                className="rounded-md border-2 border-gray-700"
              />
              <span>{dungeon.Name}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default MythicDungeon;
