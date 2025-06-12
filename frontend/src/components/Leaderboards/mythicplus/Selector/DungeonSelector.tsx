// DungeonSelector.tsx - Version harmonisée
import React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Dungeon } from "@/types/mythicPlusRuns";
import Image from "next/image";

interface DungeonSelectorProps {
  dungeons: Dungeon[];
  onDungeonChange: (dungeon: string) => void;
  selectedDungeon: string;
}

const DungeonSelector: React.FC<DungeonSelectorProps> = ({
  dungeons,
  onDungeonChange,
  selectedDungeon,
}) => {
  return (
    <Select onValueChange={onDungeonChange} value={selectedDungeon}>
      <SelectTrigger className="w-[200px] bg-slate-800/50 text-white border-slate-700 focus:ring-purple-600">
        <SelectValue placeholder="Select a dungeon" />
      </SelectTrigger>
      <SelectContent className="bg-slate-900 border-slate-700 text-white">
        <SelectItem
          key="all"
          value="all"
          className="hover:bg-slate-800 focus:bg-purple-600"
        >
          All Dungeons
        </SelectItem>
        {dungeons.map((dungeon) => {
          return (
            <SelectItem
              key={dungeon.Slug}
              value={dungeon.Slug}
              className="hover:bg-slate-800 focus:bg-purple-600"
            >
              <div className="flex items-center gap-2">
                <Image
                  src={`https://wow.zamimg.com/images/wow/icons/large/${dungeon.Icon}.jpg`}
                  alt={dungeon.Name}
                  width={30}
                  height={30}
                  unoptimized
                />
                {dungeon.Name}
              </div>
            </SelectItem>
          );
        })}
      </SelectContent>
    </Select>
  );
};

export default DungeonSelector;
