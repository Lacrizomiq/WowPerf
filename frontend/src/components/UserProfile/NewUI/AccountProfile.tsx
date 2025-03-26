// AccountProfile.tsx (WoWProfile)
"use client";

import React, { useEffect } from "react";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import { useWoWCharacters } from "@/hooks/useWowProtectedAccount";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { WoWError, WoWErrorCode } from "@/types/userCharacter/userCharacter";
import CharacterCardItem from "./CharacterCardItem";
import Image from "next/image";

export const WoWProfile: React.FC<{
  showTooltips?: boolean;
  limit?: number;
}> = ({ showTooltips = false, limit }) => {
  // Get Battle.net link status
  const { linkStatus, initiateLink } = useBattleNetLink();

  // Use the WoWCharacters hook to access locally stored characters
  // This replaces direct Blizzard API calls with cached data from our database
  const {
    userCharacters,
    isLoadingUserCharacters,
    wowProfile,
    syncCharacters,
    userCharactersError,
  } = useWoWCharacters();

  // Handle unauthorized errors by initiating Battle.net link
  useEffect(() => {
    if (
      userCharactersError instanceof WoWError &&
      userCharactersError.code === WoWErrorCode.UNAUTHORIZED
    ) {
      initiateLink();
    }
  }, [userCharactersError, initiateLink]);

  // Loading state - display spinner while characters are loading
  if (isLoadingUserCharacters) {
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500" />
      </div>
    );
  }

  // If not linked to Battle.net, show connection card
  if (!linkStatus?.linked) {
    return (
      <Card className="bg-[#131e33] border-gray-800">
        <CardHeader>
          <CardTitle>Link Battle.net Account</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="64"
              height="64"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="text-blue-500 mb-4"
            >
              <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
              <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
            </svg>
            <h3 className="text-xl font-bold mb-2">
              Connect Your Battle.net Account
            </h3>
            <p className="text-gray-400 text-center mb-6">
              Link your Battle.net account to view your WoW characters and
              access more features.
            </p>
            <Button onClick={initiateLink} className="flex items-center gap-2">
              <Image
                src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                alt="Battle.net"
                width={20}
                height={20}
              />
              Connect to Battle.net
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Use userCharacters directly instead of fetching from Blizzard API
  // These characters are already stored in our database from previous syncs
  let displayCharacters = userCharacters || [];

  // If no characters are found, prompt for synchronization
  if (displayCharacters.length === 0) {
    return (
      <Card className="bg-[#131e33] border-gray-800">
        <CardHeader>
          <CardTitle>No Characters Found</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8">
            <p className="text-gray-400 text-center mb-6">
              No characters found. Synchronize your Battle.net account to see
              your characters.
            </p>
            <Button
              onClick={() => syncCharacters()}
              className="flex items-center gap-2"
            >
              Synchronize Characters
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Limit number of characters if requested (for Overview tab)
  if (limit && displayCharacters.length > limit) {
    displayCharacters = displayCharacters.slice(0, limit);
  }

  // Render character grid
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      {displayCharacters.map((character, index) => (
        <CharacterCardItem
          key={`${character.name}-${index}`}
          character={character}
          region={wowProfile?.region || character.region || "eu"}
          onToggleDisplay={(display) => {
            /* Not used in this view */
          }}
          onSetFavorite={() => {
            /* Not used in this view */
          }}
          isTogglingDisplay={false}
          isSettingFavorite={false}
          // We can directly use the character's mythic_plus_rating from our database
          // No need to fetch it from Blizzard API again
        />
      ))}
    </div>
  );
};

export default WoWProfile;
