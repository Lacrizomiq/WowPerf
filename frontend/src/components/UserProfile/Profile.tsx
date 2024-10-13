"use client";

import React from "react";
import { useRouter } from "next/navigation";
import Sidebar from "./Sidebar";
import PersonalInfo from "@/components/UserProfile/PersonalInfo";
import ChangePassword from "@/components/UserProfile/Update/ChangePassword";
import { useUserProfile } from "@/hooks/useUserProfile";

const Profile: React.FC = () => {
  const router = useRouter();
  const {
    profile,
    isLoading,
    error,
    updateEmail,
    changePassword,
    deleteAccount,
    isUpdatingEmail,
    isChangingPassword,
    isDeletingAccount,
  } = useUserProfile();

  const handleDeleteAccount = () => {
    deleteAccount();
    router.push("/signup");
  };

  console.log(profile);

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div className="flex h-screen bg-blue-800 dark:bg-gray-900">
      <Sidebar />
      <main className="flex-1 p-8 overflow-y-auto">
        <h1 className="text-3xl font-bold mb-8 text-white">Profile Settings</h1>
        <div className="space-y-6">
          {profile && (
            <PersonalInfo
              profile={profile}
              onUpdate={updateEmail}
              isUpdating={isUpdatingEmail}
            />
          )}
          <ChangePassword
            onChangePassword={changePassword}
            isChanging={isChangingPassword}
          />
          <button
            onClick={handleDeleteAccount}
            disabled={isDeletingAccount}
            className="bg-red-500 text-white px-4 py-2 rounded disabled:bg-red-300"
          >
            {isDeletingAccount ? "Deleting..." : "Delete Account"}
          </button>
        </div>
      </main>
    </div>
  );
};

export default Profile;
