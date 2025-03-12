// ProfileHeader.tsx
// Component that displays the user's profile header with avatar and basic information
import React from "react";
import { UserProfile } from "@/libs/userService";

interface ProfileHeaderProps {
  profile: UserProfile;
}

const ProfileHeader: React.FC<ProfileHeaderProps> = ({ profile }) => {
  return (
    <div className="flex items-center gap-6 mb-8 pb-6 border-b border-gray-800">
      <div className="flex items-center justify-center w-20 h-20 bg-blue-500 rounded-full text-3xl font-bold text-white">
        {profile.username.charAt(0).toUpperCase()}
      </div>
      <div>
        <h1 className="text-2xl font-bold">{profile.username}</h1>
        <p className="text-gray-400">{profile.email}</p>
        {profile.battle_tag && (
          <p className="text-blue-400">{profile.battle_tag}</p>
        )}
      </div>
    </div>
  );
};

export default ProfileHeader;
