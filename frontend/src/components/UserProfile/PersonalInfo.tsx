"use client";

import React, { useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import Image from "next/image";
import { UserProfile } from "@/libs/userService";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import { UserErrorCode } from "@/libs/userService";

// Définition des types pour les mutations
interface MutationState {
  isError: boolean;
  error: Error | null;
  isPending: boolean;
  errorCode?: UserErrorCode;
}

interface MutationResponse {
  success: boolean;
  error?: string;
  code?: UserErrorCode;
}

interface PersonalInfoProps {
  profile: UserProfile;
  mutationStates: {
    changeUsername: MutationState;
    updateEmail: MutationState;
    changePassword?: MutationState;
    deleteAccount?: MutationState;
  };
  onUpdateEmail: (email: string) => Promise<MutationResponse>;
  onChangeUsername: (username: string) => Promise<MutationResponse>;
}

// Interface pour le statut Battle.net (si non défini dans useBattleNetLink)
interface BattleNetLinkStatus {
  linked: boolean;
  battleTag?: string;
}

const PersonalInfo: React.FC<PersonalInfoProps> = ({
  profile,
  mutationStates,
}) => {
  const queryClient = useQueryClient();
  const {
    linkStatus,
    isLoading: isLinkLoading,
    error: linkError,
    initiateLink,
    unlinkAccount,
    isUnlinking,
  } = useBattleNetLink();

  useEffect(() => {
    if (window.location.search.includes("success=true")) {
      queryClient.invalidateQueries({ queryKey: ["battleNetLinkStatus"] });
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
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["battleNetLinkStatus"] }),
        queryClient.invalidateQueries({ queryKey: ["userProfile"] }),
      ]);
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
          {/* Username */}
          <div className="flex items-center justify-between">
            <p className="block text-sm font-medium text-[#e2e8f0]">
              <span className="font-bold text-lg">Username: </span>
              {profile.username}
            </p>
            {mutationStates.changeUsername.isError && (
              <p className="text-red-500 text-sm">
                {mutationStates.changeUsername.error?.message}
              </p>
            )}
          </div>

          {/* Email */}
          <div className="flex items-center justify-between">
            <p className="block text-sm font-medium text-[#e2e8f0]">
              <span className="font-bold text-lg">Email: </span>
              {profile.email}
            </p>
            {mutationStates.updateEmail.isError && (
              <p className="text-red-500 text-sm">
                {mutationStates.updateEmail.error?.message}
              </p>
            )}
          </div>

          {/* Battle.net Tag */}
          {profile.battle_tag && (
            <div className="flex items-center justify-between">
              <p className="block text-sm font-medium text-[#e2e8f0]">
                <span className="font-bold text-lg">Battle Tag: </span>
                <span className="text-blue-400">{profile.battle_tag}</span>
              </p>
            </div>
          )}
        </div>
      </section>

      <section className="bg-deep-blue shadow rounded-lg p-6 mt-4 border border-gray-800">
        <h2 className="text-2xl font-bold mb-4 text-[#e2e8f0]">
          Battle.net Account {linkStatus?.linked ? "Linked" : "Not Linked"}
        </h2>
        {linkError && <p className="text-red-500 mb-4">{linkError.message}</p>}
        <div className="flex items-center">
          {linkStatus?.linked ? (
            <div className="space-y-4">
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
              disabled={isLinkLoading}
              className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-200 flex items-center disabled:opacity-50"
            >
              <Image
                src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                alt="Battle.net icon"
                width={24}
                height={24}
                className="mr-2"
              />
              {isLinkLoading ? "Connecting..." : "Connect to Battle.net"}
            </button>
          )}
        </div>
      </section>
    </>
  );
};

export default PersonalInfo;
