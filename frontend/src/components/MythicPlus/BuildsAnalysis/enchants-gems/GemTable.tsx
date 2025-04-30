// components/MythicPlus/BuildsAnalysis/enchants-gems/GemTable.tsx
"use client";

import Image from "next/image";
import { GemUsage } from "@/types/warcraftlogs/builds/buildsAnalysis";
import { getGemIconUrl } from "@/utils/buildsAnalysis/gemsData";

interface GemTableProps {
  slotName: string;
  gems: GemUsage[];
}

export default function GemTable({ slotName, gems }: GemTableProps) {
  if (!gems || gems.length === 0) return null;

  // Sort gems by rank
  const sortedGems = [...gems].sort((a, b) => a.rank - b.rank);

  // Calculate total count to derive usage percentage
  const totalCount = sortedGems.reduce((acc, gem) => acc + gem.usage_count, 0);

  return (
    <div className="w-full rounded-lg border border-slate-700 overflow-hidden bg-slate-900">
      <div className="overflow-x-auto">
        <table className="w-full min-w-[400px]">
          <thead className="bg-slate-800">
            <tr>
              <th className="py-2 px-3 text-left text-sm font-semibold text-white">
                {slotName} Gems
              </th>
              <th className="py-2 px-2 text-center text-sm font-semibold text-white w-20">
                Avg Key
              </th>
              <th className="py-2 px-2 text-right text-sm font-semibold text-white w-24">
                Usage
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700">
            {sortedGems.map((gem) => {
              const usagePercentage = (gem.usage_count / totalCount) * 100;

              return (
                <tr
                  key={`${gem.item_slot}-${gem.gem_ids_array.join("-")}-${
                    gem.rank
                  }`}
                  className="hover:bg-slate-800"
                >
                  <td className="py-2 px-3">
                    <div className="flex items-center">
                      {/* Gem Icons - Using utility function for consistency */}
                      <div className="flex gap-1 mr-3">
                        {gem.gem_ids_array.map((gemId, idx) => (
                          <a
                            key={`gem-${gem.item_slot}-${gem.rank}-${idx}`}
                            href={`https://www.wowhead.com/item=${gemId}`}
                            data-wowhead={`item=${gemId}`}
                            className="block"
                          >
                            <div className="w-8 h-8 border border-slate-600 rounded-full overflow-hidden bg-slate-800">
                              <Image
                                src={getGemIconUrl(gemId)}
                                alt={`Gem ${gemId}`}
                                width={32}
                                height={32}
                                className="rounded-full"
                                unoptimized
                              />
                            </div>
                          </a>
                        ))}
                      </div>

                      {/* Gem Count & Rank Info */}
                      <div>
                        <div
                          className={`text-${
                            gem.rank === 1 ? "indigo" : "slate"
                          }-${gem.rank === 1 ? "400" : "300"} ${
                            gem.rank === 1 ? "font-medium" : ""
                          }`}
                        >
                          {gem.gems_count}{" "}
                          {gem.gems_count === 1 ? "Gem" : "Gems"} Socket
                          {gem.rank === 1 && (
                            <span className="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-indigo-900 text-indigo-300">
                              Best
                            </span>
                          )}
                        </div>
                        <div className="text-xs text-slate-500 mt-0.5">
                          Item Level: {Math.round(gem.avg_item_level)}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="py-2 px-2 text-center text-white">
                    {Math.round(gem.avg_keystone_level * 10) / 10}
                  </td>
                  <td className="py-2 px-2 text-right">
                    <div className="inline-block bg-slate-800 rounded-full px-2 py-1 text-sm font-medium text-white">
                      {Math.round(usagePercentage)}%
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
