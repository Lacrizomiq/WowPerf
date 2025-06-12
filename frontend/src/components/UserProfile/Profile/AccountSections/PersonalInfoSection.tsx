// PersonalInfoSection.tsx - Section informations personnelles (lecture seule)
import React from "react";
import { UserProfile } from "@/libs/userService";

interface PersonalInfoSectionProps {
  profile: UserProfile;
}

const PersonalInfoSection: React.FC<PersonalInfoSectionProps> = ({
  profile,
}) => {
  return (
    <div className="bg-slate-800/30 border border-slate-700 rounded-lg">
      <div className="p-6 border-b border-slate-700">
        <h3 className="text-lg font-semibold">Personal Information</h3>
        <p className="text-sm text-slate-400">
          Your current account information
        </p>
      </div>
      <div className="p-6 space-y-4">
        <div className="space-y-2">
          <label
            htmlFor="current-username"
            className="text-sm font-medium text-slate-200"
          >
            Username
          </label>
          <input
            id="current-username"
            value={profile.username}
            disabled
            className="w-full px-3 py-2 bg-slate-800/50 border border-slate-700 text-slate-400 rounded-md"
          />
        </div>
        <div className="space-y-2">
          <label
            htmlFor="current-email"
            className="text-sm font-medium text-slate-200"
          >
            Email Address
          </label>
          <input
            id="current-email"
            value={profile.email}
            disabled
            className="w-full px-3 py-2 bg-slate-800/50 border border-slate-700 text-slate-400 rounded-md"
          />
        </div>
        {profile.battle_tag && (
          <div className="space-y-2">
            <label
              htmlFor="current-battle-tag"
              className="text-sm font-medium text-slate-200"
            >
              Battle Tag
            </label>
            <input
              id="current-battle-tag"
              value={profile.battle_tag}
              disabled
              className="w-full px-3 py-2 bg-slate-800/50 border border-slate-700 text-blue-400 rounded-md"
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default PersonalInfoSection;
