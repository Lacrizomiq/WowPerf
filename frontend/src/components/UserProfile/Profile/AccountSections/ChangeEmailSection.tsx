// ChangeEmailSection.tsx - Version HTML native sans shadcn/ui
import React, { useState } from "react";
import toast from "react-hot-toast";
import { useUserProfile } from "@/hooks/useUserProfile";
import { UserServiceError, UserErrorCode } from "@/libs/userService";

const ChangeEmailSection: React.FC = () => {
  const [newEmail, setNewEmail] = useState("");
  const { updateEmail, mutationStates } = useUserProfile();
  const { isPending: isUpdating } = mutationStates.updateEmail;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const result = await updateEmail(newEmail);

      if (result.success) {
        toast.success("Email updated successfully!");
        setNewEmail(""); // Reset form
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
        <h3 className="text-lg font-semibold">Change Email Address</h3>
        <p className="text-sm text-slate-400">
          Update the email address associated with your account
        </p>
      </div>
      <div className="p-6">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label
              htmlFor="new-email"
              className="text-sm font-medium text-slate-200"
            >
              New Email Address
            </label>
            <input
              id="new-email"
              type="email"
              value={newEmail}
              onChange={(e) => setNewEmail(e.target.value)}
              placeholder="Enter a new email address"
              className="w-full px-3 py-2 bg-slate-800/50 border border-slate-700 text-white placeholder-slate-400 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              required
              disabled={isUpdating}
            />
          </div>
          <button
            type="submit"
            disabled={isUpdating || !newEmail.trim()}
            className="px-4 py-2 bg-purple-600 hover:bg-purple-700 disabled:bg-indigo-600 disabled:cursor-not-allowed text-white rounded-md font-medium transition-colors"
          >
            {isUpdating ? "Updating..." : "Save Changes"}
          </button>
        </form>
      </div>
    </div>
  );
};

export default ChangeEmailSection;
