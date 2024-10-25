"use client";

import React from "react";
import { useUserProfile } from "@/hooks/useUserProfile";
import ChangeEmail from "@/components/UserProfile/Update/ChangeEmail";
import Sidebar from "@/components/UserProfile/ProfileSidebar";

const UpdateEmailPage: React.FC = () => {
  const { updateEmail, isUpdatingEmail } = useUserProfile();

  return (
    <div className="flex min-h-screen bg-gradient-to-br from-[#1a202c] to-[#2d3748] dark:bg-gray-900">
      <Sidebar />
      <main className="flex-1 p-8">
        <ChangeEmail onUpdateEmail={updateEmail} isUpdating={isUpdatingEmail} />
      </main>
    </div>
  );
};

export default UpdateEmailPage;
