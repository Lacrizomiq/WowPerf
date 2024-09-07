import React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface ExpansionSelectorProps {
  expansions: string[];
  selectedExpansion: string;
  onExpansionChange: (expansion: string) => void;
}

const ExpansionSelector: React.FC<ExpansionSelectorProps> = ({
  expansions,
  selectedExpansion,
  onExpansionChange,
}) => {
  const uniqueExpansions = Array.from(new Set(expansions));

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
