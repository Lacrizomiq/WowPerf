// CharacterCardItem.tsx
import React, { useState } from "react";
import Image from "next/image";
import { getClassIcon } from "@/utils/classandspecicons";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { UserCharacter } from "@/types/userCharacter/userCharacter";

interface CharacterCardItemProps {
  character: UserCharacter;
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
  const { name, class: className, level, realm, item_level } = character;
  const normalizedClass = className.replace(/\s+/g, "");
  const classIcon = getClassIcon(normalizedClass);
  const safeRegion = region || "eu";
  const characterName = name.toLowerCase();
  const realmSlug = realm;

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

      {/* Card content */}
      <div className="relative z-10 flex flex-col h-full p-4">
        <div className="flex items-center justify-between mb-3">
          <div className={`flex items-center gap-2 font-bold text-lg`}>
            <Image
              src={classIcon}
              alt={className}
              width={38}
              height={38}
              className="rounded-full"
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
            {className} ({character.active_spec_name || ""})
          </div>
          <div className="text-white/70">
            {realm} ({safeRegion.toUpperCase()})
          </div>
        </div>

        <div className="flex gap-2 mt-3">
          {item_level > 0 && (
            <div className="bg-black/60 px-3 py-1 rounded-md text-center">
              <div className="font-bold">{Math.round(item_level)}</div>
              <div className="text-xs text-white/70">iLvL</div>
            </div>
          )}
          {character.mythic_plus_rating > 0 && (
            <div className="bg-black/60 px-3 py-1 rounded-md text-center">
              <div className="font-bold">
                {Math.round(character.mythic_plus_rating)}
              </div>
              <div className="text-xs text-white/70">M+</div>
            </div>
          )}
        </div>

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
              variant={character.is_displayed ? "default" : "outline"}
              className={
                character.is_displayed
                  ? "bg-green-600 hover:bg-green-700 w-36"
                  : "w-36"
              }
              onClick={(e) => {
                e.stopPropagation();
                onToggleDisplay(!character.is_displayed);
              }}
              disabled={isTogglingDisplay}
            >
              {character.is_displayed ? "Hide Character" : "Show Character"}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};

export default CharacterCardItem;
