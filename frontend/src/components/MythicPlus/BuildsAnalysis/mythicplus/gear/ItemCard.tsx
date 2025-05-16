// ItemCard.tsx - Version harmonisée
import Image from "next/image";
import {
  GlobalPopularItem,
  PopularItem,
} from "@/types/warcraftlogs/builds/buildsAnalysis";
import {
  getItemQualityClass,
  getItemIconUrl,
} from "@/utils/buildsAnalysis/dataTransformer";

interface ItemCardProps {
  item: PopularItem | GlobalPopularItem;
}

export default function ItemCard({ item }: ItemCardProps) {
  // Build the parameters for the Wowhead tooltip
  const wowheadParams = `item=${item.item_id}&ilvl=${item.item_level}`;

  // Get the CSS class for the item quality
  const qualityClass = `item-quality--${item.item_quality}`;

  // Determine the border class for the icon
  const getBorderClass = () => {
    switch (item.item_quality) {
      case 4:
        return "border-purple-500"; // Epic
      case 3:
        return "border-blue-500"; // Rare
      case 2:
        return "border-green-500"; // Uncommon
      default:
        return "border-gray-500"; // Common
    }
  };

  return (
    <div className="bg-slate-800/50 border border-slate-700 rounded-lg p-3 flex items-center">
      {/* Item icon with colored border */}
      <a
        href={`https://www.wowhead.com/item=${item.item_id}`}
        data-wowhead={wowheadParams}
        className="flex-shrink-0 block"
      >
        <div
          className={`rounded-md border-2 ${getBorderClass()} overflow-hidden`}
        >
          <Image
            src={getItemIconUrl(item.item_icon)}
            alt={item.item_name}
            width={40}
            height={40}
            className="rounded"
            unoptimized
          />
        </div>
      </a>

      {/* Item information */}
      <div className="ml-3 flex-grow">
        {/* Name with color based on quality */}
        <div className={`${qualityClass} font-semibold truncate`}>
          {item.item_name}
        </div>

        {/* Item level */}
        <div className="text-sm text-gray-300">ilvl {item.item_level}</div>

        {/* Usage statistics */}
        <div className="flex items-center mt-1 text-xs">
          <div className="mr-4">
            <span className="text-gray-400">Usage: </span>
            <span className="text-white">
              {Math.round(item.usage_percentage)}%
            </span>
          </div>
          <div>
            <span className="text-gray-400">Avg Key: </span>
            <span className="text-white">
              +{item.avg_keystone_level.toFixed(1)}
            </span>
          </div>
        </div>
      </div>

      {/* Rank indicator (optional) */}
      {item.rank === 1 && (
        <div className="absolute -top-2 -right-2 bg-purple-600 text-white text-xs px-2 py-1 rounded-full">
          Best
        </div>
      )}
    </div>
  );
}
