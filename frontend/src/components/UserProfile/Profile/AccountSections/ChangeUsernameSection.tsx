// ChangeUsernameSection.tsx - Version HTML native sans shadcn/ui
import React, { useState } from "react";
import toast from "react-hot-toast";
import { useUserProfile } from "@/hooks/useUserProfile";
import { UserServiceError, UserErrorCode } from "@/libs/userService";

const ChangeUsernameSection: React.FC = () => {
  const [newUsername, setNewUsername] = useState("");
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
        setNewUsername(""); // Reset form
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
    <div className="bg-slate-800/30 border border-slate-700 rounded-lg">
      <div className="p-6 border-b border-slate-700">
        <h3 className="text-lg font-semibold">Change Username</h3>
        <p className="text-sm text-slate-400">
          Update your username. You can only change it once every 30 days.
        </p>
      </div>
      <div className="p-6">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label
              htmlFor="new-username"
              className="text-sm font-medium text-slate-200"
            >
              New Username
            </label>
            <input
              id="new-username"
              type="text"
              value={newUsername}
              onChange={(e) => setNewUsername(e.target.value)}
              placeholder="Enter a new username"
              className="w-full px-3 py-2 bg-slate-800/50 border border-slate-700 text-white placeholder-slate-400 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              required
              minLength={3}
              maxLength={50}
              disabled={isUpdating}
            />
          </div>
          <button
            type="submit"
            disabled={isUpdating || !newUsername.trim()}
            className="px-4 py-2 bg-purple-600 hover:bg-purple-700 disabled:bg-indigo-600 disabled:cursor-not-allowed text-white rounded-md font-medium transition-colors"
          >
            {isUpdating ? "Updating..." : "Save Changes"}
          </button>
        </form>
      </div>
    </div>
  );
};

export default ChangeUsernameSection;
