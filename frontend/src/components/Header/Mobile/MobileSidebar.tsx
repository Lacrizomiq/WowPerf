// components/Header/Mobile/MobileSidebar.tsx
"use client";

import React from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuth } from "@/providers/AuthContext";
import { useUserProfile } from "@/hooks/useUserProfile";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent } from "@/components/ui/sheet";
import {
  Home,
  BarChart3,
  PieChart,
  Trophy,
  ListChecks,
  LayoutDashboard,
  TrendingUp,
  LogIn,
  LogOut,
  UserPlus,
  Settings,
  ChevronRight,
} from "lucide-react";
import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import MobileSidebarSearchBar from "./MobileSidebarSearchBar";

interface MobileSidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

const MobileSidebar: React.FC<MobileSidebarProps> = ({ isOpen, onClose }) => {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated, logout } = useAuth();
  const { profile } = useUserProfile();

  const handleNavigation = (path: string) => {
    router.push(path);
    onClose();
  };

  const isActive = (path: string) => pathname === path;

  const getMenuItemClasses = (path: string) => {
    const active = isActive(path);
    return `w-full flex items-center justify-between py-3 px-4 rounded-md text-base
    ${
      active
        ? "bg-purple-900/30 text-white font-medium"
        : "text-slate-200 hover:bg-slate-800 hover:text-white"
    }`;
  };

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent
        side="left"
        className="w-full max-w-[300px] p-0 bg-[#1A1D21] border-r border-slate-800 text-slate-200"
        onCloseAutoFocus={(e) => e.preventDefault()}
        onOpenAutoFocus={(e) => e.preventDefault()}
      >
        <div className="flex items-center justify-between p-4 border-b border-slate-800">
          <div className="flex items-center gap-2">
            <div className="bg-purple-600 text-white p-1 rounded">
              <BarChart3 className="h-5 w-5" />
            </div>
            <div className="font-bold text-lg text-white">WoW Perf</div>
          </div>
        </div>

        <div className="overflow-y-auto h-full">
          {/* User Profile */}
          <div className="p-4 border-b border-slate-800">
            <div className="flex items-center">
              <Avatar className="h-10 w-10 bg-slate-700">
                <AvatarImage src={""} alt={profile?.username || "User"} />
                <AvatarFallback className="bg-slate-700 text-white">
                  {profile?.username
                    ? profile.username.charAt(0).toUpperCase()
                    : "U"}
                </AvatarFallback>
              </Avatar>
              <div className="ml-3">
                <p className="text-base font-medium">
                  {isAuthenticated ? profile?.username : "Guest"}
                </p>
                <p className="text-sm text-slate-400">
                  {isAuthenticated ? profile?.email : "Not logged in"}
                </p>
              </div>
            </div>

            <div className="mt-3 space-y-2">
              {isAuthenticated ? (
                <>
                  <Button
                    variant="outline"
                    className="w-full justify-start border-slate-700 hover:bg-slate-800"
                    onClick={() => handleNavigation("/profile")}
                  >
                    <Settings className="h-4 w-4 mr-2" />
                    Account
                  </Button>
                  <Button
                    variant="outline"
                    className="w-full justify-start border-slate-700 text-red-400 hover:text-red-300 hover:bg-slate-800"
                    onClick={() => {
                      logout();
                      onClose();
                    }}
                  >
                    <LogOut className="h-4 w-4 mr-2" />
                    Log out
                  </Button>
                </>
              ) : (
                <>
                  <Button
                    variant="outline"
                    className="w-full justify-start border-slate-700 hover:bg-slate-800"
                    onClick={() => handleNavigation("/login")}
                  >
                    <LogIn className="h-4 w-4 mr-2" />
                    Login
                  </Button>
                  <Button
                    variant="outline"
                    className="w-full justify-start border-slate-700 hover:bg-slate-800"
                    onClick={() => handleNavigation("/signup")}
                  >
                    <UserPlus className="h-4 w-4 mr-2" />
                    Register
                  </Button>
                </>
              )}
            </div>
          </div>

          {/* Search Character - Using the dedicated component */}
          <MobileSidebarSearchBar onClose={onClose} />

          {/* Navigation */}
          <div className="p-4">
            <h3 className="text-sm font-medium uppercase tracking-wider text-slate-400 mb-2 px-1">
              Navigation
            </h3>
            <nav className="space-y-1">
              <button
                onClick={() => handleNavigation("/")}
                className={getMenuItemClasses("/")}
              >
                <span className="flex items-center">
                  <Home className="h-5 w-5 mr-3" />
                  Home
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/performance-analysis")}
                className={getMenuItemClasses("/performance-analysis")}
              >
                <span className="flex items-center">
                  <BarChart3 className="h-5 w-5 mr-3" />
                  Performance Analysis
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/builds")}
                className={getMenuItemClasses("/builds")}
              >
                <span className="flex items-center">
                  <ListChecks className="h-5 w-5 mr-3" />
                  Builds
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/statistics")}
                className={getMenuItemClasses("/statistics")}
              >
                <span className="flex items-center">
                  <PieChart className="h-5 w-5 mr-3" />
                  Statistics
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/leaderboards")}
                className={getMenuItemClasses("/leaderboards")}
              >
                <span className="flex items-center">
                  <Trophy className="h-5 w-5 mr-3" />
                  Leaderboards
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/character-progress")}
                className={getMenuItemClasses("/character-progress")}
              >
                <span className="flex items-center">
                  <TrendingUp className="h-5 w-5 mr-3" />
                  Character Progress
                  <Badge
                    variant="outline"
                    className="ml-2 text-xs border-purple-600 text-purple-400"
                  >
                    Soon
                  </Badge>
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/dashboard")}
                className={getMenuItemClasses("/dashboard")}
              >
                <span className="flex items-center">
                  <LayoutDashboard className="h-5 w-5 mr-3" />
                  Dashboard
                  <Badge
                    variant="outline"
                    className="ml-2 text-xs border-purple-600 text-purple-400"
                  >
                    Soon
                  </Badge>
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>
            </nav>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
};

export default MobileSidebar;
