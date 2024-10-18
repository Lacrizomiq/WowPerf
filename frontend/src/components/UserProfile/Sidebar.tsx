"use client";

import React from "react";
import { usePathname } from "next/navigation";

interface SidebarProps {
  onSectionChange: (section: string) => void;
}

const Sidebar: React.FC<SidebarProps> = ({ onSectionChange }) => {
  const pathname = usePathname();

  const navItems = [
    { name: "Profile", section: "profile" },
    { name: "Change Username", section: "change-username" },
    { name: "Change Password", section: "change-password" },
    { name: "Change Email", section: "change-email" },
    { name: "Delete Account", section: "delete-account" },
  ];

  return (
    <aside className="w-64 bg-blue-700 dark:bg-gray-800 shadow-md">
      <div className="p-6">
        <h2 className="text-xl font-semibold mb-6 text-white dark:text-gray-200">
          Account Settings
        </h2>
        <nav>
          <ul className="space-y-2">
            {navItems.map((item) => (
              <li key={item.name}>
                <button
                  onClick={() => onSectionChange(item.section)}
                  className={`block w-full text-left px-4 py-2 rounded-md transition-colors ${
                    pathname === `/profile/${item.section}`
                      ? "bg-blue-500 text-white"
                      : "text-white dark:text-gray-300 hover:bg-blue-900 dark:hover:bg-gray-700"
                  }`}
                >
                  {item.name}
                </button>
              </li>
            ))}
          </ul>
        </nav>
      </div>
    </aside>
  );
};

export default Sidebar;
