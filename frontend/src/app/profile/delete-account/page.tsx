"use client";

import React from "react";
import { useUserProfile } from "@/hooks/useUserProfile";
import DeleteAccount from "@/components/UserProfile/Update/DeleteAccount";
import Sidebar from "@/components/UserProfile/Sidebar";

const DeleteAccountPage: React.FC = () => {
  const { deleteAccount, isDeletingAccount } = useUserProfile();

  return (
    <div className="flex min-h-screen bg-gradient-to-br from-[#1a202c] to-[#2d3748] dark:bg-gray-900">
      <Sidebar />
      <main className="flex-1 p-8">
        <DeleteAccount
          onDeleteAccount={deleteAccount}
          isDeleting={isDeletingAccount}
        />
      </main>
    </div>
  );
};

export default DeleteAccountPage;
