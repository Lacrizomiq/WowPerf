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
      gems?: number[];
      bonuses?: number[];
    };
  };
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

  return (
    <div className="p-4 bg-gradient-dark shadow-lg m-12 flex justify-center items-center glow-effect">
      <style jsx global>{`
        .wowhead-tooltip {
          scale: 1.2;
          transform-origin: top left;
          max-width: 300px;
          font-size: 14px;
        }
      `}</style>

      {/* Container for the gear layout */}
      <div className="flex flex-col items-center">
        <div className="flex justify-between w-full max-w-5xl">
          <div className="flex flex-col gap-2">
            {/* Left column (upper body gear) */}
            {orderedGear.slice(0, 8).map((item) => (
              <div key={item.slot} className="relative flex items-center">
                <a
                  href={`https://www.wowhead.com/item=${item.item_id}`}
                  data-wowhead={getWowheadParams(item, item.slot)}
                  className="block cursor-pointer"
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
                <div className="ml-2 text-white text-sm">
                  <span className={getItemQualityClass(item.item_quality)}>
                    {item.name}
                  </span>
                  <div className="text-xs">{item.item_level}</div>
                </div>
              </div>
            ))}
          </div>

          {/* Center column (character image) */}
          <div className="flex flex-col items-center justify-center">
            <Image
              src={characterProfile.main_raw_url}
              alt="Character"
              width={450}
              height={450}
              className=" scale-150"
              priority
            />
            <div className="text-blue-200 mt-4 text-center">
              {characterData.item_level_equipped} item lvl (Equipped)
            </div>
          </div>

          <div className="flex flex-col gap-2">
            {/* Right column (lower body gear) */}
            {orderedGear.slice(8).map((item) => (
              <div key={item.slot} className="relative flex items-center">
                <a
                  href={`https://www.wowhead.com/item=${item.item_id}`}
                  data-wowhead={getWowheadParams(item, item.slot)}
                  className="block cursor-pointer"
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
                <div className="ml-2 text-white text-sm">
                  <div>
                    <span className={getItemQualityClass(item.item_quality)}>
                      {item.name}
                    </span>
                  </div>
                  <div className="text-xs">{item.item_level}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
