"use client";

import React, { useState } from "react";

interface ChangeEmailProps {
  onUpdateEmail: (newEmail: string) => void;
  isUpdating: boolean;
}

const ChangeEmail: React.FC<ChangeEmailProps> = ({
  onUpdateEmail,
  isUpdating,
}) => {
  const [newEmail, setNewEmail] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onUpdateEmail(newEmail);
    alert("Email updated");
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
