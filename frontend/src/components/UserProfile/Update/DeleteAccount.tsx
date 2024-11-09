"use client";

import React, { useState } from "react";
import toast from "react-hot-toast";
import { useUserProfile } from "@/hooks/useUserProfile";
import { UserServiceError, UserErrorCode } from "@/libs/userService";
import { useAuth } from "@/providers/AuthContext";

const DeleteAccount: React.FC = () => {
  const [confirmDelete, setConfirmDelete] = useState(false);
  const { logout } = useAuth();
  const { deleteAccount, mutationStates } = useUserProfile();
  const { isPending: isDeleting } = mutationStates.deleteAccount;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!confirmDelete) {
      toast.error("Please confirm that you want to delete your account");
      return;
    }

    try {
      const result = await deleteAccount();

      if (result.success) {
        toast.success("Account deleted successfully");
        await logout(); // Utilise la fonction logout du AuthContext
      }
    } catch (error) {
      if (error instanceof UserServiceError) {
        switch (error.code) {
          case UserErrorCode.UNAUTHORIZED:
            toast.error("Session expired. Please login again");
            await logout();
            break;
          default:
            toast.error(error.message);
        }
      } else {
        toast.error("An unexpected error occurred");
      }
    }
  };

  return (
    <section className="bg-deep-blue border border-gray-800 shadow rounded-lg p-6">
      <h2 className="text-2xl font-bold mb-4 text-white">Delete Account</h2>
      <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
        <p className="font-bold">Warning!</p>
        <p>
          This action cannot be undone. All your data will be permanently
          deleted.
        </p>
      </div>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="flex items-center">
            <input
              type="checkbox"
              checked={confirmDelete}
              onChange={(e) => setConfirmDelete(e.target.checked)}
              className="form-checkbox h-5 w-5 text-red-600"
              disabled={isDeleting}
            />
            <span className="ml-2 text-white dark:text-gray-300">
              I understand that this action is irreversible and will permanently
              delete my account.
            </span>
          </label>
        </div>
        <button
          type="submit"
          disabled={isDeleting || !confirmDelete}
          className="w-full px-4 py-2 bg-red-500 text-white rounded-md hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 disabled:bg-red-300"
        >
          {isDeleting ? "Deleting Account..." : "Delete Account"}
        </button>
      </form>
    </section>
  );
};

export default DeleteAccount;
