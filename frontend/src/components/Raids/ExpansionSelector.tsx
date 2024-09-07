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
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select an expansion" />
      </SelectTrigger>
      <SelectContent>
        {availableExpansions.map((expansion) => (
          <SelectItem key={expansion} value={expansion}>
            {expansion}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default ExpansionSelector;
