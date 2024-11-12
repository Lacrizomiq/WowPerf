"use client";

import React, { useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import Image from "next/image";
import { UserProfile } from "@/libs/userService";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";

interface PersonalInfoProps {
  profile: UserProfile;
}

const PersonalInfo: React.FC<PersonalInfoProps> = ({ profile }) => {
  const queryClient = useQueryClient();
  const {
    linkStatus,
    isLoading,
    error,
    initiateLink,
    unlinkAccount,
    isUnlinking,
  } = useBattleNetLink();

  useEffect(() => {
    // Rafraîchir le statut de la liaison quand le composant est monté
    // et quand l'URL contient success=true
    if (window.location.search.includes("success=true")) {
      // Force refresh du status
      queryClient.invalidateQueries({ queryKey: ["battleNetLinkStatus"] });
      // Nettoyer l'URL
      window.history.replaceState({}, "", window.location.pathname);
    }
  }, [queryClient]);

  const handleBattleNetConnect = async () => {
    try {
      await initiateLink();
    } catch (error) {
      console.error("Failed to initiate OAuth:", error);
    }
  };

  const handleBattleNetUnlink = async () => {
    try {
      await unlinkAccount();
    } catch (error) {
      console.error("Failed to unlink Battle.net:", error);
    }
  };

  return (
    <>
      <section className="bg-deep-blue shadow rounded-lg p-6 border border-gray-800">
        <h2 className="text-2xl font-bold mb-4 text-[#e2e8f0]">
          Personal Information
        </h2>
        <div className="space-y-4">
          <div className="flex items-center">
            <p className="block text-sm font-medium text-[#e2e8f0]">
              <span className="font-bold text-lg">Username: </span>
              {profile.username}
            </p>
          </div>
          <div>
            <p className="block text-sm font-medium text-[#e2e8f0]">
              <span className="font-bold text-lg">Email: </span>
              {profile.email}
            </p>
          </div>
        </div>
      </section>

      <section className="bg-deep-blue shadow rounded-lg p-6 mt-4 border border-gray-800">
        <h2 className="text-2xl font-bold mb-4 text-[#e2e8f0]">
          Link your Battle.net account
        </h2>
        <div className="flex items-center">
          {linkStatus?.linked ? (
            <div className="space-y-4">
              <div className="text-green-500 flex items-center">
                <span className="mr-2">✓</span>
                Connected as {linkStatus.battleTag}
              </div>
              <button
                onClick={handleBattleNetUnlink}
                disabled={isUnlinking}
                className="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded transition duration-200 disabled:opacity-50"
              >
                {isUnlinking ? "Unlinking..." : "Unlink Battle.net Account"}
              </button>
            </div>
          ) : (
            <button
              onClick={handleBattleNetConnect}
              className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-200 flex items-center"
            >
              <Image
                src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                alt="Battle.net icon"
                width={24}
                height={24}
                className="mr-2"
              />
              Connect to Battle.net
            </button>
          )}
        </div>
      </section>
    </>
  );
};

export default PersonalInfo;
