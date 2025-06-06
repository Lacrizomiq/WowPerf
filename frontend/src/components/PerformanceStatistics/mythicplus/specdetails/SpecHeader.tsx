import React from "react";
import Image from "next/image";
import { Badge } from "@/components/ui/badge";
import { getSpecIcon, normalizeWowName } from "@/utils/classandspecicons";
import { getSpecBackground } from "@/utils/classandspecbackgrounds";
import StatCard from "./StatCard";
import InfoTooltip from "@/components/Shared/InfoTooltip";

interface SpecHeaderProps {
  className: string;
  specName: string;
  currentSpecData: any;
  role: string;
}

const SpecHeader: React.FC<SpecHeaderProps> = ({
  className,
  specName,
  currentSpecData,
  role,
}) => {
  // Helper function to format class name for CSS
  const formatClassNameForCSS = (className: string): string => {
    return className.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
  };

  // Get spec icon
  const specIconUrl = getSpecIcon(className, normalizeWowName(specName));

  // Get spec background class
  const backgroundClass = getSpecBackground(className, specName);

  return (
    <div
      className={`mt-6 relative pt-8 pb-6 px-4 md:px-8 rounded-lg ${backgroundClass}`}
      style={{
        backgroundSize: "cover",
        backgroundPosition: "center top",
        backgroundRepeat: "no-repeat",
      }}
    >
      {/* Dark overlay for better text readability */}
      <div className="absolute inset-0  rounded-lg" />

      <div className="relative z-10">
        <div className="flex flex-col md:flex-row md:items-center gap-6">
          {/* Spec Icon and Name */}
          <div className="flex items-center gap-4">
            <div className="relative w-12 h-12 md:w-16 md:h-16 rounded-full overflow-hidden bg-slate-800 border-2 border-purple-600">
              <Image
                src={specIconUrl}
                alt={`${specName} ${className}`}
                fill
                className="object-cover"
                unoptimized
              />
            </div>
            <div>
              <h1
                className="text-3xl md:text-4xl font-bold"
                style={{
                  color: `var(--color-${formatClassNameForCSS(className)})`,
                }}
              >
                {specName} {className}
              </h1>
              <div className="flex items-center text-slate-300 mt-1">
                <Badge
                  variant="outline"
                  className="text-xs py-0 h-5 border-slate-600"
                >
                  {role}
                </Badge>
              </div>
            </div>
          </div>

          {/* Key Summary Stats */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 md:ml-auto mt-2 md:mt-0">
            <StatCard
              title="Average Score"
              value={Math.round(
                currentSpecData.avg_global_score
              ).toLocaleString()}
              tooltip={`Average score for ${specName} ${className} specialization across the top 10 players.`}
            />
            <StatCard
              title="Overall Rank"
              value={`#${currentSpecData.overall_rank}`}
              tooltip={`Overall rank for ${specName} ${className} specialization across all specializations.`}
            />
            <StatCard
              title="Weekly Change"
              value="Coming Soon"
              isComing={true}
              tooltip={`Weekly change for ${specName} ${className} specialization rank across all specializations and average score . Coming soon.`}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default SpecHeader;
