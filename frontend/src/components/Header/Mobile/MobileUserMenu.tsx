// components/Sidebar/MobileUserMenu.tsx
import React from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/providers/AuthContext";
import { useUserProfile } from "@/hooks/useUserProfile";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { LogIn, LogOut, UserPlus, BadgeCheck } from "lucide-react";

const MobileUserMenu = () => {
  const router = useRouter();
  const { isAuthenticated, logout } = useAuth();
  const { profile } = useUserProfile();

  return (
    <div className="p-4 border-t border-sidebar-border">
      <div className="flex items-center space-x-4 mb-4">
        <Avatar className="h-10 w-10">
          <AvatarImage src="" alt="User" />
          <AvatarFallback>U</AvatarFallback>
        </Avatar>
        <div className="flex flex-col">
          <span className="text-sm font-semibold">
            {isAuthenticated ? profile?.username : "Guest"}
          </span>
          <span className="text-xs text-muted-foreground">
            {isAuthenticated ? profile?.email : "Not logged in"}
          </span>
        </div>
      </div>

      <div className="space-y-2">
        {isAuthenticated ? (
          <>
            <button
              onClick={() => router.push("/profile")}
              className="w-full flex items-center px-2 py-2 text-sm hover:bg-slate-800 rounded-md"
            >
              <BadgeCheck className="mr-2 h-4 w-4" />
              Account
            </button>
            <button
              onClick={logout}
              className="w-full flex items-center px-2 py-2 text-sm hover:bg-slate-800 rounded-md text-red-500 hover:text-red-600"
            >
              <LogOut className="mr-2 h-4 w-4" />
              Log out
            </button>
          </>
        ) : (
          <>
            <button
              onClick={() => router.push("/login")}
              className="w-full flex items-center px-2 py-2 text-sm hover:bg-slate-800 rounded-md"
            >
              <LogIn className="mr-2 h-4 w-4" />
              Login
            </button>
            <button
              onClick={() => router.push("/signup")}
              className="w-full flex items-center px-2 py-2 text-sm hover:bg-slate-800 rounded-md"
            >
              <UserPlus className="mr-2 h-4 w-4" />
              Register
            </button>
          </>
        )}
      </div>
    </div>
  );
};

export default MobileUserMenu;
