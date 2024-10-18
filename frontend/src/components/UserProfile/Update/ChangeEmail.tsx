"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";

interface ChangeEmailProps {
  onUpdateEmail: (newEmail: string) => Promise<void>;
  isUpdating: boolean;
}

const ChangeEmail: React.FC<ChangeEmailProps> = ({
  onUpdateEmail,
  isUpdating,
}) => {
  const [newEmail, setNewEmail] = useState("");
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await toast.promise(onUpdateEmail(newEmail), {
        loading: "Updating email...",
        success: "Email updated successfully!",
        error: "Failed to update email",
      });
      router.push("/profile");
    } catch (error) {
      console.error("Error updating email:", error);
    }
  };

  return (
    <section className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      <h2 className="text-2xl font-bold mb-4 text-gray-800 dark:text-gray-200">
        Change Email
      </h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="newEmail"
            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            New Email
          </label>
          <input
            type="email"
            id="newEmail"
            value={newEmail}
            onChange={(e) => setNewEmail(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-gray-800 dark:text-gray-200"
            required
          />
        </div>
        <button
          type="submit"
          disabled={isUpdating}
          className="w-full px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:bg-blue-300"
        >
          {isUpdating ? "Updating..." : "Change Email"}
        </button>
      </form>
    </section>
  );
};

export default ChangeEmail;
