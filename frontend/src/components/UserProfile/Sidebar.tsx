import React from "react";
import { usePathname, useRouter } from "next/navigation";
import { Camera, Edit2, Key, Mail, Trash2 } from "lucide-react";

const Sidebar: React.FC = () => {
  const pathname = usePathname();
  const router = useRouter();
  const navItems = [
    { name: "Profile", icon: Camera, section: "" },
    { name: "Change Username", icon: Edit2, section: "update-username" },
    { name: "Change Password", icon: Key, section: "update-password" },
    { name: "Change Email", icon: Mail, section: "update-email" },
    { name: "Delete Account", icon: Trash2, section: "delete-account" },
  ];

  return (
    <aside className="w-64 bg-[#2d3748] dark:bg-gray-800 shadow-md">
      <div className="p-6">
        <h2 className="text-xl font-semibold mb-6 text-white dark:text-gray-200">
          Account Settings
        </h2>
        <nav>
          <ul className="space-y-2">
            {navItems.map((item) => (
              <li key={item.name}>
                <button
                  onClick={() => router.push(`/profile/${item.section}`)}
                  className={`flex items-center w-full text-left px-4 py-2 rounded-md transition-colors ${
                    pathname === `/profile/${item.section}`
                      ? "bg-blue-500 text-white"
                      : "text-white dark:text-gray-300 hover:bg-blue-900 dark:hover:bg-gray-700"
                  }`}
                >
                  <item.icon className="mr-2" size={20} />
                  <span>{item.name}</span>
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
