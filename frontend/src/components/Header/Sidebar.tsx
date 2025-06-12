// components/Header/Sidebar.tsx
"use client";

import React from "react";
import { useRouter, usePathname } from "next/navigation";
import { SidebarMenu, SidebarMenuItem } from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import {
  Home,
  BarChart3,
  PieChart,
  Trophy,
  ListChecks,
  LayoutDashboard,
  TrendingUp,
  ChevronLeft,
  ChevronRight,
  LogIn,
  UserPlus,
} from "lucide-react";
import Link from "next/link";
import UserMenuOverlay from "./UserMenuOverlay";
import SidebarSearchBar from "@/components/Header/SidebarSearchBar";
import { Badge } from "@/components/ui/badge";
import { useAuth } from "@/providers/AuthContext";

interface AppSidebarProps {
  isExpanded: boolean;
  onToggleSidebar: () => void;
}

const AppSidebar: React.FC<AppSidebarProps> = ({
  isExpanded,
  onToggleSidebar,
}) => {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated } = useAuth();

  const isActive = (path: string) => pathname === path;

  // Styles améliorés pour les éléments de menu
  const getMenuItemClasses = (path: string) => {
    const active = isActive(path);
    return `w-full flex items-center py-2 px-3 rounded-md text-sm transition-colors duration-200 
    ${
      active
        ? "bg-purple-900/30 text-white"
        : "text-slate-200 hover:bg-slate-800/80 hover:text-white"
    }
    ${!isExpanded ? "justify-center" : ""}`;
  };

  // Badge "Soon" pour les fonctionnalités à venir
  const soonBadge = (
    <Badge
      variant="outline"
      className="ml-2 text-xs py-0 h-5 border-purple-600 text-purple-400"
    >
      Soon
    </Badge>
  );

  return (
    <div className="flex flex-col h-full bg-[#1A1D21] text-slate-200 border-r border-slate-800">
      {/* Header - Complètement différent selon l'état expanded/collapsed */}
      {isExpanded ? (
        <div className="p-4 flex items-center justify-between h-16 border-b border-slate-800">
          <Link href="/" className="flex items-center gap-2">
            <div className="bg-purple-600 text-white p-1 rounded flex-shrink-0">
              <BarChart3 className="h-5 w-5" />
            </div>
            <span className="font-bold text-lg">WoW Perf</span>
          </Link>
          <Button
            variant="ghost"
            size="icon"
            onClick={onToggleSidebar}
            className="bg-slate-800/80 text-white hover:bg-purple-700/80"
            aria-label="Collapse sidebar"
          >
            <ChevronLeft className="h-5 w-5" />
          </Button>
        </div>
      ) : (
        <div className="flex flex-col items-center py-3 border-b border-slate-800">
          <Button
            variant="ghost"
            size="icon"
            onClick={onToggleSidebar}
            className="bg-purple-600 text-white hover:bg-purple-700 mb-2"
            aria-label="Expand sidebar"
          >
            <ChevronRight className="h-5 w-5" />
          </Button>
          <Link href="/" className="bg-purple-600 text-white p-1 rounded">
            <BarChart3 className="h-5 w-5" />
          </Link>
        </div>
      )}

      {/* Content - Scrollable */}
      <div className="flex-1 overflow-y-auto">
        {/* User Profile Section */}
        <div className="px-4 py-2">
          <UserMenuOverlay isExpanded={isExpanded} />
          {isExpanded && !isAuthenticated && (
            <div className="mt-2 space-y-1.5">
              <Button
                variant="outline"
                size="sm"
                className="w-full justify-start border-slate-700 bg-slate-800/20 hover:bg-slate-800/90 text-slate-200 text-sm"
                onClick={() => router.push("/login")}
              >
                <LogIn className="h-4 w-4 mr-2" />
                Login
              </Button>
              <Button
                variant="outline"
                size="sm"
                className="w-full justify-start border-slate-700 bg-slate-800/20 hover:bg-slate-800/90 text-slate-200 text-sm"
                onClick={() => router.push("/signup")}
              >
                <UserPlus className="h-4 w-4 mr-2" />
                Register
              </Button>
            </div>
          )}
        </div>

        {/* Search Section - Moved to SidebarSearchBar component */}
        <SidebarSearchBar
          isExpanded={isExpanded}
          onToggleSidebar={onToggleSidebar}
        />

        {/* Navigation Section */}
        <div className="px-4 py-3">
          {isExpanded && (
            <div className="text-xs font-medium uppercase tracking-wider text-slate-400 mb-2 px-3">
              Navigation
            </div>
          )}
          <SidebarMenu>
            <SidebarMenuItem>
              <Link href="/" className={getMenuItemClasses("/")}>
                {isExpanded ? (
                  <Home className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <Home className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Home</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/performance-analysis"
                className={getMenuItemClasses("/performance-analysis")}
              >
                {isExpanded ? (
                  <BarChart3 className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <BarChart3 className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Performance Analysis</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link href="/builds" className={getMenuItemClasses("/builds")}>
                {isExpanded ? (
                  <ListChecks className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <ListChecks className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Builds</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/statistics"
                className={getMenuItemClasses("/statistics")}
              >
                {isExpanded ? (
                  <PieChart className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <PieChart className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Statistics</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/leaderboards"
                className={getMenuItemClasses("/leaderboards")}
              >
                {isExpanded ? (
                  <Trophy className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <Trophy className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Leaderboards</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/character-progress"
                className={getMenuItemClasses("/character-progress")}
              >
                {isExpanded ? (
                  <>
                    <TrendingUp className="h-4 w-4 mr-3" />
                    <div className="flex-1 flex justify-between items-center">
                      <span>Character Progress</span>
                      {soonBadge}
                    </div>
                  </>
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <TrendingUp className="w-5 h-5" />
                  </div>
                )}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/dashboard"
                className={getMenuItemClasses("/dashboard")}
              >
                {isExpanded ? (
                  <>
                    <LayoutDashboard className="h-4 w-4 mr-3" />
                    <div className="flex-1 flex justify-between items-center">
                      <span>Dashboard</span>
                      {soonBadge}
                    </div>
                  </>
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <LayoutDashboard className="w-5 h-5" />
                  </div>
                )}
              </Link>
            </SidebarMenuItem>
          </SidebarMenu>
        </div>
      </div>
    </div>
  );
};

export default AppSidebar;
