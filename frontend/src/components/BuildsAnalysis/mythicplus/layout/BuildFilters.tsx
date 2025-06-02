// BuildFilters.tsx - Version complète harmonisée
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
  onClassChange?: (className: WowClassParam) => void;
  onSpecChange?: (spec: WowSpecParam) => void;
  showDungeonSelector?: boolean;
  selectedDungeon?: string;
}

export default function BuildFilters({
  className,
  spec,
  onDungeonChange,
  onClassChange,
  onSpecChange,
  showDungeonSelector = true,
  selectedDungeon = "all",
}: BuildFiltersProps) {
  const season = "season-tww-2";
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
    if (onDungeonChange) onDungeonChange(value);
  };

  return (
    <div className="flex flex-wrap gap-2 mb-6">
      <ClassSelector
        selectedClass={className}
        onClassChange={onClassChange || (() => {})}
      />

      <SpecSelector
        selectedClass={className}
        selectedSpec={spec}
        onSpecChange={onSpecChange || (() => {})}
      />

      <div className="px-4 py-2 bg-purple-600 rounded flex items-center text-white">
        <span>Mythic+</span>
      </div>

      {showDungeonSelector && (
        <DungeonSelector
          dungeons={dungeons}
          onDungeonChange={handleDungeonChange}
          selectedDungeon={selectedDungeon}
        />
      )}
    </div>
  );
}
