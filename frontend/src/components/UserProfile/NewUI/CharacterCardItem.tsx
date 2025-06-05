// CharacterCardItem.tsx
import React, { useState } from "react";
import Image from "next/image";
import { getClassIcon } from "@/utils/classandspecicons";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { EnrichedUserCharacter } from "@/types/character/character";
import {
  isCharacterRecentlyUpdated,
  getCharacterAvatarUrl,
} from "@/utils/character/character";

interface CharacterCardItemProps {
  character: EnrichedUserCharacter;
  region: string;
  onToggleDisplay: (display: boolean) => void;
  onSetFavorite: () => void;
  isTogglingDisplay?: boolean;
  isSettingFavorite?: boolean;
  isFavorite?: boolean;
}

const CharacterCardItem: React.FC<CharacterCardItemProps> = ({
  character,
  region,
  onToggleDisplay,
  onSetFavorite,
  isTogglingDisplay = false,
  isSettingFavorite = false,
  isFavorite = false,
}) => {
  const router = useRouter();
  const [isHovered, setIsHovered] = useState(false);

  // Extract character data from enriched structure FIRST
  const {
    name,
    class: className,
    level,
    realm,
    active_spec_name,
    achievement_points,
    avatar_url,
    equipment_json,
    mythic_plus_json,
    last_api_update,
    is_displayed,
  } = character;

  // Check if character is enriched (has additional data beyond basic sync)
  const isEnriched = !!(active_spec_name || achievement_points || avatar_url);

  // Calculate derived values
  const normalizedClass = className.replace(/\s+/g, "");
  const classIcon = getClassIcon(normalizedClass);
  const safeRegion = region || "eu";
  const characterName = name.toLowerCase();
  const realmSlug = realm;

  // Get avatar URL with fallback (always use class icon if no enriched avatar)
  const avatarUrl = avatar_url || classIcon;

  // Extract item level from equipment JSON or fallback to 0
  const itemLevel =
    equipment_json?.average_item_level ||
    equipment_json?.item_level_equipped ||
    0;

  // Extract M+ rating from mythic plus JSON or fallback to 0
  const mythicPlusRating =
    mythic_plus_json?.current_rating || mythic_plus_json?.overall_rating || 0;

  // Check if character data is fresh
  const isRecentlyUpdated = isCharacterRecentlyUpdated(character, 24);

  const handleCharacterClick = () => {
    router.push(`/character/${safeRegion}/${realmSlug}/${characterName}`);
  };

  // Normalize the class name for CSS classes
  const cssClassName = className.toLowerCase().replace(/\s+/g, "-");

  // Create a background style with a gradient based on the class color
  const getBgStyle = () => {
    // Use the CSS variable corresponding to the class
    const colorVar = `var(--color-${cssClassName})`;
    return {
      background: `linear-gradient(135deg, ${colorVar} 0%, rgba(0, 0, 0, 0.7) 100%)`,
    };
  };

  return (
    <div
      className={`relative rounded-lg overflow-hidden h-56 ${
        isFavorite ? "ring-3 ring-yellow-400 border-2 border-yellow-500" : ""
      }`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Colored background according to the class */}
      <div className="absolute inset-0" style={getBgStyle()}></div>

      {/* Data status indicators - Positioned at bottom to avoid conflicts */}
      {(!isEnriched || (isEnriched && !isRecentlyUpdated)) && (
        <div className="absolute bottom-2 left-2 z-10 flex gap-2">
          {/* Basic data indicator */}
          {!isEnriched && (
            <div className="bg-yellow-500/90 text-white text-xs px-2 py-1 rounded-full flex items-center gap-1 backdrop-blur-sm">
              <svg
                className="w-3 h-3"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              Basic
            </div>
          )}

          {/* Data freshness indicator */}
          {isEnriched && !isRecentlyUpdated && (
            <div className="bg-orange-500/90 text-white text-xs px-2 py-1 rounded-full flex items-center gap-1 backdrop-blur-sm">
              <svg
                className="w-3 h-3"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              Outdated
            </div>
          )}
        </div>
      )}

      {/* Card content */}
      <div className="relative z-10 flex flex-col h-full p-4">
        <div className="flex items-center justify-between mb-3">
          <div className={`flex items-center gap-2 font-bold text-lg`}>
            {/* Always use class icon as fallback, show enriched avatar only if available */}
            <Image
              src={avatar_url || classIcon}
              alt={avatar_url ? `${name} avatar` : className}
              width={38}
              height={38}
              className="rounded-full"
              onError={(e) => {
                // Always fallback to class icon if any image fails to load
                e.currentTarget.src = classIcon;
              }}
            />
            {name}
          </div>
          <div className="flex items-center gap-2">
            <div className="bg-black/60 px-3 py-1 rounded-full text-xs flex items-center">
              <span className="h-2 w-2 rounded-full bg-green-500 mr-2"></span>
              Level {level}
              {/* Favorite character indicator */}
              {isFavorite && (
                <span
                  className="ml-2 text-yellow-400 text-xl"
                  title="Favorite character"
                >
                  ★
                </span>
              )}
            </div>
          </div>
        </div>

        <div className="text-white/90 text-sm">
          <div>
            {active_spec_name ? `${active_spec_name} ${className}` : className}
          </div>
          <div className="text-white/70">
            {realm.charAt(0).toUpperCase() + realm.slice(1).toLowerCase()} -{" "}
            {safeRegion.toUpperCase()}
          </div>

          {!isEnriched && (
            <div className="text-yellow-300 text-xs">
              Basic data only - refresh to enrich
            </div>
          )}
        </div>

        <div className="flex gap-2 mt-3">
          {itemLevel > 0 && (
            <div className="bg-black/60 px-3 py-1 rounded-md text-center">
              <div className="font-bold">{Math.round(itemLevel)}</div>
              <div className="text-xs text-white/70">iLvL</div>
            </div>
          )}
          {mythicPlusRating > 0 && (
            <div className="bg-black/60 px-3 py-1 rounded-md text-center">
              <div className="font-bold">{Math.round(mythicPlusRating)}</div>
              <div className="text-xs text-white/70">M+</div>
            </div>
          )}
        </div>

        {/* Data freshness info */}
        {last_api_update && (
          <div className="mt-auto">
            <div className="text-xs text-white/50">
              {/* Updated: {new Date(last_api_update).toLocaleDateString()} */}
            </div>
          </div>
        )}

        {/* Overlay on hover with all buttons */}
        {isHovered && (
          <div className="absolute inset-0 bg-black/70 flex flex-col items-center justify-center gap-3 transition-opacity duration-200 z-20">
            <Button
              size="sm"
              variant="default"
              className="bg-blue-500 hover:bg-blue-600 w-36"
              onClick={handleCharacterClick}
            >
              View Profile
            </Button>

            <Button
              size="sm"
              variant={isFavorite ? "default" : "outline"}
              className={
                isFavorite ? "bg-yellow-600 hover:bg-yellow-700 w-36" : "w-36"
              }
              onClick={(e) => {
                e.stopPropagation();
                onSetFavorite();
              }}
              disabled={isSettingFavorite}
            >
              {isFavorite ? "★ Favorite" : "Set as Favorite"}
            </Button>

            {/* Button to display/hide the character */}
            <Button
              size="sm"
              variant={is_displayed ? "default" : "outline"}
              className={
                is_displayed ? "bg-green-600 hover:bg-green-700 w-36" : "w-36"
              }
              onClick={(e) => {
                e.stopPropagation();
                onToggleDisplay(!is_displayed);
              }}
              disabled={isTogglingDisplay}
            >
              {is_displayed ? "Hide Character" : "Show Character"}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};

export default CharacterCardItem;
