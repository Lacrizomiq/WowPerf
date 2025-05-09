// components/Header/UserMenuOverlay.tsx
"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/providers/AuthContext";
import { useUserProfile } from "@/hooks/useUserProfile";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  LogIn,
  LogOut,
  UserPlus,
  Settings,
  ChevronsUpDown,
} from "lucide-react";

interface UserMenuOverlayProps {
  isExpanded: boolean;
}

const UserMenuOverlay: React.FC<UserMenuOverlayProps> = ({ isExpanded }) => {
  const router = useRouter();
  const { isAuthenticated, logout } = useAuth();
  const { profile } = useUserProfile();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Popover open={isOpen} onOpenChange={setIsOpen}>
      <PopoverTrigger asChild>
        <button
          className={`w-full px-2 py-2 flex items-center rounded-md transition-colors duration-200
                     text-slate-200 hover:bg-slate-800/80 hover:text-white
                     focus:outline-none focus:ring-1 focus:ring-purple-600
                     ${isExpanded ? "justify-start" : "justify-center"}`}
        >
          <Avatar className="h-8 w-8 bg-slate-700 flex-shrink-0">
            <AvatarImage src={""} alt={profile?.username || "User"} />
            <AvatarFallback className="bg-slate-700 text-white">
              {profile?.username
                ? profile.username.charAt(0).toUpperCase()
                : "U"}
            </AvatarFallback>
          </Avatar>
          {isExpanded && (
            <div className="ml-3 flex-1 min-w-0">
              <p className="text-sm font-medium truncate">
                {isAuthenticated ? profile?.username : "Guest"}
              </p>
              <p className="text-xs text-slate-400 truncate">
                {isAuthenticated ? profile?.email : "Not logged in"}
              </p>
            </div>
          )}
          {isExpanded && (
            <ChevronsUpDown className="ml-auto size-4 text-slate-400" />
          )}
        </button>
      </PopoverTrigger>
      <PopoverContent
        className="w-60 bg-slate-800 border-slate-700 text-slate-200 p-0 shadow-xl"
        align={isExpanded ? "start" : "center"}
        side={isExpanded ? "right" : "bottom"}
        sideOffset={8}
      >
        {isAuthenticated && profile ? (
          <>
            <div className="flex flex-col p-3 border-b border-slate-700">
              <p className="text-sm font-medium">{profile.username}</p>
              <p className="text-xs text-slate-400">{profile.email}</p>
            </div>
            <div className="p-1">
              <button
                className="w-full text-left flex items-center px-3 py-2 text-sm hover:bg-slate-700 rounded-sm"
                onClick={() => router.push("/profile")}
              >
                <Settings className="mr-2 h-4 w-4" />
                Account
              </button>
              <div className="h-px bg-slate-700 my-1" />
              <button
                onClick={() => {
                  logout();
                  setIsOpen(false);
                }}
                className="w-full text-left flex items-center px-3 py-2 text-sm hover:bg-slate-700 rounded-sm text-red-400 hover:text-red-300"
              >
                <LogOut className="mr-2 h-4 w-4" />
                Log out
              </button>
            </div>
          </>
        ) : (
          <div className="p-1">
            <button
              onClick={() => {
                router.push("/login");
                setIsOpen(false);
              }}
              className="w-full text-left flex items-center px-3 py-2 text-sm hover:bg-slate-700 rounded-sm"
            >
              <LogIn className="mr-2 h-4 w-4" />
              Login
            </button>
            <button
              onClick={() => {
                router.push("/signup");
                setIsOpen(false);
              }}
              className="w-full text-left flex items-center px-3 py-2 text-sm hover:bg-slate-700 rounded-sm"
            >
              <UserPlus className="mr-2 h-4 w-4" />
              Register
            </button>
          </div>
        )}
      </PopoverContent>
    </Popover>
  );
};

export default UserMenuOverlay;
