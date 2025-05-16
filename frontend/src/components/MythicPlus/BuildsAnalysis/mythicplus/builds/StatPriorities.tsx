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

  // Get stat colors for background
  const getStatColor = (statName: string) => {
    const colors: Record<string, string> = {
      Crit: "bg-pink-500 text-white",
      Haste: "bg-blue-500 text-white",
      Mastery: "bg-purple-500 text-white",
      Versatility: "bg-slate-500 text-white",
      Speed: "bg-amber-500 text-white",
      Avoidance: "bg-indigo-500 text-white",
      Leech: "bg-blue-400 text-white",
    };
    return colors[statName] || "bg-gray-600 text-white";
  };

  return (
    <div className="mb-4">
      <div className="bg-slate-900 rounded-lg border border-slate-800 p-1 mb-6">
        {/* Header with season info */}
        <div className="flex justify-between items-center mb-4 text-sm text-gray-400"></div>

        {/* Secondary Stats */}
        <div className="mb-6 flex flex-row">
          <div className="flex items-center mb-2 px-6">
            <span className="text-gray-300 mr-2">Secondary Stat Priority</span>
          </div>

          <div className="flex items-center flex-wrap">
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

        {/* Minor Stats */}
        {sortedMinor.length > 0 && (
          <div className="flex flex-row mb-2">
            <div className="flex items-center mb-2 px-6">
              <span className="text-gray-300 mr-2">Minor Stats Priority</span>
            </div>

            <div className="flex items-center flex-wrap">
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
        )}
      </div>
    </div>
  );
}
