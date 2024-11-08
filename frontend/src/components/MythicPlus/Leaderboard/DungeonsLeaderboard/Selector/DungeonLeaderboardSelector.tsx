import React from "react";
import Image from "next/image";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Dungeon } from "@/types/mythicPlusRuns";
import { DUNGEON_ENCOUNTER_MAPPING } from "@/utils/s1_tww_mapping";

interface DungeonLeaderboardSelectorProps {
  dungeons: Dungeon[];
  onDungeonChange: (encounterID: number) => void;
  selectedDungeonId: number | null;
}

const DungeonLeaderboardSelector: React.FC<DungeonLeaderboardSelectorProps> = ({
  dungeons,
  onDungeonChange,
  selectedDungeonId,
}) => {
  return (
    <Select
      onValueChange={(value) => onDungeonChange(Number(value))}
      value={selectedDungeonId?.toString() || ""}
    >
      <SelectTrigger className="w-[200px] bg-gradient-blue shadow-2xl text-white border-none">
        <SelectValue placeholder="Select a dungeon" />
      </SelectTrigger>
      <SelectContent className="bg-black text-white">
        {dungeons.map((dungeon) => {
          const dungeonInfo = DUNGEON_ENCOUNTER_MAPPING[dungeon.Slug];
          if (!dungeonInfo) return null;

          return (
            <SelectItem
              key={dungeon.Slug}
              value={dungeonInfo.id.toString()}
              className="hover:bg-gradient-purple mr-16"
            >
              <div className="flex items-center gap-2">
                <Image
                  src={`https://wow.zamimg.com/images/wow/icons/large/${dungeon.Icon}.jpg`}
                  alt={dungeonInfo.name}
                  width={30}
                  height={30}
                  unoptimized
                />
                {dungeonInfo.name}
              </div>
            </SelectItem>
          );
        })}
      </SelectContent>
    </Select>
  );
};

export default DungeonLeaderboardSelector;
