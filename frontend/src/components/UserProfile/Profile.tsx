"use client";

import React, { useState } from "react";
import Sidebar from "./Sidebar";
import PersonalInfo from "@/components/UserProfile/PersonalInfo";
import ChangePassword from "@/components/UserProfile/Update/ChangePassword";
import ChangeEmail from "@/components/UserProfile/Update/ChangeEmail";
import DeleteAccount from "@/components/UserProfile/Update/DeleteAccount";
import ChangeUsername from "@/components/UserProfile/Update/ChangeUsername";
import { useUserProfile } from "@/hooks/useUserProfile";

const Profile: React.FC = () => {
  const [activeSection, setActiveSection] = useState("profile");
  const {
    profile,
    isLoading,
    error,
    updateEmail,
    changePassword,
    deleteAccount,
    changeUsername,
    isUpdatingEmail,
    isChangingPassword,
    isDeletingAccount,
    isChangingUsername,
  } = useUserProfile();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  const renderSection = () => {
    switch (activeSection) {
      case "profile":
        return (
          profile && (
            <PersonalInfo
              profile={profile}
              onUpdate={updateEmail}
              isUpdating={isUpdatingEmail}
            />
          )
        );
      case "change-username":
        return (
          <ChangeUsername
            onUpdateUsername={changeUsername}
            isUpdating={isChangingUsername}
          />
        );
      case "change-password":
        return (
          <ChangePassword
            onChangePassword={changePassword}
            isChanging={isChangingPassword}
          />
        );
      case "change-email":
        return (
          <ChangeEmail
            onUpdateEmail={updateEmail}
            isUpdating={isUpdatingEmail}
          />
        );
      case "delete-account":
        return (
          <DeleteAccount
            onDeleteAccount={deleteAccount}
            isDeleting={isDeletingAccount}
          />
        );
      default:
        return null;
    }
  };

  return (
    <div className="flex h-screen bg-blue-800 dark:bg-gray-900">
      <Sidebar onSectionChange={setActiveSection} />
      <main className="flex-1 p-8 overflow-y-auto">
        <h1 className="text-3xl font-bold mb-8 text-white">Profile Settings</h1>
        {renderSection()}
      </main>
    </div>
  );
};

export default Profile;
