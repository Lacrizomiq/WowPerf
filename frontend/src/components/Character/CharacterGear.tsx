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
    .map((slot) => {
      const item = gear.items[slot];
      return item ? { ...item, slot } : null;
    })
    .filter(Boolean) as GearItem[];

  return (
    <div className="p-4 bg-gradient-dark">
      <style jsx global>{`
        .wowhead-tooltip {
          scale: 0.8;
          transform-origin: top left;
        }
      `}</style>
      <h2 className="text-3xl font-bold text-gradient-glow mb-4">Gear</h2>
      <p className="text-blue-200 mb-4">
        Item Level: {gear.item_level_equipped} (Equipped)
      </p>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {orderedGear.map((item) => (
          <div
            key={item.slot}
            className="bg-deep-blue bg-opacity-50 rounded-lg overflow-hidden shadow-lg hover:scale-105 transition duration-300 glow-effect p-4"
          >
            <div className="flex items-center">
              <div className="relative">
                <a
                  href={`https://www.wowhead.com/item=${item.item_id}`}
                  data-wowhead={`item=${item.item_id}&ilvl=${item.item_level}`}
                  className="block cursor-pointer"
                >
                  <Image
                    src={`https://wow.zamimg.com/images/wow/icons/large/${item.icon}.jpg`}
                    alt={item.name}
                    width={56}
                    height={56}
                    className="rounded-md mr-4"
                  />
                </a>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-gradient-glow">
                  {item.name}
                </h3>
                <p className="text-blue-200">
                  {item.slot} - iLvl: {item.item_level}
                </p>
                {item.enchant && (
                  <p className="text-green-300">Enchant: {item.enchant}</p>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
