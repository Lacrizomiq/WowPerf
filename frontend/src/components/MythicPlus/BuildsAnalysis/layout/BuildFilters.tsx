// components/builds/layout/BuildFilters.tsx
import { useState, useEffect } from "react";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import DungeonSelector from "./DungeonSelector";
import ClassSelector from "./ClassSelector";
import SpecSelector from "./SpecSelector";
import { Dungeon } from "@/types/mythicPlusRuns";

interface BuildFiltersProps {
  className: WowClassParam;
  spec: WowSpecParam;
  onDungeonChange?: (dungeon: string) => void;
  onAffixChange?: (affix: string) => void;
}

export default function BuildFilters({
  className,
  spec,
  onDungeonChange,
  onAffixChange,
}: BuildFiltersProps) {
  const season = "season-tww-2";
  const [selectedDungeon, setSelectedDungeon] = useState("all");
  const [dungeons, setDungeons] = useState<Dungeon[]>([]);

  // Fetch dungeons data
  const { data: dungeonData } = useGetBlizzardMythicDungeonPerSeason(season);

  useEffect(() => {
    if (dungeonData && dungeonData.dungeons) {
      setDungeons(dungeonData.dungeons);
    }
  }, [dungeonData]);

  // Handle changes
  const handleDungeonChange = (value: string) => {
    setSelectedDungeon(value);
    if (onDungeonChange) onDungeonChange(value);
  };

  return (
    <div className="flex flex-wrap gap-2 mb-6">
      <button className="px-4 py-2 bg-slate-800 rounded flex items-center gap-2 hover:bg-slate-700">
        <span className="text-indigo-400">Filters</span>
      </button>

      <ClassSelector
        selectedClass={className}
        onClassChange={() => {}} // Disabled in display mode
      />

      <SpecSelector
        selectedClass={className}
        selectedSpec={spec}
        onSpecChange={() => {}} // Disabled in display mode
      />

      <div className="px-4 py-2 bg-indigo-600 rounded flex items-center">
        <span>Mythic+</span>
      </div>

      <DungeonSelector
        dungeons={dungeons}
        onDungeonChange={handleDungeonChange}
        selectedDungeon={selectedDungeon}
      />
    </div>
  );
}
