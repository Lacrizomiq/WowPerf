import React, { useEffect } from "react";
import { useGetBlizzardCharacterEquipment } from "@/hooks/useBlizzardApi";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import Image from "next/image";

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

  useEffect(() => {
    if (characterData && window.$WowheadPower) {
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

  console.log("Received character data:", characterData);

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

  return (
    <div className="p-4 bg-gradient-dark shadow-lg glow-effect m-12">
      <style jsx global>{`
        .wowhead-tooltip {
          scale: 1.2;
          transform-origin: top left;
          max-width: 300px;
          font-size: 14px;
        }
      `}</style>
      <div className="flex align-center items-center justify-between mb-2 px-4">
        <h2 className="text-2xl font-bold text-gradient-glow mb-6">Gear</h2>
        <p className="text-blue-200 mb-4 ml-2">
          {characterData.item_level_equipped} item lvl (Equipped)
        </p>
      </div>
      <div className="flex flex-wrap gap-2 justify-center  ">
        {orderedGear.map((item) => (
          <div key={item.slot} className="relative">
            <a
              href={`https://www.wowhead.com/item=${item.item_id}`}
              data-wowhead={getWowheadParams(item, item.slot)}
              className="block cursor-pointer"
              data-wh-icon-size="medium"
            >
              <Image
                src={`https://wow.zamimg.com/images/wow/icons/large/${item.icon_name}.jpg`}
                alt={item.name}
                width={56}
                height={56}
                className="rounded-md border-2 border-gray-700"
              />
              <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-xs px-1 rounded">
                {item.item_level}
              </div>
            </a>
          </div>
        ))}
      </div>
    </div>
  );
}
