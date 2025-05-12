// components/performance/mythicplus/SpecCard.tsx
import Link from "next/link";
import Image from "next/image";
import { Badge } from "@/components/ui/badge";
import { SpecAverageGlobalScore } from "@/types/warcraftlogs/globalLeaderboardAnalysis";
import { getSpecIcon, normalizeWowName } from "@/utils/classandspecicons";
import { ArrowUp } from "lucide-react";

interface SpecCardProps {
  specData: SpecAverageGlobalScore;
  selectedRole: string;
}

export default function SpecCard({ specData, selectedRole }: SpecCardProps) {
  // Helper to format the class name for CSS classes (DeathKnight -> death-knight)
  const formatClassNameForCSS = (className: string): string => {
    return className.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
  };

  // Get the CSS class for the class colors
  const getClassColorClass = (): string => {
    const cssClassName = formatClassNameForCSS(specData.class);
    return `class-color--${cssClassName}`;
  };

  // Get the URL of the specialization icon
  const getSpecIconUrl = (): string => {
    try {
      const normalizedSpecName = normalizeWowName(specData.spec);
      return getSpecIcon(specData.class, normalizedSpecName);
    } catch (error) {
      console.warn(
        `Error getting spec icon for ${specData.class}-${specData.spec}:`,
        error
      );
      return "";
    }
  };

  // Format the class and specialization name for display
  const formatSpecName = (className: string, specName: string): string => {
    // Convert camelCase to normal format
    const formattedClass = className.replace(/([a-z])([A-Z])/g, "$1 $2");
    return `${specName} ${formattedClass}`;
  };

  // Determine which rank to display based on the selected role
  const displayRank =
    selectedRole === "ALL" ? specData.overall_rank : specData.role_rank;

  const specIconUrl = getSpecIconUrl();
  const classColorClass = getClassColorClass();
  const specSlug =
    specData.slug ||
    `${specData.class.toLowerCase()}-${specData.spec
      .toLowerCase()
      .replace(/ /g, "-")}`;

  // Maximum data - to implement

  return (
    <Link href={`/performance-analysis/mythic-plus/${specSlug}`}>
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 hover:border-purple-700/50 transition-all hover:shadow-md cursor-pointer">
        <div className="flex items-start gap-4">
          {/* Left section - Rank and Icon */}
          <div className="flex flex-col items-center gap-1">
            {/* Rank with empty change indicator for now */}
            <div
              className={`text-2xl font-bold ${
                displayRank <= 3 ? "text-purple-400" : "text-slate-400"
              }`}
            >
              #{displayRank}
            </div>

            {/* Space for future rank change indicator */}
            <div className="h-4 w-full flex justify-center">
              {/* Left intentionally empty as requested */}
            </div>

            {/* Icon */}
            <div className="relative w-12 h-12 rounded-full overflow-hidden bg-slate-700 mt-1">
              {specIconUrl ? (
                <Image
                  src={specIconUrl}
                  alt={`${specData.spec} icon`}
                  className="w-full h-full object-cover"
                  width={48}
                  height={48}
                  unoptimized
                />
              ) : (
                <div className="w-full h-full bg-slate-700" />
              )}
            </div>
          </div>

          {/* Information about the specialization */}
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-3">
              {/* Specialization name with class color applied */}
              <h3 className={`font-bold text-lg ${classColorClass}`}>
                {formatSpecName(specData.class, specData.spec)}
              </h3>
              <Badge
                variant="outline"
                className="text-xs py-0 h-5 border-slate-600"
              >
                {specData.role}
              </Badge>
            </div>

            {/* Score details - Layout modified */}
            <div className="grid grid-cols-2 gap-x-6 gap-y-2">
              <div>
                <div className="text-xs text-slate-400 mb-0.5">
                  Average Score
                </div>
                <div className="text-xl font-bold text-white">
                  {Math.round(specData.avg_global_score).toLocaleString()}
                </div>
              </div>

              <div>
                <div className="text-xs text-slate-400 mb-0.5">
                  Weekly Evolution{" "}
                  <Badge className="ml-2 bg-purple-600 text-[10px]">Soon</Badge>
                </div>
                <div className="flex items-center h-7">
                  {/* Left intentionally empty until we have the data */}
                  {/* <MiniSparkline /> */}
                </div>
              </div>

              <div>
                <div className="text-xs text-slate-400 mb-0.5">Max Score </div>
                <div className="text-sm font-medium">
                  {specData.max_global_score}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}

// Mini Sparkline conservée mais adapté
function MiniSparkline() {
  return (
    <div className="ml-2 flex items-end h-5 gap-0.5">
      <div
        className="w-1 bg-purple-900/50 rounded-sm"
        style={{ height: "30%" }}
      ></div>
      <div
        className="w-1 bg-purple-900/50 rounded-sm"
        style={{ height: "50%" }}
      ></div>
      <div
        className="w-1 bg-purple-900/50 rounded-sm"
        style={{ height: "40%" }}
      ></div>
      <div
        className="w-1 bg-purple-900/50 rounded-sm"
        style={{ height: "60%" }}
      ></div>
      <div
        className="w-1 bg-purple-600 rounded-sm"
        style={{ height: "80%" }}
      ></div>
    </div>
  );
}
