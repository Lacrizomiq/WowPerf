"use client";

import React from "react";
import ProfileSidebar from "./ProfileSidebar";
import PersonalInfo from "@/components/UserProfile/PersonalInfo";
import { useUserProfile } from "@/hooks/useUserProfile";

const Profile: React.FC = () => {
  const { profile, isLoading, error, updateEmail, isUpdatingEmail } =
    useUserProfile();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div className="flex min-h-screen w-full bg-black">
      <ProfileSidebar />
      <main className="flex-1 p-8">
        {profile && (
          <PersonalInfo
            profile={profile}
            onUpdate={updateEmail}
            isUpdating={isUpdatingEmail}
          />
        )}
      </main>
    </div>
  );
};

export default Profile;
