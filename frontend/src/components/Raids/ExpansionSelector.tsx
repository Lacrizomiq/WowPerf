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
      <SelectTrigger className="w-[200px] bg-deep-blue shadow-2xl text-white border-none">
        <SelectValue placeholder="Select an expansion" />
      </SelectTrigger>
      <SelectContent className="bg-deep-blue  text-white">
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
