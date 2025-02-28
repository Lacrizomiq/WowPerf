// AccountProfile.tsx (WoWProfile)
"use client";

import React, { useEffect } from "react";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import { useWoWProfile } from "@/hooks/useWowProtectedAccount";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { WoWError, WoWErrorCode } from "@/libs/wowProtectedAccountService";
import { toast } from "react-hot-toast";
import CharacterCard from "./CharacterCard";
import Image from "next/image";

export const WoWProfile: React.FC<{
  showTooltips?: boolean;
  limit?: number;
}> = ({ showTooltips = false, limit }) => {
  const { linkStatus, initiateLink } = useBattleNetLink();
  const { wowProfile, isLoading, error, refetch } = useWoWProfile();

  useEffect(() => {
    if (error instanceof WoWError && error.code === WoWErrorCode.UNAUTHORIZED) {
      initiateLink();
    }
  }, [error, initiateLink]);

  if (isLoading) {
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

  // Find all characters from wow profile
  let characters = wowProfile
    ? wowProfile.wow_accounts.flatMap((account) => account.characters)
    : [];

  // Limit number of characters if requested (for Overview tab)
  if (limit && characters.length > limit) {
    characters = characters.slice(0, limit);
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      {characters.length > 0 ? (
        characters.map((character, index) => (
          <CharacterCard
            key={`${character.name}-${index}`}
            character={character}
            region={wowProfile?.region}
            showTooltip={showTooltips}
            isMainCard={limit === 1}
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
