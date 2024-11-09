"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { useUserProfile } from "@/hooks/useUserProfile";
import { UserServiceError, UserErrorCode } from "@/libs/userService";

const ChangePassword: React.FC = () => {
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const router = useRouter();

  const { changePassword, mutationStates } = useUserProfile();
  const { isPending: isChanging } = mutationStates.changePassword;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (newPassword !== confirmPassword) {
      toast.error("New passwords do not match");
      return;
    }

    // Validation du mot de passe
    const passwordRegex = /^(?=.*[!@#$%^&*()_+]).{8,}$/;
    if (!passwordRegex.test(newPassword)) {
      toast.error(
        "Password must be at least 8 characters and contain at least one special character"
      );
      return;
    }

    try {
      const result = await changePassword(currentPassword, newPassword);

      if (result.success) {
        toast.success("Password changed successfully!");
        router.push("/profile");
      }
    } catch (error) {
      if (error instanceof UserServiceError) {
        switch (error.code) {
          case UserErrorCode.INVALID_PASSWORD:
            toast.error("Current password is incorrect");
            break;
          case UserErrorCode.UNAUTHORIZED:
            toast.error("Session expired. Please login again");
            router.push("/login");
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
      <h2 className="text-2xl font-bold mb-4 text-white">Change Password</h2>
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
            required
            disabled={isChanging}
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
            required
            pattern="(?=.*[!@#$%^&*()_+]).{8,}"
            disabled={isChanging}
          />
          <p className="mt-1 text-sm text-gray-400">
            Must be at least 8 characters with at least one special character
          </p>
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
            required
            disabled={isChanging}
          />
        </div>
        <button
          type="submit"
          disabled={isChanging}
          className="w-full px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:bg-blue-300"
        >
          {isChanging ? "Changing Password..." : "Change Password"}
        </button>
      </form>
    </section>
  );
};

export default ChangePassword;
