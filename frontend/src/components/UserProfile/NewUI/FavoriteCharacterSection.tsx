import React from "react";
import { useWoWCharacters } from "@/hooks/useWowProtectedAccount";
import { useUserProfile } from "@/hooks/useUserProfile";
import CharacterCardItem from "./CharacterCardItem";
import { Button } from "@/components/ui/button";

export const FavoriteCharacterSection: React.FC = () => {
  // Get the favorite character ID from user profile
  const { profile } = useUserProfile();
  const favoriteCharacterId = profile?.favorite_character_id;

  const {
    userCharacters,
    isLoadingUserCharacters,
    syncCharacters,
    isSyncing,
    setFavoriteCharacter,
    isSettingFavorite,
    toggleCharacterDisplay,
    isTogglingDisplay,
    wowProfile,
  } = useWoWCharacters();

  // Loading state
  if (isLoadingUserCharacters) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500" />
      </div>
    );
  }

  // No characters or not connected to Battle.net
  if (!userCharacters || userCharacters.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-400 mb-4">
          You don&apos;t have any synchronized characters yet.
        </p>
        <Button
          onClick={() => syncCharacters()}
          disabled={isSyncing}
          className="bg-blue-500 hover:bg-blue-600"
        >
          {isSyncing ? "Synchronizing..." : "Synchronize characters"}
        </Button>
      </div>
    );
  }

  // Find the favorite character by its ID
  const favoriteCharacter = favoriteCharacterId
    ? userCharacters.find((char) => char.id === favoriteCharacterId)
    : null;

  // If no favorite character found, use a visible character as a fallback
  const displayCharacter =
    favoriteCharacter || userCharacters.find((char) => char.is_displayed);

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
        region={wowProfile?.region || "eu"}
        onToggleDisplay={(display) =>
          toggleCharacterDisplay({ characterId: displayCharacter.id, display })
        }
        onSetFavorite={() => setFavoriteCharacter(displayCharacter.id)}
        isTogglingDisplay={isTogglingDisplay}
        isSettingFavorite={isSettingFavorite}
        isFavorite={isFavorite}
      />
    </div>
  );
};

export default FavoriteCharacterSection;
