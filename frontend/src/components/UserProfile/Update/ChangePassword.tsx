"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";

interface ChangePasswordProps {
  onChangePassword: (
    currentPassword: string,
    newPassword: string
  ) => Promise<void>;
  isChanging: boolean;
}

const ChangePassword: React.FC<ChangePasswordProps> = ({
  onChangePassword,
  isChanging,
}) => {
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newPassword === confirmPassword) {
      try {
        await toast.promise(onChangePassword(currentPassword, newPassword), {
          loading: "Changing password...",
          success: "Password changed successfully!",
          error: "Failed to change password",
        });
        router.push("/profile");
      } catch (error) {
        console.error("Error changing password:", error);
      }
    } else {
      // Handle password mismatch error
      alert("New passwords do not match");
    }
  };

  return (
    <section className="bg-[#374151] dark:bg-gray-800 shadow rounded-lg p-6">
      <h2 className="text-2xl font-bold mb-4 text-white dark:text-gray-200">
        Change Password
      </h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="currentPassword"
            className="block text-sm font-medium text-white dark:text-gray-300 mb-2"
          >
            Current Password
          </label>
          <input
            type="password"
            id="currentPassword"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-gray-800 dark:text-gray-200"
          />
        </div>
        <div>
          <label
            htmlFor="newPassword"
            className="block text-sm font-medium text-white dark:text-gray-300 mb-2"
          >
            New Password
          </label>
          <input
            type="password"
            id="newPassword"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-gray-800 dark:text-gray-200"
          />
        </div>
        <div>
          <label
            htmlFor="confirmPassword"
            className="block text-sm font-medium text-white dark:text-gray-300 mb-2"
          >
            Confirm New Password
          </label>
          <input
            type="password"
            id="confirmPassword"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-gray-800 dark:text-gray-200"
          />
        </div>
        <button
          type="submit"
          className="w-full px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
        >
          Change Password
        </button>
      </form>
    </section>
  );
};

export default ChangePassword;
