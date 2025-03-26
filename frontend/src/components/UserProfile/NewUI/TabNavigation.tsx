// TabNavigation.tsx
import React from "react";

interface TabNavigationProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

const TabNavigation: React.FC<TabNavigationProps> = ({
  activeTab,
  setActiveTab,
}) => {
  return (
    <div className="flex gap-1 border-b border-gray-800 mb-6">
      <div
        className={`flex items-center gap-2 px-6 py-4 cursor-pointer border-b-2 font-semibold ${
          activeTab === "overview"
            ? "border-blue-500 text-blue-500"
            : "border-transparent hover:text-blue-400 hover:border-blue-400"
        }`}
        onClick={() => setActiveTab("overview")}
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
          <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
          <polyline points="9 22 9 12 15 12 15 22"></polyline>
        </svg>
        Overview
      </div>

      <div
        className={`flex items-center gap-2 px-6 py-4 cursor-pointer border-b-2 font-semibold ${
          activeTab === "characters"
            ? "border-blue-500 text-blue-500"
            : "border-transparent hover:text-blue-400 hover:border-blue-400"
        }`}
        onClick={() => setActiveTab("characters")}
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
          <path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48"></path>
        </svg>
        Characters
      </div>

      <div
        className={`flex items-center gap-2 px-6 py-4 cursor-pointer border-b-2 font-semibold ${
          activeTab === "account"
            ? "border-blue-500 text-blue-500"
            : "border-transparent hover:text-blue-400 hover:border-blue-400"
        }`}
        onClick={() => setActiveTab("account")}
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
          <circle cx="12" cy="8" r="5"></circle>
          <path d="M20 21a8 8 0 0 0-16 0"></path>
        </svg>
        Account
      </div>

      <div
        className={`flex items-center gap-2 px-6 py-4 cursor-pointer border-b-2 font-semibold ${
          activeTab === "security"
            ? "border-blue-500 text-blue-500"
            : "border-transparent hover:text-blue-400 hover:border-blue-400"
        }`}
        onClick={() => setActiveTab("security")}
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
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"></path>
        </svg>
        Security
      </div>

      <div
        className={`flex items-center gap-2 px-6 py-4 cursor-pointer border-b-2 font-semibold ${
          activeTab === "connections"
            ? "border-blue-500 text-blue-500"
            : "border-transparent hover:text-blue-400 hover:border-blue-400"
        }`}
        onClick={() => setActiveTab("connections")}
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
          <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
          <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
        </svg>
        Connections
      </div>
    </div>
  );
};

export default TabNavigation;
