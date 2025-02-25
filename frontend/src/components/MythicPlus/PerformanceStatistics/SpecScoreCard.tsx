import React from "react";
import { SpecAverageGlobalScore } from "@/types/warcraftlogs/globalLeaderboardAnalysis";
import { specMapping } from "@/utils/specmapping";
import { getSpecIcon, normalizeWowName } from "@/utils/classandspecicons";
import Image from "next/image";
import Link from "next/link";

interface SpecScoreCardProps {
  specData: SpecAverageGlobalScore;
}

const SpecScoreCard: React.FC<SpecScoreCardProps> = ({ specData }) => {
  // Get role from the spec mapping
  const getRole = (
    className: string,
    specName: string
  ): "TANK" | "HEALER" | "DPS" | undefined => {
    if (specMapping[className] && specMapping[className][specName]) {
      return specMapping[className][specName].role;
    }
    return undefined;
  };

  const role = getRole(specData.class, specData.spec);

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

  const specIconUrl = getSpecIconUrl();
  const classColorClass = getClassColorClass();

  // Create the spec slug directly - just lowercase both values and join with hyphen
  const specSlug = `${specData.class.toLowerCase()}-${specData.spec.toLowerCase()}`;

  return (
    <div>
      <Link href={`/mythic-plus/analysis/${specSlug}`}>
        <div
          className="bg-[#112240] rounded-md p-4 border-l-4 transition-all hover:transform hover:-translate-y-1 hover:shadow-lg cursor-pointer"
          style={{
            borderLeftColor: `var(--color-${formatClassNameForCSS(
              specData.class
            )})`,
          }}
        >
          <div className="flex items-center">
            {/* Spec icon */}
            {specIconUrl ? (
              <Image
                src={specIconUrl}
                alt={`${specData.spec} icon`}
                className="rounded mr-3"
                width={42}
                height={42}
                unoptimized // Add this for external images
              />
            ) : (
              <div className="w-8 h-8 bg-gray-700 rounded mr-3"></div>
            )}

            <div>
              <h3 className={`font-bold ${classColorClass}`}>
                {formatSpecName(specData.class, specData.spec)}
              </h3>
              <p className="text-xs text-gray-400 capitalize">
                {role?.toLowerCase()}
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
