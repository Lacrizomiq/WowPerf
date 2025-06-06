import React from "react";
import { useCharacters } from "@/hooks/useCharacters";
import { useUserProfile } from "@/hooks/useUserProfile";
import CharacterCardItem from "./CharacterCardItem";
import { Button } from "@/components/ui/button";

export const FavoriteCharacterSection: React.FC = () => {
  // Get the favorite character ID from user profile
  const { profile } = useUserProfile();
  const favoriteCharacterId = profile?.favorite_character_id;

  const {
    characters,
    isLoadingCharacters,
    actions,
    isLoading,
    ui,
    rateLimitState,
    region,
  } = useCharacters();

  // Loading state
  if (isLoadingCharacters) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500" />
      </div>
    );
  }

  // No characters found - show appropriate sync/refresh button
  if (!Array.isArray(characters) || characters.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-400 mb-4">
          You don&apos;t have any synchronized characters yet.
        </p>

        {/* Show sync button (for first time) */}
        <div className="space-y-3">
          <Button
            onClick={actions.syncAndEnrich}
            disabled={ui.isDisabled.sync}
            className="bg-blue-500 hover:bg-blue-600"
          >
            {isLoading.sync ? "Synchronizing..." : "Synchronize Characters"}
          </Button>

          {/* Rate limit message */}
          {ui.showRateLimit && (
            <div className="text-sm text-orange-400">
              ⏱️ {rateLimitState.message}
            </div>
          )}
        </div>
      </div>
    );
  }

  // Ensure characters is an array before using array methods
  const charactersArray = Array.isArray(characters) ? characters : [];

  // Find the favorite character by its ID
  const favoriteCharacter = favoriteCharacterId
    ? charactersArray.find((char) => char.id === favoriteCharacterId)
    : null;

  // If no favorite character found, use a displayed character as a fallback
  const displayCharacter =
    favoriteCharacter || charactersArray.find((char) => char.is_displayed);

  if (!displayCharacter) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-400">
          No favorite character set. Visit the &quot;Characters&quot; tab to set
          one.
        </p>
      </div>
    );
  }

  const isFavorite = favoriteCharacterId === displayCharacter.id;

  return (
    <div className="max-w-md mx-auto">
      <CharacterCardItem
        character={displayCharacter}
        region={region}
        onToggleDisplay={() => {
          /* Not used in favorite view */
        }}
        onSetFavorite={() => {
          /* Not used in favorite view */
        }}
        isTogglingDisplay={false}
        isSettingFavorite={false}
        isFavorite={isFavorite}
      />
    </div>
  );
};

export default FavoriteCharacterSection;
