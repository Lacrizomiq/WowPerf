// AccountTab.tsx - Nouveau tab account avec toutes les sections inline
import React from "react";
import { UserProfile } from "@/libs/userService";
import PersonalInfoSection from "./AccountSections/PersonalInfoSection";
import ChangeUsernameSection from "./AccountSections/ChangeUsernameSection";
import ChangeEmailSection from "./AccountSections/ChangeEmailSection";
import ChangePasswordSection from "./AccountSections/ChangePasswordSection";
import DeleteAccountSection from "./AccountSections/DeleteAccountSection";

interface AccountTabProps {
  profile: UserProfile;
  isActive: boolean;
}

const AccountTab: React.FC<AccountTabProps> = ({ profile, isActive }) => {
  if (!isActive) return null;

  return (
    <div className="space-y-6 grid grid-cols-1 md:grid-cols-2 gap-6">
      {/* Personal Information (read-only) */}
      <PersonalInfoSection profile={profile} />

      {/* Change Username */}
      <ChangeUsernameSection />

      {/* Change Email */}
      <ChangeEmailSection />

      {/* Change Password */}
      <ChangePasswordSection />

      {/* Delete Account (Danger Zone) */}
      <DeleteAccountSection />
    </div>
  );
};

export default AccountTab;
