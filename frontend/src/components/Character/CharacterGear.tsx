import React, { useEffect } from "react";
import { useGetBlizzardCharacterEquipment } from "@/hooks/useBlizzardApi";
import { useGetBlizzardCharacterProfile } from "@/hooks/useBlizzardApi";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import Image from "next/image";
/// <reference path="../../global.d.ts" />

declare global {
  interface Window {
    $WowheadPower?: {
      refreshLinks: () => void;
    };
    whTooltips?: {
      renameLinks: boolean;
      iconSize: string;
      hideSpecs: string;
    };
  }
}

interface CharacterGearProps {
  region: string;
  realm: string;
  name: string;
  namespace: string;
  locale: string;
}

interface EquipmentData {
  item_level_equipped: number;
  item_level_total: number;
  items: {
    [key: string]: {
      item_id: number;
      item_level: number;
      item_quality: number;
      icon_name: string;
      icon_url: string;
      name: string;
      slot: string;
      enchant?: number;
      enchant_name?: string;
      stats?: ItemStat[];
      gems?: number[];
      bonuses?: number[];
    };
  };
}

interface ItemStat {
  Type: string;
  Value: number;
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
  "finger_1",
  "finger_2",
  "trinket_1",
  "trinket_2",
  "main_hand",
  "off_hand",
];

export default function CharacterGear({
  region,
  realm,
  name,
  namespace,
  locale,
}: CharacterGearProps) {
  const {
    data: characterData,
    isLoading,
    error,
  } = useGetBlizzardCharacterEquipment(region, realm, name, namespace, locale);
  useWowheadTooltips();

  const { data: characterProfile } = useGetBlizzardCharacterProfile(
    region,
    realm,
    name,
    namespace,
    locale
  );

  console.log(characterData);

  useEffect(() => {
    if (
      characterData &&
      typeof window !== "undefined" &&
      window.$WowheadPower
    ) {
      window.$WowheadPower.refreshLinks();
    }
  }, [characterData]);

  if (isLoading) return <div className="text-white">Loading gear data...</div>;
  if (error)
    return (
      <div className="text-red-500">
        Error loading gear data:{" "}
        {error instanceof Error ? error.message : "Unknown error"}
      </div>
    );
  if (!characterData || !characterData.items)
    return <div className="text-yellow-500">No gear data found</div>;

  const orderedGear = gearOrder
    .map((slot) =>
      characterData.items[slot] ? { ...characterData.items[slot], slot } : null
    )
    .filter(Boolean) as (EquipmentData["items"][string] & { slot: string })[];

  const getWowheadParams = (
    item: EquipmentData["items"][string],
    slot: string
  ) => {
    let params = `item=${item.item_id}&ilvl=${item.item_level}`;
    if (item.bonuses?.length) params += `&bonus=${item.bonuses.join(":")}`;
    if (item.gems?.length) params += `&gems=${item.gems.join(":")}`;
    if (item.enchant) params += `&ench=${item.enchant}`;
    return params;
  };

  const getItemQualityClass = (quality: number) => {
    switch (quality) {
      case 4:
        return "item-quality--4";
      case 3:
        return "item-quality--3";
      case 2:
        return "item-quality--2";
      default:
        return "";
    }
  };

  function ItemDisplay({
    item,
  }: {
    item: EquipmentData["items"][string] & { slot: string };
  }) {
    return (
      <div className="relative flex items-center bg-deep-blue border-2  p-2 h-20 w-full">
        <a
          href={`https://www.wowhead.com/item=${item.item_id}`}
          data-wowhead={getWowheadParams(item, item.slot)}
          className="block cursor-pointer flex-shrink-0"
          data-wh-icon-size="medium"
        >
          <Image
            src={item.icon_url}
            alt={item.name}
            width={48}
            height={48}
            className="rounded-md border-2 border-gray-700"
          />
        </a>
        <div className="ml-2 text-white text-sm flex-grow overflow-hidden">
          <span
            className={`${getItemQualityClass(
              item.item_quality
            )} truncate block`}
          >
            {item.name}
          </span>
          <div className="truncate">
            <span>{item.enchant_name}</span>
          </div>
          <div className="text-xs">{item.item_level}</div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 bg-[#002440] rounded-xl shadow-lg m-4">
      <style jsx global>{`
        .wowhead-tooltip {
          scale: 1.2;
          transform-origin: top left;
          max-width: 300px;
          font-size: 14px;
        }
      `}</style>

      <div className="text-blue-200 mb-4">
        {characterData.item_level_equipped.toFixed(1)} item lvl (Equipped)
      </div>

      {/* Armor */}
      <h2 className="text-xl font-bold text-white mb-4">Armor</h2>
      <div className="flex flex-col md:flex-row space-x-0 md:space-x-4 space-y-4 md:space-y-0 mb-8">
        <div className="flex flex-col w-full md:w-1/2">
          {orderedGear.slice(0, 5).map((item) => (
            <ItemDisplay key={item.slot} item={item} />
          ))}
        </div>
        <div className="flex flex-col w-full md:w-1/2">
          {orderedGear.slice(5, 10).map((item) => (
            <ItemDisplay key={item.slot} item={item} />
          ))}
        </div>
      </div>

      {/* Fingers */}
      <h2 className="text-xl font-bold text-white mb-4">Fingers</h2>
      <div className="flex flex-col md:flex-row space-x-0 md:space-x-4 space-y-4 md:space-y-0 mb-8">
        {orderedGear.slice(10, 12).map((item) => (
          <ItemDisplay key={item.slot} item={item} />
        ))}
      </div>

      {/* Trinkets */}
      <h2 className="text-xl font-bold text-white mb-4">Trinkets</h2>
      <div className="flex flex-col md:flex-row space-x-0 md:space-x-4 space-y-4 md:space-y-0 mb-8">
        {orderedGear.slice(12, 14).map((item) => (
          <ItemDisplay key={item.slot} item={item} />
        ))}
      </div>

      {/* Weapon */}
      <h2 className="text-xl font-bold text-white mb-4">Weapon</h2>
      <div className="flex flex-col md:flex-row space-x-0 md:space-x-4 space-y-4 md:space-y-0">
        {orderedGear.slice(14).map((item) => (
          <ItemDisplay key={item.slot} item={item} />
        ))}
      </div>
    </div>
  );
}
