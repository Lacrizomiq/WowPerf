// Profile.tsx
// Main component that coordinates the profile interface and tab navigation
"use client";

import React, { useState, useRef } from "react";
import { useUserProfile } from "@/hooks/useUserProfile";
import { useAuth } from "@/providers/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import TabNavigation from "./TabNavigation";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import ChangePassword from "../Update/ChangePassword";
import ChangeEmail from "../Update/ChangeEmail";
import ChangeUsername from "../Update/ChangeUsername";
import DeleteAccount from "../Update/DeleteAccount";
import ProfileHeader from "./ProfileHeader";
import OverviewTab from "./Tabs/OverviewTab";
import CharactersTab from "./Tabs/CharactersTab";
import AccountTab from "./Tabs/AccountTab";
import SecurityTab from "./Tabs/SecurityTab";
import ConnectionsTab from "./Tabs/ConnectionsTab";
import { showError, TOAST_IDS } from "@/utils/toastManager";
import FavoriteCharacterSection from "./FavoriteCharacterSection";

const Profile: React.FC = () => {
  // State for active tab navigation
  const [activeTab, setActiveTab] = useState("overview");

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
        <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
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
    <div className="max-w-7xl mx-auto px-4 py-8">
      {/* Profile header with avatar and basic info */}
      <ProfileHeader profile={profile} />

      {/* Tab navigation */}
      <TabNavigation activeTab={activeTab} setActiveTab={handleNavigate} />

      {/* Tab content container */}
      <div className="mt-6">
        {/* Main tabs */}
        {activeTab === "overview" && (
          <OverviewTab
            profile={profile}
            linkStatus={linkStatus}
            isLinkLoading={isLinkLoading}
            isUnlinking={isUnlinking}
            onBattleNetLink={handleBattleNetLink}
            onBattleNetUnlink={handleBattleNetUnlink}
            onNavigate={handleNavigate}
          />
        )}

        {/* Modified condition for CharactersTab */}
        {activeTab === "characters" && linkStatus?.linked ? (
          <CharactersTab />
        ) : activeTab === "characters" ? (
          <ConnectionsTab
            profile={profile}
            linkStatus={linkStatus}
            isLinkLoading={isLinkLoading}
            isUnlinking={isUnlinking}
            onBattleNetLink={handleBattleNetLink}
            onBattleNetUnlink={handleBattleNetUnlink}
          />
        ) : null}

        {activeTab === "account" && (
          <AccountTab profile={profile} onNavigate={handleNavigate} />
        )}

        {activeTab === "security" && (
          <SecurityTab onNavigate={handleNavigate} />
        )}

        {activeTab === "connections" && (
          <ConnectionsTab
            profile={profile}
            linkStatus={linkStatus}
            isLinkLoading={isLinkLoading}
            isUnlinking={isUnlinking}
            onBattleNetLink={handleBattleNetLink}
            onBattleNetUnlink={handleBattleNetUnlink}
          />
        )}

        {/* Action tabs */}
        {activeTab === "username" && <ChangeUsername />}
        {activeTab === "email" && <ChangeEmail />}
        {activeTab === "password" && <ChangePassword />}
        {activeTab === "delete" && <DeleteAccount />}
      </div>
    </div>
  );
};

export default Profile;
