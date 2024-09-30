import React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface RegionSelectorProps {
  regions: { value: string; label: string }[];
  onRegionChange: (region: string) => void;
  selectedRegion: string;
}

const RegionSelector: React.FC<RegionSelectorProps> = ({
  regions,
  onRegionChange,
  selectedRegion,
}) => {
  return (
    <Select onValueChange={onRegionChange} value={selectedRegion}>
      <SelectTrigger className="w-[200px] bg-gradient-blue shadow-2xl text-white border-none">
        <SelectValue placeholder="Select a region" />
      </SelectTrigger>
      <SelectContent className="bg-black text-white">
        <SelectItem
          key="world"
          value="world"
          className="hover:bg-gradient-purple"
        >
          World
        </SelectItem>
        {regions.map((region) => (
          <SelectItem
            key={region.value}
            value={region.value}
            className="hover:bg-gradient-purple"
          >
            {region.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default RegionSelector;
