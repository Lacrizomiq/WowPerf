"use client";

import React from "react";
import ProfileSidebar from "./ProfileSidebar";
import PersonalInfo from "@/components/UserProfile/PersonalInfo";
import { useUserProfile } from "@/hooks/useUserProfile";
import { useAuth } from "@/providers/AuthContext";
import { useRequireAuth } from "@/providers/AuthContext";

const Profile: React.FC = () => {
  // Protection de la route
  useRequireAuth();

  const { user } = useAuth();
  const { profile, isLoading, error, mutationStates } = useUserProfile();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        Loading...
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen text-red-500">
        {error instanceof Error ? error.message : "An error occurred"}
      </div>
    );
  }

  return (
    <div className="flex min-h-screen w-full bg-black">
      <ProfileSidebar />
      <main className="flex-1 p-8">
        {profile && (
          <PersonalInfo
            profile={profile}
            isBattleNetLinked={!!user?.battlenet_id}
          />
        )}
      </main>
    </div>
  );
};

export default Profile;
