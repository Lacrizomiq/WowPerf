// StatPriorities.tsx - Version mise à jour avec divs séparées et couleurs harmonisées
import { StatPriority } from "@/types/warcraftlogs/builds/buildsAnalysis";
import { groupStatsByCategory } from "@/utils/buildsAnalysis/dataTransformer";

interface StatPrioritiesProps {
  stats: StatPriority[];
}

export default function StatPriorities({ stats }: StatPrioritiesProps) {
  const { secondary, minor } = groupStatsByCategory(stats);

  // Sort by priority rank
  const sortedSecondary = [...secondary].sort(
    (a, b) => a.priority_rank - b.priority_rank
  );
  const sortedMinor = [...minor].sort(
    (a, b) => a.priority_rank - b.priority_rank
  );

  // Get stat colors for background - harmonisées avec le design système
  const getStatColor = (statName: string) => {
    const colors: Record<string, string> = {
      Crit: "bg-pink-600 text-white",
      Haste: "bg-blue-600 text-white",
      Mastery: "bg-purple-600 text-white",
      Versatility: "bg-gray-600 text-white",
      Speed: "bg-amber-600 text-white",
      Avoidance: "bg-indigo-600 text-white",
      Leech: "bg-teal-600 text-white",
    };
    return colors[statName] || "bg-gray-600 text-white";
  };

  return (
    <div className="mb-4">
      {/* Secondary Stats */}
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-6 mb-4">
        <div className="flex items-center mb-4">
          <h3 className="text-lg font-medium text-white">Stat Priority</h3>
        </div>

        <div className="flex items-center flex-wrap gap-2">
          {sortedSecondary.map((stat, index) => (
            <div key={stat.stat_name} className="flex items-center">
              <div
                className={`px-3 py-1 rounded-md flex items-center ${getStatColor(
                  stat.stat_name
                )}`}
              >
                <span className="mr-1.5 bg-black bg-opacity-20 w-5 h-5 rounded-full flex items-center justify-center font-bold text-sm">
                  {index + 1}
                </span>
                <span>{stat.stat_name}</span>
                <span className="ml-2 bg-black bg-opacity-20 px-2 py-1 rounded-md text-xs">
                  {stat.avg_value.toFixed(0)}
                </span>
              </div>
              {index < sortedSecondary.length - 1 && (
                <span className="mx-2 text-gray-500">→</span>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Minor Stats  
      {sortedMinor.length > 0 && (
        <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-6">
          <div className="flex items-center mb-4">
            <h3 className="text-lg font-medium text-white">
              Minor Stats Priority
            </h3>
          </div>

          <div className="flex items-center flex-wrap gap-2">
            {sortedMinor.map((stat, index) => (
              <div key={stat.stat_name} className="flex items-center">
                <div
                  className={`px-3 py-1 rounded-md flex items-center ${getStatColor(
                    stat.stat_name
                  )}`}
                >
                  <span className="mr-1.5 bg-black bg-opacity-20 w-5 h-5 rounded-full flex items-center justify-center font-bold text-sm">
                    {index + 1}
                  </span>
                  <span>{stat.stat_name}</span>
                  <span className="ml-2 bg-black bg-opacity-20 px-2 py-1 rounded-md text-xs">
                    {stat.avg_value.toFixed(0)}
                  </span>
                </div>
                {index < sortedMinor.length - 1 && (
                  <span className="mx-2 text-gray-500">→</span>
                )}
              </div>
            ))}
          </div>
        </div>
      )} */}
    </div>
  );
}
