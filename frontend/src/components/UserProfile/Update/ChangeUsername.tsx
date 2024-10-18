import React, { useState } from "react";

interface ChangeUsernameProps {
  onUpdateUsername: (newUsername: string) => void;
  isUpdating: boolean;
}

const ChangeUsername: React.FC<ChangeUsernameProps> = ({
  onUpdateUsername,
  isUpdating,
}) => {
  const [newUsername, setNewUsername] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onUpdateUsername(newUsername);
  };

  return (
    <section className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      <h2 className="text-2xl font-bold mb-4 text-gray-800 dark:text-gray-200">
        Change Username
      </h2>
      <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
        Please note that you can only change your username once every 30 days.
      </p>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="newUsername"
            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            New Username
          </label>
          <input
            type="text"
            id="newUsername"
            value={newUsername}
            onChange={(e) => setNewUsername(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-gray-800 dark:text-gray-200"
          />
        </div>
        <button
          type="submit"
          disabled={isUpdating}
          className="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 transition-colors"
        >
          {isUpdating ? "Updating..." : "Update Username"}
        </button>
      </form>
    </section>
  );
};

export default ChangeUsername;
