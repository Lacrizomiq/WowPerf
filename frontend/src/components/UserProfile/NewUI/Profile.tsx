// Profile.tsx (refactorisé)
"use client";

import React, { useState } from "react";
import { useUserProfile } from "@/hooks/useUserProfile";
import { useAuth } from "@/providers/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import TabNavigation from "./TabNavigation";
import { WoWProfile } from "./AccountProfile";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import Image from "next/image";
import ChangePassword from "../Update/ChangePassword";
import ChangeEmail from "../Update/ChangeEmail";
import ChangeUsername from "../Update/ChangeUsername";
import DeleteAccount from "../Update/DeleteAccount";

const Profile: React.FC = () => {
  const [activeTab, setActiveTab] = useState("overview");
  const { isAuthenticated } = useAuth();
  const router = useRouter();
  const {
    profile,
    isLoading,
    error,
    mutationStates,
    updateEmail,
    changeUsername,
  } = useUserProfile();

  const {
    linkStatus,
    isLoading: isLinkLoading,
    error: linkError,
    initiateLink,
    unlinkAccount,
    isUnlinking,
  } = useBattleNetLink();

  // Route protection
  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-red-500 bg-red-100 p-4 rounded-lg">
          {error instanceof Error ? error.message : "An error occurred"}
        </div>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-yellow-500">No profile data available</div>
      </div>
    );
  }

  const handleBattleNetUnlink = async () => {
    try {
      await unlinkAccount();
    } catch (error) {
      console.error("Failed to unlink Battle.net:", error);
    }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      {/* Header with avatar and basic info */}
      <div className="flex items-center gap-6 mb-8 pb-6 border-b border-gray-800">
        <div className="flex items-center justify-center w-20 h-20 bg-blue-500 rounded-full text-3xl font-bold text-white">
          {profile.username.charAt(0).toUpperCase()}
        </div>
        <div>
          <h1 className="text-2xl font-bold">{profile.username}</h1>
          <p className="text-gray-400">{profile.email}</p>
          {profile.battle_tag && (
            <p className="text-blue-400">{profile.battle_tag}</p>
          )}
        </div>
      </div>

      {/* Tab Navigation */}
      <TabNavigation activeTab={activeTab} setActiveTab={setActiveTab} />

      {/* Tab content */}
      <div className="mt-6">
        {/* Overview Tab */}
        {activeTab === "overview" && (
          <div className="space-y-6">
            <Card className="bg-[#131e33] border-gray-800 p-6">
              <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-blue-500"
                >
                  <circle cx="12" cy="8" r="5" />
                  <path d="M20 21a8 8 0 0 0-16 0" />
                </svg>
                Personal Information
              </h2>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-gray-400">Username</p>
                  <p className="font-medium">{profile.username}</p>
                </div>
                <div>
                  <p className="text-gray-400">Email</p>
                  <p className="font-medium">{profile.email}</p>
                </div>
                {profile.battle_tag && (
                  <div>
                    <p className="text-gray-400">Battle Tag</p>
                    <p className="font-medium text-blue-400">
                      {profile.battle_tag}
                    </p>
                  </div>
                )}
                <div>
                  <p className="text-gray-400">Member Since</p>
                  <p className="font-medium">February 28, 2025</p>
                </div>
              </div>
            </Card>

            <Card className="bg-[#131e33] border-gray-800 p-6">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-bold flex items-center gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="text-blue-500"
                  >
                    <path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48" />
                  </svg>
                  Your Characters
                </h2>
                <button
                  className="text-blue-500 flex items-center gap-1 hover:underline"
                  onClick={() => setActiveTab("characters")}
                >
                  View all
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <path d="M5 12h14" />
                    <path d="m12 5 7 7-7 7" />
                  </svg>
                </button>
              </div>

              <WoWProfile limit={1} />
            </Card>

            <Card className="bg-[#131e33] border-gray-800 p-6">
              <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-blue-500"
                >
                  <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
                  <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
                </svg>
                Connected Accounts
              </h2>

              <div className="flex justify-between items-center">
                <div className="flex items-center gap-4">
                  <Image
                    src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                    alt="Battle.net"
                    width={40}
                    height={40}
                  />
                  <div>
                    <h3 className="font-semibold">Battle.net</h3>
                    {linkStatus?.linked ? (
                      <p>Connected as: {profile.battle_tag}</p>
                    ) : (
                      <p className="text-gray-400">Not connected</p>
                    )}
                  </div>
                </div>

                {linkStatus?.linked ? (
                  <Button
                    variant="destructive"
                    onClick={handleBattleNetUnlink}
                    disabled={isUnlinking}
                  >
                    {isUnlinking ? "Unlinking..." : "Disconnect"}
                  </Button>
                ) : (
                  <Button
                    onClick={initiateLink}
                    disabled={isLinkLoading}
                    className="flex items-center gap-2"
                  >
                    <Image
                      src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                      alt="Battle.net"
                      width={20}
                      height={20}
                    />
                    {isLinkLoading ? "Connecting..." : "Connect"}
                  </Button>
                )}
              </div>
            </Card>
          </div>
        )}

        {/* Characters Tab */}
        {activeTab === "characters" && (
          <Card className="bg-[#131e33] border-gray-800 p-6">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="text-blue-500"
              >
                <path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48" />
              </svg>
              Your Characters
            </h2>

            <p className="text-gray-400 mb-6">
              Only characters from the same region as your Battle.net account
              will be displayed, and only level 80 characters will be shown.
            </p>

            <WoWProfile showTooltips={true} />
          </Card>
        )}

        {/* Account Tab */}
        {activeTab === "account" && (
          <Card className="bg-[#131e33] border-gray-800 p-6">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="text-blue-500"
              >
                <circle cx="12" cy="8" r="5" />
                <path d="M20 21a8 8 0 0 0-16 0" />
              </svg>
              Account Information
            </h2>

            <div className="grid grid-cols-2 gap-4 mb-6">
              <div>
                <p className="text-gray-400">Username</p>
                <p className="font-medium">{profile.username}</p>
              </div>
              <div>
                <p className="text-gray-400">Email</p>
                <p className="font-medium">{profile.email}</p>
              </div>
              {profile.battle_tag && (
                <div>
                  <p className="text-gray-400">Battle Tag</p>
                  <p className="font-medium text-blue-400">
                    {profile.battle_tag}
                  </p>
                </div>
              )}
            </div>

            <div className="flex gap-4">
              <button
                className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md flex items-center gap-2"
                onClick={() => setActiveTab("username")}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M12 20h9" />
                  <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
                </svg>
                Change Username
              </button>

              <button
                className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md flex items-center gap-2"
                onClick={() => setActiveTab("email")}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <rect width="20" height="16" x="2" y="4" rx="2" />
                  <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" />
                </svg>
                Change Email
              </button>
            </div>
          </Card>
        )}

        {/* Security Tab */}
        {activeTab === "security" && (
          <Card className="bg-[#131e33] border-gray-800 p-6">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="text-blue-500"
              >
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10" />
              </svg>
              Security Settings
            </h2>

            <div className="mb-8">
              <h3 className="font-semibold text-lg mb-2">Change Password</h3>
              <p className="text-gray-400 mb-4">
                It is a good idea to use a strong password that you do not use
                elsewhere
              </p>
              <button
                className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md flex items-center gap-2"
                onClick={() => setActiveTab("password")}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
                Change Password
              </button>
            </div>

            <div>
              <h3 className="font-semibold text-lg mb-2">Delete Account</h3>
              <p className="text-gray-400 mb-4">
                This action is permanent and cannot be undone. All your data
                will be deleted.
              </p>
              <button
                className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-md flex items-center gap-2"
                onClick={() => setActiveTab("delete")}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M3 6h18" />
                  <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
                  <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
                  <path d="M10 11v6" />
                  <path d="M14 11v6" />
                </svg>
                Delete Account
              </button>
            </div>
          </Card>
        )}

        {/* Connections Tab */}
        {activeTab === "connections" && (
          <Card className="bg-[#131e33] border-gray-800 p-6">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="text-blue-500"
              >
                <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
                <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
              </svg>
              Connected Accounts
            </h2>

            <div className="flex justify-between items-center mb-6">
              <div className="flex items-center gap-4">
                <Image
                  src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                  alt="Battle.net"
                  width={40}
                  height={40}
                />
                <div>
                  <h3 className="font-semibold">Battle.net</h3>
                  {linkStatus?.linked ? (
                    <p>Connected as: {profile.battle_tag}</p>
                  ) : (
                    <p className="text-gray-400">Not connected</p>
                  )}
                </div>
              </div>

              {linkStatus?.linked ? (
                <Button
                  variant="destructive"
                  onClick={handleBattleNetUnlink}
                  disabled={isUnlinking}
                >
                  {isUnlinking ? "Unlinking..." : "Disconnect"}
                </Button>
              ) : (
                <Button
                  onClick={initiateLink}
                  disabled={isLinkLoading}
                  className="flex items-center gap-2"
                >
                  <Image
                    src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                    alt="Battle.net"
                    width={20}
                    height={20}
                  />
                  {isLinkLoading ? "Connecting..." : "Connect"}
                </Button>
              )}
            </div>

            <p className="text-gray-400">
              Disconnecting your Battle.net account will remove access to all
              character data. You will need to reconnect your account to view
              your characters again.
            </p>
          </Card>
        )}

        {/* Additional Tabs for Specific Actions */}
        {activeTab === "username" && <ChangeUsername />}
        {activeTab === "email" && <ChangeEmail />}
        {activeTab === "password" && <ChangePassword />}
        {activeTab === "delete" && <DeleteAccount />}
      </div>
    </div>
  );
};

export default Profile;
