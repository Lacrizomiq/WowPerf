"use client";

import React from "react";
import { useUserProfile } from "@/hooks/useUserProfile";
import ChangeUsername from "@/components/UserProfile/Update/ChangeUsername";
import Sidebar from "@/components/UserProfile/ProfileSidebar";

const UpdateUsernamePage: React.FC = () => {
  const { changeUsername, isChangingUsername } = useUserProfile();

  return (
    <div className="flex min-h-screen bg-black">
      <Sidebar />
      <main className="flex-1 p-8">
        <ChangeUsername
          onUpdateUsername={changeUsername}
          isUpdating={isChangingUsername}
        />
      </main>
    </div>
  );
};

export default UpdateUsernamePage;
