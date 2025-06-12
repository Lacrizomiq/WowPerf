// Profile.tsx - Version HTML native complète sans restriction navigation
"use client";

import React, { useState, useRef, useEffect } from "react";
import { useUserProfile } from "@/hooks/useUserProfile";
import { useAuth } from "@/providers/AuthContext";
import { useRouter } from "next/navigation";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import { showError, TOAST_IDS } from "@/utils/toastManager";

// Composants natifs
import ProfileHeader from "./ProfileHeader";
import TabNavigation from "./TabNavigation";
import AccountTab from "./AccountTab";
import CharactersTab from "./CharactersTab";
import ConnectionsTab from "./ConnectionsTab";

const Profile: React.FC = () => {
  // State for active tab navigation
  const [activeTab, setActiveTab] = useState("account");

  // Authentication and routing
  const { isAuthenticated } = useAuth();
  const router = useRouter();

  // User profile data and mutations
  const { profile, isLoading, error } = useUserProfile();

  // Battle.net connection state
  const {
    linkStatus,
    isLoading: isLinkLoading,
    error: linkError,
    initiateLink,
    unlinkAccount,
    isUnlinking,
  } = useBattleNetLink();

  // Route protection - redirect to login if not authenticated
  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  // Loading state
  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-purple-500"></div>
      </div>
    );
  }

  // Error handling
  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-red-500 bg-red-100 p-4 rounded-lg">
          {error instanceof Error ? error.message : "An error occurred"}
        </div>
      </div>
    );
  }

  // Check if profile exists
  if (!profile) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-yellow-500">No profile data available</div>
      </div>
    );
  }

  // Handle Battle.net unlink action
  const handleBattleNetUnlink = async () => {
    try {
      await unlinkAccount();
      // Si on est dans l'onglet Characters, on peut y rester
      // L'onglet affichera les personnages sauvegardés avec un warning
    } catch (error) {
      console.error("Failed to unlink Battle.net:", error);
    }
  };

  // Handle Battle.net link action
  const handleBattleNetLink = async () => {
    try {
      await initiateLink();
    } catch (error) {
      console.error("Failed to link Battle.net:", error);
    }
  };

  // Navigation libre vers tous les onglets
  const handleNavigate = (tab: string) => {
    setActiveTab(tab);
  };

  return (
    <div className="flex flex-col min-h-screen bg-[#1A1D21]">
      {/* Profile header with summary */}
      <div className="">
        <ProfileHeader profile={profile} />
      </div>

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-4 md:px-8 py-6">
        {/* Tab navigation */}
        <TabNavigation activeTab={activeTab} setActiveTab={handleNavigate} />

        {/* Tab content */}
        <div className="mt-6">
          {/* Account Tab */}
          <AccountTab profile={profile} isActive={activeTab === "account"} />

          {/* Characters Tab - TOUJOURS accessible */}
          <CharactersTab isActive={activeTab === "characters"} />

          {/* Connections Tab */}
          <ConnectionsTab
            profile={profile}
            linkStatus={linkStatus}
            isLinkLoading={isLinkLoading}
            isUnlinking={isUnlinking}
            onBattleNetLink={handleBattleNetLink}
            onBattleNetUnlink={handleBattleNetUnlink}
            isActive={activeTab === "connections"}
          />
        </div>
      </main>
    </div>
  );
};

export default Profile;
