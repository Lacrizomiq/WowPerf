import { useState } from "react";
import ClassCard from "./ClassCard";
import { WowClassParam } from "@/types/warcraftlogs/builds/classSpec";

export default function ClassesGrid() {
  // Organized classes by row by alphabetical order
  const classRows: WowClassParam[][] = [
    ["priest", "demonhunter", "druid", "evoker"],
    ["hunter", "mage", "monk", "paladin"],
    ["priest", "deathknight", "rogue", "shaman"],
    ["warlock", "warrior"],
  ];

  return (
    <div className="min-h-screen bg-black text-slate-100 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-4xl font-bold mb-2">Data-Driven Mythic+ Builds</h1>
        <p className="text-slate-400 mb-10">
          Select a specialization to view builds, talents, gear, and more.
        </p>

        <div className="space-y-6">
          {classRows.map((classRow, rowIndex) => (
            <div
              key={rowIndex}
              className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"
            >
              {classRow.map((wowClass) => (
                <ClassCard key={wowClass} className={wowClass} />
              ))}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
