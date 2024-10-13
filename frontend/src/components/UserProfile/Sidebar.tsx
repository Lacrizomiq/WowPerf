"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";

const Sidebar: React.FC = () => {
  const pathname = usePathname();

  const navItems = [
    { name: "Profile", href: "/profile" },
    { name: "Security", href: "/profile/security" },
    { name: "Notifications", href: "/profile/notifications" },
    { name: "Linked Accounts", href: "/profile/linked-accounts" },
    { name: "Privacy", href: "/profile/privacy" },
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
                <Link
                  href={item.href}
                  className={`block px-4 py-2 rounded-md transition-colors ${
                    pathname === item.href
                      ? "bg-blue-500 text-white"
                      : "text-white dark:text-gray-300 hover:bg-blue-900 dark:hover:bg-gray-700"
                  }`}
                >
                  {item.name}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      </div>
    </aside>
  );
};

export default Sidebar;
