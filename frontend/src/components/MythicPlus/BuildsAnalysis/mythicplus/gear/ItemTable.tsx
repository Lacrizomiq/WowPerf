// ItemTable.tsx - Version harmonisée
"use client";

import Image from "next/image";
import {
  GlobalPopularItem,
  PopularItem,
} from "@/types/warcraftlogs/builds/buildsAnalysis";
import { getItemIconUrl } from "@/utils/buildsAnalysis/dataTransformer";

interface ItemTableProps {
  slotName: string;
  items: (PopularItem | GlobalPopularItem)[];
}

export default function ItemTable({ slotName, items }: ItemTableProps) {
  if (!items || items.length === 0) return null;

  // Sort items by rank
  const sortedItems = [...items].sort((a, b) => a.rank - b.rank);

  return (
    <div className="w-full rounded-lg border border-slate-700 overflow-hidden bg-slate-800/30">
      <div className="overflow-x-auto">
        <table className="w-full min-w-[400px]">
          <thead className="bg-slate-800/50">
            <tr>
              <th className="py-2 px-3 text-left text-sm font-semibold text-white">
                {slotName}
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
            {sortedItems.map((item) => (
              <tr
                key={`${item.item_id}-${item.rank}`}
                className="hover:bg-slate-800/70"
              >
                <td className="py-2 px-3">
                  <a
                    href={`https://www.wowhead.com/item=${item.item_id}`}
                    data-wowhead={`item=${item.item_id}&ilvl=${item.item_level}`}
                    className="flex items-center"
                  >
                    <div className="mr-2 flex-shrink-0">
                      <Image
                        src={getItemIconUrl(item.item_icon)}
                        alt={item.item_name}
                        width={32}
                        height={32}
                        className="rounded"
                        unoptimized
                      />
                    </div>
                    <span
                      className={`item-quality--${item.item_quality} truncate`}
                    >
                      {item.item_name}
                    </span>
                  </a>
                </td>
                <td className="py-2 px-2 text-center text-white">
                  {Math.round(item.avg_keystone_level)}
                </td>
                <td className="py-2 px-2 text-right">
                  <div className="inline-block bg-slate-800/70 rounded-full px-2 py-1 text-sm font-medium text-white">
                    {Math.round(item.usage_percentage)}%
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
