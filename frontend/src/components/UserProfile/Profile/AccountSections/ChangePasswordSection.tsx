// ChangePasswordSection.tsx - Version HTML native sans shadcn/ui
import React, { useState } from "react";
import toast from "react-hot-toast";
import { useUserProfile } from "@/hooks/useUserProfile";
import { UserServiceError, UserErrorCode } from "@/libs/userService";

const ChangePasswordSection: React.FC = () => {
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

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
        // Reset form
        setCurrentPassword("");
        setNewPassword("");
        setConfirmPassword("");
      }
    } catch (error) {
      if (error instanceof UserServiceError) {
        switch (error.code) {
          case UserErrorCode.INVALID_PASSWORD:
            toast.error("Current password is incorrect");
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
        <h3 className="text-lg font-semibold">Change Password</h3>
        <p className="text-sm text-slate-400">
          Change your password to secure your account
        </p>
      </div>
      <div className="p-6">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label
              htmlFor="current-password"
              className="text-sm font-medium text-slate-200"
            >
              Current Password
            </label>
            <input
              id="current-password"
              type="password"
              value={currentPassword}
              onChange={(e) => setCurrentPassword(e.target.value)}
              placeholder="Enter your current password"
              className="w-full px-3 py-2 bg-slate-800/50 border border-slate-700 text-white placeholder-slate-400 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              required
              disabled={isChanging}
            />
          </div>
          <div className="space-y-2">
            <label
              htmlFor="new-password"
              className="text-sm font-medium text-slate-200"
            >
              New Password
            </label>
            <input
              id="new-password"
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              placeholder="Enter a new password"
              className="w-full px-3 py-2 bg-slate-800/50 border border-slate-700 text-white placeholder-slate-400 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              required
              pattern="(?=.*[!@#$%^&*()_+]).{8,}"
              disabled={isChanging}
            />
            <p className="text-sm text-slate-400">
              Must be at least 8 characters with at least one special character
            </p>
          </div>
          <div className="space-y-2">
            <label
              htmlFor="confirm-password"
              className="text-sm font-medium text-slate-200"
            >
              Confirm New Password
            </label>
            <input
              id="confirm-password"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="Confirm your new password"
              className="w-full px-3 py-2 bg-slate-800/50 border border-slate-700 text-white placeholder-slate-400 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              required
              disabled={isChanging}
            />
          </div>
          <button
            type="submit"
            disabled={
              isChanging || !currentPassword || !newPassword || !confirmPassword
            }
            className="px-4 py-2 bg-purple-600 hover:bg-purple-700 disabled:bg-indigo-600 disabled:cursor-not-allowed text-white rounded-md font-medium transition-colors"
          >
            {isChanging ? "Changing Password..." : "Change Password"}
          </button>
        </form>
      </div>
    </div>
  );
};

export default ChangePasswordSection;
