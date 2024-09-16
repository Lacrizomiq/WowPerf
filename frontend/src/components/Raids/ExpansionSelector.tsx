import React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface ExpansionSelectorProps {
  currentExpansion: string;
  onExpansionChange: (expansion: string) => void;
}

const ExpansionSelector: React.FC<ExpansionSelectorProps> = ({
  currentExpansion,
  onExpansionChange,
}) => {
  const availableExpansions = ["TWW", "DF"];

  return (
    <Select onValueChange={onExpansionChange} value={currentExpansion}>
      <SelectTrigger className="w-[200px] bg-gradient-purple  text-white border-none">
        <SelectValue placeholder="Select an expansion" />
      </SelectTrigger>
      <SelectContent className="bg-gradient-purple  text-white">
        {availableExpansions.map((expansion) => (
          <SelectItem
            key={expansion}
            value={expansion}
            className="hover:bg-gradient-purple"
          >
            {expansion}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default ExpansionSelector;
