"use client";

import React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Raid } from "@/types/raids";

interface ExpansionSelectorProps {
  raids: Raid[];
  onExpansionChange: (expansion: string) => void;
  selectedExpansion: string;
}

// Select an expansion to see the raids
const ExpansionSelector: React.FC<ExpansionSelectorProps> = ({
  raids,
  onExpansionChange,
  selectedExpansion,
}) => {
  const uniqueExpansions = Array.from(
    new Set(raids.map((raid) => raid.Expansion))
  );

  return (
    <Select onValueChange={onExpansionChange} value={selectedExpansion}>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select an expansion" />
      </SelectTrigger>
      <SelectContent>
        {uniqueExpansions.map((expansion) => (
          <SelectItem key={expansion} value={expansion}>
            {expansion}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default ExpansionSelector;
