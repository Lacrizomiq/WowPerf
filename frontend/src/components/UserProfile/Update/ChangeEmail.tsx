"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { useUserProfile } from "@/hooks/useUserProfile";
import { UserServiceError, UserErrorCode } from "@/libs/userService";

const ChangeEmail: React.FC = () => {
  const [newEmail, setNewEmail] = useState("");
  const router = useRouter();

  const { updateEmail, mutationStates } = useUserProfile();
  const { isPending: isUpdating } = mutationStates.updateEmail;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const result = await updateEmail(newEmail);

      if (result.success) {
        toast.success("Email updated successfully!");
        router.push("/profile");
      } else {
        toast.error(result.error || "Failed to update email");
      }
    } catch (error) {
      if (error instanceof UserServiceError) {
        switch (error.code) {
          case UserErrorCode.INVALID_EMAIL:
            toast.error("Please enter a valid email address");
            break;
          case UserErrorCode.EMAIL_EXISTS:
            toast.error("This email is already in use");
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
      <h2 className="text-2xl font-bold mb-4 text-white">Change Email</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="newEmail"
            className="block text-sm font-medium text-white dark:text-gray-300 mb-2"
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
