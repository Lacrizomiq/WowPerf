"use client";

import React, { useState } from "react";

interface DeleteAccountProps {
  onDeleteAccount: () => void;
  isDeleting: boolean;
}

const DeleteAccount: React.FC<DeleteAccountProps> = ({
  onDeleteAccount,
  isDeleting,
}) => {
  const [confirmDelete, setConfirmDelete] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (confirmDelete) {
      onDeleteAccount();
    } else {
      alert("Please confirm that you want to delete your account.");
    }
  };

  return (
    <section className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      <h2 className="text-2xl font-bold mb-4 text-gray-800 dark:text-gray-200">
        Delete Account
      </h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="flex items-center">
            <input
              type="checkbox"
              checked={confirmDelete}
              onChange={(e) => setConfirmDelete(e.target.checked)}
              className="form-checkbox h-5 w-5 text-red-600"
            />
            <span className="ml-2 text-gray-700 dark:text-gray-300">
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
