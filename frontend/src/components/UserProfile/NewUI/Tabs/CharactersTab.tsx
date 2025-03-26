// CharactersTab.tsx
import React, { useState } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useWoWCharacters } from "@/hooks/useWowProtectedAccount";
import CharacterCardItem from "../CharacterCardItem";
import { showSuccess, TOAST_IDS } from "@/utils/toastManager";
import { useUserProfile } from "@/hooks/useUserProfile";

const CharactersTab: React.FC = () => {
  // Get user profile to access favorite character id
  const { profile } = useUserProfile();

  // Utilization of the combined hook to manage characters
  const {
    userCharacters,
    isLoadingUserCharacters,
    syncCharacters,
    isSyncing,
    refreshCharacters,
    isRefreshing,
    setFavoriteCharacter,
    isSettingFavorite,
    toggleCharacterDisplay,
    isTogglingDisplay,
    wowProfile,
  } = useWoWCharacters();

  // Local state for filtering and sorting
  const [classFilter, setClassFilter] = useState<string | null>(null);
  const [showHidden, setShowHidden] = useState(false);

  // Handle initial sync action if no characters
  const handleSync = () => {
    syncCharacters();
    showSuccess(
      "Synchronization of characters in progress...",
      TOAST_IDS.CHARACTERS_SYNC
    );
  };

  // Handle refresh action
  const handleRefresh = () => {
    refreshCharacters();
    showSuccess(
      "Refreshing characters in progress...",
      TOAST_IDS.CHARACTERS_REFRESH
    );
  };

  // Apply filters to characters
  const getFilteredCharacters = () => {
    if (!userCharacters) return [];

    return userCharacters.filter((char) => {
      // Filter by class
      if (classFilter && char.class !== classFilter) return false;

      // Display filter
      if (!showHidden && !char.is_displayed) return false;

      return true;
    });
  };

  // Get unique classes for filters
  const getUniqueClasses = () => {
    if (!userCharacters) return [];

    return Array.from(new Set(userCharacters.map((char) => char.class)));
  };

  // Loading state
  if (isLoadingUserCharacters) {
    return (
      <Card className="bg-[#131e33] border-gray-800 p-6">
        <div className="flex justify-center items-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500" />
        </div>
      </Card>
    );
  }

  // No characters found
  if (!userCharacters || userCharacters.length === 0) {
    return (
      <Card className="bg-[#131e33] border-gray-800 p-6">
        <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="text-blue-500"
          >
            <path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48" />
          </svg>
          Your characters
        </h2>

        <div className="text-center py-12">
          <p className="text-gray-400 mb-6">
            No characters found. Synchronize your Battle.net account to see your
            characters.
          </p>
          <Button
            onClick={handleSync}
            disabled={isSyncing}
            className="bg-blue-500 hover:bg-blue-600"
          >
            {isSyncing ? "Synchronization..." : "Synchronize characters"}
          </Button>
        </div>
      </Card>
    );
  }

  // Get favorite character ID from user profile
  const favoriteCharacterId = profile?.favorite_character_id;

  // Filter characters
  const filteredCharacters = getFilteredCharacters();
  const uniqueClasses = getUniqueClasses();

  return (
    <Card className="bg-[#131e33] border-gray-800 p-6">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-bold flex items-center gap-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="text-blue-500"
          >
            <path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48" />
          </svg>
          Your characters
        </h2>

        <Button
          onClick={handleRefresh}
          disabled={isRefreshing}
          className="bg-blue-500 hover:bg-blue-600"
        >
          {isRefreshing ? "Refreshing..." : "Refresh"}
        </Button>
      </div>

      <p className="text-gray-400 mb-4">
        You can manage the display of your characters and define your favorite
        character.
      </p>

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-2 mb-6">
        <div className="text-sm text-gray-400 mr-2">Filter by class:</div>
        <Button
          onClick={() => setClassFilter(null)}
          variant={classFilter === null ? "default" : "secondary"}
          size="sm"
          className={`rounded-full ${
            classFilter === null ? "bg-blue-500" : "bg-gray-800"
          }`}
        >
          All
        </Button>

        {uniqueClasses.map((className) => (
          <Button
            key={className}
            onClick={() => setClassFilter(className)}
            variant={classFilter === className ? "default" : "secondary"}
            size="sm"
            className={`rounded-full ${
              classFilter === className ? "bg-blue-500" : "bg-gray-800"
            }`}
          >
            {className}
          </Button>
        ))}

        <div className="ml-auto">
          <Button
            onClick={() => setShowHidden(!showHidden)}
            variant="outline"
            size="sm"
            className="rounded-full"
          >
            {showHidden ? (
              <span className="flex items-center gap-1">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24"></path>
                  <path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"></path>
                  <path d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"></path>
                  <line x1="2" x2="22" y1="2" y2="22"></line>
                </svg>
                Hide hidden
              </span>
            ) : (
              <span className="flex items-center gap-1">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
                Show all
              </span>
            )}
          </Button>
        </div>
      </div>

      {/* List of characters */}
      {filteredCharacters.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {filteredCharacters.map((character) => (
            <CharacterCardItem
              key={character.id}
              character={character}
              region={wowProfile?.region || "eu"}
              onToggleDisplay={(display) =>
                toggleCharacterDisplay({ characterId: character.id, display })
              }
              onSetFavorite={() => setFavoriteCharacter(character.id)}
              isTogglingDisplay={isTogglingDisplay}
              isSettingFavorite={isSettingFavorite}
              isFavorite={character.id === favoriteCharacterId}
            />
          ))}
        </div>
      ) : (
        <div className="text-center py-8 text-gray-400">
          <p>No character matches your filtering criteria.</p>
        </div>
      )}
    </Card>
  );
};

export default CharactersTab;
