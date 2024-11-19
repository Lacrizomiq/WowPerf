"use client";

import React from "react";
import ProfileSidebar from "./ProfileSidebar";
import PersonalInfo from "@/components/UserProfile/PersonalInfo";
import { useUserProfile } from "@/hooks/useUserProfile";
import { useAuth } from "@/providers/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

const Profile: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const router = useRouter();

  // Protection de route
  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  const {
    profile,
    isLoading,
    error,
    mutationStates,
    updateEmail,
    changeUsername,
  } = useUserProfile();

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

  const personalInfoMutationStates = {
    changeUsername: mutationStates.changeUsername,
    updateEmail: mutationStates.updateEmail,
    changePassword: mutationStates.changePassword,
    deleteAccount: mutationStates.deleteAccount,
  };

  return (
    <div className="flex min-h-screen w-full bg-black">
      <ProfileSidebar />
      <main className="flex-1 p-8">
        <PersonalInfo
          profile={profile}
          mutationStates={personalInfoMutationStates}
          onUpdateEmail={updateEmail}
          onChangeUsername={changeUsername}
        />
      </main>
    </div>
  );
};

export default Profile;
