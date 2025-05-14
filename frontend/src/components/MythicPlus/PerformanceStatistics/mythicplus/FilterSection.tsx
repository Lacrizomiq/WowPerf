// components/performance/mythicplus/FilterSection.tsx
import { RefreshCw } from "lucide-react";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import { useEffect, useState } from "react";
import { Dungeon } from "@/types/mythicPlusRuns";

import ClassSelector from "@/components/MythicPlus/PerformanceStatistics/selector/ClassSelector";
import DungeonSelector from "@/components/MythicPlus/PerformanceStatistics/selector/DungeonSelector";

interface FilterSectionProps {
  selectedRole: string;
  selectedClass: string | null;
  selectedDungeon: string;
  availableClasses: string[];
  onRoleChange: (role: any) => void;
  onClassChange: (className: string | null) => void;
  onDungeonChange: (dungeon: string) => void;
  onResetFilters: () => void;
  isFiltered: boolean;
}

export default function FilterSection({
  selectedRole,
  selectedClass,
  selectedDungeon,
  availableClasses,
  onRoleChange,
  onClassChange,
  onDungeonChange,
  onResetFilters,
  isFiltered,
}: FilterSectionProps) {
  const [dungeons, setDungeons] = useState<Dungeon[]>([]);

  const season = "season-tww-2"; // Saison courante
  const { data: dungeonData } = useGetBlizzardMythicDungeonPerSeason(season);

  useEffect(() => {
    if (dungeonData && dungeonData.dungeons) {
      setDungeons(dungeonData.dungeons);
    }
  }, [dungeonData]);

  return (
    <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-4 md:p-5">
      {/* Responsive layout */}
      <div className="md:flex md:justify-between md:space-x-8 md:items-end">
        {/* Grouped filters for medium and wide screens */}
        <div className="md:flex md:space-x-8 space-y-6 md:space-y-0 md:items-end">
          {/* Role Section */}
          <div>
            <label className="text-base font-medium mb-3 block text-white">
              Role
            </label>
            <div className="flex flex-wrap gap-2">
              <button
                className={`px-4 py-2 rounded-md text-sm ${
                  selectedRole === "ALL"
                    ? "bg-purple-600 text-white"
                    : "bg-slate-800 text-slate-300 hover:bg-slate-700"
                }`}
                onClick={() => onRoleChange("ALL")}
              >
                All
              </button>
              <button
                className={`px-4 py-2 rounded-md text-sm ${
                  selectedRole === "Tank"
                    ? "bg-purple-600 text-white"
                    : "bg-slate-800 text-slate-300 hover:bg-slate-700"
                }`}
                onClick={() => onRoleChange("Tank")}
              >
                Tank
              </button>
              <button
                className={`px-4 py-2 rounded-md text-sm ${
                  selectedRole === "Healer"
                    ? "bg-purple-600 text-white"
                    : "bg-slate-800 text-slate-300 hover:bg-slate-700"
                }`}
                onClick={() => onRoleChange("Healer")}
              >
                Healer
              </button>
              <button
                className={`px-4 py-2 rounded-md text-sm ${
                  selectedRole === "DPS"
                    ? "bg-purple-600 text-white"
                    : "bg-slate-800 text-slate-300 hover:bg-slate-700"
                }`}
                onClick={() => onRoleChange("DPS")}
              >
                DPS
              </button>
            </div>
          </div>

          {/* Class Section */}
          <div>
            <label className="text-base font-medium mb-3 block text-white">
              Class
            </label>
            <ClassSelector
              selectedClass={selectedClass}
              onClassChange={onClassChange}
              availableClasses={availableClasses}
            />
          </div>

          {/* Dungeon Section */}
          <div>
            <label className="text-base font-medium mb-3 block text-white">
              Dungeon
            </label>
            <DungeonSelector
              dungeons={dungeons}
              onDungeonChange={onDungeonChange}
              selectedDungeon={selectedDungeon}
            />
          </div>
        </div>

        {/* Reset Button - different responsive versions */}
        <div className="mt-6 md:mt-0">
          {/* Mobile version (full width, with text) */}
          <button
            onClick={onResetFilters}
            className={`md:hidden flex items-center justify-center w-full px-4 py-2 border border-slate-700 rounded-md hover:bg-slate-700 transition-colors ${
              !isFiltered ? "opacity-50 cursor-not-allowed" : ""
            }`}
            disabled={!isFiltered}
            aria-label="Reset filters"
          >
            <RefreshCw className="h-4 w-4 mr-2" />
            <span>Reset Filters</span>
          </button>

          {/* Desktop version (only icon) */}
          <button
            onClick={onResetFilters}
            className={`hidden md:flex md:items-center md:justify-center p-2 border border-slate-700 rounded-md hover:bg-slate-700 transition-colors ${
              !isFiltered ? "opacity-50 cursor-not-allowed" : ""
            }`}
            disabled={!isFiltered}
            aria-label="Reset filters"
          >
            <RefreshCw className="h-5 w-5" />
          </button>
        </div>
      </div>
    </div>
  );
}
