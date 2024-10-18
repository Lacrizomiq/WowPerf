"use client";

import React, { useState } from "react";

interface PersonalInfoProps {
  profile: any;
  onUpdate: (email: string) => void;
  isUpdating: boolean;
}

const PersonalInfo: React.FC<PersonalInfoProps> = ({
  profile,
  onUpdate,
  isUpdating,
}) => {
  const [email, setEmail] = useState(profile?.email || "");

  return (
    <>
      <section className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h2 className="text-2xl font-bold mb-4 text-gray-800 dark:text-gray-200">
          Personal Information
        </h2>

        <div className="flex items-center">
          <p className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            <span className="font-bold text-lg">Username : </span>
            {profile?.username}
          </p>
        </div>
        <div>
          <p className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            <span className="font-bold text-lg">Email : </span>
            {profile?.email}
          </p>
        </div>
      </section>

      <section className="bg-white dark:bg-gray-800 shadow rounded-lg p-6 mt-4">
        <h2 className="text-2xl font-bold mb-4 text-gray-800 dark:text-gray-200">
          Link your WowPerf account to your battle.net account
        </h2>
        <div className="flex items-center"></div>
      </section>
    </>
  );
};

export default PersonalInfo;
