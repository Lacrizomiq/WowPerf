// CharactersTab.tsx
import React, { useState } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useCharacters } from "@/hooks/useCharacters";
import CharacterCardItem from "./CharacterCardItem";
import { useUserProfile } from "@/hooks/useUserProfile";
import { isCharacterEnriched } from "@/utils/character/character";
import { EnrichmentNotice } from "./Tabs/EnrichmentNotice";
import { EnrichedUserCharacter } from "@/types/character/character";

const CharactersTab: React.FC = () => {
  // Get user profile to access favorite character id
  const { profile } = useUserProfile();

  // Use the new useCharacters hook
  const {
    characters,
    isLoadingCharacters,
    actions,
    isLoading,
    rateLimitState,
    ui,
    region,
  } = useCharacters();

  console.log("characters", characters);

  const charactersArray: EnrichedUserCharacter[] = characters || [];

  console.log("charactersArray", charactersArray);
  console.log("charactersArray.length", charactersArray?.length);
  console.log("Array.isArray(charactersArray)", Array.isArray(charactersArray));

  // Local state for filtering and sorting
  const [classFilter, setClassFilter] = useState<string | null>(null);

  // Determine if user has characters (for conditional button logic)
  const hasCharacters =
    Array.isArray(charactersArray) && charactersArray.length > 0;

  // Handle sync action (first time users)
  const handleSync = () => {
    actions.syncAndEnrich();
  };

  // Handle refresh action (existing users)
  const handleRefresh = () => {
    actions.refreshAndEnrich();
  };

  // Apply filters to characters - SIMPLIFIED
  const getFilteredCharacters = (): EnrichedUserCharacter[] => {
    if (!Array.isArray(charactersArray)) return [];

    return charactersArray.filter((char: EnrichedUserCharacter) => {
      // Filter by class if selected
      if (classFilter && char.class !== classFilter) return false;

      // For now, show all displayed characters (keep it simple)
      return char.is_displayed; // Only show displayed characters
    });
  };

  // Get unique classes for filters - SIMPLIFIED
  const getUniqueClasses = (): string[] => {
    if (!Array.isArray(charactersArray)) return [];
    return Array.from(
      new Set(charactersArray.map((char: EnrichedUserCharacter) => char.class))
    );
  };

  // Loading state
  if (isLoadingCharacters) {
    return (
      <Card className="bg-[#131e33] border-gray-800 p-6">
        <div className="flex justify-center items-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500" />
        </div>
      </Card>
    );
  }

  // No characters found - show sync button
  if (!hasCharacters) {
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
          <div className="mb-6">
            <svg
              className="w-16 h-16 text-gray-500 mx-auto mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>

            <p className="text-gray-400 mb-4">
              {isLoading.sync
                ? "Synchronizing your characters..."
                : "Sync not done, click on sync button to display your characters"}
            </p>

            {isLoading.sync && (
              <div className="text-sm text-gray-500 space-y-1">
                <div className="flex items-center justify-center gap-2">
                  <div className="animate-pulse h-2 w-2 bg-blue-500 rounded-full"></div>
                  <span>Connecting to Battle.net</span>
                </div>
                <div className="flex items-center justify-center gap-2">
                  <div className="animate-pulse h-2 w-2 bg-blue-500 rounded-full"></div>
                  <span>Fetching character data</span>
                </div>
                <div className="flex items-center justify-center gap-2">
                  <div className="animate-pulse h-2 w-2 bg-blue-500 rounded-full"></div>
                  <span>Enriching character information</span>
                </div>
              </div>
            )}
          </div>

          {/* Sync button */}
          <Button
            onClick={handleSync}
            disabled={ui.isDisabled.sync}
            className="bg-blue-500 hover:bg-blue-600 min-w-[200px]"
          >
            {isLoading.sync ? "Synchronizing..." : "Sync Characters"}
          </Button>

          {/* Rate limit message */}
          {ui.showRateLimit && (
            <div className="mt-4 p-3 bg-orange-900/30 border border-orange-500/50 rounded-lg">
              <div className="text-orange-400 font-semibold flex items-center justify-center gap-2">
                <svg
                  className="w-5 h-5"
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
                {rateLimitState.formattedTime && (
                  <span className="font-mono">
                    {rateLimitState.formattedTime}
                  </span>
                )}
              </div>
              <p className="text-sm text-orange-300 mt-1">
                {rateLimitState.message}
              </p>
            </div>
          )}
        </div>
      </Card>
    );
  }

  // Get favorite character ID from user profile
  const favoriteCharacterId = profile?.favorite_character_id;

  // Filter characters
  const filteredCharacters = getFilteredCharacters();
  console.log("filteredCharacters", filteredCharacters);
  const uniqueClasses = getUniqueClasses();

  // Check if there are unenriched characters - SIMPLIFIED
  const hasUnenrichedCharacters =
    Array.isArray(charactersArray) &&
    charactersArray.some(
      (char: EnrichedUserCharacter) => !isCharacterEnriched(char)
    );

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

        {/* Refresh button (only shown when user has characters) */}
        <div className="flex items-center gap-3">
          {/* Rate limit indicator */}
          {ui.showRateLimit && (
            <div className="text-sm text-orange-400 flex items-center gap-2">
              <svg
                className="w-4 h-4"
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
              {rateLimitState.formattedTime}
            </div>
          )}

          <Button
            onClick={handleRefresh}
            disabled={ui.isDisabled.refresh}
            className="bg-blue-500 hover:bg-blue-600"
          >
            {isLoading.refresh ? "Refreshing..." : "Refresh"}
          </Button>
        </div>
      </div>

      <p className="text-gray-400 mb-4">
        You can manage the display of your characters and define your favorite
        character. Click refresh to update character data from Battle.net.
      </p>

      {/* Rate limit message (detailed) */}
      {ui.showRateLimit && (
        <div className="mb-4 p-3 bg-orange-900/30 border border-orange-500/50 rounded-lg">
          <p className="text-orange-300 text-sm">⏱️ {rateLimitState.message}</p>
        </div>
      )}

      {/* Filters - SIMPLIFIED */}
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

        {uniqueClasses.map((className: string) => (
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
      </div>

      {/* Enrichment notice */}
      <EnrichmentNotice
        hasUnenrichedCharacters={hasUnenrichedCharacters}
        onRefresh={handleRefresh}
        isRefreshing={isLoading.refresh}
      />

      {/* List of characters */}
      {filteredCharacters.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {filteredCharacters.map((character: EnrichedUserCharacter) => (
            <CharacterCardItem
              key={character.id}
              character={character}
              region={region}
              onToggleDisplay={(display) => {
                /* TODO: Implement toggle display for individual characters */
              }}
              onSetFavorite={() => {
                /* TODO: Implement set favorite for individual characters */
              }}
              isTogglingDisplay={false}
              isSettingFavorite={false}
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
