"use client";

import React, { useEffect } from "react";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import { useWoWProfile } from "@/hooks/useWowProtectedAccount";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { WoWError, WoWErrorCode } from "@/libs/wowProtectedAccountService";
import { toast } from "react-hot-toast";
import { useRouter } from "next/navigation";
import Image from "next/image";
import { getClassIcon } from "@/utils/classandspecicons";

export const WoWProfile: React.FC = () => {
  const router = useRouter();
  const { linkStatus, initiateLink, unlinkAccount } = useBattleNetLink();
  const { wowProfile, isLoading, error, refetch } = useWoWProfile();

  useEffect(() => {
    if (error instanceof WoWError && error.code === WoWErrorCode.UNAUTHORIZED) {
      initiateLink();
    }
  }, [error, initiateLink]);

  const handleCharacterClick = (character: {
    name: string;
    realm: { slug: string };
  }) => {
    if (!wowProfile?.region) {
      toast.error("Cannot determine region");
      return;
    }

    const realmSlug = character.realm.slug;
    const characterName = character.name.toLowerCase();
    router.push(
      `/character/${wowProfile.region}/${realmSlug}/${characterName}`
    );
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500" />
      </div>
    );
  }

  if (!linkStatus?.linked) {
    return (
      <Card className="w-full max-w-md mx-auto mt-8">
        <CardHeader>
          <CardTitle>Link Battle.net Account</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="mb-4">
            Link your Battle.net account to view your WoW profile
          </p>
          <Button onClick={initiateLink}>Connect Battle.net</Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="mt-4">
      {" "}
      {/* Suppression de la Card englobante */}
      {wowProfile ? (
        <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-6">
          {wowProfile.wow_accounts.flatMap((account) =>
            account.characters.map((character, charIndex) => {
              const className = character.playable_class.name;
              const normalizedClassName = className.replace(/\s+/g, "");
              const classIcon = getClassIcon(normalizedClassName);

              return (
                <Card
                  key={`${character.name}-${charIndex}`}
                  className="bg-black border-gray-800 cursor-pointer hover:bg-gray-800 transition-all duration-200 transform hover:-translate-y-1"
                  onClick={() => handleCharacterClick(character)}
                >
                  <CardContent className="p-6">
                    <div className="flex items-center gap-3 mb-3">
                      <Image
                        src={classIcon}
                        alt={className}
                        width={32}
                        height={32}
                        className="rounded"
                      />
                      <h4
                        className={`text-lg font-bold class-color--${character.playable_class.id}`}
                      >
                        {character.name}
                      </h4>
                    </div>
                    <div className="space-y-2 text-sm">
                      <div className="flex items-center gap-2">
                        <span
                          className="inline-block w-2 h-2 rounded-full"
                          style={{
                            backgroundColor:
                              character.level >= 80 ? "#4ade80" : "#ef4444",
                          }}
                        />
                        <p className="text-gray-200">Level {character.level}</p>
                      </div>
                      <p className="text-gray-300">
                        {character.playable_class.name}
                      </p>
                      <p className="text-gray-400">
                        {character.realm.name} (
                        {wowProfile.region?.toUpperCase()})
                      </p>
                    </div>
                  </CardContent>
                </Card>
              );
            })
          )}
        </div>
      ) : (
        <div className="text-center py-6 text-gray-400">
          <p>No profile data available</p>
        </div>
      )}
    </div>
  );
};
