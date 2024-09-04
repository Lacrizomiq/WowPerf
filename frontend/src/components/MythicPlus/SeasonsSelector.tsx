"use client";

import React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Season {
  slug: string;
  name: string;
  short_name: string;
}

interface SeasonsSelectorProps {
  seasons: Season[];
  onSeasonChange: (seasonSlug: string) => void;
}

const SeasonsSelector: React.FC<SeasonsSelectorProps> = ({
  seasons,
  onSeasonChange,
}) => {
  return (
    <Select onValueChange={onSeasonChange} defaultValue={seasons[0].slug}>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select a season" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value={seasons[0].slug}>{seasons[0].name}</SelectItem>
        {seasons.slice(1).map((season) => (
          <SelectItem key={season.slug} value={season.slug}>
            {season.name}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default SeasonsSelector;
