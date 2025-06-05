// AccountProfile.tsx (WoWProfile)
"use client";

import React from "react";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import { useCharacters } from "@/hooks/useCharacters";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import CharacterCardItem from "./CharacterCardItem";
import Image from "next/image";

export const WoWProfile: React.FC<{
  showTooltips?: boolean;
  limit?: number;
}> = ({ showTooltips = false, limit }) => {
  // Get Battle.net link status
  const { linkStatus, initiateLink } = useBattleNetLink();

  // Use the new useCharacters hook for enriched data from BDD
  const {
    characters,
    isLoadingCharacters,
    actions,
    isLoading,
    ui,
    rateLimitState,
    region,
  } = useCharacters();

  // Loading state - display spinner while characters are loading
  if (isLoadingCharacters) {
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

  // Use characters from our enriched database instead of direct API calls
  let displayCharacters = Array.isArray(characters) ? characters : [];

  // If no characters are found, show sync or appropriate message
  if (displayCharacters.length === 0) {
    return (
      <Card className="bg-[#131e33] border-gray-800">
        <CardHeader>
          <CardTitle>No Characters Found</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8">
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

              <p className="text-gray-400 text-center mb-4">
                {isLoading.sync
                  ? "Synchronizing your characters..."
                  : "Sync not done, click on sync button to display your characters"}
              </p>

              {isLoading.sync && (
                <div className="text-sm text-gray-500 space-y-1">
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

            {/* Sync Characters Button */}
            <Button
              onClick={actions.syncAndEnrich}
              disabled={ui.isDisabled.sync}
              className="flex items-center gap-2 mb-4"
            >
              {isLoading.sync ? "Synchronizing..." : "Sync Characters"}
            </Button>

            {/* Rate limit message */}
            {ui.showRateLimit && (
              <div className="mt-4 p-3 bg-orange-900/30 border border-orange-500/50 rounded-lg max-w-md">
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
                <p className="text-sm text-orange-300 mt-1 text-center">
                  {rateLimitState.message}
                </p>
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    );
  }

  // Limit number of characters if requested (for Overview tab)
  if (limit && displayCharacters.length > limit) {
    displayCharacters = displayCharacters.slice(0, limit);
  }

  // Render character grid using enriched data from our database
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      {displayCharacters.map((character, index) => (
        <CharacterCardItem
          key={`${character.name}-${character.id}-${index}`}
          character={character}
          region={region}
          onToggleDisplay={() => {
            /* Not used in this overview */
          }}
          onSetFavorite={() => {
            /* Not used in this overview */
          }}
          isTogglingDisplay={false}
          isSettingFavorite={false}
          // Character data is already enriched from our database
          // including mythic_plus_rating, item_level, active_spec_name, etc.
        />
      ))}
    </div>
  );
};

export default WoWProfile;
