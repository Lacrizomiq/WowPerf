import React, { useEffect } from "react";
import { useGetRaiderIoCharacterGear } from "@/hooks/useRaiderioApi";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import Image from "next/image";

interface CharacterGearProps {
  region: string;
  realm: string;
  name: string;
}

interface GearItem {
  item_id: number;
  item_level: number;
  item_quality: number;
  icon: string;
  name: string;
  enchant?: number;
  slot: string;
  bonuses?: number[];
  gems?: number[];
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
];

export default function CharacterGear({
  region,
  realm,
  name,
}: CharacterGearProps) {
  const {
    data: characterData,
    isLoading,
    error,
  } = useGetRaiderIoCharacterGear(region, realm, name);
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
  if (!characterData || !characterData.gear)
    return <div className="text-yellow-500">No gear data found</div>;

  const { gear } = characterData;

  const orderedGear = gearOrder
    .map((slot) => (gear.items[slot] ? { ...gear.items[slot], slot } : null))
    .filter(Boolean) as GearItem[];

  const getWowheadParams = (item: GearItem) => {
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
      <div className="flex align-center items-center justify-center">
        <h2 className="text-2xl font-bold text-gradient-glow mb-4 ">Gear</h2>
        <p className="text-blue-200 mb-4 ml-2">
          {gear.item_level_equipped} Item Level (Equipped)
        </p>
      </div>
      <div className="flex flex-wrap gap-2 justify-center  ">
        {orderedGear.map((item) => (
          <div key={item.slot} className="relative">
            <a
              href={`https://www.wowhead.com/item=${item.item_id}`}
              data-wowhead={getWowheadParams(item)}
              className="block cursor-pointer"
              data-wh-icon-size="medium"
            >
              <Image
                src={`https://wow.zamimg.com/images/wow/icons/large/${item.icon}.jpg`}
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
