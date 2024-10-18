"use client";

import React from "react";
import { useUserProfile } from "@/hooks/useUserProfile";
import ChangePassword from "@/components/UserProfile/Update/ChangePassword";
import Sidebar from "@/components/UserProfile/Sidebar";

const UpdatePasswordPage: React.FC = () => {
  const { changePassword, isChangingPassword } = useUserProfile();

  return (
    <div className="flex min-h-screen bg-gradient-to-br from-[#1a202c] to-[#2d3748] dark:bg-gray-900">
      <Sidebar />
      <main className="flex-1 p-8">
        <ChangePassword
          onChangePassword={changePassword}
          isChanging={isChangingPassword}
        />
      </main>
    </div>
  );
};

export default UpdatePasswordPage;
