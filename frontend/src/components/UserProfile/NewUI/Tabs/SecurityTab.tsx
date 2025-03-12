// SecurityTab.tsx
// Security settings including password change and account deletion
import React from "react";
import { Card } from "@/components/ui/card";

interface SecurityTabProps {
  onNavigate: (tab: string) => void;
}

const SecurityTab: React.FC<SecurityTabProps> = ({ onNavigate }) => {
  return (
    <Card className="bg-[#131e33] border-gray-800 p-6">
      <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          className="text-blue-500"
        >
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10" />
        </svg>
        Security Settings
      </h2>

      <div className="mb-8">
        <h3 className="font-semibold text-lg mb-2">Change Password</h3>
        <p className="text-gray-400 mb-4">
          It is a good idea to use a strong password that you do not use
          elsewhere
        </p>
        <button
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md flex items-center gap-2"
          onClick={() => onNavigate("password")}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
            <path d="M7 11V7a5 5 0 0 1 10 0v4" />
          </svg>
          Change Password
        </button>
      </div>

      <div>
        <h3 className="font-semibold text-lg mb-2">Delete Account</h3>
        <p className="text-gray-400 mb-4">
          This action is permanent and cannot be undone. All your data will be
          deleted.
        </p>
        <button
          className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-md flex items-center gap-2"
          onClick={() => onNavigate("delete")}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M3 6h18" />
            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
            <path d="M10 11v6" />
            <path d="M14 11v6" />
          </svg>
          Delete Account
        </button>
      </div>
    </Card>
  );
};

export default SecurityTab;
