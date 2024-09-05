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
    <Select onValueChange={onSeasonChange} value={selectedSeason.slug}>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select a season" />
      </SelectTrigger>
      <SelectContent>
        {seasons.map((season) => (
          <SelectItem key={season.slug} value={season.slug}>
            {season.name}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default SeasonsSelector;
