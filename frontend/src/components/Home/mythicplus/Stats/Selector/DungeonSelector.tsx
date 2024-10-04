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
  dungeons: (Dungeon & { Slug: string })[];
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
      <SelectTrigger className="w-[200px] bg-gradient-blue shadow-2xl text-white border-none">
        <SelectValue placeholder="Select a dungeon" />
      </SelectTrigger>
      <SelectContent className="bg-black text-white">
        <SelectItem key="all" value="all" className="hover:bg-gradient-purple">
          All Dungeons
        </SelectItem>
        {dungeons.map((dungeon) => (
          <SelectItem
            key={dungeon.Slug}
            value={dungeon.Slug}
            className="hover:bg-gradient-purple mr-16"
          >
            <div className="flex items-center gap-2">
              <Image
                src={`https://wow.zamimg.com/images/wow/icons/large/${dungeon.Icon}.jpg`}
                alt={dungeon.Name}
                width={30}
                height={30}
              />
              {dungeon.Name}
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default DungeonSelector;
