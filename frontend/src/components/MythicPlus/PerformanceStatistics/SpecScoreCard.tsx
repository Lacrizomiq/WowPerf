import React from "react";
import { SpecAverageGlobalScore } from "@/types/warcraftlogs/globalLeaderboardAnalysis";
import { getSpecIcon, normalizeWowName } from "@/utils/classandspecicons";
import Image from "next/image";
import Link from "next/link";

interface SpecScoreCardProps {
  specData: SpecAverageGlobalScore;
  selectedRole: "Tank" | "Healer" | "DPS" | "ALL";
}

const SpecScoreCard: React.FC<SpecScoreCardProps> = ({
  specData,
  selectedRole,
}) => {
  // Helper to format class name for CSS classes (DeathKnight -> death-knight)
  const formatClassNameForCSS = (className: string): string => {
    return className.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
  };

  // Get class colors CSS class for the card
  const getClassColorClass = (): string => {
    const cssClassName = formatClassNameForCSS(specData.class);
    return `class-color--${cssClassName}`;
  };

  // Get spec icon URL
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

  // Format the class and spec name for display
  const formatSpecName = (className: string, specName: string): string => {
    // Convert camelCase class names to normal format
    const formattedClass = className.replace(/([a-z])([A-Z])/g, "$1 $2");
    return `${specName} ${formattedClass}`;
  };

  // Determine which rank to display based on selected role
  const displayRank =
    selectedRole === "ALL" ? specData.overall_rank : specData.role_rank;

  const specIconUrl = getSpecIconUrl();
  const classColorClass = getClassColorClass();
  const specSlug =
    specData.slug ||
    `${specData.class.toLowerCase()}-${specData.spec
      .toLowerCase()
      .replace(/ /g, "-")}`;

  return (
    <div>
      <Link href={`/mythic-plus/analysis/${specSlug}`}>
        <div
          className="bg-[#112240] rounded-md p-4 transition-all hover:transform hover:-translate-y-1 hover:shadow-lg cursor-pointer relative overflow-hidden"
          style={{
            borderLeftColor: `var(--color-${formatClassNameForCSS(
              specData.class
            )})`,
            borderLeftWidth: "4px",
          }}
        >
          {/* Rank Badge - Positioned at the right side */}
          <div className="absolute right-0 top-0 bottom-0 flex items-center justify-center w-16">
            <div className="flex items-center justify-center w-12 h-12 rounded-full bg-blue-700 text-white font-bold shadow-md">
              #{displayRank}
            </div>
          </div>

          <div className="flex items-center pr-14">
            {" "}
            {/* Add padding right to make room for the badge */}
            {/* Spec icon */}
            {specIconUrl ? (
              <Image
                src={specIconUrl}
                alt={`${specData.spec} icon`}
                className="rounded mr-3"
                width={42}
                height={42}
                unoptimized
              />
            ) : (
              <div className="w-12 h-12 bg-gray-700 rounded mr-3"></div>
            )}
            <div>
              <h3 className={`text-xl font-bold ${classColorClass}`}>
                {formatSpecName(specData.class, specData.spec)}
              </h3>
              <p className="text-sm text-gray-400 capitalize">
                {specData.role?.toLowerCase()}
              </p>
            </div>
            <div className="ml-auto text-right">
              <p className="text-2xl font-bold text-white">
                {Math.round(specData.avg_global_score).toLocaleString()}
              </p>
            </div>
          </div>
        </div>
      </Link>
    </div>
  );
};

export default SpecScoreCard;
