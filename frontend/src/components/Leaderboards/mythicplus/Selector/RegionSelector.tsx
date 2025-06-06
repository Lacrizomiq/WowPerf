// RegionSelector.tsx - Version harmonisée
import React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface RegionSelectorProps {
  regions: string[];
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
      <SelectTrigger className="w-[200px] bg-slate-800/50 text-white border-slate-700 focus:ring-purple-600">
        <SelectValue placeholder="Select a region" />
      </SelectTrigger>
      <SelectContent className="bg-slate-900 border-slate-700 text-white">
        <SelectItem
          key="world"
          value="world"
          className="hover:bg-slate-800 focus:bg-purple-600"
        >
          World
        </SelectItem>
        {regions.map((region) => (
          <SelectItem
            key={region}
            value={region}
            className="hover:bg-slate-800 focus:bg-purple-600"
          >
            {region}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default RegionSelector;
