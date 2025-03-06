// AccountTab.tsx
// Account information and settings
import React from "react";
import { Card } from "@/components/ui/card";
import { UserProfile } from "@/libs/userService";

interface AccountTabProps {
  profile: UserProfile;
  onNavigate: (tab: string) => void;
}

const AccountTab: React.FC<AccountTabProps> = ({ profile, onNavigate }) => {
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
          <circle cx="12" cy="8" r="5" />
          <path d="M20 21a8 8 0 0 0-16 0" />
        </svg>
        Account Information
      </h2>

      <div className="grid gap-4 mb-6">
        <div>
          <p className="text-gray-400">Username</p>
          <p className="font-medium">{profile.username}</p>
        </div>
        <div>
          <p className="text-gray-400">Email</p>
          <p className="font-medium">{profile.email}</p>
        </div>
        {profile.battle_tag && (
          <div>
            <p className="text-gray-400">Battle Tag</p>
            <p className="font-medium text-blue-400">{profile.battle_tag}</p>
          </div>
        )}
      </div>

      <div className="flex gap-4">
        <button
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md flex items-center gap-2"
          onClick={() => onNavigate("username")}
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
            <path d="M12 20h9" />
            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
          </svg>
          Change Username
        </button>

        <button
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md flex items-center gap-2"
          onClick={() => onNavigate("email")}
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
            <rect width="20" height="16" x="2" y="4" rx="2" />
            <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" />
          </svg>
          Change Email
        </button>
      </div>
    </Card>
  );
};

export default AccountTab;
