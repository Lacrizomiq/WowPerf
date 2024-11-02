"use client";

import React, { useState } from "react";
import toast from "react-hot-toast";

interface DeleteAccountProps {
  onDeleteAccount: () => Promise<void>;
  isDeleting: boolean;
}

const DeleteAccount: React.FC<DeleteAccountProps> = ({
  onDeleteAccount,
  isDeleting,
}) => {
  const [confirmDelete, setConfirmDelete] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (confirmDelete) {
      try {
        await toast.promise(onDeleteAccount(), {
          loading: "Deleting account...",
          success: "Account deleted successfully!",
          error: "Failed to delete account",
        });
      } catch (error) {
        console.error("Error deleting account:", error);
      }
    } else {
      alert("Please confirm that you want to delete your account.");
    }
  };

  return (
    <section className="bg-deep-blue border border-gray-800 shadow rounded-lg p-6">
      <h2 className="text-2xl font-bold mb-4 text-white">Delete Account</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="flex items-center">
            <input
              type="checkbox"
              checked={confirmDelete}
              onChange={(e) => setConfirmDelete(e.target.checked)}
              className="form-checkbox h-5 w-5 text-red-600"
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
          {isDeleting ? "Deleting..." : "Delete Account"}
        </button>
      </form>
    </section>
  );
};

export default DeleteAccount;
