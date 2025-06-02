// BuildCard.tsx - Version mise à jour avec les couleurs harmonisées
"use client";

import { useState } from "react";
import { TopTalentBuild } from "@/types/warcraftlogs/builds/buildsAnalysis";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import {
  formatDisplayClassName,
  formatDisplaySpecName,
} from "@/utils/classandspecicons";
import { Copy, Check } from "lucide-react";
import Link from "next/link";

interface BuildCardProps {
  build: TopTalentBuild;
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function BuildCard({ build, className, spec }: BuildCardProps) {
  const [isCopied, setIsCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(build.talent_import);
    setIsCopied(true);
    setTimeout(() => setIsCopied(false), 2000);
  };

  const displayClassName = formatDisplayClassName(className);
  const displaySpecName = formatDisplaySpecName(spec);

  return (
    <div className="bg-slate-800/30 rounded-lg border border-slate-700 overflow-hidden">
      <div className="p-6">
        {/* Header section with title and stats */}
        <div className="flex flex-col md:flex-row justify-between mb-6">
          {/* Title and Description */}
          <div className="mb-4 md:mb-0 md:pr-6 md:max-w-[50%]">
            <h3 className="text-2xl font-bold text-white mb-2">
              {displaySpecName} {displayClassName} Talents Build
            </h3>
            <p className="text-slate-400 text-sm">
              Most Popular {displaySpecName} {displayClassName} Talents Build
              for Mythic+ in The War Within Season 2.
            </p>
          </div>

          {/* Stats and Buttons */}
          <div className="flex flex-wrap gap-4 items-start justify-end">
            {/* Average Key */}
            <div className="bg-slate-800/50 rounded-md p-4 text-center w-[100px] h-[100px] flex flex-col items-center justify-center">
              <div className="text-2xl font-bold text-white">
                +{Math.round(build.avg_keystone_level)}
              </div>
              <div className="text-slate-400 text-xs uppercase mt-1">
                AVERAGE KEY
              </div>
            </div>

            {/* Talent Popularity */}
            <div className="bg-slate-800/50 rounded-md p-4 text-center w-[100px] h-[100px] flex flex-col items-center justify-center">
              <div className="text-2xl font-bold text-white">
                {Math.round(build.avg_usage_percentage)}%
              </div>
              <div className="text-slate-400 text-xs uppercase mt-1">
                TALENT POPULARITY
              </div>
            </div>

            {/* Buttons */}
            <div className="w-[100px] h-[100px] flex flex-col items-center justify-between">
              <button
                onClick={handleCopy}
                className={`w-full h-[48px] rounded-md transition-colors flex items-center justify-center ${
                  isCopied
                    ? "bg-green-600 text-white"
                    : "bg-purple-600 hover:bg-purple-700 text-white"
                }`}
              >
                {isCopied ? (
                  <div className="flex items-center gap-2">
                    <Check className="w-4 h-4" />
                    <p>Copy</p>
                  </div>
                ) : (
                  <div className="flex items-center gap-2">
                    <Copy className="w-4 h-4" />
                    <p>Copy</p>
                  </div>
                )}
              </button>

              <Link
                href={`/builds/mythic-plus/${className}/${spec}/talents`}
                className="w-full h-[48px] bg-slate-800/50 hover:bg-slate-700 text-white rounded-md flex items-center justify-center"
              >
                View More
              </Link>
            </div>
          </div>
        </div>

        {/* Talents iframe - displayed directly */}
        <div className="bg-black bg-opacity-30 rounded-lg p-2 border border-slate-700 shadow-xl">
          <iframe
            src={`https://www.raidbots.com/simbot/render/talents/${build.talent_import}?width=1000&level=80`}
            width="100%"
            height="600px"
            className="w-full"
          ></iframe>
        </div>
      </div>
    </div>
  );
}
