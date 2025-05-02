// components/MythicPlus/BuildsAnalysis/enchants-gems/EnchantTable.tsx
"use client";

import { EnchantUsage } from "@/types/warcraftlogs/builds/buildsAnalysis";

interface EnchantTableProps {
  slotName: string;
  enchants: EnchantUsage[];
}

export default function EnchantTable({
  slotName,
  enchants,
}: EnchantTableProps) {
  if (!enchants || enchants.length === 0) return null;

  // Sort enchants by rank
  const sortedEnchants = [...enchants].sort((a, b) => a.rank - b.rank);

  // Calculate total count to derive usage percentage
  const totalCount = sortedEnchants.reduce(
    (acc, enchant) => acc + enchant.usage_count,
    0
  );

  return (
    <div className="w-full rounded-lg border border-slate-700 overflow-hidden bg-slate-900">
      <div className="overflow-x-auto">
        <table className="w-full min-w-[400px]">
          <thead className="bg-slate-800">
            <tr>
              <th className="py-2 px-3 text-left text-sm font-semibold text-white">
                {slotName} Enchants
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
            {sortedEnchants.map((enchant) => {
              const usagePercentage = (enchant.usage_count / totalCount) * 100;

              return (
                <tr
                  key={`${enchant.item_slot}-${enchant.permanent_enchant_id}-${enchant.rank}`}
                  className="hover:bg-slate-800"
                >
                  <td className="py-2 px-3">
                    <span
                      className={`text-${
                        enchant.rank === 1 ? "indigo" : "slate"
                      }-${enchant.rank === 1 ? "400" : "300"} ${
                        enchant.rank === 1 ? "font-medium" : ""
                      }`}
                    >
                      {enchant.permanent_enchant_name}
                      {enchant.rank === 1 && (
                        <span className="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-indigo-900 text-indigo-300">
                          Best
                        </span>
                      )}
                    </span>
                  </td>
                  <td className="py-2 px-2 text-center text-white">
                    {Math.round(enchant.avg_keystone_level * 10) / 10}
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
