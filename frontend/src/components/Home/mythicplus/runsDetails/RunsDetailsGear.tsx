import React from "react";
import Image from "next/image";
import { EquippedItem } from "@/types/runsDetails";

interface RunsDetailsGearProps {
  items: { [key: string]: EquippedItem };
}

const gearOrder = [
  "head",
  "neck",
  "shoulder",
  "back",
  "chest",
  "wrist",
  "hands",
  "waist",
  "legs",
  "feet",
  "finger1",
  "finger2",
  "trinket1",
  "trinket2",
  "mainhand",
  "offhand",
  "ranged",
];

const RunsDetailsGear: React.FC<RunsDetailsGearProps> = ({ items }) => {
  const getWowheadParams = (item: EquippedItem) => {
    let params = `item=${item.item_id}&ilvl=${item.item_level}`;
    if (item.bonuses?.length) params += `&bonus=${item.bonuses.join(":")}`;
    if (item.gems?.length) params += `&gems=${item.gems.join(":")}`;
    if (item.enchant) params += `&ench=${item.enchant}`;
    return params;
  };

  const getItemQualityClass = (quality: number) => {
    switch (quality) {
      case 4:
        return "text-purple-400";
      case 3:
        return "text-blue-400";
      case 2:
        return "text-green-400";
      default:
        return "text-white";
    }
  };

  return (
    <div className="flex flex-row gap-2">
      <div className="mt-4 flex flex-wrap flex-col gap-2 px-2 mb-2">
        {gearOrder.slice(0, 8).map((slot) => {
          const item = items[slot];
          if (!item) return null;
          return (
            <div key={slot} className="flex items-center space-x-2">
              <a
                href={`https://www.wowhead.com/item=${item.item_id}`}
                data-wowhead={getWowheadParams(item)}
              >
                <Image
                  src={`https://wow.zamimg.com/images/wow/icons/large/${item.icon}.jpg`}
                  alt={item.name}
                  width={32}
                  height={32}
                  className="rounded-md border border-gray-600"
                />
              </a>
              <span
                className={`text-sm ${getItemQualityClass(item.item_quality)}`}
              >
                {item.name} ({item.item_level})
              </span>
            </div>
          );
        })}
      </div>
      <div className="mt-4 flex flex-wrap flex-col gap-2 px-2 mb-2">
        {gearOrder.slice(8, 16).map((slot) => {
          const item = items[slot];
          if (!item) return null;
          return (
            <div key={slot} className="flex items-center space-x-2">
              <a
                href={`https://www.wowhead.com/item=${item.item_id}`}
                data-wowhead={getWowheadParams(item)}
              >
                <Image
                  src={`https://wow.zamimg.com/images/wow/icons/large/${item.icon}.jpg`}
                  alt={item.name}
                  width={32}
                  height={32}
                  className="rounded-md border border-gray-600"
                />
              </a>
              <span
                className={`text-sm ${getItemQualityClass(item.item_quality)}`}
              >
                {item.name} ({item.item_level})
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default RunsDetailsGear;
