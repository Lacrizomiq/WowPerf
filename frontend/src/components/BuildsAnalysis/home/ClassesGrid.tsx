import { WowClassParam } from "@/types/warcraftlogs/builds/classSpec";
import ClassCard from "./ClassCard";

export default function ClassesGrid() {
  // Organized classes by row by alphabetical order
  const classRows: WowClassParam[][] = [
    ["deathknight", "demonhunter", "druid", "evoker"],
    ["hunter", "mage", "monk", "paladin"],
    ["priest", "rogue", "shaman", "warlock"],
    ["warrior"],
  ];

  return (
    <div className="flex flex-col min-h-screen">
      {/* Page Header */}
      <header className="pt-8 pb-6 px-4 md:px-8 border-b border-slate-800">
        <div className="container mx-auto">
          <h1 className="text-3xl md:text-4xl font-bold mb-2">
            Data-Driven Builds
          </h1>
          <p className="text-muted-foreground text-base md:text-lg">
            Select a specialization to view builds, talents, gear, and more.
          </p>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-4 md:px-8 py-6">
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
      </main>
    </div>
  );
}
