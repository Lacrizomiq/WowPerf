"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { useUserProfile } from "@/hooks/useUserProfile";
import { UserServiceError, UserErrorCode } from "@/libs/userService";

const ChangeUsername: React.FC = () => {
  const [newUsername, setNewUsername] = useState("");
  const router = useRouter();

  const { changeUsername, mutationStates } = useUserProfile();
  const { isPending: isUpdating } = mutationStates.changeUsername;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (newUsername.length < 3 || newUsername.length > 50) {
      toast.error("Username must be between 3 and 50 characters");
      return;
    }

    try {
      const result = await changeUsername(newUsername);

      if (result.success) {
        toast.success("Username updated successfully!");
        router.push("/profile");
      }
    } catch (error) {
      if (error instanceof UserServiceError) {
        switch (error.code) {
          case UserErrorCode.USERNAME_EXISTS:
            toast.error("This username is already taken");
            break;
          case UserErrorCode.USERNAME_CHANGE_LIMIT:
            toast.error("You can only change your username once every 30 days");
            break;
          case UserErrorCode.INVALID_USERNAME:
            toast.error("Invalid username format");
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
      <h2 className="text-2xl font-bold mb-4 text-white">Change Username</h2>
      <p className="text-sm text-white dark:text-gray-400 mb-4">
        Please note that you can only change your username once every 30 days.
      </p>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="newUsername"
            className="block text-sm font-medium text-white dark:text-gray-300 mb-2"
          >
            New Username
          </label>
          <input
            type="text"
            id="newUsername"
            value={newUsername}
            onChange={(e) => setNewUsername(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-gray-800 dark:text-gray-200"
            required
            minLength={3}
            maxLength={50}
            disabled={isUpdating}
          />
        </div>
        <button
          type="submit"
          disabled={isUpdating}
          className="w-full px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:bg-blue-300"
        >
          {isUpdating ? "Updating Username..." : "Update Username"}
        </button>
      </form>
    </section>
  );
};

export default ChangeUsername;
