// CharacterCard.tsx
import React from "react";
import Image from "next/image";
import { getClassIcon } from "@/utils/classandspecicons";
import { useRouter } from "next/navigation";
import { EnrichedUserCharacter } from "@/types/character/character";
import { getCharacterAvatarUrl } from "@/utils/character/character";

interface CharacterCardProps {
  character: EnrichedUserCharacter;
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

  // Extract data from enriched character (from BDD)
  const {
    name,
    class: className,
    level,
    realm,
    active_spec_name,
    avatar_url,
    equipment_json,
    mythic_plus_json,
  } = character;

  const normalizedClass = className.replace(/\s+/g, "");
  const classIcon = getClassIcon(normalizedClass);
  const safeRegion = region || character.region || "eu";
  const characterName = name.toLowerCase();
  const realmSlug = realm;

  // Get enriched data from our database instead of API calls
  const avatarUrl = getCharacterAvatarUrl(character);

  // Extract item level from equipment JSON (already enriched)
  const itemLevel =
    equipment_json?.average_item_level ||
    equipment_json?.item_level_equipped ||
    0;

  // Extract M+ score from mythic plus JSON (already enriched)
  const mPlusScore =
    mythic_plus_json?.current_rating || mythic_plus_json?.overall_rating || 0;

  const handleCharacterClick = () => {
    router.push(`/character/${safeRegion}/${realmSlug}/${characterName}`);
  };

  // Normalize class name for CSS classes
  const cssClassName = className.toLowerCase().replace(/\s+/g, "-");

  // Create a background style with gradient based on class color
  const getBgStyle = () => {
    // Use the CSS variable corresponding to the class
    const colorVar = `var(--color-${cssClassName})`;
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
            className={`flex items-center gap-2 font-bold text-lg text-white`}
          >
            {/* Use enriched avatar if available, fallback to class icon */}
            {avatar_url ? (
              <Image
                src={avatar_url}
                alt={`${name} avatar`}
                width={38}
                height={38}
                className="rounded-full"
                onError={(e) => {
                  // Fallback to class icon if avatar fails to load
                  e.currentTarget.src = classIcon;
                }}
              />
            ) : (
              <Image
                src={classIcon}
                alt={className}
                width={38}
                height={38}
                className="rounded-full"
              />
            )}
            {name}
          </div>
          <div className="flex items-center gap-2 bg-black/60 px-3 py-1 rounded-full text-xs">
            <span className="h-2 w-2 rounded-full bg-green-500"></span>
            Level {level}
          </div>
        </div>

        <div className="text-white/90 text-sm">
          <div>
            {/* Use enriched spec name from our database */}
            {active_spec_name || className}
          </div>
          <div className="text-white/70">
            {realm} ({safeRegion.toUpperCase()})
          </div>
        </div>

        <div className="flex gap-2 mt-3">
          {/* Item level from enriched data */}
          <div className="bg-black/60 px-3 py-1 rounded-md text-center">
            <div className="font-bold">
              {itemLevel > 0 ? Math.round(itemLevel) : "-"}
            </div>
            <div className="text-xs text-white/70">iLvL</div>
          </div>

          {/* M+ score from enriched data */}
          <div className="bg-black/60 px-3 py-1 rounded-md text-center">
            <div className="font-bold">
              {mPlusScore > 0 ? Math.round(mPlusScore) : "0"}
            </div>
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
