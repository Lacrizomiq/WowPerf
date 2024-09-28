import React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { StaticRaid } from "@/types/raids";
import Image from "next/image";

interface RaidsSelectorProps {
  raids: StaticRaid[];
  onRaidChange: (raid: StaticRaid) => void;
  selectedRaid: StaticRaid | null;
}

const RaidsSelector: React.FC<RaidsSelectorProps> = ({
  raids,
  onRaidChange,
  selectedRaid,
}) => {
  return (
    <Select
      onValueChange={(value) => {
        const raid = raids.find((r) => r.Slug === value);
        if (raid) {
          onRaidChange(raid);
        }
      }}
      value={selectedRaid?.Slug || ""}
    >
      <SelectTrigger className="w-[200px] bg-gradient-blue shadow-2xl text-white border-none">
        <SelectValue placeholder="Select a raid" />
      </SelectTrigger>
      <SelectContent className="bg-black text-white">
        {raids.map((raid) => (
          <SelectItem
            key={raid.Slug}
            value={raid.Slug}
            className="hover:bg-gradient-purple mr-16"
          >
            <div className="flex items-center gap-2">
              <Image
                src={`https://wow.zamimg.com/images/wow/icons/large/${raid.Icon}.jpg`}
                alt={raid.Name}
                width={30}
                height={30}
              />
              {raid.Name}
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default RaidsSelector;
