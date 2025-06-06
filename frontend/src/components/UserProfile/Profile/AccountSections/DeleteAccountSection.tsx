// DeleteAccountSection.tsx - Version HTML native sans shadcn/ui
import React, { useState } from "react";
import toast from "react-hot-toast";
import { useUserProfile } from "@/hooks/useUserProfile";
import { UserServiceError, UserErrorCode } from "@/libs/userService";
import { useAuth } from "@/providers/AuthContext";

const DeleteAccountSection: React.FC = () => {
  const [confirmDelete, setConfirmDelete] = useState(false);
  const { logout } = useAuth();
  const { deleteAccount, mutationStates } = useUserProfile();
  const { isPending: isDeleting } = mutationStates.deleteAccount;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!confirmDelete) {
      toast.error("Please confirm that you want to delete your account");
      return;
    }

    try {
      const result = await deleteAccount();

      if (result.success) {
        toast.success("Account deleted successfully");
        await logout(); // use the logout function from the AuthContext
      }
    } catch (error) {
      if (error instanceof UserServiceError) {
        switch (error.code) {
          case UserErrorCode.UNAUTHORIZED:
            toast.error("Session expired. Please login again");
            await logout();
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
    <div className="bg-red-900/20 border border-red-700 rounded-lg">
      <div className="p-6 border-b border-red-700">
        <h3 className="text-lg font-semibold text-red-400 flex items-center gap-2">
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 16.5c-.77.833.192 2.5 1.732 2.5z"
            />
          </svg>
          Danger Zone
        </h3>
        <p className="text-sm text-red-300">
          This action is permanent and cannot be undone. All your data will be
          deleted.
        </p>
      </div>
      <div className="p-6">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="bg-red-900/30 border border-red-700 rounded-lg p-4">
            <p className="text-red-300 text-sm font-medium">Warning!</p>
            <p className="text-red-200 text-sm mt-1">
              This action cannot be undone. All your data will be permanently
              deleted.
            </p>
          </div>

          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              id="confirm-delete"
              checked={confirmDelete}
              onChange={(e) => setConfirmDelete(e.target.checked)}
              disabled={isDeleting}
              className="h-4 w-4 text-red-600 rounded border-slate-600 bg-slate-800 focus:ring-red-500"
            />
            <label
              htmlFor="confirm-delete"
              className="text-sm text-slate-300 cursor-pointer"
            >
              I understand that this action is irreversible and will permanently
              delete my account.
            </label>
          </div>

          <button
            type="submit"
            disabled={isDeleting || !confirmDelete}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 disabled:bg-red-300 disabled:cursor-not-allowed text-white rounded-md font-medium transition-colors"
          >
            {isDeleting ? "Deleting Account..." : "Delete Account"}
          </button>
        </form>
      </div>
    </div>
  );
};

export default DeleteAccountSection;
