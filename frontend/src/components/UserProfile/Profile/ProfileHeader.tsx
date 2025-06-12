// ProfileHeader.tsx - Version HTML native sans shadcn/ui
import React from "react";
import { UserProfile } from "@/libs/userService";

interface ProfileHeaderProps {
  profile: UserProfile;
}

const ProfileHeader: React.FC<ProfileHeaderProps> = ({ profile }) => {
  return (
    <>
      {/* Page Header */}
      <header className="pt-8 pb-6 px-4 md:px-8 border-b border-slate-800">
        <div className="container mx-auto">
          <h1 className="text-3xl md:text-4xl font-bold mb-2">
            Account Management
          </h1>
          <p className="text-slate-400 text-base md:text-lg">
            Manage your profile, characters and connections
          </p>
        </div>
      </header>

      {/* Profile Summary Card */}
      <div className="bg-slate-800/30 border border-slate-700 rounded-lg mb-6 mt-6 mx-8">
        <div className="p-6">
          <div className="flex flex-col md:flex-row items-center md:items-start gap-6">
            {/* Avatar */}
            <div className="relative w-24 h-24 rounded-full overflow-hidden bg-slate-700 border-2 border-purple-600 flex items-center justify-center">
              <span className="text-3xl font-bold text-white">
                {profile.username.charAt(0).toUpperCase()}
              </span>
            </div>

            {/* User Information */}
            <div className="flex-1 text-center md:text-left">
              <h2 className="text-2xl font-bold">{profile.username}</h2>
              <p className="text-slate-400">{profile.email}</p>
              {profile.battle_tag && (
                <p className="text-blue-400">{profile.battle_tag}</p>
              )}
              <p className="text-sm text-slate-500 mt-1">
                Member since: {new Date().toLocaleDateString()}{" "}
                {/* TODO: Use real created_at */}
              </p>
              <div className="mt-2 flex flex-wrap gap-2 justify-center md:justify-start">
                {profile.battle_tag && (
                  <span className="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-slate-800 text-slate-400 border border-slate-600">
                    Battle.net Linked
                  </span>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default ProfileHeader;
