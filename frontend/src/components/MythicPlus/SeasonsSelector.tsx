"use client";

import React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Dungeon } from "@/types/mythicPlusRuns";

interface Season {
  slug: string;
  name: string;
  shortName: string;
  id: number;
  Dungeons: null | Dungeon[];
}

interface SeasonsSelectorProps {
  seasons: Season[];
  onSeasonChange: (seasonSlug: string) => void;
  selectedSeason: Season;
}

const SeasonsSelector: React.FC<SeasonsSelectorProps> = ({
  seasons,
  onSeasonChange,
  selectedSeason,
}) => {
  return (
    <>
      <Select onValueChange={onSeasonChange} value={selectedSeason.slug}>
        <SelectTrigger className="w-[200px] bg-gradient-purple text-white border-none">
          <SelectValue placeholder="Select a season" />
        </SelectTrigger>
        <SelectContent className="bg-gradient-purple text-white">
          {seasons.map((season) => (
            <SelectItem
              key={season.slug}
              value={season.slug}
              className="hover:bg-gradient-purple"
            >
              {season.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </>
  );
};

export default SeasonsSelector;
