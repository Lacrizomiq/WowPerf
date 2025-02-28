// CharacterCard.tsx
import React from "react";
import Image from "next/image";
import { getClassIcon } from "@/utils/classandspecicons";
import { useRouter } from "next/navigation";
import {
  useGetBlizzardCharacterProfile,
  useGetBlizzardCharacterEquipment,
  useGetBlizzardCharacterMythicPlusBestRuns,
} from "@/hooks/useBlizzardApi";

interface CharacterCardProps {
  character: {
    name: string;
    playable_class: { name: string; id: number };
    level: number;
    realm: { name: string; slug: string };
  };
  region?: string;
  showTooltip?: boolean;
  isMainCard?: boolean;
}

const CharacterCard: React.FC<CharacterCardProps> = ({
  character,
  region,
  showTooltip = false,
  isMainCard = false,
}) => {
  const router = useRouter();
  const { name, playable_class, level, realm } = character;
  const normalized = playable_class.name.replace(/\s+/g, "");
  const classIcon = getClassIcon(normalized);
  const safeRegion = region || "eu";
  const characterName = name.toLowerCase();
  const realmSlug = realm.slug;

  // Dynamically build namespace based on the region
  const staticNamespace = `static-${safeRegion}`;
  const dynamicNamespace = `dynamic-${safeRegion}`;
  const profileNamespace = `profile-${safeRegion}`;

  // Fetch character equipment data
  const { data: equipmentData } = useGetBlizzardCharacterEquipment(
    safeRegion,
    realmSlug,
    characterName,
    staticNamespace,
    "en_GB"
  );

  // Fetch mythic+ data
  const { data: mythicPlusData } = useGetBlizzardCharacterMythicPlusBestRuns(
    safeRegion,
    realmSlug,
    characterName,
    dynamicNamespace,
    "en_GB",
    "13"
  );

  // Get item level and round it
  const itemLevel = equipmentData?.data?.item_level_equipped
    ? Math.round(equipmentData.data.item_level_equipped)
    : "-";

  // Calculate M+ score
  const mPlusScore = mythicPlusData?.OverallMythicRating || "-";

  // Fetch additional character profile information
  const { data: characterData } = useGetBlizzardCharacterProfile(
    safeRegion,
    realmSlug,
    characterName,
    profileNamespace,
    "en_GB"
  );

  const handleCharacterClick = () => {
    router.push(`/character/${safeRegion}/${realmSlug}/${characterName}`);
  };

  // Normalize class name for CSS classes
  const className = playable_class.name.toLowerCase().replace(/\s+/g, "-");

  // Create appropriate color class based on class name
  const colorClass = `class-color--${className}`;

  // Create a background style with gradient based on class color
  const getBgStyle = () => {
    // Use the CSS variable corresponding to the class
    const colorVar = `var(--color-${className})`;
    return {
      background: `linear-gradient(135deg, ${colorVar} 0%, rgba(0, 0, 0, 0.7) 100%)`,
    };
  };

  return (
    <div
      className={`relative rounded-lg overflow-hidden cursor-pointer transition-all duration-300 ${
        isMainCard ? "h-48" : "h-56"
      } hover:-translate-y-1 hover:shadow-xl`}
      onClick={handleCharacterClick}
    >
      {/* Class-colored background */}
      <div className="absolute inset-0" style={getBgStyle()}></div>

      {/* Card content */}
      <div className="relative z-10 flex flex-col justify-end h-full p-4">
        <div className="flex items-center justify-between mb-3">
          <div
            className={`flex items-center gap-2 ${colorClass} font-bold text-lg`}
          >
            <Image
              src={classIcon}
              alt={playable_class.name}
              width={24}
              height={24}
              className="rounded-full"
            />
            {name}
          </div>
          <div className="flex items-center gap-2 bg-black/60 px-3 py-1 rounded-full text-xs">
            <span className="h-2 w-2 rounded-full bg-green-500"></span>
            Level {level}
          </div>
        </div>

        <div className="text-white/90 text-sm">
          <div>
            {playable_class.name} ({characterData?.data?.active_spec_name || ""}
            )
          </div>
          <div className="text-white/70">
            {realm.name} ({safeRegion.toUpperCase()})
          </div>
        </div>

        <div className="flex gap-2 mt-3">
          <div className="bg-black/60 px-3 py-1 rounded-md text-center">
            <div className="font-bold">{itemLevel}</div>
            <div className="text-xs text-white/70">iLvL</div>
          </div>
          <div className="bg-black/60 px-3 py-1 rounded-md text-center">
            <div className="font-bold">{mPlusScore}</div>
            <div className="text-xs text-white/70">M+</div>
          </div>
        </div>
      </div>

      {/* Hover tooltip */}
      {showTooltip && (
        <div className="absolute inset-0 bg-black/90 flex items-center justify-center opacity-0 transition-opacity duration-300 hover:opacity-100 z-20">
          <button
            className="bg-blue-500 hover:bg-blue-600 px-4 py-2 rounded-md font-semibold text-white"
            onClick={handleCharacterClick}
          >
            View Profile
          </button>
        </div>
      )}
    </div>
  );
};

export default CharacterCard;
