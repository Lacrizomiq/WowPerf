"use client";

import React, { useState } from "react";
import Image from "next/image";

interface PersonalInfoProps {
  profile: any;
}

const PersonalInfo: React.FC<PersonalInfoProps> = ({ profile }) => {
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
        <div className="flex items-center">
          <button className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded transition duration-200 flex items-center">
            <Image
              src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
              alt="Battle.net icon"
              width={24}
              height={24}
              className="mr-2"
            />
            Connect to your Blizzard account
          </button>
        </div>
      </section>
    </>
  );
};

export default PersonalInfo;
