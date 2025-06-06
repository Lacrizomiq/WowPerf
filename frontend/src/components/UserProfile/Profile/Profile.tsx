// Profile.tsx - Version HTML native complète sans shadcn/ui
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

  // Ref to track if Battle.net warning was already shown
  const battleNetWarningShown = useRef(false);

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

  // Check if the user can access the Characters tab - ONE TIME useEffect
  useEffect(() => {
    if (
      activeTab === "characters" &&
      !linkStatus?.linked &&
      !battleNetWarningShown.current
    ) {
      battleNetWarningShown.current = true;
      setActiveTab("connections");

      // Use the centralized function with ID
      showError(
        "You must link your Battle.net account to access your characters",
        TOAST_IDS.BATTLENET_LINKING
      );

      setTimeout(() => {
        battleNetWarningShown.current = false;
      }, 2000);
    }
  }, [activeTab, linkStatus]);

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
      // If in the Characters tab, redirect to Connections after disconnection
      if (activeTab === "characters") {
        setActiveTab("connections");
      }
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

  // Handle tab navigation
  const handleNavigate = (tab: string) => {
    // Check if the user can access the Characters tab
    if (tab === "characters" && !linkStatus?.linked) {
      // The notification will be handled by the useEffect, no need to put it here
      setActiveTab("connections"); // Redirect to the Connections tab
    } else {
      setActiveTab(tab);
    }
  };

  return (
    <div className="flex flex-col min-h-screen">
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

          {/* Characters Tab - only if Battle.net is linked */}
          {activeTab === "characters" && linkStatus?.linked && (
            <CharactersTab isActive={true} />
          )}

          {/* Show connections tab if trying to access characters without link */}
          {activeTab === "characters" && !linkStatus?.linked && (
            <ConnectionsTab
              profile={profile}
              linkStatus={linkStatus}
              isLinkLoading={isLinkLoading}
              isUnlinking={isUnlinking}
              onBattleNetLink={handleBattleNetLink}
              onBattleNetUnlink={handleBattleNetUnlink}
              isActive={true}
            />
          )}

          {/* Connections Tab */}
          {activeTab === "connections" && (
            <ConnectionsTab
              profile={profile}
              linkStatus={linkStatus}
              isLinkLoading={isLinkLoading}
              isUnlinking={isUnlinking}
              onBattleNetLink={handleBattleNetLink}
              onBattleNetUnlink={handleBattleNetUnlink}
              isActive={true}
            />
          )}
        </div>
      </main>
    </div>
  );
};

export default Profile;
