// AccountProfile.tsx (WoWProfile)
"use client";

import React, { useEffect, useState } from "react";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import { useWoWProfile } from "@/hooks/useWowProtectedAccount";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { WoWError, WoWErrorCode } from "@/libs/wowProtectedAccountService";
import { toast } from "react-hot-toast";
import CharacterCard from "./CharacterCard";
import Image from "next/image";

// Define a type for the character with M+ score
interface CharacterWithScore {
  name: string;
  playable_class: { name: string; id: number };
  level: number;
  realm: { name: string; slug: string };
  mPlusScore?: number;
  // Add any other properties from your character object
}

export const WoWProfile: React.FC<{
  showTooltips?: boolean;
  limit?: number;
}> = ({ showTooltips = false, limit }) => {
  const { linkStatus, initiateLink } = useBattleNetLink();
  const { wowProfile, isLoading, error, refetch } = useWoWProfile();
  // Properly type the state variable
  const [sortedCharacters, setSortedCharacters] = useState<
    CharacterWithScore[]
  >([]);
  const [isCharactersLoading, setIsCharactersLoading] = useState(true);

  useEffect(() => {
    if (error instanceof WoWError && error.code === WoWErrorCode.UNAUTHORIZED) {
      initiateLink();
    }
  }, [error, initiateLink]);

  // New effect to handle character sorting
  // New approach in the useEffect
  useEffect(() => {
    // First set the characters from the profile without fetching scores
    if (wowProfile && wowProfile.wow_accounts) {
      const characters = wowProfile.wow_accounts.flatMap(
        (account) => account.characters
      ) as CharacterWithScore[];

      // Initially set characters without sorting
      setSortedCharacters(characters);

      // Then fetch scores and update
      const fetchAllScores = async () => {
        setIsCharactersLoading(true);

        try {
          // Create a map to store scores by character name
          const scoreMap = new Map();

          // Fetch scores for each character
          for (const character of characters) {
            const safeRegion = wowProfile.region || "eu";
            const profileNamespace = `profile-${safeRegion}`;

            try {
              const response = await fetch(
                `/api/blizzard/character/${safeRegion}/${
                  character.realm.slug
                }/${character.name.toLowerCase()}/mythic-plus-best-runs?namespace=${profileNamespace}&locale=en_GB&season=13`
              );

              if (response.ok) {
                const data = await response.json();
                const score = data?.OverallMythicRating || 0;
                scoreMap.set(character.name, score);
                console.log(`Fetched score for ${character.name}: ${score}`);
              } else {
                console.log(
                  `Error fetching score for ${character.name}: ${response.status}`
                );
                scoreMap.set(character.name, 0);
              }
            } catch (error) {
              console.error(`Error for ${character.name}:`, error);
              scoreMap.set(character.name, 0);
            }
          }

          // Now update all characters with their scores
          const charactersWithScores = characters.map((character) => ({
            ...character,
            mPlusScore: scoreMap.get(character.name) || 0,
          }));

          // Log before sorting
          console.log(
            "Before sorting:",
            charactersWithScores.map((c) => `${c.name}: ${c.mPlusScore}`)
          );

          // Sort characters by M+ score (highest first)
          const sortedChars = [...charactersWithScores].sort((a, b) => {
            const scoreA = Number(a.mPlusScore) || 0;
            const scoreB = Number(b.mPlusScore) || 0;
            return scoreB - scoreA;
          });

          // Log after sorting
          console.log(
            "After sorting:",
            sortedChars.map((c) => `${c.name}: ${c.mPlusScore}`)
          );

          // Update state with sorted characters
          setSortedCharacters(sortedChars);
        } catch (error) {
          console.error("Error in score fetching:", error);
        } finally {
          setIsCharactersLoading(false);
        }
      };

      fetchAllScores();
    }
  }, [wowProfile]);

  if (isLoading || isCharactersLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500" />
      </div>
    );
  }

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

  // Log the final list of characters being displayed
  console.log(
    "Final display characters with scores:",
    sortedCharacters.map((c) => ({
      name: c.name,
      score: c.mPlusScore,
    }))
  );

  // Limit number of characters if requested (for Overview tab)
  let displayCharacters = sortedCharacters;
  if (limit && displayCharacters.length > limit) {
    displayCharacters = displayCharacters.slice(0, limit);
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      {displayCharacters.length > 0 ? (
        displayCharacters.map((character, index) => (
          <CharacterCard
            key={`${character.name}-${index}`}
            character={character}
            region={wowProfile?.region}
            showTooltip={showTooltips}
            isMainCard={limit === 1}
            mythicPlusScore={character.mPlusScore} // Pass M+ score to avoid refetching
          />
        ))
      ) : (
        <div className="col-span-full text-center py-8 text-gray-400">
          <p>No characters found for your account.</p>
        </div>
      )}
    </div>
  );
};
